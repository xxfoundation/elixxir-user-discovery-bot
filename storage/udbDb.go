////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Handles the database backend for udb storage

package storage

import (
	"fmt"
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"git.xx.network/xx_network/primitives/id"
	"time"
)

// Check if a username is available
func (db *DatabaseImpl) CheckUser(username string, id *id.ID, rsaPem string) error {
	var err error
	var facts []*Fact
	var count int
	err = db.db.Where("type = ? AND fact = ?", Username, username).Find(&facts).Count(&count).Error
	if err != nil {
		return errors.WithMessage(err, "Failed to check facts for desired username")
	}
	if count > 0 {
		return errors.New("error: username in use")
	}

	err = db.db.Where("type = ? AND user_id = ?", Username, id.Marshal()).Find(&facts).Count(&count).Error
	if err != nil {
		return errors.WithMessage(err, "Failed to check facts for usernames registered to user")
	}
	if count > 0 {
		return errors.New("error: user has already registered a username")
	}

	return nil
}

// Insert a new user object
func (db *DatabaseImpl) InsertUser(user *User) error {
	return db.db.Create(user).Error
}

// Retrieve a user by ID
func (db *DatabaseImpl) GetUser(id []byte) (*User, error) {
	result := &User{}
	err := db.db.First(&result, "id = ?", id).Error
	return result, err
}

// Delete a user by ID
func (db *DatabaseImpl) DeleteUser(id []byte) error {
	return db.db.Delete(&User{
		Id: id,
	}).Error
}

// Insert a new fact
func (db *DatabaseImpl) InsertFact(fact *Fact) error {
	return db.db.Create(fact).Error
}

// Retreive a fact by confirmation ID
func (db *DatabaseImpl) MarkFactVerified(factHash []byte) error {
	return db.db.Model(&Fact{}).Where("hash = ?", factHash).UpdateColumn("verified", true).Error
}

// Delete a fact by confirmation ID
func (db *DatabaseImpl) DeleteFact(factHash []byte) error {
	return db.db.Delete(&Fact{
		Hash: factHash,
	}).Error
}

// Insert a twilio-verified fact
func (db *DatabaseImpl) InsertFactTwilio(userID, factHash, signature []byte, factType uint, fact, confirmationID string) error {
	f := &Fact{
		Hash:      factHash,
		UserId:    userID,
		Fact:      "fact",
		Type:      uint8(factType),
		Signature: signature,
		Verified:  false,
		Timestamp: time.Now(),
	}

	tv := &TwilioVerification{
		ConfirmationId: confirmationID,
		FactHash:       factHash,
	}

	tf := func(tx *gorm.DB) error {
		var err error
		if err = tx.Set("gorm:insert_option", "ON CONFLICT (hash) DO UPDATE SET timestamp = NOW()").Create(f).Error; err != nil {
			return err
		}

		if err = tx.Set("gorm:insert_option", fmt.Sprintf("ON CONFLICT (fact_hash) DO UPDATE SET confirmation_id = '%+v'", confirmationID)).Create(tv).Error; err != nil {
			return err
		}
		return nil
	}

	return db.db.Transaction(tf)
}

// Verify a fact through twilio
func (db *DatabaseImpl) MarkTwilioFactVerified(confirmationId string) error {
	tf := func(tx *gorm.DB) error {
		var err error
		tv := &TwilioVerification{}
		err = tx.Where("confirmation_id = ?", confirmationId).First(tv).Error
		if err != nil {
			return err
		}
		err = tx.Model(&Fact{}).Where("hash = ?", tv.FactHash).UpdateColumn("verified", true).Error
		if err != nil {
			return err
		}
		err = tx.Delete(tv, "confirmation_id = ?", confirmationId).Error
		return err
	}
	return db.db.Transaction(tf)
}

// Search for users by facts
func (db *DatabaseImpl) Search(factHashes [][]byte) ([]*User, error) {
	var facts []*Fact
	err := db.db.Where("hash in (?) and verified", factHashes).Find(&facts).Error
	if err != nil {
		return nil, err
	}

	var found = make(map[id.ID][]Fact)
	var usernames = make(map[id.ID]Fact)
	for _, f := range facts {
		// Unmarshal uid for this fact
		uid, err := id.Unmarshal(f.UserId)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to unmarshal uid")
		}

		// Add user if not hit already, add to list of facts otherwise
		if fl, ok := found[*uid]; ok {
			found[*uid] = append(fl, *f)
		} else {
			found[*uid] = []Fact{*f}
		}

		// Username handling
		if f.Type == uint8(Username) {
			if _, ok := usernames[*uid]; ok {
				continue
			} else {
				usernames[*uid] = *f
			}
		}
	}
	var users []*User
	for uid, fl := range found {
		u := &User{}
		err = db.db.Preload("Facts", "type = 0").Take(u, "id = ?", uid.Marshal()).Error
		if err != nil {
			return nil, err
		}
		if _, ok := usernames[uid]; ok {
			u.Facts = fl
		} else {
			u.Facts = append(u.Facts, fl...)
		}
		users = append(users, u)
	}

	return users, nil
}

func (db *DatabaseImpl) StartFactManager(i time.Duration) chan chan bool {
	stopChan := make(chan chan bool)
	go func() {
		interval := time.NewTicker(i)
		select {
		case <-interval.C:
			tf := func(tx *gorm.DB) error {
				var err error
				var facts []*Fact
				err = db.db.Where(&facts, "verified = false AND timestamp <= (NOW() - INTERVAL '5 minutes')").Error
				if err != nil {
					return errors.Errorf("error retrieving out of date unverified facts: %+v", err)
				}
				for _, f := range facts {
					err = db.db.Delete(f, "hash = ?", f.Hash).Error
					if err != nil {
						return errors.Errorf("error deleting out of date fact %+v: %+v", f.Hash, err)
					}
				}
				return err
			}
			err := db.db.Transaction(tf)
			if err != nil {
				jww.ERROR.Print(err)
			}
		case kc := <-stopChan:
			kc <- true
			return
		}
	}()
	return stopChan
}
