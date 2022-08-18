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
	"time"
)

var (
	errorChannelsNotActive = errors.New("error channels are not enabled in user discovery")
	errorUserBanned        = errors.New("error user is banned from channels")
	errorPubkeyMismatch    = errors.New("error cannot register second public key for user")
)

// Handle a request for channel authentication from a user.
// Accepts pb.ChannelLeaseRequest, UD ed25519 private key,
// & UD storage interface.  Returns a ChannelLeaseResponse.
func authorizeChannelUser(req *pb.ChannelLeaseRequest, s *storage.Storage,
	param params.Channels) (*pb.ChannelLeaseResponse, error) {
	// Return error if not configured to run this endpoint
	if !param.Enabled {
		return nil, errorChannelsNotActive
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

	requestTimestamp := time.Unix(0, req.Timestamp)
	now := time.Now()

	// Verify channel request authenticity based on public key from database
	err = channel.VerifyChannelIdentityRequest(req.UserPubKeyRSASignature,
		req.UserEd25519PubKey, now, requestTimestamp, pubKey)
	if err != nil {
		return nil, err
	}

	// Check for previous registration, return error if banned or if attempting
	// to use public key other than the one stored
	prevChanId, err := s.GetChannelIdentity(req.UserID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	} else if prevChanId != nil && !bytes.Equal(prevChanId.PublicKey,
		req.UserEd25519PubKey) {
		return nil, errorPubkeyMismatch
	} else if prevChanId != nil && prevChanId.Banned {
		return nil, errorUserBanned
	}

	// If no lease, or if lease expired, issue new lease
	// If lease unexpired, but within grace period, issue new lease
	// Otherwise, use stored lease
	var lease time.Time

	if prevChanId != nil && requestTimestamp.Before(time.Unix(0,
		prevChanId.Lease).Add(-1*param.LeaseGracePeriod)) {
		lease = time.Unix(0, prevChanId.Lease)
	} else {
		lease = now.Add(param.LeaseTime)
	}

	// Sign lease + user's public key
	udSig := channel.SignChannelLease(req.UserEd25519PubKey, u.Username,
		lease, param.Ed25519Key)

	// Insert identity to database
	err = s.InsertChannelIdentity(&storage.ChannelIdentity{
		UserId:    req.UserID,
		PublicKey: req.UserEd25519PubKey,
		Lease:     lease.UnixNano(),
	})
	if err != nil {
		return nil, err
	}

	// Return lease and signature
	return &pb.ChannelLeaseResponse{
		Lease:                   lease.UnixNano(),
		UserEd25519PubKey:       req.UserEd25519PubKey,
		UDLeaseEd25519Signature: udSig,
	}, nil
}
