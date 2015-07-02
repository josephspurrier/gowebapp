package mysql

// Helpful post: http://stackoverflow.com/a/17384204

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Connect will set the connection information so mysql.DB can be used
func Connect(dsn string) (*sqlx.DB, error) {
	DB, err := sqlx.Connect("mysql", dsn)
	return DB, err
}
