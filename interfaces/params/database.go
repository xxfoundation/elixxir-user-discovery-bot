////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
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
