////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Mock implementation of a verification service for testing fact verification

package twilio

import (
	"math/rand"
	"strconv"
)

type mockVerifier struct {
	Codes map[string]int
	index int
}

func (v *mockVerifier) Verification(to, channel string) (string, error) {
	cid := strconv.Itoa(v.index)
	v.index++
	v.Codes[cid] = rand.Int()
	return cid, nil
}

func (v *mockVerifier) VerificationCheck(code string, to string) (bool, error) {
	_, ok := v.Codes[to]
	if !ok {
		return false, nil
	}
	return ok, nil
}

func newMockVerifier() *mockVerifier {
	return &mockVerifier{
		Codes: make(map[string]int),
		index: 0,
	}
}
