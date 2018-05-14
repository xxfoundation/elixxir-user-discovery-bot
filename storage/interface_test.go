////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package storage

import (
	"testing"
)

func TestValueTypeString(t *testing.T) {
	vtStr := Email.String()
	if vtStr != "Email" {
		t.Errorf("Could not convert ValueType to string!")
	}
	unknownStr := ValueType(4).String()
	if unknownStr != "Unknown" {
		t.Errorf("ValueType string returned: %s, Expected: Unknown", unknownStr)
	}
}
