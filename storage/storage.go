////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Handles the high level storage API.
// This layer merges the business logic layer and the database layer

package storage

import "gitlab.com/elixxir/user-discovery-bot/interfaces/params"

// API for the storage layer
type Storage struct {
	// Stored Database interface
	database
}

// Create a new Storage object wrapping a database interface
// Returns a Storage object, close function, and error
func NewStorage(p params.Database) (*Storage, func() error, error) {
	return newDatabase(p.DbUsername, p.DbPassword, p.DbName, p.DbAddress, p.DbPort)
}
