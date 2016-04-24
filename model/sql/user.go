package model

import (
	"database/sql"
	"errors"
	"fmt"
	"time"

	"github.com/josephspurrier/gowebapp/shared/database"
)

// *****************************************************************************
// User
// *****************************************************************************

// User table contains the information for each user
type User struct {
	ID        uint32    `db:"id"`
	FirstName string    `db:"first_name"`
	LastName  string    `db:"last_name"`
	Email     string    `db:"email"`
	Password  string    `db:"password"`
	StatusID  uint8     `db:"status_id"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Deleted   uint8     `db:"deleted"`
}

// UserStatus table contains every possible user status (active/inactive)
type UserStatus struct {
	ID        uint8     `db:"id"`
	Status    string    `db:"status"`
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
	Deleted   uint8     `db:"deleted"`
}

var (
	// ErrCode is a config or an internal error
	ErrCode = errors.New("Case statement in code is not correct.")
	// ErrNoResult is a not results error
	ErrNoResult = errors.New("Result not found.")
	// ErrUnavailable is a database not available error
	ErrUnavailable = errors.New("Database is unavailable.")
)

// UserID returns the user id
func (u *User) UserID() string {
	return fmt.Sprintf("%v", u.Id)
}

// standardizeErrors returns the same error regardless of the database used
func standardizeError(err error) error {
	if err == sql.ErrNoRows {
		return ErrNoResult
	}

	return err
}

// UserByEmail gets user information from email
func UserByEmail(email string) (User, error) {
	result := User{}
	err := database.Sql.Get(&result, "SELECT id, password, status_id, first_name FROM user WHERE email = ? LIMIT 1", email)
	return result, err
}

// UserIDByEmail gets user id from email
func UserIDByEmail(email string) (User, error) {
	result := User{}
	err := database.Sql.Get(&result, "SELECT id FROM user WHERE email = ? LIMIT 1", email)
	return result, err
}

// UserCreate creates user
func UserCreate(firstName, lastName, email, password string) error {
	_, err := database.Sql.Exec("INSERT INTO user (first_name, last_name, email, password) VALUES (?,?,?,?)", firstName, lastName, email, password)
	return err
}
