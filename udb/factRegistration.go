package udb

import (
	"fmt"
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/crypto/hash"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/elixxir/user-discovery-bot/twilio"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"gitlab.com/xx_network/primitives/id"
)

var (
	invalidFactRegisterRequestError = errors.New("Unable to parse required fields in fact registration request.")
	factExistsError                 = errors.New("Cannot register fact that already exists.")
	noUserError                     = errors.New("User associated with fact not registered.")
	invalidUserKeyError             = errors.New("Could not parse user's key.")
	invalidFactSigError             = errors.New("Failed to verify fact signature.")
	invalidUserIdError              = errors.New("Failed to parse user ID.")
	twilioRegFailureError           = errors.New("Failed to register fact with Twilio.")
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
		return &pb.FactRegisterResponse{}, invalidFactRegisterRequestError
	}

	// Return an error if the fact is already registered
	hashedFact := request.Fact.Digest()
	if len(storage.UserDiscoveryDB.Search([][]byte{hashedFact})) != 0 {
		return &pb.FactRegisterResponse{}, factExistsError
	}
	test := uint(5)
	fmt.Println(test)

	// Return an error if the fact's user is not registered
	user, err := store.GetUser(request.UID)
	if err != nil {
		return &pb.FactRegisterResponse{}, noUserError
	}

	// Parse the client's public key
	clientPubKey, err := rsa.LoadPublicKeyFromPem([]byte(user.RsaPub))
	if err != nil {
		return &pb.FactRegisterResponse{}, invalidUserKeyError
	}

	// Return an error if the fact signature cannot be verified
	err = rsa.Verify(clientPubKey, hash.CMixHash, hashedFact, request.FactSig, nil)
	if err != nil {
		return &pb.FactRegisterResponse{}, invalidFactSigError
	}

	// Marshal user ID
	userID, err := id.Unmarshal(request.UID)
	if err != nil {
		return &pb.FactRegisterResponse{}, invalidUserIdError
	}

	// Register fact with Twilio to get confirmation ID
	confirmationID, err := twilio.RegisterFact(userID, request.Fact.Fact,
		uint8(request.Fact.FactType), request.FactSig, nil)
	if err != nil {
		return &pb.FactRegisterResponse{}, twilioRegFailureError
	}

	// Create response
	response := &pb.FactRegisterResponse{
		ConfirmationID: confirmationID,
	}

	return response, nil
}
