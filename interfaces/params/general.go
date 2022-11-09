////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

// Handles Params-related functionality for the UserDiscovery layer

package params

type General struct {
	SessionPath     string
	ProtoUserJson   []byte
	Ndf             string
	PermCert        []byte
	BannedUserList  string
	BannedRegexList string

	Database
	IO
	Twilio
}
