package io

import (
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/crypto/hash"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/elixxir/user-discovery-bot/twilio"
	"gitlab.com/xx_network/comms/connect"
	"gitlab.com/xx_network/comms/messages"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"gitlab.com/xx_network/primitives/id"
	"strconv"
)

var (
	invalidFactRegisterRequestError = "Unable to parse required fields in FactRegisterRequest."
	factExistsError                 = "Cannot register fact that already exists."
	noUserError                     = "User associated with fact not registered."
	invalidUserKeyError             = "Could not parse user's key."
	invalidFactSigError             = "Failed to verify fact signature."
	getUserFailureError             = "Failed to find user"
	invalidUserIdError              = "Failed to parse user ID."
	twilioRegFailureError           = "Failed to register fact with Twilio."

	invalidFactConfirmRequestError = "Unable to parse required fields in FactConfirmRequest."
	invalidFactCodeError           = "Failed to parse the FactConfirmRequest code."
	twilioConfirmFailureError      = "Failed to confirm fact with Twilio"
	twilioVerificationFailureError = "Twilio verification failed."
)

// registerFact is an endpoint that attempts to register a user's fact.
func registerFact(request *pb.FactRegisterRequest, verifier *twilio.Manager, store *storage.Storage,
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
	ret, err := store.Search([][]byte{hashedFact})
	if err != nil {
		return &pb.FactRegisterResponse{}, err
	}
	if len(ret) != 0 {
		return &pb.FactRegisterResponse{}, errors.New(factExistsError)
	}

	// Return an error if the fact's user is not registered
	user, err := store.GetUser(request.UID)
	if err != nil {
		return &pb.FactRegisterResponse{}, errors.Errorf(getUserFailureError+": %+v", err)
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
	confirmationID, err := verifier.RegisterFact(userID, request.Fact.Fact,
		uint8(request.Fact.FactType), request.FactSig)
	if err != nil {
		return &pb.FactRegisterResponse{}, errors.New(twilioRegFailureError)
	}

	// Create response
	response := &pb.FactRegisterResponse{
		ConfirmationID: confirmationID,
	}

	return response, nil
}

// confirmFact verifies the fact via Twilio and sets the fact in the database as
// confirmed.
func confirmFact(request *pb.FactConfirmRequest, verifier *twilio.Manager, store *storage.Storage,
	auth *connect.Auth) (*messages.Ack, error) {

	// Ensure client is properly authenticated
	if !auth.IsAuthenticated || auth.Sender.IsDynamicHost() {
		return &messages.Ack{}, connect.AuthError(auth.Sender.GetId())
	}

	// Return an error if the request is nil
	if request == nil {
		return &messages.Ack{}, errors.New(invalidFactConfirmRequestError)
	}

	// Convert fact code to integer
	code, err := strconv.Atoi(request.Code)
	if err != nil {
		return &messages.Ack{}, errors.New(invalidFactCodeError)
	}

	valid, err := verifier.ConfirmFact(request.ConfirmationID, code)
	if err != nil {
		return &messages.Ack{}, errors.Errorf(twilioConfirmFailureError+": %+v", err)
	} else if !valid {
		return &messages.Ack{}, errors.New(twilioVerificationFailureError)
	}

	return &messages.Ack{}, nil
}
