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

func (m *MapImpl) CheckUser(username string, id *id.ID, rsaPem string) error {
	return nil
}

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
	factid := factId{}
	copy(factid[:], fact.FactHash)
	m.facts[factid] = fact
	return nil
}

func (m *MapImpl) VerifyFact(factHash []byte) error {
	return nil
}

// Delete a fact from the map
func (m *MapImpl) DeleteFact(confirmationId []byte) error {
	factid := factId{}
	copy(factid[:], confirmationId)
	delete(m.facts, factid)
	return nil
}

func (m *MapImpl) InsertFactTwilio(userID, factHash, signature []byte, fact string, factType uint, confirmationID string) error {
	return nil
}
func (m *MapImpl) VerifyFactTwilio(confirmationId string) error {
	return nil
}

func (m *MapImpl) Search(factHashs [][]byte) []*User {
	return nil
}
