package main

import (
	"log"
	"os"
	"testing"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"github.com/leftovers-2025/kds_backend/internal/kds/datasource"
	migrate "github.com/rubenv/sql-migrate"
)

var globalDB *sqlx.DB
var globalDialect = "mysql"

func setupTest() {
	globalDB = datasource.GetMySqlConnection()
	globalDialect = "mysql"

	migrations := &migrate.FileMigrationSource{
		Dir: "../../migrations",
	}
	n, err := migrate.Exec(globalDB.DB, globalDialect, migrations, migrate.Down)
	if err != nil {
		panic(err.Error())
	}
	log.Printf("Applied %d migrations down to testDB\n", n)
	n, err = migrate.Exec(globalDB.DB, globalDialect, migrations, migrate.Up)
	if err != nil {
		panic(err.Error())
	}
	log.Printf("Applied %d migrations up to testDB\n", n)
}

func teardownTest() {
	if globalDB != nil {
		globalDB.Close()
	}
}

func TestMain(m *testing.M) {
	setupTest()
	code := m.Run()
	teardownTest()
	os.Exit(code)
}
