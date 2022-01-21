////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Handles Params-related functionality for the UserDiscovery layer

package params

type General struct {
	SessionPath             string
	ProtoUserJson           []byte
	Ndf                     string
	PermCert                []byte
	RestrictedUserListPath  string // Path to list of line-seperated usernames
	RestrictedRegexListPath string // Path to list of line-seperated regexes

	Database
	IO
	Twilio
}
