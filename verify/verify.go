package verify

import (
	"encoding/json"
	"fmt"
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

type VerificationService interface {
	Verification(to, channel string) (string, error)
	VerificationCheck(code int, to string) (bool, error)
}

type TwilioVerifier struct {
	accountSid      string
	authToken       string
	verificationSid string
}

type Channel int

const (
	SMS Channel = iota
	Email
	Voice
)

func (c Channel) String() string {
	return [...]string{"sms", "email", "call"}[c]
}

func (v *TwilioVerifier) Verification(to, channel string) (string, error) {
	verificationURL := fmt.Sprintf("https://verify.twilio.com/v2/Services/%s/Verifications", v.verificationSid)
	payload := url.Values{}
	payload.Set("To", to)
	payload.Set("Channel", channel)

	data, err := v.twilioRequest(payload, verificationURL)
	if err != nil {
		return "", err
	}
	sid := fmt.Sprint(data["sid"])

	return sid, err
}

func (v *TwilioVerifier) VerificationCheck(code int, to string) (bool, error) {
	checkUrl := fmt.Sprintf("https://verify.twilio.com/v2/Services/%s/VerificationCheck", v.verificationSid)
	payload := url.Values{}
	payload.Set("To", to)
	payload.Set("Code", strconv.Itoa(code))

	data, err := v.twilioRequest(payload, checkUrl)
	if err != nil {
		return false, errors.WithMessage(err, "Failed to make verification check request")
	}
	jww.INFO.Println(data)
	valid, err := strconv.ParseBool(fmt.Sprint(data["valid"]))
	if err != nil {
		return false, errors.WithMessage(err, "Failed to pares verification check response")
	}
	return valid, nil
}

func (v *TwilioVerifier) twilioRequest(payload url.Values, url string) (map[string]interface{}, error) {
	client := &http.Client{} // TODO: this may need special configurations.  See Transport object

	req, err := http.NewRequest("POST", url, strings.NewReader(payload.Encode()))
	if err != nil {
		return nil, errors.WithStack(err)
	}
	req.SetBasicAuth(v.accountSid, v.authToken)
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
		return nil, errors.Errorf("error: request failed with status %d: %+v", resp.StatusCode, resp.Status)
	}
}
