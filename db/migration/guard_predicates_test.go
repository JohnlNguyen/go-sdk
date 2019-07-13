package migration

import (
	"context"
	"database/sql"
	"fmt"
	"testing"

	"go-sdk/assert"
	"go-sdk/db"
)

func TestGuard(t *testing.T) {
	assert := assert.New(t)
	tx, err := db.Default().Begin()
	assert.Nil(err)
	defer tx.Rollback()

	tableName := randomName()
	err = createTestTable(tableName, tx)
	assert.Nil(err)

	err = insertTestValue(tableName, 4, "test", tx)
	assert.Nil(err)

	var didRun bool
	action := Actions(func(ctx context.Context, c *db.Connection, itx *sql.Tx) error {
		didRun = true
		return nil
	})

	err = Guard("test", func(c *db.Connection, itx *sql.Tx) (bool, error) {
		return c.QueryInTx(fmt.Sprintf("select * from %s", tableName), itx).Any()
	})(
		context.Background(),
		db.Default(),
		tx,
		action,
	)
	assert.Nil(err)
	assert.True(didRun)
}
