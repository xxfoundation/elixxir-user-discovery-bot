package cmix

import (
	"github.com/golang/protobuf/proto"
	"gitlab.com/elixxir/client/interfaces/message"
	"gitlab.com/elixxir/client/switchboard"
	"gitlab.com/elixxir/client/ud"
	"gitlab.com/elixxir/crypto/e2e"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/xx_network/primitives/id"
	"testing"
	"time"
)

func TestManager_LookupProcess(t *testing.T) {
	c := &MockClient{
		sw: switchboard.New(),
	}
	db := storage.NewTestDB(t)
	manager := &Manager{
		lookupChan: make(chan message.Receive, 1),
		db:         db,
		client:     c,
	}
	uid := id.NewIdFromString("zezima", id.User, t)
	err := db.InsertUser(&storage.User{
		Id:     uid.Marshal(),
		RsaPub: "rsapub",
	})
	if err != nil {
		t.Errorf("Failed to insert dummy user: %+v", err)
	}
	go manager.LookupProcess()
	payload, err := proto.Marshal(&ud.LookupSend{
		UserID: uid.Marshal(),
		CommID: 1,
	})
	if err != nil {
		t.Errorf("Failed to marshal payload: %+v", err)
	}
	manager.lookupChan <- message.Receive{
		ID:          e2e.NewMessageID([]byte("test"), uint64(420)),
		Payload:     payload,
		MessageType: message.UdLookup,
		Sender:      id.NewIdFromString("zezima", id.User, t),
		Timestamp:   time.Now(),
		Encryption:  0,
	}
	time.Sleep(time.Second * 5)
	if !c.r {
		t.Error("Failed to receive lookup message")
	}
}

func TestManager_handleLookup(t *testing.T) {
	db := storage.NewTestDB(t)
	manager := &Manager{
		db: db,
	}
	uid := id.NewIdFromString("zezima", id.User, t)
	err := db.InsertUser(&storage.User{
		Id:     uid.Marshal(),
		RsaPub: "rsapub",
	})
	if err != nil {
		t.Errorf("Failed to insert dummy user: %+v", err)
	}
	resp := manager.handleLookup(&ud.LookupSend{
		UserID: uid.Marshal(),
		CommID: 0,
	}, uid)
	if resp.Error != "" {
		t.Errorf("Failed to handle lookup: %+v", resp.Error)
	}
}
