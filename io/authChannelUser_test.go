package io

import (
	"crypto/ed25519"
	"gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/crypto/channel"
	"gitlab.com/elixxir/user-discovery-bot/interfaces/params"
	"gitlab.com/elixxir/user-discovery-bot/storage"
	"gitlab.com/xx_network/crypto/csprng"
	"gitlab.com/xx_network/crypto/signature/rsa"
	"testing"
	"time"
)

func TestAuthChannelUser(t *testing.T) {
	// Initialize client and storage
	clientId, clientKey := initClientFields(t)
	store := storage.NewTestDB(t)
	leaseTime := 100
	leaseGracePeriod := 10

	err := store.InsertUser(&storage.User{
		Id:     clientId.Bytes(),
		RsaPub: string(rsa.CreatePublicKeyPem(clientKey.GetPublic())),
	})
	if err != nil {
		t.Fatalf("Failed to insert user: %+v", err)
	}

	rng := csprng.NewSystemRNG()
	rng.SetSeed([]byte("seed"))

	ts := int64(1000)

	udPub, udPriv, err := ed25519.GenerateKey(rng)
	if err != nil {
		t.Fatalf("Failed to generate ud ed25519 key: %+v", err)
	}

	userPub, _, err := ed25519.GenerateKey(rng)
	if err != nil {
		t.Fatalf("Failed to generate ud ed25519 key: %+v", err)

	}

	sig, err := channel.SignChannelIdentityRequest(userPub, ts, clientKey, rng)
	if err != nil {
		t.Fatalf("Failed to sign user pub key: %+v", err)
	}

	resp, err := authorizeChannelUser(&mixmessages.ChannelLeaseRequest{
		UserID:                 clientId.Bytes(),
		UserEd25519PubKey:      userPub,
		Timestamp:              ts,
		UserPubKeyRSASignature: sig,
	}, store, params.Channels{
		Enabled:          true,
		LeaseTime:        time.Duration(leaseTime),
		LeaseGracePeriod: time.Duration(leaseGracePeriod),
		Ed25519Key:       udPriv,
	})
	if err != nil {
		t.Fatalf("Failed to authorizeChannelUser: %+v", err)
	}

	ok := channel.VerifyChannelLease(resp.UDLeaseEd25519Signature, userPub, uint64(resp.Lease), udPub)
	if !ok {
		t.Fatal("Failed to verify ud signature returned by authorizeChannelUser")
	}

	if resp.Lease != ts+int64(leaseTime) {
		t.Errorf("Lease not calculated as expected: %+v", err)
	}

	ts2 := resp.Lease - 15
	sig2, err := channel.SignChannelIdentityRequest(userPub, ts2, clientKey, rng)
	if err != nil {
		t.Fatalf("Failed to sign user pub key: %+v", err)
	}
	resp2, err := authorizeChannelUser(&mixmessages.ChannelLeaseRequest{
		UserID:                 clientId.Bytes(),
		UserEd25519PubKey:      userPub,
		Timestamp:              ts2,
		UserPubKeyRSASignature: sig2,
	}, store, params.Channels{
		Enabled:          true,
		LeaseTime:        time.Duration(leaseTime),
		LeaseGracePeriod: time.Duration(leaseGracePeriod),
		Ed25519Key:       udPriv,
	})
	if err != nil {
		t.Fatalf("Failed to authorizeChannelUser: %+v", err)
	}

	if resp2.Lease != resp.Lease {
		t.Errorf("Lease should not have changed\n\tExpected: %d\n\tReceived: %d\n", resp.Lease, resp2.Lease)
	}

	ts3 := resp.Lease - 3
	sig3, err := channel.SignChannelIdentityRequest(userPub, ts3, clientKey, rng)
	if err != nil {
		t.Fatalf("Failed to sign user pub key: %+v", err)
	}
	resp3, err := authorizeChannelUser(&mixmessages.ChannelLeaseRequest{
		UserID:                 clientId.Bytes(),
		UserEd25519PubKey:      userPub,
		Timestamp:              ts3,
		UserPubKeyRSASignature: sig3,
	}, store, params.Channels{
		Enabled:          true,
		LeaseTime:        time.Duration(leaseTime),
		LeaseGracePeriod: time.Duration(leaseGracePeriod),
		Ed25519Key:       udPriv,
	})
	if err != nil {
		t.Fatalf("Failed to authorizeChannelUser: %+v", err)
	}
	if resp3.Lease != ts3+int64(leaseTime) {
		t.Errorf("Did not receive expected lease\n\tExpected: %d\n\tReceived: %d\n", ts3+int64(leaseTime), resp3.Lease)
	}
}
