////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Twilio verification service code, using POST requests

package twilio

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/user-discovery-bot/interfaces/params"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

const VERIFICATION_URL = "https://verify.twilio.com/v2/Services/%s/Verifications"
const VERIFICATION_CHECK_URL = "https://verify.twilio.com/v2/Services/%s/VerificationCheck"
const PAYLOAD_TO = "To"
const PAYLOAD_SID = "VerificationSid"
const PAYLOAD_CODE = "Code"
const PAYLOAD_CHAN = "Channel"

// Interface for verification service
type VerificationService interface {
	Verification(to, channel string) (string, error)
	VerificationCheck(code string, to string) (bool, error)
}

// Channels that can be passed into twilio
type Channel int

const (
	SMS Channel = iota
	Email
)

func (c Channel) String() string {
	return [...]string{"sms", "email"}[c]
}

type verifier struct {
	p params.Twilio
}

// Posts to the verification endpoint of twilio, returns confirmation id
func (v *verifier) Verification(to, channel string) (string, error) {
	jww.INFO.Printf("Attempting to verify %s via %s", to, channel)
	verificationURL := fmt.Sprintf(VERIFICATION_URL, v.p.VerificationSid)
	payload := url.Values{}
	payload.Set(PAYLOAD_TO, to)
	payload.Set(PAYLOAD_CHAN, channel)

	data, err := v.twilioRequest(payload, verificationURL)
	if err != nil {
		if strings.Contains(err.Error(), "Account is not active") {
			jww.FATAL.Printf("Failed to submit verification request: %+v", err)
		}
		return "", err
	}
	jww.INFO.Printf("Response data: %+v", data)
	sid := fmt.Sprint(data["sid"])

	return sid, err
}

// Posts to the verificationcheck endpoint of twilio, returns verification status (bool)
func (v *verifier) VerificationCheck(code string, to string) (bool, error) {
	checkUrl := fmt.Sprintf(VERIFICATION_CHECK_URL, v.p.VerificationSid)
	payload := url.Values{}
	payload.Set(PAYLOAD_SID, to)
	payload.Set(PAYLOAD_CODE, code)

	data, err := v.twilioRequest(payload, checkUrl)
	if err != nil {
		// https://www.twilio.com/docs/api/errors/20404
		// If twilio cannot find the verification check, it probably expired
		if strings.Contains(err.Error(), "errors/20404") {
			return false, errors.New("Your verification code may have expired; please resubmit")
		}
		jww.ERROR.Printf("Failed to submit verification check request: %+v", err)
		return false, errors.WithMessage(err, "Failed to make verification check request")
	}
	jww.INFO.Println(data)
	valid, err := strconv.ParseBool(fmt.Sprint(data["valid"]))
	if err != nil {
		return false, errors.WithMessage(err, "Failed to pares verification check response")
	}
	return valid, nil
}

// Helper function for sending post requests to twilio
func (v *verifier) twilioRequest(payload url.Values, url string) (map[string]interface{}, error) {
	client := &http.Client{} // TODO: this may need special configurations.  See Transport object

	req, err := http.NewRequest("POST", url, strings.NewReader(payload.Encode()))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	req.SetBasicAuth(v.p.AccountSid, v.p.AuthToken)
	req.Header.Add("Accept", "application/json")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := client.Do(req)
	if err != nil {
		return nil, errors.WithStack(err)
	}
	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&data)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return data, nil
	} else {
		var data map[string]interface{}
		decoder := json.NewDecoder(resp.Body)
		err = decoder.Decode(&data)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		return nil, errors.Errorf("error: request failed with status %d (%+v): %+v", resp.StatusCode, resp.Status, data)
	}
}
