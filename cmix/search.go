package cmix

import (
	"github.com/golang/protobuf/proto"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/client/interfaces/message"
	"gitlab.com/elixxir/client/interfaces/params"
	"gitlab.com/elixxir/client/ud"
	"gitlab.com/xx_network/primitives/id"
)

func (m *Manager) SearchProcess() {
	for true {
		request := <-m.searchChan

		if request.Encryption != message.E2E {
			jww.ERROR.Printf("Ignoring improperly encrypted search "+
				"request from %s", request.Sender)
		}

		searchMsg := &ud.SearchSend{}
		if err := proto.Unmarshal(request.Payload, searchMsg); err != nil {
			jww.ERROR.Printf("failed to unmarshal search "+
				"request from %s: %+v", request.Sender, err)
			continue
		}

		response := m.handleSearch(searchMsg, request.Sender)

		marshaledResponse, err := proto.Marshal(response)
		if err != nil {
			jww.ERROR.Printf("failed to marshal responce "+
				"to search request from %s: %+v", request.Sender, err)
			continue
		}

		responseMsg := message.Send{
			Recipient:   request.Sender,
			Payload:     marshaledResponse,
			MessageType: message.UdSearchResponse,
		}

		_, _, err = m.client.SendE2E(responseMsg, params.GetDefaultE2E())

		if err != nil {
			jww.ERROR.Printf("failed to send responce "+
				"to search request from %s: %+v", request.Sender, err)
		}
	}
}

func (m *Manager) handleSearch(msg *ud.SearchSend, requestor *id.ID) (response *ud.SearchResponse) {
	response = &ud.SearchResponse{
		Contacts: nil,
		CommID:   msg.CommID,
		Error:    "",
	}

	var factHashs [][]byte
	facts := msg.GetFact()
	for _, f := range facts {
		factHashs = append(factHashs, f.Hash)
	}

	users, err := m.db.Search(factHashs)
	if err != nil {
		response.Error = errors.WithMessage(err, "handleSearch error: failed to execute search").Error()
		return
	}
	for _, u := range users {
		var ufacts []*ud.HashFact
		for _, f := range u.Facts {
			ufacts = append(ufacts, &ud.HashFact{
				Hash: f.Hash,
				Type: int32(f.Type),
			})
		}
		response.Contacts = append(response.Contacts, &ud.Contact{
			UserID:    u.Id,
			PubKey:    u.DhPub,
			TrigFacts: ufacts,
		})
	}

	return
}