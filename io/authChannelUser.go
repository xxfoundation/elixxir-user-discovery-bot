package io

import (
	"bytes"
	"crypto/ed25519"
	"github.com/pkg/errors"
	pb "gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/crypto/channel"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"gorm.io/gorm"
	"time"
)

// TODO make these configurable
const leaseTime = time.Hour * 500
const leaseGraceTime = time.Hour * 24
const channelEndpointActive = true

const (
	errorChannelsNotActive = ""
	errorUserBanned        = ""
	errorPubkeyMismatch    = ""
)

// Handle a request for channel authentication from a user.
// Accepts pb.ChannelAuthenticationRequest, UD ed25519 private key, & UD storage interface
// Returns a ChannelAuthenticationResponser
func authorizeChannelUser(req *pb.ChannelAuthenticationRequest, udEd25519PrivKey ed25519.PrivateKey, s *storage.Storage) (*pb.ChannelAuthenticationResponse, error) {
	if !channelEndpointActive {
		return nil, errors.New(errorChannelsNotActive)
	}

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

	prevChanId, err := s.GetChannelIdentity(req.UserID)
	if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
		return nil, err
	} else if !bytes.Equal(prevChanId.Ed25519Pub, req.UserEd25519PubKey) {
		return nil, errors.New(errorPubkeyMismatch)
	} else if prevChanId.Banned {
		return nil, errors.New(errorUserBanned)
	}

	var lease int64
	if prevChanId != nil {
		if prevChanId.Lease < req.Timestamp {
			if prevChanId.Lease < req.Timestamp+leaseGraceTime.Nanoseconds() {
				lease = prevChanId.Lease
			} else {
				lease = req.Timestamp + leaseTime.Nanoseconds()
			}
		} else {
			lease = req.Timestamp + leaseTime.Nanoseconds()
		}
	} else {
		lease = req.Timestamp + leaseTime.Nanoseconds()
	}

	udSig := channel.SignResponse(req.UserEd25519PubKey, uint64(lease), udEd25519PrivKey)
	if err != nil {
		return nil, err
	}

	err = s.InsertChannelIdentity(&storage.ChannelIdentity{
		UserId:     req.UserID,
		Ed25519Pub: req.UserEd25519PubKey,
		Lease:      lease,
	})
	if err != nil {
		return nil, err
	}

	return &pb.ChannelAuthenticationResponse{
		Lease:            lease,
		UDSignedEdPubKey: udSig,
	}, nil
}
