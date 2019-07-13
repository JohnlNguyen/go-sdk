package migration

import (
	"context"
	"database/sql"

	"go-sdk/db"
	"go-sdk/exception"
)

// Group creates a new GroupedActions from a given list of actionable.
func Group(actions ...Actionable) GroupedActions {
	return GroupedActions(actions)
}

// GroupedActions is an atomic series of migration actions.
// It uses transactions to apply these actions as an atomic unit.
type GroupedActions []Actionable

// Action runs the groups actions within a transaction.
func (ga GroupedActions) Action(ctx context.Context, c *db.Connection) (err error) {
	var tx *sql.Tx
	tx, err = c.Begin()
	if err != nil {
		return
	}
	// commit or rollback the transaction.
	defer func() {
		if err != nil {
			if txErr := tx.Rollback(); txErr != nil {
				err = exception.Nest(err, txErr)
			}
		} else {
			if txErr := tx.Commit(); txErr != nil {
				err = exception.Nest(err, txErr)
			}
		}
	}()

	for _, a := range ga {
		err = a.Action(ctx, c, tx)
		if err != nil {
			return
		}
	}

	return
}
