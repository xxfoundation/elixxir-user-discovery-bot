package io

import (
	"bytes"
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/crypto/channel"
	"gitlab.com/elixxir/user-discovery-bot/interfaces/params"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"gorm.io/gorm"
)

const (
	errorChannelsNotActive = "error channels are not enabled in user discovery"
	errorUserBanned        = "error user is banned from channels"
	errorPubkeyMismatch    = "error cannot register second public key for user"
)

// Handle a request for channel authentication from a user.
// Accepts pb.ChannelAuthenticationRequest, UD ed25519 private key, & UD storage interface
// Returns a ChannelAuthenticationResponser
func authorizeChannelUser(req *pb.ChannelAuthenticationRequest, s *storage.Storage, param params.Channels) (*pb.ChannelAuthenticationResponse, error) {
	// Return error if not configured to run this endpoint
	if !param.Enabled {
		return nil, errors.New(errorChannelsNotActive)
	}

	// Fetch user RSA public key from database
	u, err := s.GetUser(req.UserID)
	if err != nil {
		return nil, err
	}
	pubKey, err := rsa.LoadPublicKeyFromPem([]byte(u.RsaPub))
	if err != nil {
		return nil, err
	}

	// Verify channel request authenticity based on public key from database
	err = channel.VerifyChannelIdentityRequest(req.UserSignedEdPubKey, req.UserEd25519PubKey, req.Timestamp, pubKey)
	if err != nil {
		return nil, err
	}

	// Check for previous registration, return error if banned or if attempting
	// to use public key other than the one stored
	prevChanId, err := s.GetChannelIdentity(req.UserID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	} else if prevChanId != nil && !bytes.Equal(prevChanId.Ed25519Pub, req.UserEd25519PubKey) {
		return nil, errors.New(errorPubkeyMismatch)
	} else if prevChanId != nil && prevChanId.Banned {
		return nil, errors.New(errorUserBanned)
	}

	// If no lease, or if lease expired, issue new lease
	// If lease unexpired, but within grace period, issue new lease
	// Otherwise, use stored lease
	var lease int64
	if prevChanId != nil && prevChanId.Lease-param.LeaseGracePeriod.Nanoseconds() > req.Timestamp {
		lease = prevChanId.Lease
	} else {
		lease = req.Timestamp + param.LeaseTime.Nanoseconds()
	}

	// Sign lease + user's public key
	udSig := channel.SignChannelLease(req.UserEd25519PubKey, uint64(lease), param.Ed25519Key)

	// Insert identity to database
	err = s.InsertChannelIdentity(&storage.ChannelIdentity{
		UserId:     req.UserID,
		Ed25519Pub: req.UserEd25519PubKey,
		Lease:      lease,
	})
	if err != nil {
		return nil, err
	}

	// Return lease and signature
	return &pb.ChannelAuthenticationResponse{
		Lease:            lease,
		UDSignedEdPubKey: udSig,
	}, nil
}
