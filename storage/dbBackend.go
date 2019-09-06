////////////////////////////////////////////////////////////////////////////////
// Copyright © 2019 Privategrity Corporation                                   /
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
func (m *DatabaseImpl) GetUser(id []byte) (*User, error) {
	user := &User{Id: id}
	err := m.db.Select(&user)
	if err != nil {
		// If there was an error, no user for the given ID was found
		return nil, errors.New(fmt.Sprintf("unable to get user with id %v: %+v",
			string(id),
			errors.New(err.Error())))
	}
	// If we found a user for the given ID, return it
	return user, nil
}

// Fetch a User from the database by Value
func (m *DatabaseImpl) GetUserByValue(value string) (*User, error) {
	user := new(User)
	err := m.db.Model(user).Where("value = ?", value).Select()
	if err != nil {
		// If there was an error, no user for the given Value was found
		return nil, errors.New(fmt.Sprintf(
			"unable to get user with value %s: %+v", value,
			errors.New(err.Error())))
	}
	// If we found a user for the given Value, return it
	return user, nil
}

// Fetch a User from the database by KeyId
func (m *DatabaseImpl) GetUserByKeyId(keyId string) (*User, error) {
	user := new(User)
	err := m.db.Model(user).Where("key_id = ?", keyId).Select()
	if err != nil {
		// If there was an error, no user for the given KeyId was found
		return nil, errors.New(fmt.Sprintf(
			"unable to get user with keyId %s: %+v", keyId,
			errors.New(err.Error())))
	}
	// If we found a user for the given KeyId, return it
	return user, nil
}

//Delete a User from the database by the userID
func (m *DatabaseImpl) DeleteUser(id []byte) error {
	user := &User{Id: id}
	err := m.db.Delete(user)
	if err != nil {
		// If there was an error, no user for the given id was found
		return errors.New(fmt.Sprintf(
			"unable to delete user with keyId %s: %+v", id,
			errors.New(err.Error())))
	}

	return nil
}
