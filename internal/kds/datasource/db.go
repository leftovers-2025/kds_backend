package datasource

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	migrate "github.com/rubenv/sql-migrate"
)

// MySQLに接続する
func GetMySqlConnection() *sqlx.DB {
	db, err := sqlx.Connect("mysql", getMySqlDSN())
	if err != nil {
		panic("failed to connect MySQL database: " + err.Error())
	}
	mysqlMigration(db)
	return db
}

func mysqlMigration(db *sqlx.DB) error {
	doMigrate, ok := os.LookupEnv("KDS_RUNTIME_MIGRATION")
	if !ok || doMigrate != "1" {
		fmt.Println("migration skipped. runtime migration is disabled")
		return nil
	}
	migration := &migrate.FileMigrationSource{
		Dir: "migrations",
	}
	n, err := migrate.Exec(db.DB, "mysql", migration, migrate.Up)
	if err != nil {
		return err
	}
	fmt.Printf("Applied %d migrations to mysql\n", n)
	return nil
}

// MySQLのDSNを環境変数から取得する
func getMySqlDSN() string {
	host := mustEnv("MYSQL_HOST")
	user := mustEnv("MYSQL_USER")
	password := mustEnv("MYSQL_PASSWORD")
	database := mustEnv("MYSQL_DATABASE")
	return fmt.Sprintf("%s:%s@tcp(%s)/%s?parseTime=true", user, password, host, database)
}

// 環境変数を読み込み、存在しない場合はpanicを発生
func mustEnv(envName string) string {
	env, ok := os.LookupEnv(envName)
	if !ok {
		panic(fmt.Sprintf("\"%s\" is not set", envName))
	}
	return env
}
