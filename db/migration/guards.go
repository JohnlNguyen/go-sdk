package migration

import (
	"context"
	"database/sql"
	"fmt"

	"go-sdk/db"
	"go-sdk/stringutil"
)

const (
	verbCreate = "create"
	verbAlter  = "alter"
	verbRun    = "run"

	nounColumn     = "column"
	nounTable      = "table"
	nounIndex      = "index"
	nounConstraint = "constraint"
	nounRole       = "role"

	adverbAlways    = "always"
	adverbExists    = "exists"
	adverbNotExists = "not exists"
)

// GuardFunc is a control for migration steps.
type GuardFunc func(context.Context, *db.Connection, *sql.Tx, Action) error

// --------------------------------------------------------------------------------
// Guards
// --------------------------------------------------------------------------------

// Guard returns a function that determines if a step in a group should run.
func Guard(description string, predicate func(c *db.Connection, tx *sql.Tx) (bool, error)) GuardFunc {
	return func(ctx context.Context, c *db.Connection, tx *sql.Tx, step Action) error {
		proceed, err := predicate(c, tx)
		if err != nil {
			if suite := GetContextSuite(ctx); suite != nil {
				return suite.Error(WithLabel(ctx, description), err)
			}
			return err
		}

		if !proceed {
			if suite := GetContextSuite(ctx); suite != nil {
				suite.Skipf(ctx, description)
			}
			return nil
		}

		err = step(ctx, c, tx)
		if err != nil {
			if suite := GetContextSuite(ctx); suite != nil {
				return suite.Error(WithLabel(ctx, description), err)
			}
			return err
		}
		if suite := GetContextSuite(ctx); suite != nil {
			suite.Applyf(ctx, description)
		}
		return nil
	}
}

// Always always runs a step.
func Always() GuardFunc {
	return Guard("always run", func(_ *db.Connection, _ *sql.Tx) (bool, error) { return true, nil })
}

// IfExists only runs the statement if the given item exists.
func IfExists(statement string) GuardFunc {
	return Guard("if exists run", func(c *db.Connection, tx *sql.Tx) (bool, error) {
		return Exists(c, tx, statement)
	})
}

// IfNotExists only runs the statement if the given item doesn't exist.
func IfNotExists(statement string) GuardFunc {
	return Guard("if not exists run", func(c *db.Connection, tx *sql.Tx) (bool, error) {
		return NotExists(c, tx, statement)
	})
}

// ColumnNotExistsWithPredicate creates a table on the given connection if it does not exist.
func ColumnNotExistsWithPredicate(predicate Predicate2, tableName, columnName string) GuardFunc {
	return Guard(fmt.Sprintf("create column `%s` on `%s`", columnName, tableName), func(c *db.Connection, tx *sql.Tx) (bool, error) {
		return Not(predicate(c, tx, tableName, columnName))
	})
}

// ConstraintNotExistsWithPredicate creates a table on the given connection if it does not exist.
func ConstraintNotExistsWithPredicate(predicate Predicate, constraintName string) GuardFunc {
	return Guard(fmt.Sprintf("create constraint `%s`", constraintName), func(c *db.Connection, tx *sql.Tx) (bool, error) {
		return Not(predicate(c, tx, constraintName))
	})
}

// TableNotExistsWithPredicate creates a table on the given connection if it does not exist.
func TableNotExistsWithPredicate(predicate Predicate, tableName string) GuardFunc {
	return Guard(fmt.Sprintf("create table `%s`", tableName), func(c *db.Connection, tx *sql.Tx) (bool, error) {
		return Not(predicate(c, tx, tableName))
	})
}

// IndexNotExistsWithPredicate creates a index on the given connection if it does not exist.
func IndexNotExistsWithPredicate(predicate Predicate2, tableName, indexName string) GuardFunc {
	return Guard(fmt.Sprintf("create index `%s` on `%s`", indexName, tableName), func(c *db.Connection, tx *sql.Tx) (bool, error) {
		return Not(predicate(c, tx, tableName, indexName))
	})
}

// RoleNotExistsWithPredicate creates a new role if it doesn't exist.
func RoleNotExistsWithPredicate(predicate Predicate, roleName string) GuardFunc {
	return Guard(fmt.Sprintf("create role `%s`", roleName), func(c *db.Connection, tx *sql.Tx) (bool, error) {
		return Not(predicate(c, tx, roleName))
	})
}

// ColumnExistsWithPredicate alters an existing column, erroring if it doesn't exist
func ColumnExistsWithPredicate(predicate Predicate2, tableName, columnName string) GuardFunc {
	return Guard(fmt.Sprintf("alter column `%s` on `%s`", columnName, tableName), func(c *db.Connection, tx *sql.Tx) (bool, error) {
		return predicate(c, tx, tableName, columnName)
	})
}

// ConstraintExistsWithPredicate alters an existing constraint, erroring if it doesn't exist
func ConstraintExistsWithPredicate(predicate Predicate, constraintName string) GuardFunc {
	return Guard(fmt.Sprintf("alter constraint `%s`", constraintName), func(c *db.Connection, tx *sql.Tx) (bool, error) {
		return predicate(c, tx, constraintName)
	})
}

// TableExistsWithPredicate alters an existing table, erroring if it doesn't exist
func TableExistsWithPredicate(predicate Predicate, tableName string) GuardFunc {
	return Guard(fmt.Sprintf("alter table `%s`", tableName), func(c *db.Connection, tx *sql.Tx) (bool, error) {
		return predicate(c, tx, tableName)
	})
}

// IndexExistsWithPredicate alters an existing index, erroring if it doesn't exist
func IndexExistsWithPredicate(predicate Predicate2, tableName, indexName string) GuardFunc {
	return Guard(fmt.Sprintf("alter index `%s` on `%s`", indexName, tableName), func(c *db.Connection, tx *sql.Tx) (bool, error) {
		return predicate(c, tx, tableName, indexName)
	})
}

// RoleExistsWithPredicate alters an existing role in the db
func RoleExistsWithPredicate(predicate Predicate, roleName string) GuardFunc {
	return Guard(fmt.Sprintf("alter role `%s`", roleName), func(c *db.Connection, tx *sql.Tx) (bool, error) {
		return predicate(c, tx, roleName)
	})
}

// Predicate is a function that evaluates based on a string param.
type Predicate func(*db.Connection, *sql.Tx, string) (bool, error)

// Predicate2 is a function that evaluates based on two string params.
type Predicate2 func(*db.Connection, *sql.Tx, string, string) (bool, error)

// Not inverts the output of a predicate.
func Not(proceed bool, err error) (bool, error) {
	return !proceed, err
}

// --------------------------------------------------------------------------------
// Guard Helpers
// --------------------------------------------------------------------------------

// Exists returns if a statement has results.
func Exists(c *db.Connection, tx *sql.Tx, selectStatement string) (bool, error) {
	if !stringutil.HasPrefixCaseless(selectStatement, "select") {
		return false, fmt.Errorf("statement must be a `SELECT`")
	}
	return c.QueryInTx(selectStatement, tx).Any()
}

// NotExists returns if a statement doesnt have results.
func NotExists(c *db.Connection, tx *sql.Tx, selectStatement string) (bool, error) {
	if !stringutil.HasPrefixCaseless(selectStatement, "select") {
		return false, fmt.Errorf("statement must be a `SELECT`")
	}
	return c.QueryInTx(selectStatement, tx).None()
}
