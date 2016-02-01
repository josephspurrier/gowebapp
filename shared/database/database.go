package database

import (
	"encoding/json"
	"fmt"
	"log"
	"time"

	"github.com/boltdb/bolt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/jmoiron/sqlx"
	"gopkg.in/mgo.v2"
)

var (
	BoltDB    *bolt.DB     // Bolt wrapper
	Mongo     *mgo.Session // Mongo wrapper
	Sql       *sqlx.DB     // SQL wrapper
	databases DatabaseInfo // Database info
)

type DatabaseType string

const (
	TypeBolt    DatabaseType = "Bolt"
	TypeMongoDB DatabaseType = "MongoDB"
	TypeMySQL   DatabaseType = "MySQL"
)

type DatabaseInfo struct {
	Type    DatabaseType
	MySQL   MySQLInfo
	Bolt    BoltInfo
	MongoDB MongoDBInfo
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

// BoltInfo is the details for the database connection
type BoltInfo struct {
	Path string
}

// MongoDBInfo is the details for the database connection
type MongoDBInfo struct {
	URL      string
	Database string
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
func Connect(d DatabaseInfo) {
	var err error

	// Store the config
	databases = d

	switch d.Type {
	case TypeMySQL:
		// Connect to MySQL
		if Sql, err = sqlx.Connect("mysql", DSN(d.MySQL)); err != nil {
			log.Println("SQL Driver Error", err)
		}

		// Check if is alive
		if err = Sql.Ping(); err != nil {
			log.Println("Database Error", err)
		}
	case TypeBolt:
		// Connect to Bolt
		if BoltDB, err = bolt.Open(d.Bolt.Path, 0600, nil); err != nil {
			log.Println("Bolt Driver Error", err)
		}
	case TypeMongoDB:
		// Connect to MongoDB
		if Mongo, err = mgo.DialWithTimeout(d.MongoDB.URL, 5); err != nil {
			log.Println("MongoDB Driver Error", err)
			return
		}

		// Prevents these errors: read tcp 127.0.0.1:27017: i/o timeout
		Mongo.SetSocketTimeout(1 * time.Second)

		// Check if is alive
		if err = Mongo.Ping(); err != nil {
			log.Println("Database Error", err)
		}
	default:
		log.Println("No registered database in config")
	}
}

// Update makes a modification to Bolt
func Update(bucket_name string, key string, dataStruct interface{}) error {
	err := BoltDB.Update(func(tx *bolt.Tx) error {
		// Create the bucket
		bucket, e := tx.CreateBucketIfNotExists([]byte(bucket_name))
		if e != nil {
			return e
		}

		// Encode the record
		encoded_record, e := json.Marshal(dataStruct)
		if e != nil {
			return e
		}

		// Store the record
		if e = bucket.Put([]byte(key), encoded_record); e != nil {
			return e
		}
		return nil
	})
	return err
}

// View retrieves a record in Bolt
func View(bucket_name string, key string, dataStruct interface{}) error {
	err := BoltDB.View(func(tx *bolt.Tx) error {
		// Get the bucket
		b := tx.Bucket([]byte(bucket_name))
		if b == nil {
			return bolt.ErrBucketNotFound
		}

		// Retrieve the record
		v := b.Get([]byte(key))
		if len(v) < 1 {
			return bolt.ErrInvalid
		}

		// Decode the record
		e := json.Unmarshal(v, &dataStruct)
		if e != nil {
			return e
		}

		return nil
	})

	return err
}

// CheckConnection returns true if MongoDB is available
func CheckConnection() bool {
	if Mongo == nil {
		Connect(databases)
	}

	if Mongo != nil {
		return true
	}

	return false
}

// ReadConfig returns the database information
func ReadConfig() DatabaseInfo {
	return databases
}
