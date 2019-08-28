////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package storage

import (
	"fmt"
	"gitlab.com/elixxir/primitives/id"
	"reflect"
	"testing"
)

//Happy path
func TestMap_UpsertUser(t *testing.T) {
	m := &MapImpl{
		Users: make(map[*id.User]*User),
	}

	usr := NewUser()
	usr.Id = make([]byte, 8)

	err := m.UpsertUser(usr)

	if err != nil {
		t.Errorf("Expected to successfully upsert user, recieved err: %+v", err)
	}
}

//Test that map updates a new user being inserted with same id
func TestMap_UpsertDuplicate(t *testing.T) {
	m := &MapImpl{
		Users: make(map[*id.User]*User),
	}

	usr := NewUser()
	usr.Id = make([]byte, 8)

	_ = m.UpsertUser(usr)

	usr2 := usr
	usr2.Value = "email"

	_ = m.UpsertUser(usr2)

	observedUser, _ := m.GetUser(usr)

	if observedUser.Value != usr.Value {
		t.Errorf("Failed to update a user with new information")
	}
}

//Happy path
func TestMapImpl_GetUser(t *testing.T) {
	m := &MapImpl{
		Users: make(map[*id.User]*User),
	}

	//Populate the user
	usr := NewUser()
	usr.Key = make([]byte, 8)
	usr.Id = make([]byte, 8)
	usr.Value = "email"
	usr.ValueType = 1

	_ = m.UpsertUser(usr)
	retrievedUser, _ := m.GetUser(usr)

	if !reflect.DeepEqual(retrievedUser, usr) {
		t.Errorf("Expected to retrieve %+v, recieved: %+v", usr, retrievedUser)
	}

}

//Error path: pull a nonexistant user
func TestMapImpl_GetUser_EmptyMap(t *testing.T) {
	m := &MapImpl{
		Users: make(map[*id.User]*User),
	}
	//Create user, never insert in map
	usr := NewUser()
	usr.Id = make([]byte, 8)
	usr.Value = "email"

	//Search for usr in empty map
	retrievedUser, _ := m.GetUser(usr)

	//Check that no user is obtained from an empty map
	if !reflect.DeepEqual(retrievedUser, NewUser()) {
		t.Errorf("Expected to not find user in empty map. Map: %+v", m)
	}

}

//Error path: request a value that doesn't exist in the map
func TestMapImpl_GetUser_NilValue(t *testing.T) {
	m := &MapImpl{
		Users: make(map[*id.User]*User),
	}

	//Make a user with a value set
	usr := NewUser()
	usr.Value = "email"
	_ = m.UpsertUser(usr)

	//Search for a user with no value set
	usr2 := NewUser()
	usr2.Id = make([]byte, 8)
	retrievedUser,  err:= m.GetUser(usr2)
	fmt.Println(err)
	//Should return an empty user, as map doesn't have a user with id set
	if !reflect.DeepEqual(retrievedUser, NewUser()) || err == nil {
		t.Errorf("Should have retrieved: %+v: Recieved: %+v", NewUser(), retrievedUser)
	}

}

//Happy path: Insert and get a user for every user attribute
func TestMapImpl_GetUser_AddAndGet(t *testing.T) {
	m := &MapImpl{
		Users: make(map[*id.User]*User),
	}

	//Insert user with ID and get user
	usrID := NewUser()
	usrID.Id = make([]byte, 8)
	_ = m.UpsertUser(usrID)
	retrievedUser, _ := m.GetUser(usrID)
	if !reflect.DeepEqual(retrievedUser, usrID) {
		t.Errorf("Inserted and pulled an id. "+
			"Should have retrieved: %+v, recieved: %+v", retrievedUser, usrID)
	}

	//Insert user with val and get user
	usrVal := NewUser()
	usrVal.Value = "email"
	_ = m.UpsertUser(usrVal)
	retrievedUser, _ = m.GetUser(usrVal)
	if !reflect.DeepEqual(retrievedUser, usrVal) {
		t.Errorf("Inserted and pulled a value. "+
			"Should have retrieved: %+v, recieved: %+v", retrievedUser, usrVal)
	}

	//Insert user with val type and then get user
	usrValType := NewUser()
	usrValType.ValueType = 1
	_ = m.UpsertUser(usrValType)
	retrievedUser, _ = m.GetUser(usrValType)
	if !reflect.DeepEqual(retrievedUser, usrValType) {
		t.Errorf("Inserted and pulled a value type. "+
			"Should have retrieved: %+v, recieved: %+v", usrValType, usrValType)
	}

	//Insert a user with key and then get user
	usrKey := NewUser()
	usrKey.Key = make([]byte, 8)
	_ = m.UpsertUser(usrKey)
	retrievedUser, _ = m.GetUser(usrKey)
	if !reflect.DeepEqual(retrievedUser, usrKey) {
		t.Errorf("Inserted and pulled a key. "+
			"Should have retrieved: %+v, recieved: %+v", retrievedUser, usrKey)
	}

	//Insert a user with key id and then get user
	usrKeyId := NewUser()
	usrKeyId.KeyId = "test"
	_ = m.UpsertUser(usrKeyId)
	retrievedUser, _ = m.GetUser(usrKeyId)
	if !reflect.DeepEqual(retrievedUser, usrKeyId) {
		t.Errorf("Inserted and pulled a keyID. "+
			"Should have retrieved: %+v, recieved: %+v", retrievedUser, usrKeyId)
	}

}
