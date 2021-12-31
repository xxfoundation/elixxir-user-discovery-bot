////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Handles Params-related functionality for the UserDiscovery layer

package params

type General struct {
	SessionPath       string
	ProtoUserJsonPath string
	ProtoUserJson     []byte
	Ndf               string
	PermCert          []byte

	Database
	IO
	Twilio
}
