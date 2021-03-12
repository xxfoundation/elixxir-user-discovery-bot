package twilio

import (
	"gitlab.com/elixxir/user-discovery-bot/interfaces/params"
	"testing"
)

var v = verifier{
	p: params.Twilio{
		AccountSid:      "AC93abc3f04f849213a4cbdd59046e48e5",
		AuthToken:       "062a0d90e0eac18b702a4e89837d8f73",
		VerificationSid: "VA03c19b46b20a054d7b5c3eb3a2bec1bf",
	},
}

// This test cannot pass in a non-error state as twilio
// doesn't have testing credentials for verification
func TestTwilioVerifier_Verification(t *testing.T) {
	sid, err := v.Verification("+17813151633", "sms")
	if err == nil {
		t.Errorf("Test verification not enabled: %+v", err)
	}
	println(sid)
}

func TestTwilioVerifier_VerificationCheck(t *testing.T) {
	code := 222410
	ok, err := v.VerificationCheck(code, "+17813151633")
	if err == nil {
		t.Errorf("Test verification not enabled: %+v", err)
	}
	t.Log(ok)
}
