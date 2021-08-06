// Package schema contains the database schema, migrations and seeding data.
package schema

import (
	"context"
	_ "embed" // Calls init function.

	"github.com/ardanlabs/darwin"
	"github.com/ardanlabs/service/foundation/database"
	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

var (
	//go:embed sql/schema.sql
	schemaDoc string
)

// Migrate attempts to bring the schema for db up to date with the migrations
// defined in this package.
func Migrate(ctx context.Context, db *sqlx.DB) error {
	if err := database.StatusCheck(ctx, db); err != nil {
		return errors.Wrap(err, "status check database")
	}

	driver, err := darwin.NewGenericDriver(db.DB, darwin.PostgresDialect{})
	if err != nil {
		return errors.Wrap(err, "construct darwin driver")
	}

	d := darwin.New(driver, darwin.ParseMigrations(schemaDoc))
	return d.Migrate()
}
