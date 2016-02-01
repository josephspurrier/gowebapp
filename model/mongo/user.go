package model

import (
	"errors"
	"time"

	"github.com/josephspurrier/gowebapp/shared/database"

	"gopkg.in/mgo.v2"
	"gopkg.in/mgo.v2/bson"
)

// *****************************************************************************
// User
// *****************************************************************************

// User table contains the information for each user
type User struct {
	ObjectId   bson.ObjectId `bson:"_id"`
	First_name string        `db:"first_name" bson:"first_name"`
	Last_name  string        `db:"last_name" bson:"last_name"`
	Email      string        `db:"email" bson:"email"`
	Password   string        `db:"password" bson:"password"`
	Status_id  uint8         `db:"status_id" bson:"status_id"`
	Created_at time.Time     `db:"created_at" bson:"created_at"`
	Updated_at time.Time     `db:"updated_at" bson:"updated_at"`
	Deleted    uint8         `db:"deleted" bson:"deleted"`
}

var (
	ErrCode        = errors.New("Case statement in code is not correct.")
	ErrNoResult    = errors.New("Result not found.")
	ErrUnavailable = errors.New("Database is unavailable.")
)

// Id returns the user id
func (u *User) ID() string {
	return u.ObjectId.Hex()
}

// standardizeErrors returns the same error regardless of the database used
func standardizeError(err error) error {
	if err == mgo.ErrNotFound {
		return ErrNoResult
	}

	return err
}

// UserByEmail gets user information from email
func UserByEmail(email string) (User, error) {
	var err error

	result := User{}

	if database.CheckConnection() {
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("user")
		err = c.Find(bson.M{"email": email}).One(&result)
	} else {
		err = ErrUnavailable
	}

	return result, standardizeError(err)
}

// UserCreate creates user
func UserCreate(first_name, last_name, email, password string) error {
	var err error

	now := time.Now()

	if database.CheckConnection() {
		session := database.Mongo.Copy()
		defer session.Close()
		c := session.DB(database.ReadConfig().MongoDB.Database).C("user")

		user := &User{
			ObjectId:   bson.NewObjectId(),
			First_name: first_name,
			Last_name:  last_name,
			Email:      email,
			Password:   password,
			Status_id:  1,
			Created_at: now,
			Updated_at: now,
			Deleted:    0,
		}
		err = c.Insert(user)
	} else {
		err = ErrUnavailable
	}

	return standardizeError(err)
}
