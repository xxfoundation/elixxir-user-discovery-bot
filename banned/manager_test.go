///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package banned

import (
	"reflect"
	"regexp"
	"testing"
)

func TestNewManager(t *testing.T) {
	expectedManager := &Manager{
		bannedUserList: map[string]struct{}{
			"Privategrity":      {},
			"Privategrity_Corp": {},
		},
		bannedRegexList: []*regexp.Regexp{
			regexp.MustCompile("xx"),
			regexp.MustCompile("xx.*?network"),
		},
	}

	bannedUserList := ""
	for key := range expectedManager.bannedUserList {
		bannedUserList += key + "\n"
	}

	bannedRegexList := ""
	for _, regex := range expectedManager.bannedRegexList {
		bannedRegexList += regex.String() + "\n"
	}

	m, err := NewManager(bannedUserList, bannedRegexList)
	if err != nil {
		t.Fatalf("NewManager error: %v", err)
	}

	if !reflect.DeepEqual(m, expectedManager) {
		t.Errorf("Constructed manager does not match expected output."+
			"\nExpected: %+v"+
			"\nReceived: %+v", expectedManager, m)
	}
}

func TestManager_IsBanned_GoodUsername(t *testing.T) {
	bannedUserList := "Privategrity\nPrivategrity_Corp"
	bannedRegexList := "xx\nxx.*?network"

	m, err := NewManager(bannedUserList, bannedRegexList)
	if err != nil {
		t.Fatalf("NewManager error: %v", err)
	}

	goodUsernames := []string{
		"john_doe",
		"private",
		"network",
		"Privategrity!??!Corporation",
	}

	for _, goodUsername := range goodUsernames {
		if m.IsBanned(goodUsername) {
			t.Errorf("Username %q was recognized as banned when it should not be", goodUsername)
		}

	}

}

//
func TestManager_IsBanned_BadUsername(t *testing.T) {
	bannedUserList := "Privategrity\nPrivategrity_Corp"
	bannedRegexList := "xx\nxx.*?network"

	m, err := NewManager(bannedUserList, bannedRegexList)
	if err != nil {
		t.Fatalf("NewManager error: %v", err)
	}

	badUsernames := []string{
		"xxfsdfsdfsdklfjnetwork",
		"Privategrity",
		"Privategrity_Corp",
		"exxplostion",
	}

	for _, badUsername := range badUsernames {
		if !m.IsBanned(badUsername) {
			t.Errorf("Username %q was not recognized as banned when it should be", badUsername)
		}

	}

}
