////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Handles Params-related functionality for the Twilio layer

package params

type Twilio struct {
	AccountSid      string
	AuthToken       string
	VerificationSid string
}
