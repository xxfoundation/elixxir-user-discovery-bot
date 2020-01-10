package udb

import (
	"encoding/base64"
	"gitlab.com/elixxir/client/cmixproto"
	"gitlab.com/elixxir/primitives/id"
	fingerprint2 "gitlab.com/elixxir/user-discovery-bot/fingerprint"
	"gitlab.com/elixxir/user-discovery-bot/testutil"
	"strings"
	"testing"
	"time"
)

func TestSearchListener_Hear(t *testing.T) {
	ch := make(chan string, 1)
	mockSender := testutil.NewMockSender(ch)
	mockDB := testutil.GetMockDatabase("testid", "testval", "testkeyid", "testkey", false, false)
	var listener = SearchListener{
		Sender: mockSender,
		db:     mockDB,
	}
	var senderID = id.NewUserFromUint(5, t)

	msg := NewMessage("EMAIL rick@elixxir.io", cmixproto.Type_UDB_SEARCH, senderID)

	listener.Hear(msg, false)

	select {
	case msgStr := <-ch:
		if strings.Contains(msgStr, "SEARCH") && strings.Contains(msgStr, "FOUND") {
			t.Log("Received found message for search")
		} else {
			t.Error("Did not receive found message for search")
		}
	case <-time.After(5 * time.Second):
		t.Error("Test timed out")
	}

}

func TestGetKeyListener_Hear(t *testing.T) {
	ch := make(chan string, 1)
	mockSender := testutil.NewMockSender(ch)
	mockDB := testutil.GetMockDatabase("testid", "testval", "testkeyid", "testkey", false, false)
	pubKeyBits := "S8KXBczy0jins9uS4LgBPt0bkFl8t00MnZmExQ6GcOcu8O7DKgAsNzLU7a" +
		"+gMTbIsS995IL/kuFF8wcBaQJBY23095PMSQ/nMuetzhk9HdXxrGIiKBo3C/n4SClpq4H+PoF9XziEVKua8JxGM2o83KiCK3tNUpaZbAAElkjueY7wuD96h4oaA+WV5Nh87cnIZ+fAG0uLve2LSHZ0FBZb3glOpNAOv7PFWkvN2BO37ztOQCXTJe72Y5ReoYn7nWVNxGUh0ilal+BRuJt1GZ7whOGDRE0IXfURIoK2yjyAnyZJWWMhfGsL5S6iL4aXUs03mc8BHKRq3HRjvTE10l3YFA=="

	pubKey := make([]byte, 256)
	pubKey, _ = base64.StdEncoding.DecodeString(pubKeyBits)

	fingerprint := fingerprint2.Fingerprint(pubKey)
	var listener = GetKeyListener{
		Sender: mockSender,
		db:     mockDB,
	}
	var senderID = id.NewUserFromUint(5, t)

	msg := NewMessage(fingerprint, cmixproto.Type_UDB_GET_KEY, senderID)
	listener.Hear(msg, false)

	time.Sleep(3 * time.Second)

	select {
	case msgStr := <-ch:
		if strings.Contains(msgStr, "GETKEY") && !strings.Contains(msgStr, "NOTFOUND") {
			t.Log("Received found message for getkey")
		} else {
			t.Error("Did not receive found message for search")
		}
	case <-time.After(5 * time.Second):
		t.Error("Test timed out")
	}
}

func TestPushKeyListener_Hear(t *testing.T) {
	ch := make(chan string, 1)
	mockSender := testutil.NewMockSender(ch)
	mockDB := testutil.GetMockDatabase("testid", "testval", "testkeyid", "testkey", false, true)
	pubKeyBits := "S8KXBczy0jins9uS4LgBPt0bkFl8t00MnZmExQ6GcOcu8O7DKgAsNzLU7a" +
		"+gMTbIsS995IL/kuFF8wcBaQJBY23095PMSQ/nMuetzhk9HdXxrGIiKBo3C/n4SClpq4H+PoF9XziEVKua8JxGM2o83KiCK3tNUpaZbAAElkjueY7wuD96h4oaA+WV5Nh87cnIZ+fAG0uLve2LSHZ0FBZb3glOpNAOv7PFWkvN2BO37ztOQCXTJe72Y5ReoYn7nWVNxGUh0ilal+BRuJt1GZ7whOGDRE0IXfURIoK2yjyAnyZJWWMhfGsL5S6iL4aXUs03mc8BHKRq3HRjvTE10l3YFA=="

	var listener = PushKeyListener{
		Sender: mockSender,
		db:     mockDB,
	}
	var senderID = id.NewUserFromUint(5, t)

	msg := NewMessage("myKeyId "+pubKeyBits, cmixproto.Type_UDB_PUSH_KEY, senderID)
	listener.Hear(msg, false)

	time.Sleep(3 * time.Second)

	select {
	case msgStr := <-ch:
		if strings.Contains(msgStr, "PUSHKEY COMPLETE") {
			t.Log("Received complete message for pushkey")
		} else {
			t.Error("Did not receive pushkey complete")
		}
	case <-time.After(5 * time.Second):
		t.Error("Test timed out")
	}
}

func TestRegisterListener_Hear(t *testing.T) {
	ch := make(chan string, 1)
	mockSender := testutil.NewMockSender(ch)
	mockDB := testutil.GetMockDatabase("testid", "testval", "testkeyid", "testkey", true, false)
	pubKeyBits := "S8KXBczy0jins9uS4LgBPt0bkFl8t00MnZmExQ6GcOcu8O7DKgAsNzLU7a" +
		"+gMTbIsS995IL/kuFF8wcBaQJBY23095PMSQ/nMuetzhk9HdXxrGIiKBo3C/n4SClpq4H+PoF9XziEVKua8JxGM2o83KiCK3tNUpaZbAAElkjueY7wuD96h4oaA+WV5Nh87cnIZ+fAG0uLve2LSHZ0FBZb3glOpNAOv7PFWkvN2BO37ztOQCXTJe72Y5ReoYn7nWVNxGUh0ilal+BRuJt1GZ7whOGDRE0IXfURIoK2yjyAnyZJWWMhfGsL5S6iL4aXUs03mc8BHKRq3HRjvTE10l3YFA=="

	pubKey := make([]byte, 256)
	pubKey, _ = base64.StdEncoding.DecodeString(pubKeyBits)

	fingerprint := fingerprint2.Fingerprint(pubKey)

	var listener = RegisterListener{
		Sender:    mockSender,
		db:        mockDB,
		blacklist: *InitBlackList("./blacklists/bannedNames.txt"),
	}
	var senderID = id.NewUserFromUint(5, t)

	msg := NewMessage("EMAIL rick@elixxir.io "+fingerprint, cmixproto.Type_UDB_REGISTER, senderID)
	listener.Hear(msg, false)

	time.Sleep(3 * time.Second)

	select {
	case msgStr := <-ch:
		if strings.Contains(msgStr, "REGISTRATION COMPLETE") {
			t.Log("Received registration complete message")
		} else {
			t.Error("Did not receive registration complete")
		}
	case <-time.After(5 * time.Second):
		t.Error("Test timed out")
	}
}
