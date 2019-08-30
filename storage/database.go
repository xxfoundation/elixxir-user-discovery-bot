////////////////////////////////////////////////////////////////////////////////
// Copyright © 2019 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package storage

import (
	"fmt"
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"gitlab.com/elixxir/client/globals"
	"gitlab.com/elixxir/primitives/id"
)

// Struct implementing the Database Interface with an underlying DB
type DatabaseImpl struct {
	db *pg.DB // Stored database connection
}

var UserDiscoveryDb Database

type Database interface {
	// Insert or Update a User into the database
	UpsertUser(user *User) error
	// Fetch a User from the database by ID
	GetUser(id []byte) (*User, error)
	// Fetch a User from the database by Value
	GetUserByValue(value string) (*User, error)
	// Fetch a User from the database by KeyId
	GetUserByKeyId(keyId string) (*User, error)
}

// Struct representing the udb_users table in the database
type User struct {
	// Overwrite table name
	tableName struct{} `sql:"udb_users,alias:udb_users"`

	// User Id
	Id []byte `sql:",pk"`
	// Identifying information
	Value string
	// Type of identifying information as denoted by the ValueType type
	ValueType int `sql:"type:value_type"`
	// Hash of the User public key
	KeyId string
	// User public key
	Key []byte
}

// Initialize a new User object
func NewUser() *User {
	return &User{
		Id:        make([]byte, 0),
		Value:     "",
		ValueType: -1,
		KeyId:     "",
		Key:       make([]byte, 0),
	}
}

func (u *User) SetID(id []byte) {
	u.Id = id
}

func (u *User) SetValue(val string) {
	u.Value = val
}

func (u *User) SetValueType(valType int) {
	u.ValueType = valType
}

func (u *User) SetKeyID(keyID string) {
	u.KeyId = keyID
}

func (u *User) SetKey(key []byte) {
	u.Key = key
}

// Initialize the Database interface with database backend
func NewDatabase(username, password, database, address string) Database {
	// Create the database connection
	db := pg.Connect(&pg.Options{
		User:         username,
		Password:     password,
		Database:     database,
		Addr:         address,
		MaxRetries:   10,
		MinIdleConns: 1,
	})

	// Initialize the schema
	err := createSchema(db)
	if err != nil {
		// If an error is thrown with the database, run with a map backend
		globals.Log.INFO.Println("Using map backend for User Discovery!")
		return &MapImpl{
			Users: make(map[*id.User]*User),
		}
	}

	// Create the ValueType enum in the database
	_, err = db.Exec(fmt.Sprintf(`CREATE TYPE value_type AS ENUM (%d,%d);`,
		Email, Nick))
	if err != nil {
		globals.Log.FATAL.Panicf("Unable to create enum: %+v", err)
	}

	// Return the database-backed Database interface
	globals.Log.INFO.Println("Using database backend for User Discovery!")
	return &DatabaseImpl{
		db: db,
	}
}

// Create the database schema
func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{&User{}} {
		err := db.CreateTable(model, &orm.CreateTableOptions{
			// Ignore create table if already exists?
			IfNotExists: true,
			// Create temporary table?
			Temp: false,
			// FKConstraints causes CreateTable to create foreign key constraints
			// for has one relations. ON DELETE hook can be added using tag
			// `sql:"on_delete:RESTRICT"` on foreign key field.
			FKConstraints: true,
			// Replaces PostgreSQL data type `text` with `varchar(n)`
			// Varchar: 255
		})
		if err != nil {
			// Return error if one comes up
			return err
		}
	}
	// No error, return nil
	return nil
}
