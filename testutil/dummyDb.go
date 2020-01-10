package testutil

import (
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/user-discovery-bot/storage"
)

type MockDatabase struct {
	User    *storage.User
	newuser bool
	newkey  bool
}

func (d MockDatabase) UpsertUser(user *storage.User) error {
	jww.INFO.Printf("Called UpsertUser with params: { user: %+v }", user)
	return nil
}

func (d MockDatabase) GetUser(id []byte) (*storage.User, error) {
	jww.INFO.Printf("Called GetUser with params: { id: %+v }", id)
	return d.User, nil
}

func (d MockDatabase) GetUserByValue(value string) (*storage.User, error) {
	jww.INFO.Printf("Called GetUserByValue with params: { value: %+v }", value)
	if d.newuser {
		return nil, errors.New("Unable to find any user with that value (mocked)")
	}
	return d.User, nil
}

func (d MockDatabase) GetUserByKeyId(keyId string) (*storage.User, error) {
	jww.INFO.Printf("Called GetUserByKeyId with params: { keyId: %+v }", keyId)
	if d.newkey {
		return nil, errors.New("Mock: new key")
	} else if d.newuser {
		return &storage.User{
			Id:        nil,
			Value:     "",
			ValueType: 0,
			KeyId:     "",
			Key:       nil,
		}, nil
	}
	return d.User, nil
}

func (d MockDatabase) DeleteUser(id []byte) error {
	jww.INFO.Printf("Called DeleteUser with params: { id: %+v }", id)
	return nil
}

func GetMockDatabase(id, value, keyid, key string, newuser, newkey bool) *MockDatabase {
	db := MockDatabase{User: &storage.User{
		Id:        []byte(id),
		Value:     value,
		ValueType: 0,
		KeyId:     keyid,
		Key:       []byte(key),
	},
		newuser: newuser,
		newkey:  newkey}
	return &db
}
