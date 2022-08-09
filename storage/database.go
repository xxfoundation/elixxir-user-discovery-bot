////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Handles low level Database control and interfaces

package storage

import (
	"fmt"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/xx_network/primitives/id"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"gorm.io/gorm/logger"
	"sync"
	"time"
)

// Interface declaration for Storage methods
type database interface {
	CheckUser(username string, id *id.ID) error

	InsertUser(user *User) error
	GetUser(id []byte) (*User, error)
	DeleteUser(id []byte) error

	InsertFact(fact *Fact) error
	MarkFactVerified(factHash []byte) error
	DeleteFact(factHash []byte) error

	InsertFactTwilio(userID, factHash, signature []byte, factType uint, confirmationID string) error
	MarkTwilioFactVerified(confirmationId string) error

	Search(factHashes [][]byte) ([]*User, error)

	StartFactManager(i time.Duration) chan chan bool

	InsertChannelIdentity(identity *ChannelIdentity) error
	GetChannelIdentity(id []byte) (*ChannelIdentity, error)
}

// Struct implementing the Database Interface with an underlying DB
type DatabaseImpl struct {
	db *gorm.DB // Stored database connection
}

// ID type for facts map
type factId [32]byte

// Struct implementing the Database Interface with an underlying Map
type MapImpl struct {
	users             map[id.ID]*User
	usernames         map[id.ID]*Fact
	facts             map[factId]*Fact
	channelIdentities map[id.ID]*ChannelIdentity
	sync.RWMutex
}

// Struct defining the users table for the database
type User struct {
	Id        []byte `gorm:"primaryKey"`
	Username  string `gorm:"not null;unique"`
	RsaPub    string `gorm:"not null;unique"`
	DhPub     []byte `gorm:"not null;unique"`
	Salt      []byte `gorm:"not null"`
	Signature []byte `gorm:"not null"`
	// Time in which user registered with the network (ie permissioning)
	RegistrationTimestamp time.Time `gorm:"not null"`
	Facts                 []Fact    `gorm:"constraint:OnDelete:CASCADE"`
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
	Hash           []byte `gorm:"primaryKey"`
	UserId         []byte `gorm:"not null"`
	Fact           string
	Type           uint8 `gorm:"not null"`
	Signature      []byte
	Verified       bool      `gorm:"not null"`
	Timestamp      time.Time `gorm:"not null"`
	ConfirmationId string
}

// ChannelIdentity represents the data which is stored by user discovery on a User's channel registration
type ChannelIdentity struct {
	UserId    []byte `gorm:"primaryKey"`
	PublicKey []byte `gorm:"not null"`
	Lease     int64  `gorm:"not null"`
	Banned    bool   `gorm:"default:false"`
	User      User   `gorm:"foreignKey:UserId;constraint:OnDelete:CASCADE;"`
}

// Initialize the Database interface with database backend
// Returns a database interface and error
func newDatabase(username, password, dbName, address,
	port string) (database, error) {

	var err error
	var db *gorm.DB
	// Connect to the database if the correct information is provided
	if address != "" && port != "" {
		// Create the database connection
		connectString := fmt.Sprintf(
			"host=%s port=%s user=%s dbname=%s sslmode=disable",
			address, port, username, dbName)
		// Handle empty database password
		if len(password) > 0 {
			connectString += fmt.Sprintf(" password=%s", password)
		}
		db, err = gorm.Open(postgres.Open(connectString), &gorm.Config{
			Logger: logger.New(jww.TRACE, logger.Config{LogLevel: logger.Info}),
		})
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
			users:             map[id.ID]*User{},
			facts:             map[factId]*Fact{},
			usernames:         map[id.ID]*Fact{},
			channelIdentities: map[id.ID]*ChannelIdentity{},
		}

		return database(mapImpl), nil
	}

	// Get and configure the internal database ConnPool
	sqlDb, err := db.DB()
	if err != nil {
		return database(&DatabaseImpl{}), errors.Errorf("Unable to configure database connection pool: %+v", err)
	}
	// SetMaxIdleConns sets the maximum number of connections in the idle connection pool.
	sqlDb.SetMaxIdleConns(10)
	// SetMaxOpenConns sets the maximum number of open connections to the Database.
	sqlDb.SetMaxOpenConns(50)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be idle.
	sqlDb.SetConnMaxIdleTime(10 * time.Minute)
	// SetConnMaxLifetime sets the maximum amount of time a connection may be reused.
	sqlDb.SetConnMaxLifetime(12 * time.Hour)

	// Initialize the database schema
	// WARNING: Order is important. Do not change without database testing
	models := []interface{}{User{}, Fact{}, ChannelIdentity{}}
	for _, model := range models {
		err = db.AutoMigrate(model)
		if err != nil {
			return database(&DatabaseImpl{}), err
		}
	}

	jww.INFO.Println("Database backend initialized successfully!")
	return &DatabaseImpl{db: db}, nil
}
