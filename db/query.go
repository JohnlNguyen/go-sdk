package db

import (
	"context"
	"database/sql"
	"reflect"

	"github.com/blend/go-sdk/ex"
)

// --------------------------------------------------------------------------------
// Query Result
// --------------------------------------------------------------------------------

// Query is the intermediate result of a query.
type Query struct {
	Context       context.Context
	Statement     string
	CachedPlanKey string
	Args          []interface{}

	Rows *sql.Rows
	Err  error

	Conn       *Connection
	Invocation *Invocation
	Tx         *sql.Tx
}

// Execute runs a given query, yielding the raw results.
func (q *Query) Execute() (rows *sql.Rows, err error) {
	defer func() { err = q.finish(recover(), err) }()
	rows, err = q.query()
	return
}

// Any returns if there are any results for the query.
func (q *Query) Any() (found bool, err error) {
	defer func() { err = q.finish(recover(), err) }()

	q.Rows, q.Err = q.query()
	if q.Err != nil {
		err = q.Err
		return
	}
	defer func() { err = ex.Nest(err, q.Rows.Close()) }()

	found = q.Rows.Next()
	return
}

// None returns if there are no results for the query.
func (q *Query) None() (notFound bool, err error) {
	defer func() { err = q.finish(recover(), err) }()

	q.Rows, q.Err = q.query()
	if q.Err != nil {
		err = q.Err
		return
	}
	defer func() { err = ex.Nest(err, Error(q.Rows.Close())) }()
	notFound = !q.Rows.Next()
	return
}

// Scan writes the results to a given set of local variables.
// It returns if the query produced a row, and returns `ErrTooManyRows` if there
// are multiple row results.
func (q *Query) Scan(args ...interface{}) (found bool, err error) {
	defer func() { err = q.finish(recover(), err) }()

	q.Rows, q.Err = q.query()
	if q.Err != nil {
		err = q.Err
		return
	}
	defer func() { err = ex.Nest(err, Error(q.Rows.Close())) }()

	if q.Rows.Next() {
		found = true
		if err = q.Rows.Scan(args...); err != nil {
			err = Error(err)
			return
		}
	}
	if q.Rows.Next() {
		err = Error(ErrTooManyRows)
	}

	return
}

// ReadInto writes the query result to a single object via. reflection mapping. If there is more than one result, the first
// result is mapped to to object, and ErrTooManyRows is returned. Unlike Out(), if a field on the stuct is not present
// or is an "empty" value in the result set, ReadInto() clears the field on object it is populating. In short ReadInto() maps the
// output of your query into object as "exactly" as possible. Where you can, prefer ReadInto() over Out().
// If you need multiple queries to fill up your object struct, you should be using Out().
func (q *Query) ReadInto(object interface{}) (found bool, err error) {
	return q.populate(object, true)
}

// Out writes the query result to a single object via. reflection mapping. If there is more than one result, the first
// result is mapped to to object, and ErrTooManyRows is returned. Out() is different than ReadInto(), in that is DOES NOT change
// struct fields on object that are empty in the result set. If a result field is null the value that was present in
// object is maintained.
func (q *Query) Out(object interface{}) (found bool, err error) {
	return q.populate(object, false)
}

// OutMany writes the query results to a slice of objects.
func (q *Query) OutMany(collection interface{}) (err error) {
	defer func() { err = q.finish(recover(), err) }()

	q.Rows, q.Err = q.query()
	if q.Err != nil {
		err = q.Err
		return
	}
	defer func() { err = ex.Nest(err, q.Rows.Close()) }()

	sliceType := ReflectType(collection)
	if sliceType.Kind() != reflect.Slice {
		err = Error(ErrCollectionNotSlice)
		return
	}

	sliceInnerType := ReflectSliceType(collection)
	collectionValue := ReflectValue(collection)
	v := makeNew(sliceInnerType)
	meta := CachedColumnCollectionFromType(newColumnCacheKey(sliceInnerType), sliceInnerType)

	isPopulatable := IsPopulatable(v)

	var didSetRows bool
	for q.Rows.Next() {
		newObj := makeNew(sliceInnerType)
		if isPopulatable {
			err = AsPopulatable(newObj).Populate(q.Rows)
		} else {
			err = PopulateByName(newObj, q.Rows, meta, true)
		}
		if err != nil {
			return
		}

		newObjValue := ReflectValue(newObj)
		collectionValue.Set(reflect.Append(collectionValue, newObjValue))
		didSetRows = true
	}

	// this initializes the slice if we didn't add elements to it.
	if !didSetRows {
		collectionValue.Set(reflect.MakeSlice(sliceType, 0, 0))
	}
	return
}

// Each executes the consumer for each result of the query (one to many).
func (q *Query) Each(consumer RowsConsumer) (err error) {
	defer func() { err = q.finish(recover(), err) }()

	q.Rows, q.Err = q.query()
	if q.Err != nil {
		err = q.Err
		return
	}
	defer func() { err = ex.Nest(err, Error(q.Rows.Close())) }()

	for q.Rows.Next() {
		if err = consumer(q.Rows); err != nil {
			err = Error(err)
			return
		}
	}
	return
}

// First executes the consumer for the first result of a query.
// It returns `ErrTooManyRows` if more than one result is returned.
func (q *Query) First(consumer RowsConsumer) (err error) {
	defer func() { err = q.finish(recover(), err) }()

	q.Rows, q.Err = q.query()
	if q.Err != nil {
		err = q.Err
		return
	}
	defer func() { err = ex.Nest(err, Error(q.Rows.Close())) }()

	if q.Rows.Next() {
		if err = consumer(q.Rows); err != nil {
			return
		}
	}
	return
}

// --------------------------------------------------------------------------------
// helpers
// --------------------------------------------------------------------------------

func (q *Query) populate(object interface{}, resetRowEmpty bool) (found bool, err error) {
	defer func() { err = q.finish(recover(), err) }()
	q.Rows, q.Err = q.query()
	if q.Err != nil {
		err = q.Err
		return
	}
	defer func() { err = ex.Nest(err, Error(q.Rows.Close())) }()

	sliceType := ReflectType(object)
	if sliceType.Kind() != reflect.Struct {
		err = Error(ErrDestinationNotStruct)
		return
	}

	columnMeta := CachedColumnCollectionFromInstance(object)
	if q.Rows.Next() {
		found = true
		if populatable, ok := object.(Populatable); ok {
			err = populatable.Populate(q.Rows)
		} else {
			err = PopulateByName(object, q.Rows, columnMeta, resetRowEmpty)
		}
		if err != nil {
			return
		}
	} else if _, ok := object.(Populatable); !ok {
		PopulateEmpty(object, columnMeta)
	}

	if q.Rows.Next() {
		err = Error(ErrTooManyRows)
	}

	return
}

func (q *Query) query() (rows *sql.Rows, err error) {
	if q.Err != nil {
		err = q.Err
		return
	}

	stmt, stmtErr := q.Invocation.Prepare(q.Statement)
	if stmtErr != nil {
		err = Error(stmtErr)
		return
	}
	defer func() { err = q.Invocation.CloseStatement(stmt, err) }()

	rows, err = stmt.QueryContext(q.Context, q.Args...)
	if err != nil && !ex.Is(err, sql.ErrNoRows) {
		err = Error(err)
	}
	return
}

func (q *Query) finish(r interface{}, err error) error {
	return q.Invocation.Finish(q.Statement, r, err)
}
