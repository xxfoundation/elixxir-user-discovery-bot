////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Handles the database backend for udb storage

package storage

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
func (db *DatabaseImpl) GetFact(confirmationId []byte) (*Fact, error) {
	result := &Fact{}
	err := db.db.First(&result, "confirmation_id = ?", confirmationId).Error
	return result, err
}

// Delete a fact by confirmation ID
func (db *DatabaseImpl) DeleteFact(confirmationId []byte) error {
	return db.db.Delete(&Fact{
		ConfirmationId: confirmationId,
	}).Error
}

// Confirm a fact by confirmation ID
func (db *DatabaseImpl) ConfirmFact(confirmationId []byte) error {
	return db.db.Model(&Fact{}).Where("confirmation_id = ?", confirmationId).UpdateColumn("verification_status", 1).Error
}
