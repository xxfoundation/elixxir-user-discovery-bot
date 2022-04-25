package cmix

import (
	"fmt"
	"github.com/golang/protobuf/proto"
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

type lookupManager struct {
	m *Manager
}

func (lm *lookupManager) Callback(req *single.Request, eid receptionID.EphemeralIdentity, rids []rounds.Round) {
	jww.INFO.Printf("Received lookup request from %s [%+v] on rids %+v", req.GetPartner(), eid, rids)
	resp := lm.handleLookup(req)
	marshaledResponse, err := proto.Marshal(resp)
	if err != nil {
		jww.ERROR.Printf("Failed to marshal request to lookup request from "+
			"%s: %+v", req.GetPartner(), err)
		return
	}
	rid, err := req.Respond(marshaledResponse, cmix.GetDefaultCMIXParams(), time.Minute)
	jww.INFO.Printf("Responded to lookup request from %s over round %d", req.GetPartner(), rid)
}

func (lm *lookupManager) handleLookup(req *single.Request) *ud.LookupResponse {
	response := &ud.LookupResponse{}

	msg := &ud.LookupSend{}
	if err := proto.Unmarshal(req.GetPayload(), msg); err != nil {
		jww.ERROR.Printf("Failed to unmarshal lookup request from %s: %+v",
			req.GetPartner(), err)
		response.Error = err.Error()
		return response
	}

	// Decode the ID to lookup
	lookupID, err := id.Unmarshal(msg.UserID)
	if err != nil {
		response.Error = fmt.Sprintf("failed to unmarshal lookup ID in "+
			"request from %s: %+v", req.GetPartner(), err)
		jww.WARN.Printf("Failed to handle lookup response: %+v", response.Error)
		response.Error = err.Error()
		return response
	}

	// Lookup the ID
	usr, err := lm.m.db.GetUser(lookupID.Marshal())
	if err != nil {
		response.Error = fmt.Sprintf("failed to lookup ID %s in request from "+
			"%s: %+v", lookupID, req.GetPartner(), err)
		jww.WARN.Printf("Failed to handle lookup response: %+v", response.Error)
		response.Error = err.Error()
		return response
	}
	if len(usr.Facts) > 0 && usr.Facts[0].Type == uint8(fact.Username) {
		response.Username = usr.Username
	}

	response.PubKey = usr.DhPub
	return response
}
