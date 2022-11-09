////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package twilio

import (
	"gitlab.com/elixxir/user-discovery-bot/interfaces/params"
	"testing"
)

func TestNewManager(t *testing.T) {
	m := NewManager(params.Twilio{}, nil)
	if m == nil {
		t.Error("This should not happen")
	}
}
