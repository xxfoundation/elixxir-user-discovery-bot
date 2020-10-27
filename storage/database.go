////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Handles low level Database control and interfaces

package storage

import (
	"fmt"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/xx_network/primitives/id"
	"sync"
	"time"
)

// Interface declaration for Storage methods
type database interface {
	CheckUser(username string, id *id.ID, rsaPem string) error

	InsertUser(user *User) error
	GetUser(id []byte) (*User, error)
	DeleteUser(id []byte) error

	InsertFact(fact *Fact) error
	MarkFactVerified(factHash []byte) error
	DeleteFact(factHash []byte) error

	InsertFactTwilio(userID, factHash, signature []byte, factType uint, fact, confirmationID string) error
	MarkTwilioFactVerified(confirmationId string) error

	Search(factHashes [][]byte) []*User

	StartFactManager(i time.Duration) chan chan bool
}

// Struct implementing the Database Interface with an underlying DB
type DatabaseImpl struct {
	db *gorm.DB // Stored database connection
}

// ID type for facts map
type factId [32]byte

// Struct implementing the Database Interface with an underlying Map
type MapImpl struct {
	users               map[id.ID]*User
	facts               map[factId]*Fact
	twilioVerifications map[string]*TwilioVerification
	sync.RWMutex
}

// Struct defining the users table for the database
type User struct {
	Id        []byte `gorm:"primary_key"`
	RsaPub    string `gorm:"NOT NULL"`
	DhPub     []byte `gorm:"NOT NULL"`
	Salt      []byte `gorm:"NOT NULL"`
	Signature []byte `gorm:"NOT NULL"`
	Facts     []Fact `gorm:"foreignkey:UserId;association_foreignkey:Id"`
}

// Fact type enum
type FactType uint8

const (
	Username FactType = iota
	SMS
	Email
)

func (f FactType) String() string {
	return [...]string{"Username", "SMS", "Email"}[f]
}

// Struct defining the facts table in the database
type Fact struct {
	Hash         []byte             `gorm:"primary_key"`
	UserId       []byte             `gorm:"NOT NULL;type:bytea REFERENCES users(Id)"`
	Fact         string             `gorm:"NOT NULL"`
	Type         uint8              `gorm:"NOT NULL"`
	Signature    []byte             `gorm:"NOT NULL"`
	Verified     bool               `gorm:"NOT NULL"`
	Timestamp    time.Time          `gorm:"NOT NULL"`
	Verification TwilioVerification `gorm:"foreignkey:FactHash;association_foreignkey:Hash"`
}

// Struct defining twilio_verifications table
type TwilioVerification struct {
	ConfirmationId string `gorm:"primary_key"`
	FactHash       []byte `gorm:"NOT NULL;type:bytea REFERENCES facts(Hash)"`
}

// Initialize the Database interface with database backend
// Returns a Storage interface, Close function, and error
func newDatabase(username, password, database, address,
	port string) (*Storage, func() error, error) {
	var err error
	var db *gorm.DB
	//connect to the database if the correct information is provided
	if address != "" && port != "" {
		// Create the database connection
		connectString := fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s sslmode=disable",
			address, port, username, database)
		// Handle empty database password
		if len(password) > 0 {
			connectString += fmt.Sprintf(" password=%s", password)
		}
		db, err = gorm.Open("postgres", connectString)
	}

	// Return the map-backend interface
	// in the event there is a database error or information is not provided
	if (address == "" || port == "") || err != nil {

		if err != nil {
			jww.WARN.Printf("Unable to initialize database backend: %+v", err)
		} else {
			jww.WARN.Printf("Database backend connection information not provided")
		}

		defer jww.INFO.Println("Map backend initialized successfully!")

		mapImpl := &MapImpl{
			users:               map[id.ID]*User{},
			facts:               map[factId]*Fact{},
			twilioVerifications: map[string]*TwilioVerification{},
		}

		return &Storage{mapImpl}, func() error { return nil }, nil
	}

	// Initialize the database logger
	db.SetLogger(jww.TRACE)
	db.LogMode(true)

	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	db.DB().SetMaxIdleConns(10)
	// SetMaxOpenConns sets the maximum number of open connections to the database.
	db.DB().SetMaxOpenConns(100)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	db.DB().SetConnMaxLifetime(24 * time.Hour)

	// Initialize the database schema
	// WARNING: Order is important. Do not change without database testing
	models := []interface{}{User{}, Fact{}, TwilioVerification{}}
	for _, model := range models {
		err = db.AutoMigrate(model).Error
		if err != nil {
			return nil, func() error { return nil }, err
		}
	}

	jww.INFO.Println("Database backend initialized successfully!")
	return &Storage{&DatabaseImpl{db: db}}, db.Close, nil
}
