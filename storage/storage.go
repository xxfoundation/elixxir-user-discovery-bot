////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Handles the high level storage API.
// This layer merges the business logic layer and the database layer

package storage

import (
	jww "github.com/spf13/jwalterweatherman"
	"testing"
)

// API for the storage layer
type Storage struct {
	// Stored Database interface
	database
}

// Create a new Storage object wrapping a database interface
// Returns a Storage object and error
func NewStorage(username, password, dbName, address, port string) (*Storage, error) {
	db, err := newDatabase(username, password, dbName, address, port)
	storage := &Storage{db}
	return storage, err
}

func NewTestDB(t *testing.T) *Storage {
	if t == nil {
		jww.FATAL.Panic("Cannot use this outside of testing")
	}
	mockDb, err := NewStorage("", "", "", "", "11")
	if err != nil {
		jww.FATAL.Panicf("Failed to init mock db: %+v", err)
	}
	return mockDb
}
