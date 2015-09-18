package database

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	_ "github.com/mattn/go-sqlite3"
)

var (
	// Database wrapper
	DB *sqlx.DB
)

type Databases struct {
	Type   string
	MySQL  MySQLInfo
	SQLite SQLiteInfo
}

// MySQLInfo is the details for the database connection
type MySQLInfo struct {
	Username  string
	Password  string
	Name      string
	Hostname  string
	Port      int
	Parameter string
}

// SQLiteInfo is the details for the database connection
type SQLiteInfo struct {
	Parameter string
}

// DSN returns the Data Source Name
func DSN(ci MySQLInfo) string {
	// Example: root:@tcp(localhost:3306)/test
	return ci.Username +
		":" +
		ci.Password +
		"@tcp(" +
		ci.Hostname +
		":" +
		fmt.Sprintf("%d", ci.Port) +
		")/" +
		ci.Name + ci.Parameter
}

// Connect to the database
func Connect(d Databases) {
	var err error

	switch d.Type {
	case "MySQL":
		// Connect to MySQL
		if DB, err = sqlx.Connect("mysql", DSN(d.MySQL)); err != nil {
			log.Println("SQL Driver Error", err)
		}

		// Check if is alive
		if err = DB.Ping(); err != nil {
			log.Println("Database Error", err)
		}
	case "SQLite":
		// Connect to SQLite
		if DB, err = sqlx.Connect("sqlite3", d.SQLite.Parameter); err != nil {
			log.Println("SQL Driver Error", err)
		}

		// Check if is alive
		if err = DB.Ping(); err != nil {
			log.Println("Database Error", err)
		}
	default:
		log.Println("No registered database in config")
	}
}
