////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////
package udb

import (
	"github.com/pkg/errors"
	"reflect"
	"testing"
)

// Tests if the file parser parses blacklist correctly.
func TestBlackListFileParse(t *testing.T) {
	fileList, err := BlacklistFileParse("./blacklists/bannedNames.txt")

	if err != nil {
		t.Errorf("BlacklistFileParse() produced an unexpected error\n\treceived: %#v\n\texpected: %#v", err, nil)
	}

	if !reflect.DeepEqual(fileList, bannedNames) {
		t.Errorf("BlacklistFileParse() did not read the correct values from file\n\treceived: %#v\n\texpected: %#v", fileList, bannedNames)
	}
}

// Tests if the file parser throws and error when no file is present.
func TestBlackistFileParse_NoFile(t *testing.T) {
	fileList, err := BlacklistFileParse("./blacklists/no_name.txt")

	if err == nil {
		t.Errorf("BlacklistFileParse() did not produce an error when it should have\n\treceived: %#v\n\texpected: %#v", err, errors.New("some error"))
	}

	if !reflect.DeepEqual(fileList, []string{}) {
		t.Errorf("BlacklistFileParse() did not read the correct values from file\n\treceived: %#v\n\texpected: %#v", fileList, []string{})
	}
}

// Tests if the file parser parses an empty file correctly without error.
func TestBlacklistFileParse_Empty(t *testing.T) {
	fileList, err := BlacklistFileParse("./blacklists/empty.txt")

	if err != nil {
		t.Errorf("BlacklistFileParse() produced an unexpected error\n\treceived: %#v\n\texpected: %#v", err, nil)
	}

	if !reflect.DeepEqual(fileList, []string{}) {
		t.Errorf("BlacklistFileParse() did not read the correct values from empty file\n\treceived: %#v\n\texpected: %#v", fileList, []string{})
	}
}

//Test that it parses and stores the blacklisted names correctly
func TestBlackList_Exists(t *testing.T) {
	testBlackList := InitBlackList("./blacklists/bannedNames.txt")

	if testBlackList.Exists("DavidChaum") != true {
		t.Error("Failed to detect a banned name")
	}

	if testBlackList.Exists("UnicornKitty") != false {
		t.Error("False positive: Detected a non-blacklisted username")
	}

	if testBlackList.Exists("JakeTayl0r") != false {
		t.Error("False positive: Detected a non-blacklisted username")
	}
}
