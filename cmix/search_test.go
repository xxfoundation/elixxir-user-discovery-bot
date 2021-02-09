package cmix

import (
	"github.com/golang/protobuf/proto"
	"gitlab.com/elixxir/client/interfaces/message"
	"gitlab.com/elixxir/client/switchboard"
	"gitlab.com/elixxir/client/ud"
	"gitlab.com/elixxir/crypto/e2e"
	"gitlab.com/elixxir/primitives/fact"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/xx_network/primitives/id"
	"testing"
	"time"
)

func TestManager_SearchProcess(t *testing.T) {
	c := &MockClient{
		sw: switchboard.New(),
	}
	db := storage.NewTestDB(t)
	manager := &Manager{
		searchChan: make(chan message.Receive, 1),
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
	fid := []byte("facthash")
	err = db.InsertFact(&storage.Fact{
		Hash:     fid,
		UserId:   uid.Marshal(),
		Fact:     "water is wet",
		Type:     0,
		Verified: true,
	})
	if err != nil {
		t.Errorf("Failed to insert dummy fact: %+v", err)
	}

	go manager.SearchProcess()
	payload, err := proto.Marshal(&ud.SearchSend{
		Fact: []*ud.HashFact{
			{
				Hash: fid,
				Type: 0,
			},
		},
		CommID: 0,
	})
	manager.searchChan <- message.Receive{
		ID:          e2e.NewMessageID([]byte("test"), 420),
		Payload:     payload,
		MessageType: message.UdSearch,
		Sender:      uid,
		Timestamp:   time.Now(),
		Encryption:  0,
	}
	time.Sleep(time.Second * 5)
	if !c.r {
		t.Error("Failed to receive search message")
	}
}

func TestManager_handleSearch(t *testing.T) {
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
	fid := []byte("facthash")
	err = db.InsertFact(&storage.Fact{
		Hash:     fid,
		UserId:   uid.Marshal(),
		Fact:     "water is wet",
		Type:     0,
		Verified: true,
	})
	if err != nil {
		t.Errorf("Failed to insert dummy fact: %+v", err)
	}
	resp := manager.handleSearch(&ud.SearchSend{
		Fact: []*ud.HashFact{
			{
				Hash: fid,
				Type: 0,
			},
		},
		CommID: 0,
	}, id.NewIdFromString("zezima", id.User, t))
	if resp.Error != "" {
		t.Errorf("failed to handle search: %+v", resp.Error)
	}
	if len(resp.Contacts) != 1 {
		t.Errorf("Did not receive expected number of contacts")
	}

	resp = manager.handleSearch(&ud.SearchSend{
		Fact: []*ud.HashFact{
			{
				Hash: fid,
				Type: int32(fact.Nickname),
			},
		},
		CommID: 0,
	}, id.NewIdFromString("zezima", id.User, t))
	if len(resp.Contacts) != 0 {
		t.Errorf("Should not be able to search with nickname")
	}
}
