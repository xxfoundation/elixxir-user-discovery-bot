package cmix

import (
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/client/single"
	"gitlab.com/elixxir/client/ud"
	"gitlab.com/elixxir/primitives/fact"
	"gitlab.com/xx_network/primitives/id"
	"time"
)

func (m *Manager) searchCallback(payload []byte, c single.Contact) {
	searchMsg := &ud.SearchSend{}
	if err := proto.Unmarshal(payload, searchMsg); err != nil {
		jww.ERROR.Printf("Failed to unmarshal search request from %s: %+v",
			c.GetPartner(), err)
		return
	}

	jww.INFO.Printf("Search request from %s: %v", c.GetPartner(), payload)

	response := m.handleSearch(searchMsg, c)

	jww.INFO.Printf("Search for %+v completed & found %+v", searchMsg.Fact, response.Contacts)

	marshaledResponse, err := proto.Marshal(response)
	if err != nil {
		jww.ERROR.Printf("Failed to marshal request to search request from "+
			"%s: %+v", c.GetPartner(), err)
		return
	}

	// TODO: make timeout come from config file, default to 1 minute
	err = m.singleUse.RespondSingleUse(c, marshaledResponse, 1*time.Minute)
	if err != nil {
		jww.ERROR.Printf("Failed to send single-use response to to search "+
			"request from %s: %+v", c.GetPartner(), err)
		return
	}
}

func (m *Manager) handleSearch(msg *ud.SearchSend, c single.Contact) *ud.SearchResponse {
	response := &ud.SearchResponse{}

	var factHashes [][]byte
	facts := msg.GetFact()
	for _, f := range facts {
		if fact.FactType(f.Type) == fact.Nickname {
			jww.WARN.Printf("Cannot search by nickname; fact hash %+v rejected.",
				f.Hash)
			continue
		}
		factHashes = append(factHashes, f.Hash)
	}

	users, err := m.db.Search(factHashes)
	if err != nil {
		response.Error = errors.WithMessage(err, "failed to execute search").Error()
		jww.WARN.Printf("Failed to handle search response: %+v", response.Error)
		return response
	}

	for _, u := range users {
		uid, _ := id.Unmarshal(u.Id)
		jww.DEBUG.Printf("User found in search by %s: %s", c.GetPartner(), uid)
		var contact = &ud.Contact{
			UserID: u.Id,
			PubKey: u.DhPub,
		}

		var uFacts []*ud.HashFact
		for _, f := range u.Facts {
			if f.Type == uint8(fact.Username) {
				contact.Username = f.Fact
			}
			uFacts = append(uFacts, &ud.HashFact{
				Hash: f.Hash,
				Type: int32(f.Type),
			})
		}
		contact.TrigFacts = uFacts

		response.Contacts = append(response.Contacts, contact)
	}

	if len(response.Contacts) == 0 {
		response.Error = "NO RESULTS FOUND"
	}

	return response
}
