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
	copy(factid[:], fact.Hash)
	m.facts[factid] = fact

	fact.Timestamp = time.Now()

	return nil
}

// Verify fact in mapimpl
func (m *MapImpl) MarkFactVerified(factHash []byte) error {
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

// Insert a twilio-verified fact
func (m *MapImpl) InsertFactTwilio(userID, factHash, signature []byte, factType uint, fact, confirmationID string) error {
	f := Fact{
		Hash:      factHash,
		UserId:    userID,
		Fact:      fact,
		Type:      uint8(factType),
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
	m.fhToVerification[fid] = &tv
	return nil
}

// Verify a twilio fact
func (m *MapImpl) MarkTwilioFactVerified(confirmationId string) error {
	fid := factId{}
	copy(fid[:], m.twilioVerifications[confirmationId].FactHash)
	m.facts[fid].Verified = true
	delete(m.twilioVerifications, confirmationId)
	delete(m.fhToVerification, fid)
	return nil
}

// Search for users by fact hashes
func (m *MapImpl) Search(factHashes [][]byte) ([]*User, error) {
	users := map[id.ID]*User{}
	for _, h := range factHashes {
		fid := factId{}
		copy(fid[:], h)
		if f, ok := m.facts[fid]; ok {
			uid, err := id.Unmarshal(f.UserId)
			if err != nil {
				jww.ERROR.Printf("Failed to decode uid %+v: %+v", f.UserId, err)
			}
			users[*uid] = m.users[*uid]
		}
	}
	var result []*User
	for _, u := range users {
		result = append(result, u)
	}
	return result, nil
}

func (m *MapImpl) StartFactManager(i time.Duration) chan chan bool {
	stopChan := make(chan chan bool)
	go func() {
		interval := time.NewTicker(i)
		select {
		case <-interval.C:
			for factId, f := range m.facts {
				if !f.Verified && f.Timestamp.Before(time.Now().Add(-5*time.Minute)) {
					delete(m.facts, factId)
				}
			}
		case kc := <-stopChan:
			kc <- true
			return
		}
	}()
	return stopChan
}
