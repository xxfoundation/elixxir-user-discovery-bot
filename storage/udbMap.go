////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Handles the Map backend for udb storage

package storage

import (
	"errors"
	"gitlab.com/xx_network/primitives/id"
)

// Insert a new user
func (m *MapImpl) InsertUser(user *User) error {
	uid, err := id.Unmarshal(user.Id)
	if err != nil {
		return err
	}
	m.users[*uid] = user
	return nil
}

// Get a user from the map
func (m *MapImpl) GetUser(uid []byte) (*User, error) {
	nuid, err := id.Unmarshal(uid)
	if err != nil {
		return nil, err
	}
	u, _ := m.users[*nuid]
	return u, nil
}

// Delete a user from the map
func (m *MapImpl) DeleteUser(uid []byte) error {
	nuid, err := id.Unmarshal(uid)
	if err != nil {
		return err
	}
	delete(m.users, *nuid)
	return nil
}

// Insert a new fact
func (m *MapImpl) InsertFact(fact *Fact) error {
	uid, err := id.Unmarshal(fact.UserId)
	if err != nil {
		return err
	}
	if _, ok := m.users[*uid]; !ok {
		return errors.New("error: associated user not found")
	}
	factid := ConfirmationId{}
	copy(factid[:], fact.ConfirmationId)
	m.facts[factid] = fact
	return nil
}

// Get a fact from the map
func (m *MapImpl) GetFact(confirmationId []byte) (*Fact, error) {
	factid := ConfirmationId{}
	copy(factid[:], confirmationId)
	f, _ := m.facts[factid]
	return f, nil
}

// Delete a fact from the map
func (m *MapImpl) DeleteFact(confirmationId []byte) error {
	factid := ConfirmationId{}
	copy(factid[:], confirmationId)
	delete(m.facts, factid)
	return nil
}

// Confirm a fact in the map
func (m *MapImpl) ConfirmFact(confirmationId string) error {
	factid := ConfirmationId{}
	copy(factid[:], confirmationId)
	if _, ok := m.facts[factid]; !ok {
		return errors.New("specified fact not found")
	}
	m.facts[factid].VerificationStatus = 1
	return nil
}
