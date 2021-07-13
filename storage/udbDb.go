////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Handles the database backend for udb storage

package storage

import (
	"github.com/jinzhu/gorm"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/xx_network/primitives/id"
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
	return db.db.Model(&Fact{}).Where("hash = ?", factHash).UpdateColumn("verified", "true").Error
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
	}

	tv := &TwilioVerification{
		ConfirmationId: confirmationID,
		FactHash:       factHash,
	}

	tf := func(tx *gorm.DB) error {
		var err error
		if err = tx.Create(f).Error; err != nil {
			return err
		}
		if err = tx.Create(tv).Error; err != nil {
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

	var found map[id.ID]bool
	found = make(map[id.ID]bool)
	var users []*User
	for _, f := range facts {
		uid, err := id.Unmarshal(f.UserId)
		if err != nil {
			return nil, errors.WithMessage(err, "failed to unmarshal uid")
		}
		if _, ok := found[*uid]; ok {
			continue
		}
		u := &User{}
		err = db.db.Preload("Facts").Take(u, "id = ?", f.UserId).Error
		if err != nil {
			return nil, err
		}
		found[*uid] = true
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
