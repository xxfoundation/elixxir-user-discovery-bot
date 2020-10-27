////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Handles Params-related functionality for the Database layer

package params

type Database struct {
	DbUsername string
	DbPassword string
	DbName     string
	DbAddress  string
	DbPort     string
}
