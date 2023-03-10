////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Handles the Map backend for udb storage

package storage

import (
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/primitives/fact"
	"gitlab.com/xx_network/primitives/id"
	"strings"
	"time"
)

func (m *MapImpl) CheckUser(username string, id *id.ID) error {
	for _, f := range m.facts {
		if f.Type == uint8(Username) && strings.ToUpper(f.Fact) == strings.ToUpper(username) {
			return errors.New("Username already exists")
		}
	}
	return nil
}

// Insert a new user
func (m *MapImpl) InsertUser(user *User) error {
	uid, err := id.Unmarshal(user.Id)
	if err != nil {
		return err
	}
	m.users[*uid] = user
	for _, f := range user.Facts {
		err = m.InsertFact(&f)
		if err != nil {
			return err
		}
	}
	return nil
}

// Get a user from the map
func (m *MapImpl) GetUser(uid []byte) (*User, error) {
	nuid, err := id.Unmarshal(uid)
	if err != nil {
		return nil, err
	}
	u, _ := m.users[*nuid]
	un, ok := m.usernames[*nuid]
	if ok {
		u.Facts = []Fact{*un}
	}
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
	if fact.Type == uint8(Username) {
		m.usernames[*uid] = fact
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
func (m *MapImpl) InsertFactTwilio(userID, factHash, signature []byte, factType uint, confirmationID string) error {
	f := Fact{
		Hash:           factHash,
		UserId:         userID,
		Type:           uint8(factType),
		Signature:      signature,
		Verified:       false,
		Timestamp:      time.Now(),
		ConfirmationId: confirmationID,
	}
	fid := factId{}
	copy(fid[:], factHash)
	m.facts[fid] = &f
	return nil
}

// Verify a twilio fact
func (m *MapImpl) MarkTwilioFactVerified(confirmationId string) error {
	for _, myFact := range m.facts {
		if myFact.ConfirmationId == confirmationId {
			myFact.Verified = true
		}
	}
	return nil
}

// Search for users by fact hashes
func (m *MapImpl) Search(factHashes [][]byte) ([]*User, error) {
	users := map[id.ID]*User{}
	unames := map[id.ID]Fact{}
	for _, h := range factHashes {
		fid := factId{}
		copy(fid[:], h)
		if f, ok := m.facts[fid]; ok {
			uid, err := id.Unmarshal(f.UserId)
			if err != nil {
				return nil, errors.WithMessagef(err, "Failed to decode uid %+v", f.UserId)
			}
			if u, found := users[*uid]; found {
				u.Facts = append(u.Facts, *f)
			} else {
				u, ok := m.users[*uid]
				if !ok {
					return nil, errors.New("no user associated with hash, this should not be possible")
				}
				users[*uid] = &User{
					Username:              u.Username,
					Id:                    u.Id,
					RsaPub:                u.RsaPub,
					DhPub:                 u.DhPub,
					Salt:                  u.Salt,
					Signature:             u.Signature,
					RegistrationTimestamp: u.RegistrationTimestamp,
					Facts:                 []Fact{*f},
				}
			}
			if f.Type == uint8(fact.Username) {
				unames[*uid] = *f
			}
		}
	}
	var result []*User
	for i, u := range users {
		if _, ok := unames[i]; !ok {
			uid, err := id.Unmarshal(u.Id)
			if err != nil {
				return nil, errors.WithMessagef(err, "Failed to decode uid %+v", u.Id)
			}
			if uname, ok := m.usernames[*uid]; !ok {
				jww.WARN.Println("No username associated with user")
			} else {
				u.Facts = append(u.Facts, *uname)
			}
		}
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
