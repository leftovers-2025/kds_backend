package datasource

import (
	"fmt"
	"os"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// MySQLに接続する
func GetMySqlConnection() *sqlx.DB {
	db, err := sqlx.Connect("mysql", getMySqlDSN())
	if err != nil {
		panic("failed to connect MySQL database: " + err.Error())
	}
	return db
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
