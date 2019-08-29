////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Handles the Map backend for the user discovery bot

package storage

import (
	"bytes"
	"fmt"
	"gitlab.com/elixxir/primitives/id"
	"golang.org/x/tools/go/ssa/interp/testdata/src/errors"
	"strings"
	"sync"
)

// Struct implementing the Database Interface with an underlying Map
type MapImpl struct {
	Users map[*id.User]*User
	lock  sync.Mutex
}

// Insert or Update a User into the map backend
func (m *MapImpl) UpsertUser(user *User) error {
	m.lock.Lock()
	//Insert or update the user in the map
	tmpIndx := id.NewUserFromBytes(user.Id)
	m.Users[tmpIndx] = user

	m.lock.Unlock()
	return nil
}

// Fetch a User from the database by ID
func (m *MapImpl) GetUser(id []byte) (*User, error) {
	m.lock.Lock()

	//Iterate through the list of users and find matching values
	for _, u := range m.Users {

		if bytes.Compare(u.Id, id) == 0 && bytes.Compare(u.Id, make([]byte, 0)) != 0 {
			m.lock.Unlock()
			return u, nil
		}

	}
	m.lock.Unlock()
	return NewUser(), errors.New("Unable to find any user with that ID")
}

// Fetch a User from the database by Value
func (m *MapImpl) GetUserByValue(value string) (*User, error) {
	m.lock.Lock()
	for _, u := range m.Users {
		if strings.Compare(u.Value, value) == 0 && u.Value != "" {
			m.lock.Unlock()
			fmt.Println(m)
			return u, nil
		}
	}

	m.lock.Unlock()
	return NewUser(), errors.New("Unable to find any user with that value")
}

// Fetch a User from the database by KeyId
func (m *MapImpl) GetUserByKeyId(keyId string) (*User, error) {
	m.lock.Lock()

	for _, u := range m.Users {
		if strings.Compare(u.KeyId, keyId) == 0 && u.KeyId != "" {
			m.lock.Unlock()
			return u, nil
		}
	}
	m.lock.Unlock()
	return NewUser(), errors.New("Unable to find any user with that keyID")
}
