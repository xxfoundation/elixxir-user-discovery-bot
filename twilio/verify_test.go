package twilio

import "gitlab.com/elixxir/user-discovery-bot/interfaces/params"

var v = verifier{
	p: params.Twilio{
		AccountSid:      "AC9b3d1637b738d3ab5dacf85d76ec866b",
		AuthToken:       "ea059036ea7a38abe7f058f69cf9d0d8",
		VerificationSid: "VA03c19b46b20a054d7b5c3eb3a2bec1bf",
	},
}

//func TestTwilioVerifier_Verification(t *testing.T) {
//	sid, err := v.Verification("+17813151633", "sms")
//	if err != nil {
//		t.Error(err)
//	}
//	println(sid)
//}
//
//func TestTwilioVerifier_VerificationCheck(t *testing.T) {
//	code := 222410
//	ok, err := v.VerificationCheck(code, "+17813151633")
//	if err != nil {
//		t.Error(err)
//	}
//	t.Log(ok)
//}
