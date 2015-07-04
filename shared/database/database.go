package database

import (
	"fmt"
	"log"

	"github.com/josephspurrier/gowebapp/shared/database/mysql"

	"github.com/jmoiron/sqlx"
)

var (
	// Database wrapper
	DB *sqlx.DB
)

// ConnectionInfo is the details for the database connection
type ConnectionInfo struct {
	Username  string
	Password  string
	Name      string
	Hostname  string
	Port      int
	Parameter string
}

// DSN returns the Data Source Name
func DSN(ci ConnectionInfo) string {
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
func Connect(ci ConnectionInfo) {
	var err error

	// Connect to MySQL
	if DB, err = mysql.Connect(DSN(ci)); err != nil {
		log.Fatalln("SQL Driver Error", err)
	}

	// Check if MySQL is alive
	if err := DB.Ping(); err != nil {
		log.Fatalln("Database Error", err)
	}
}
