////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package fingerprint

import (
	"testing"
)

func TestFingerprint(t *testing.T) {
	expected := "lsXBLOhwWxCWIgDL31s5Hxo/qc+8mKZP4kT9tD/6iTM="
	testVal := []byte{
		'T', 'h', 'i', 's', ' ', 'i', 's', ' ', 't', 'h', 'e', ' ',
		't', 'e', 's', 't', 'v', 'a', 'l',
	}
	retVal := Fingerprint(testVal)
	if retVal != expected {
		t.Errorf("Fingerprint failed, Expected: %s, Got: %s", expected, retVal)
	}
}
