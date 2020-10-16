package udb

import (
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/crypto/hash"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"gitlab.com/xx_network/primitives/id"
)

// RegisterFact is an endpoint that attempts to register a user's fact.
func RegisterFact(request *pb.FactRegisterRequest, store storage.Storage,
	auth *connect.Auth) (*pb.FactRegisterResponse, error) {

	// Ensure client is properly authenticated
	if !auth.IsAuthenticated || auth.Sender.IsDynamicHost() {
		return &pb.FactRegisterResponse{}, connect.AuthError(auth.Sender.GetId())
	}

	// Return an error if the request is invalid
	if request == nil || request.Fact == nil {
		return &pb.FactRegisterResponse{}, errors.New("Unable to parse required fields in fact registration request.")
	}

	// Return an error if the fact is already registered
	hashedFact := request.Fact.Digest()
	if len(storage.UserDiscoveryDB.Search([][]byte{hashedFact})) > 0 {
		return &pb.FactRegisterResponse{}, errors.New("Unable to parse required fields in fact registration request.")
	}

	// Return an error if the fact already exists
	user, err := store.GetUser(request.UID)
	if err != nil {
		return &pb.FactRegisterResponse{}, errors.New("User not registered.")
	}

	// Parse the client's public key
	clientPubKey, err := rsa.LoadPublicKeyFromPem([]byte(user.RsaPub))
	if err != nil {
		return &pb.FactRegisterResponse{}, errors.New("Could not parse user's key.")
	}

	// Return an error if the fact signature cannot be verified
	err = rsa.Verify(clientPubKey, hash.CMixHash, hashedFact, request.FactSig, nil)
	if err != nil {
		return &pb.FactRegisterResponse{}, errors.New("Failed to verify fact signature.")
	}

	_ = store.InsertFactTwilio(request.UID, hashedFact, request.FactSig, request.Fact.Fact, uint(request.Fact.FactType), "")

	RegisterFact(uid *id.ID, fact string, factType uint8, signature []byte, verifier VerificationService) (string, error)

}