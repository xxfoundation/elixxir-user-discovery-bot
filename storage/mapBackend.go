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
	"github.com/pkg/errors"
	" gitlab.com/xx_network/primitives/id"
	"strings"
	"sync"
)

// Struct implementing the Database Interface with an underlying Map
type MapImpl struct {
	Users map[*id.ID]*User
	lock  sync.Mutex
}

// Insert or Update a User into the map backend
func (m *MapImpl) UpsertUser(user *User) error {
	m.lock.Lock()
	defer m.lock.Unlock()

	//Insert or update the user in the map
	tempIndex, err := id.Unmarshal(user.Id)
	if err != nil {
		return err
	}
	m.Users[tempIndex] = user

	return nil
}

// Fetch a User from the database by ID
func (m *MapImpl) GetUser(userID *id.ID) (*User, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	//Iterate through the list of users and find matching values
	for _, u := range m.Users {

		if bytes.Compare(u.Id, userID.Bytes()) == 0 && bytes.Compare(u.Id, make([]byte, 0)) != 0 {
			return u, nil
		}

	}

	return NewUser(), errors.New("Unable to find any user with that ID")
}

// Fetch a User from the database by Value
func (m *MapImpl) GetUserByValue(value string) (*User, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, u := range m.Users {
		if strings.Compare(u.Value, value) == 0 && u.Value != "" {
			fmt.Println(m)
			return u, nil
		}
	}

	return NewUser(), errors.New("Unable to find any user with that value")
}

// Fetch a User from the database by KeyId
func (m *MapImpl) GetUserByKeyId(keyId string) (*User, error) {
	m.lock.Lock()
	defer m.lock.Unlock()

	for _, u := range m.Users {
		if strings.Compare(u.KeyId, keyId) == 0 && u.KeyId != "" {
			return u, nil
		}
	}

	return NewUser(), errors.New("Unable to find any user with that keyID")
}

//Delete user by user id
func (m *MapImpl) DeleteUser(userID *id.ID) error {
	m.lock.Lock()

	delete(m.Users, userID)
	m.lock.Unlock()
	return nil

}
