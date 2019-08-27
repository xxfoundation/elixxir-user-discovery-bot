////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Handles the database ORM for the user discovery bot

package storage

import (
	"fmt"
	"github.com/pkg/errors"
)

// Insert or Update a User into the database
func (m *DatabaseImpl) UpsertUser(user *User) error {
	// Perform the upsert
	_, err := m.db.Model(user).
		// On conflict, update the user's fields
		OnConflict("(id) DO UPDATE").
		Set("value = EXCLUDED.value," +
			"value_type = EXCLUDED.value_type," +
			"key_id = EXCLUDED.key_id," +
			"key = EXCLUDED.key").
		// Otherwise, insert the new user
		Insert()
	if err != nil {
		return errors.New(fmt.Sprintf("unable to upsert user %s: %+v",
			string(user.Id),
			errors.New(err.Error())))
	}
	return nil
}

// Fetch a User from the database
func (m *DatabaseImpl) GetUser(user *User) (*User, error) {
	err := m.db.Select(&user)
	if err != nil {
		// If there was an error, no user for the given ID was found
		return nil, errors.New(fmt.Sprintf("unable to find user %v: %+v",
			string(user.Id),
			errors.New(err.Error())))
	}
	// If we found a user for the given ID, return it
	return user, nil
}
