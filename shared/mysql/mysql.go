package mysql

import (
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
)

// Connection provides access to the database wrapper
type Connection struct {
	Link *sqlx.DB
}

// ConnectionInfo is the details for the database connection
type ConnectionInfo struct {
	Username  string
	Password  string
	Database  string
	Hostname  string
	Port      int
	Parameter string
}

var connInfo ConnectionInfo

func connectionInfo() string {
	// Example: root:@tcp(localhost:3306)/test
	return connInfo.Username +
		":" +
		connInfo.Password +
		"@tcp(" +
		connInfo.Hostname +
		":" +
		fmt.Sprintf("%d", connInfo.Port) +
		")/" +
		connInfo.Database + connInfo.Parameter
}

// Config will set the connection information
func Config(ci ConnectionInfo) {
	connInfo = ci
}

// Instance will return a connection from the pool
func Instance() (*Connection, error) {
	c := &Connection{}

	conn, err := sqlx.Connect("mysql", connectionInfo())
	c.Link = conn
	if err != nil {
		log.Println("SQL Driver Error", err)
	}

	return c, err
}
