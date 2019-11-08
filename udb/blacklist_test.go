////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package udb

import (
	"testing"
)

func TestBlackList_Exists(t *testing.T) {
	testBlackList := InitBlackList("./whitelists/bannedNames.txt")

	if testBlackList.Exists("DavidChaum") != true {
		t.Error("Failed to detect a banned name")
	}

	if testBlackList.Exists("UnicornKitty") != false {
		t.Error("False positive: Detected a non-blacklisted username")
	}

	if testBlackList.Exists("JakeTayl0r") {
		t.Error("False positive: Detected a non-blacklisted username")
	}
}
