////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Handles Params-related functionality for the Channels endpoint

package params

import "time"

type Channels struct {
	Enabled          bool
	LeaseTime        time.Duration
	LeaseGracePeriod time.Duration
	Ed25519Key       []byte
}
