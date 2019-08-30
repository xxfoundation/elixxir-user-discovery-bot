////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package storage

import (
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
	usr.SetID(make([]byte, 8))

	_ = m.UpsertUser(usr)

	usr2 := usr
	usr2.SetValue("email@website.com")

	_ = m.UpsertUser(usr2)

	observedUser, _ := m.GetUser(usr.Id)

	if observedUser.Value != usr.Value {
		t.Errorf("Failed to update a user with new information")
	}
}

//TODO: Duplicate/add for all the new get user funcs
//Happy path
func TestMapImpl_GetUser(t *testing.T) {
	m := &MapImpl{
		Users: make(map[*id.User]*User),
	}

	//Populate the user
	usr := NewUser()
	usr.SetKeyID("testKeyFP")
	usr.SetID(make([]byte, 8))
	usr.SetValue("email@website.com")
	usr.SetValueType(1)

	_ = m.UpsertUser(usr)
	retrievedUser, _ := m.GetUser(usr.Id)

	if !reflect.DeepEqual(retrievedUser, usr) {
		t.Errorf("Failed to retrieve by user ID. Expected to retrieve %+v, recieved: %+v", usr, retrievedUser)
	}

	retrievedUser, _ = m.GetUserByKeyId(usr.KeyId)
	if !reflect.DeepEqual(usr, retrievedUser) {
		t.Errorf("Failed to retrieve by key ID. Expected to retrieve %+v, recieved: %+v", usr, retrievedUser)
	}

	retrievedUser, _ = m.GetUserByValue(usr.Value)
	if !reflect.DeepEqual(usr, retrievedUser) {
		t.Errorf("Failed to retrieve by value. Expected to retrieve %+v, recieved: %+v", usr, retrievedUser)
	}

}

//Error path: pull a nonexistant user
func TestMapImpl_GetUser_EmptyMap(t *testing.T) {
	m := &MapImpl{
		Users: make(map[*id.User]*User),
	}
	//Create user, never insert in map
	usr := NewUser()
	usr.SetID(make([]byte, 8))
	usr.SetValue("email@website.com")
	usr.SetKeyID("testKeyFP")
	//Search for usr in empty map
	retrievedUser, _ := m.GetUser(usr.Id)

	//Check that no user is obtained from an empty map
	if !reflect.DeepEqual(retrievedUser, NewUser()) {
		t.Errorf("Expected to not find user in empty map. Map: %+v", m)
	}

	retrievedUser, _ = m.GetUserByValue(usr.Value)
	//Check that no user is obtained from an empty map
	if !reflect.DeepEqual(retrievedUser, NewUser()) {
		t.Errorf("Expected to not find user in empty map. Map: %+v", m)
	}

	retrievedUser, _ = m.GetUserByKeyId(usr.KeyId)
	//Check that no user is obtained from an empty map
	if !reflect.DeepEqual(retrievedUser, NewUser()) {
		t.Errorf("Expected to not find user in empty map. Map: %+v", m)
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
	retrievedUser, _ := m.GetUser(usrID.Id)
	if !reflect.DeepEqual(retrievedUser, usrID) {
		t.Errorf("Inserted and pulled an id. "+
			"Should have retrieved: %+v, recieved: %+v", retrievedUser, usrID)
	}

	//Insert user with val and get user
	usrVal := NewUser()
	usrVal.Value = "email"
	_ = m.UpsertUser(usrVal)
	retrievedUser, _ = m.GetUserByValue(usrVal.Value)
	if !reflect.DeepEqual(retrievedUser, usrVal) {
		t.Errorf("Inserted and pulled a value. "+
			"Should have retrieved: %+v, recieved: %+v", retrievedUser, usrVal)
	}

	//Insert a user with key id and then get user
	usrKeyId := NewUser()
	usrKeyId.KeyId = "testKeyFP"
	_ = m.UpsertUser(usrKeyId)
	retrievedUser, _ = m.GetUserByKeyId(usrKeyId.KeyId)
	if !reflect.DeepEqual(retrievedUser, usrKeyId) {
		t.Errorf("Inserted and pulled a keyID. "+
			"Should have retrieved: %+v, recieved: %+v", retrievedUser, usrKeyId)
	}

}
