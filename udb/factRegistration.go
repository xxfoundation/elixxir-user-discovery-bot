package udb

import (
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
	invalidFactRegisterRequestError = "Unable to parse required fields in fact registration request."
	factExistsError                 = "Cannot register fact that already exists."
	noUserError                     = "User associated with fact not registered."
	invalidUserKeyError             = "Could not parse user's key."
	invalidFactSigError             = "Failed to verify fact signature."
	invalidUserIdError              = "Failed to parse user ID."
	twilioRegFailureError           = "Failed to register fact with Twilio."
)

// RegisterFact is an endpoint that attempts to register a user's fact.
func RegisterFact(request *pb.FactRegisterRequest, verifier twilio.VerificationService, store storage.Storage,
	auth *connect.Auth) (*pb.FactRegisterResponse, error) {

	// Ensure client is properly authenticated
	if !auth.IsAuthenticated || auth.Sender.IsDynamicHost() {
		return &pb.FactRegisterResponse{}, connect.AuthError(auth.Sender.GetId())
	}

	// Return an error if the request is invalid
	if request == nil || request.Fact == nil {
		return &pb.FactRegisterResponse{}, errors.New(invalidFactRegisterRequestError)
	}

	// Return an error if the fact is already registered
	hashedFact := request.Fact.Digest()
	if len(store.Search([][]byte{hashedFact})) != 0 {
		return &pb.FactRegisterResponse{}, errors.New(factExistsError)
	}

	// Return an error if the fact's user is not registered
	user, err := store.GetUser(request.UID)
	if err != nil {
		return &pb.FactRegisterResponse{}, errors.Errorf("Failed to parse user ID: %+v", err)
	} else if user == nil {
		return &pb.FactRegisterResponse{}, errors.New(noUserError)
	}

	// Parse the client's public key
	clientPubKey, err := rsa.LoadPublicKeyFromPem([]byte(user.RsaPub))
	if err != nil {
		return &pb.FactRegisterResponse{}, errors.New(invalidUserKeyError)
	}

	// Return an error if the fact signature cannot be verified
	err = rsa.Verify(clientPubKey, hash.CMixHash, hashedFact, request.FactSig, nil)
	if err != nil {
		return &pb.FactRegisterResponse{}, errors.New(invalidFactSigError)
	}

	// Marshal user ID
	userID, err := id.Unmarshal(request.UID)
	if err != nil {
		return &pb.FactRegisterResponse{}, errors.New(invalidUserIdError)
	}

	// Register fact with Twilio to get confirmation ID
	confirmationID, err := twilio.RegisterFact(userID, request.Fact.Fact,
		uint8(request.Fact.FactType), request.FactSig, verifier, store)
	if err != nil {
		return &pb.FactRegisterResponse{}, errors.New(twilioRegFailureError)
	}

	// Create response
	response := &pb.FactRegisterResponse{
		ConfirmationID: confirmationID,
	}

	return response, nil
}
