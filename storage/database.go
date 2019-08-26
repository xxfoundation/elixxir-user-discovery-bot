////////////////////////////////////////////////////////////////////////////////
// Copyright © 2019 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package storage

import (
	"github.com/go-pg/pg"
	jww "github.com/spf13/jwalterweatherman"
	"pg/orm"
	"time"
)

// Struct implementing the Database Interface with an underlying DB
type DatabaseImpl struct {
	db *pg.DB // Stored database connection
}


type UserDb Database

type Database interface {
	// TODO: Fill me in
}

// Initialize the Database interface with database backend
func NewDatabase(username, password, database, address string) Database {

	// Create the database connection
	db := pg.Connect(&pg.Options{
		User:        username,
		Password:    password,
		Database:    database,
		Addr:        address,
		PoolSize:    1,
		MaxRetries:  10,
		PoolTimeout: time.Duration(2) * time.Minute,
		IdleTimeout: time.Duration(10) * time.Minute,
		MaxConnAge:  time.Duration(1) * time.Hour,
	})

	// Ensure an empty NodeInformation table
	err := db.DropTable(&NodeInformation{},
		&orm.DropTableOptions{IfExists: true})
	if err != nil {
		// If an error is thrown with the database, run with a map backend
		jww.INFO.Println("Using map backend for User Discovery!")
		return Database(&MapImpl{
			/*client: make(map[string]*RegistrationCode),
			node:   make(map[string]*NodeInformation),*/
		})
	}

	// Initialize the schema
	jww.INFO.Println("Using database backend for User Discovery!")
	err = createSchema(db)
	if err != nil {
		jww.FATAL.Panicf("Unable to initialize database backend for UDB: %+v", err)
	}

	// Return the database-backed Database interface
	jww.INFO.Println("Database backend initialized successfully!")
	return Database(&DatabaseImpl{
		db: db,
	})
}

// Create the database schema
func createSchema(db *pg.DB) error {
	for _, model := range []interface{}{&RegistrationCode{}, &NodeInformation{}} {
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
}'