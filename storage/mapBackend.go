////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Handles the Map backend for the user discovery bot

package storage

import (
	"bytes"
	"gitlab.com/elixxir/primitives/id"
	"golang.org/x/tools/go/ssa/interp/testdata/src/errors"
	"strings"
)

// Insert or Update a User into the map backend
func (m *MapImpl) UpsertUser(user *User) error {
	m.lock.Lock()
	//Insert or update the user in the map
	tmpIndx := id.NewUserFromBytes(user.Id)
	m.Users[tmpIndx] = user

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

	//Iterate through the list of users and find matching values
	for _, u := range m.Users {

		if bytes.Compare(u.Id, user.Id) == 0 && bytes.Compare(u.Id, make([]byte, 0)) != 0 {
			m.lock.Unlock()
			return u, nil
		}

		if strings.Compare(u.Value, user.Value) == 0  && u.Value != "" {
			m.lock.Unlock()
			return u, nil
		}

		if u.ValueType == user.ValueType && u.ValueType != -1 {
			m.lock.Unlock()
			return u, nil
		}

		if (bytes.Compare(u.Key, user.Key) == 0) && bytes.Compare(u.Key, make([]byte, 0)) != 0 {
			m.lock.Unlock()
			return u, nil
		}

		if strings.Compare(u.KeyId, user.KeyId) == 0 && u.KeyId != "" {
			m.lock.Unlock()
			return u, nil
		}
	}
	m.lock.Unlock()
	return NewUser(), errors.New("Unable to find any user with those values")
}
