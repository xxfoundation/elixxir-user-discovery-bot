package cmix

import (
	"gitlab.com/elixxir/client/ud"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/xx_network/primitives/id"
	"testing"
)

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
}
