////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Mock implementation of a verification service for testing fact verification

package twilio

import (
	"math/rand"
	"strconv"
)

// Mocked verifier implementation
var MV = &MockVerifier{Codes: map[string]int{}, index: 0}

type MockVerifier struct {
	Codes map[string]int
	index int
}

func (v *MockVerifier) Verification(to, channel string) (string, error) {
	cid := strconv.Itoa(v.index)
	v.index++
	v.Codes[cid] = rand.Int()
	return cid, nil
}

func (v *MockVerifier) VerificationCheck(code int, to string) (bool, error) {
	c, ok := v.Codes[to]
	if !ok || c != code {
		return false, nil
	}
	return ok, nil
}
