package io

import (
	"crypto/ed25519"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/crypto/channel"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"time"
)

// Handle a request for channel authentication from a user.
// Accepts pb.ChannelAuthenticationRequest, UD ed25519 private key, & UD storage interface
// Returns a ChannelAuthenticationResponser
func authorizeChannelUser(req *pb.ChannelAuthenticationRequest, udEd25519PrivKey ed25519.PrivateKey, s *storage.Storage) (*pb.ChannelAuthenticationResponse, error) {
	u, err := s.GetUser(req.UserID)
	if err != nil {
		return nil, err
	}
	pubKey, err := rsa.LoadPublicKeyFromPem([]byte(u.RsaPub))
	if err != nil {
		return nil, err
	}

	err = channel.VerifyRequest(req.UserSignedEdPubKey, req.UserEd25519PubKey, req.Timestamp, pubKey)
	if err != nil {
		return nil, err
	}

	lease := req.Timestamp + (time.Hour * 24 * 7 * 3).Nanoseconds()

	udSig := channel.SignResponse(req.UserEd25519PubKey, uint64(lease), udEd25519PrivKey)
	if err != nil {
		return nil, err
	}

	return &pb.ChannelAuthenticationResponse{
		Lease:            lease,
		UDSignedEdPubKey: udSig,
	}, nil
}
