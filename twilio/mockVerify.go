package twilio

import (
	"math/rand"
	"strconv"
)

var MV *MockVerifier = &MockVerifier{codes: map[string]int{}, index: 0}

type MockVerifier struct {
	codes map[string]int
	index int
}

func (v *MockVerifier) Verification(to, channel string) (string, error) {
	cid := strconv.Itoa(v.index)
	v.index++
	v.codes[cid] = rand.Int()
	return cid, nil
}

func (v *MockVerifier) VerificationCheck(code int, to string) (bool, error) {
	c, ok := v.codes[to]
	if !ok || c != code {
		return false, nil
	}
	return ok, nil
}