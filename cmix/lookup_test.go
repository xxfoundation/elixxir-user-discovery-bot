package cmix

import (
	"gitlab.com/elixxir/client/ud"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/xx_network/primitives/id"
	"testing"
)

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
