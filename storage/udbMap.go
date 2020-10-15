////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Handles the Map backend for udb storage

package storage

import (
	"errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/xx_network/primitives/id"
	"time"
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
	fid := &factId{}
	copy(fid[:], factHash)
	m.facts[*fid].Verified = true
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
	f := Fact{
		FactHash:  factHash,
		UserId:    userID,
		Fact:      fact,
		FactType:  uint8(factType),
		Signature: signature,
		Verified:  false,
		Timestamp: time.Now(),
	}
	tv := TwilioVerification{
		ConfirmationId: confirmationID,
		FactHash:       factHash,
	}
	fid := factId{}
	copy(fid[:], factHash)
	m.facts[fid] = &f
	m.twilioVerifications[confirmationID] = &tv
	return nil
}
func (m *MapImpl) VerifyFactTwilio(confirmationId string) error {
	fid := factId{}
	copy(fid[:], m.twilioVerifications[confirmationId].FactHash)
	m.facts[fid].Verified = true
	delete(m.twilioVerifications, confirmationId)
	return nil
}

func (m *MapImpl) Search(factHashs [][]byte) []*User {
	users := map[id.ID]User{}
	for _, h := range factHashs {
		fid := factId{}
		copy(fid[:], h)
		if f, ok := m.facts[fid]; ok {
			uid, err := id.Unmarshal(f.UserId)
			if err != nil {
				jww.ERROR.Print("Failed to decode uid %+v: %+v", f.UserId, err)
			}
			users[*uid] = User{
				Id: uid.Marshal(),
			}
		}
	}
	var result []*User
	for _, u := range users {
		result = append(result, &u)
	}
	return result
}
