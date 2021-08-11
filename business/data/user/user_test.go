package user_test

import (
	"testing"

	"github.com/ardanlabs/service/business/data/tests"
)

var dbc = tests.DBContainer{
	Image: "postgres:13-alpine",
	Port:  "5432",
	Args:  []string{"-e", "POSTGRES_PASSWORD=postgres"},
}

func TestUser(t *testing.T) {
	_, _, teardown := tests.NewUnit(t, dbc)
	t.Cleanup(teardown)
}
