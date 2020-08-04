////////////////////////////////////////////////////////////////////////////////
// Copyright © 2019 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package storage

import (
	"github.com/go-pg/pg"
	"github.com/go-pg/pg/orm"
	"gitlab.com/elixxir/client/globals"
	" gitlab.com/xx_network/primitives/id"
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
	GetUser(id *id.ID) (*User, error)
	// Fetch a User from the database by Value
	GetUserByValue(value string) (*User, error)
	// Fetch a User from the database by KeyId
	GetUserByKeyId(keyId string) (*User, error)
	//Delete a user
	DeleteUser(id *id.ID) error
}

// Struct representing the udb_users table in the database
type User struct {
	// Overwrite table name
	tableName struct{} `sql:"udb_users,alias:udb_users"`

	// User Id
	Id []byte `sql:",pk,unique"`
	// Identifying informationgo-pg
	Value string `sql:",unique"`
	// Type of identifying information as denoted by the ValueType type
	ValueType int
	// Hash of the User public key
	KeyId string `sql:",unique"`
	// User public key
	Key []byte `sql:",unique"`
}

// Initialize a new User object
func NewUser() *User {
	return &User{
		Id:        make([]byte, id.ArrIDLen),
		Value:     "",
		ValueType: -1,
		KeyId:     "",
		Key:       make([]byte, 0),
	}
}

func (u *User) SetID(id *id.ID) {
	u.Id = id.Marshal()
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
		globals.Log.ERROR.Printf("Unable to initalize database backend: %+v", err)
		globals.Log.INFO.Println("Using map backend for User Discovery!")
		return &MapImpl{
			Users: make(map[*id.ID]*User),
		}
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
			FKConstraints: false,
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
