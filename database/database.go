package database

import (
	"time"
)

type User struct {
	Id         int
	First_name string
	Last_name  string
	Email      string
	Password   string
	Status_id  int
	Created_at time.Time
	Updated_at time.Time
	Deleted    int
}

type User_status struct {
	Id         int
	Status     string
	Created_at time.Time
	Updated_at time.Time
	Deleted    int
}
