////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Handles the Map backend for the user discovery bot

package storage

import (
	"gitlab.com/elixxir/primitives/id"
	"reflect"
)

// Insert or Update a User into the map backend
func (m *MapImpl) UpsertUser(user *User) error {
	m.lock.Lock()
	//Insert or update the user in the map
	tmpIndx := id.NewUserFromBytes(user.Id)
	m.users[tmpIndx] = user

	m.lock.Unlock()
	return nil
}

// Fetch a User from the map backend. Pass in a user with any attribute values
// you want to search and we will search for them
func (m *MapImpl) GetUser(user *User) (*User, error) {
	m.lock.Lock()
	/*
		var err error
		retUser, ok := m.users[user.Id]
		if ok {
			err = errors.New(fmt.Sprintf(
				"User %+v has not been added!", user))
		}*/
	retrievedUser := NewUser()

	//Flatten map into a list
	users := make([]*User, 0)
	for _, value := range m.users {
		users = append(users, value)
	}

	//Iterate through the list of users and find matching values
	for _, u := range users {
		if reflect.DeepEqual(u.Id, user.Id) && u.Id != nil {
			retrievedUser.Id = u.Id
		}

		if reflect.DeepEqual(u.Value, user.Value) && u.Value != "" {
			retrievedUser.Value = u.Value
		}

		if u.ValueType == user.ValueType {
			retrievedUser.ValueType = u.ValueType
		}

		if reflect.DeepEqual(u.Key, user.Key) && u.Key != nil {
			retrievedUser.Key = u.Key
		}

		if reflect.DeepEqual(u.KeyId, user.KeyId) && u.KeyId != "" {
			retrievedUser.KeyId = u.KeyId
		}

	}

	m.lock.Unlock()
	return retrievedUser, nil
}
