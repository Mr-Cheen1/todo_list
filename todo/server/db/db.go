package db

import (
	"database/sql"
	"fmt"
	"os"

	_ "github.com/lib/pq" // Импорт драйвера PostgreSQL для использования с database/sql.
)

var DB *sql.DB

func InitDB() {
	host := os.Getenv("DB_HOST")
	port := os.Getenv("DB_PORT")
	user := os.Getenv("DB_USER")
	password := os.Getenv("DB_PASSWORD")
	dbname := os.Getenv("DB_NAME")

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
