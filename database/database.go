package database

import (
	"fmt"
	"log"
	"time"

	"github.com/josephspurrier/gowebapp/shared/mysql"

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

// *****************************************************************************
// User
// *****************************************************************************

type User struct {
	Id         int       `db:"id"`
	First_name string    `db:"first_name"`
	Last_name  string    `db:"last_name"`
	Email      string    `db:"email"`
	Password   string    `db:"password"`
	Status_id  int       `db:"status_id"`
	Created_at time.Time `db:"created_at"`
	Updated_at time.Time `db:"updated_at"`
	Deleted    int       `db:"deleted"`
	User_status
}

type User_status struct {
	Id         int       `db:"id"`
	Status     string    `db:"status"`
	Created_at time.Time `db:"created_at"`
	Updated_at time.Time `db:"updated_at"`
	Deleted    int       `db:"deleted"`
}

// UserByEmail gets user information from email
func UserByEmail(email string) (User, error) {
	result := User{}
	err := DB.Get(&result, "SELECT id, password, status_id, first_name FROM user WHERE email = ? LIMIT 1", email)
	return result, err
}

// UserIdByEmail gets user id from email
func UserIdByEmail(email string) (User, error) {
	result := User{}
	err := DB.Get(&result, "SELECT id FROM user WHERE email = ? LIMIT 1", email)
	return result, err
}

// CreateUser creates user
func CreateUser(first_name, last_name, email, password string) error {
	_, err := DB.Exec("INSERT INTO user (first_name, last_name, email, password) VALUES (?,?,?,?)", first_name, last_name, email, password)
	return err
}
