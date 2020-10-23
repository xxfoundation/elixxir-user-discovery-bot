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
	"gitlab.com/xx_network/primitives/id"
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
		return errors.WithMessage(err, "Failed to check facts for usernames registerd to user")
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

// Retreive a user by ID
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
func (db *DatabaseImpl) Search(factHashs [][]byte) []*User {
	var facts []*Fact
	db.db.Select(&Fact{}, "hash in ?", factHashs).Find(&facts)

	var users []*User
	for _, f := range facts {
		users = append(users, &User{
			Id: f.UserId,
		})
	}

	return users
}
