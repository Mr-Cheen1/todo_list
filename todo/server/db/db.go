package db

import (
	"database/sql"
	"fmt"

	_ "github.com/lib/pq" // Импорт драйвера PostgreSQL для использования с database/sql.
)

var DB *sql.DB

func InitDB(host, port, user, password, dbname string) {
	connStr := fmt.Sprintf("postgres://%s:%s@%s:%s/%s?sslmode=disable", user, password, host, port, dbname)

	var err error
	DB, err = sql.Open("postgres", connStr)
	if err != nil {
		panic(err)
	}
}

func CloseDB() {
	DB.Close()
}
