////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Handles Params-related functionality for the IO layer

package params

type IO struct {
	Cert []byte
	Key  []byte
	Port string
}
