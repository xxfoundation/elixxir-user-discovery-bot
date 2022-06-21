package cmix

import (
	"bytes"
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/client/cmix"
	"gitlab.com/elixxir/client/cmix/identity/receptionID"
	"gitlab.com/elixxir/client/cmix/rounds"
	"gitlab.com/elixxir/client/single"
	"gitlab.com/elixxir/client/ud"
	"gitlab.com/elixxir/primitives/fact"
	"gitlab.com/xx_network/primitives/id"
	"time"
)

type searchManager struct {
	m *Manager
}

func (sm *searchManager) Callback(req *single.Request, eid receptionID.EphemeralIdentity, rids []rounds.Round) {
	jww.INFO.Printf("Received search request from %s [%+v] on rids %+v", req.GetPartner(), eid, rids)
	resp := sm.handleSearch(req)
	marshaledResponse, err := proto.Marshal(resp)
	if err != nil {
		jww.ERROR.Printf("Failed to marshal request to lookup request from "+
			"%s: %+v", req.GetPartner(), err)
		return
	}
	rid, err := req.Respond(marshaledResponse, cmix.GetDefaultCMIXParams(), time.Minute)
	jww.INFO.Printf("Responded to search request from %s over round %d", req.GetPartner(), rid)
}

func (sm *searchManager) handleSearch(req *single.Request) *ud.SearchResponse {
	response := &ud.SearchResponse{}

	msg := &ud.SearchSend{}
	if err := proto.Unmarshal(req.GetPayload(), msg); err != nil {
		jww.ERROR.Printf("Failed to unmarshal search request from %s: %+v",
			req.GetPartner(), err)
		return response
	}

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

	users, err := sm.m.db.Search(factHashes)
	if err != nil {
		response.Error = errors.WithMessage(err, "failed to execute search").Error()
		jww.WARN.Printf("Failed to handle search response: %+v", response.Error)
		return response
	}

	jww.DEBUG.Printf("Raw search returned %+v", users)

	for _, u := range users {
		if bytes.Compare(u.Id, id.DummyUser[:]) == 0 {
			jww.DEBUG.Printf("Don't return dummy user")
			continue
		}
		var contact = &ud.Contact{
			UserID: u.Id,
			PubKey: u.DhPub,
		}

		var uFacts []*ud.HashFact
		for _, f := range u.Facts {
			if f.Type == uint8(fact.Username) {
				contact.Username = u.Username
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
