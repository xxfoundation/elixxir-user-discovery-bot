package io

import (
	"crypto/ed25519"
	"gitlab.com/elixxir/comms/mixmessages"
	"gitlab.com/elixxir/crypto/channel"
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

	err := store.InsertUser(&storage.User{
		Id:     clientId.Bytes(),
		RsaPub: string(rsa.CreatePublicKeyPem(clientKey.GetPublic())),
	})
	if err != nil {
		t.Fatalf("Failed to insert user: %+v", err)
	}

	rng := csprng.NewSystemRNG()
	rng.SetSeed([]byte("seed"))

	ts := time.Now().UnixNano()

	udPub, udPriv, err := ed25519.GenerateKey(rng)
	if err != nil {
		t.Fatalf("Failed to generate ud ed25519 key: %+v", err)
	}

	userPub, _, err := ed25519.GenerateKey(rng)
	if err != nil {
		t.Fatalf("Failed to generate ud ed25519 key: %+v", err)

	}

	sig, err := channel.SignRequest(userPub, ts, clientKey, rng)
	if err != nil {
		t.Fatalf("Failed to sign user pub key: %+v", err)
	}

	resp, err := authorizeChannelUser(&mixmessages.ChannelAuthenticationRequest{
		UserID:             clientId.Bytes(),
		UserEd25519PubKey:  userPub,
		Timestamp:          ts,
		UserSignedEdPubKey: sig,
	}, udPriv, store)
	if err != nil {
		t.Fatalf("Failed to authorizeChannelUser: %+v", err)
	}

	ok := channel.VerifyResponse(resp.UDSignedEdPubKey, userPub, uint64(resp.Lease), udPub)
	if !ok {
		t.Fatal("Failed to verify ud signature returned by authorizeChannelUser")
	}
}
