package model

import (
	"errors"
	"time"

	"github.com/josephspurrier/gowebapp/shared/database"

	"gopkg.in/mgo.v2/bson"
)

// *****************************************************************************
// User
// *****************************************************************************

// User table contains the information for each user
type User struct {
	ObjectID  bson.ObjectId `bson:"_id"`
	FirstName string        `db:"first_name" bson:"first_name"`
	LastName  string        `db:"last_name" bson:"last_name"`
	Email     string        `db:"email" bson:"email"`
	Password  string        `db:"password" bson:"password"`
	StatusID  uint8         `db:"status_id" bson:"status_id"`
	CreatedAt time.Time     `db:"created_at" bson:"created_at"`
	UpdatedAt time.Time     `db:"updated_at" bson:"updated_at"`
	Deleted   uint8         `db:"deleted" bson:"deleted"`
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
	return u.ObjectID.Hex()
}

// standardizeErrors returns the same error regardless of the database used
func standardizeError(err error) error {
	return err
}

// UserByEmail gets user information from email
func UserByEmail(email string) (User, error) {
	var err error

	result := User{}

	err = database.View("user", email, &result)
	if err != nil {
		err = ErrNoResult
	}

	return result, standardizeError(err)
}

// UserCreate creates user
func UserCreate(firstName, lastName, email, password string) error {
	var err error

	now := time.Now()

	user := &User{
		ObjectID:  bson.NewObjectId(),
		FirstName: firstName,
		LastName:  lastName,
		Email:     email,
		Password:  password,
		StatusID:  1,
		CreatedAt: now,
		UpdatedAt: now,
		Deleted:   0,
	}

	err = database.Update("user", user.Email, &user)

	return standardizeError(err)
}
