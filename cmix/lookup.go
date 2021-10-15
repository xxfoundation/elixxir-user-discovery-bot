package cmix

import (
	"fmt"
	"github.com/golang/protobuf/proto"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/client/single"
	"gitlab.com/elixxir/client/ud"
	"gitlab.com/elixxir/primitives/fact"
	"gitlab.com/xx_network/primitives/id"
	"time"
)

func (m *Manager) lookupCallback(payload []byte, c single.Contact) {
	lookupMsg := &ud.LookupSend{}
	if err := proto.Unmarshal(payload, lookupMsg); err != nil {
		jww.ERROR.Printf("Failed to unmarshal lookup request from %s: %+v",
			c.GetPartner(), err)
		return
	}

	jww.INFO.Printf("Lookup request from %s: %v", c.GetPartner(), payload)

	response := m.handleLookup(lookupMsg, c)

	marshaledResponse, err := proto.Marshal(response)
	if err != nil {
		jww.ERROR.Printf("Failed to marshal request to lookup request from "+
			"%s: %+v", c.GetPartner(), err)
		return
	}

	// TODO: make timeout come from config file, default to 1 minute
	err = m.singleUse.RespondSingleUse(c, marshaledResponse, 1*time.Minute)
	if err != nil {
		jww.ERROR.Printf("Failed to send single-use response to to lookup "+
			"request from %s: %+v", c.GetPartner(), err)
		return
	}
}

func (m *Manager) handleLookup(msg *ud.LookupSend, c single.Contact) *ud.LookupResponse {
	response := &ud.LookupResponse{}

	// Decode the ID to lookup
	lookupID, err := id.Unmarshal(msg.UserID)
	if err != nil {
		response.Error = fmt.Sprintf("failed to unmarshal lookup ID in "+
			"request from %s: %+v", c.GetPartner(), err)
		jww.WARN.Printf("Failed to handle lookup response: %+v", response.Error)
		return response
	}

	// Lookup the ID
	usr, err := m.db.GetUser(lookupID.Marshal())
	if err != nil {
		response.Error = fmt.Sprintf("failed to lookup ID %s in request from "+
			"%s: %+v", lookupID, c.GetPartner(), err)
		jww.WARN.Printf("Failed to handle lookup response: %+v", response.Error)
		return response
	}
	if len(usr.Facts) > 0 && usr.Facts[0].Type == uint8(fact.Username) {
		response.Username = usr.Facts[0].Fact
	}

	response.PubKey = usr.DhPub
	return response
}
