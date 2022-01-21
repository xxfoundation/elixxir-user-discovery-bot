///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package restricted

import (
	"gitlab.com/elixxir/user-discovery-bot/validation"
	"gitlab.com/xx_network/primitives/utils"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"testing"
	"time"
)

// Tests that NewManager returns a new Manager with the expected values.
func TestNewManager(t *testing.T) {
	usernamesPath := "restrictedUsernames.txt"
	regexPath := "restrictedRegex.txt"
	expectedManager := &Manager{
		usernamePath: usernamesPath,
		regexPath:    regexPath,
		usernames: map[string]struct{}{
			"privategrity":      {},
			"privategrity_corp": {},
		},
		regexes: []*regexp.Regexp{
			regexp.MustCompile("xx"),
			regexp.MustCompile("xx.*?network"),
		},
	}

	usernameList := usernamesToList(expectedManager.usernames)
	regexList := regexesToList(expectedManager.regexes)

	err := utils.WriteFile(
		usernamesPath, []byte(usernameList), utils.FilePerms, utils.DirPerms)
	if err != nil {
		t.Errorf("Failed to write file: %+v", err)
	}
	err = utils.WriteFile(
		regexPath, []byte(regexList), utils.FilePerms, utils.DirPerms)
	if err != nil {
		t.Errorf("Failed to write file: %+v", err)
	}
	defer func() {
		err = os.RemoveAll(usernamesPath)
		if err != nil {
			t.Errorf("Error deleting test file %q: %+v", usernamesPath, err)
		}
		err = os.RemoveAll(regexPath)
		if err != nil {
			t.Errorf("Error deleting test file %q: %+v", regexPath, err)
		}
	}()

	quit := make(chan struct{})
	m, err := NewManager(usernamesPath, regexPath, quit)
	if err != nil {
		t.Errorf("NewManager returned an error: %+v", err)
	}

	if !reflect.DeepEqual(m, expectedManager) {
		t.Errorf("New manager does not match expected."+
			"\nexpected: %+v\nreceived: %+v", expectedManager, m)
	}

	quit <- struct{}{}
}

// Tests that Manager.IsRestricted returns false for a list of known non-
// restricted usernames.
func TestManager_IsRestricted_GoodUsernames(t *testing.T) {
	usernameList := "Privategrity\nPrivategrity_Corp"
	regexList := "xx\nxx.*?network"

	m, deleteFunc := newTestManager(usernameList, regexList, t)
	defer deleteFunc()

	usernames := []string{
		"john_doe",
		"private",
		"network",
		"Privategrity!??!Corporation",
	}

	for _, username := range usernames {
		if m.IsRestricted(username) {
			t.Errorf("Username %q was recognized as restricted when it "+
				"should not be.", username)
		}
	}
}

// Tests that Manager.IsRestricted returns true for a list of known restricted
// usernames.
func TestManager_IsRestricted_BadUsernames(t *testing.T) {
	usernameList := "Privategrity\nPrivategrity_Corp"
	regexList := "xx\nxx.*?network"

	m, deleteFunc := newTestManager(usernameList, regexList, t)
	defer deleteFunc()

	usernames := []string{
		"xxfsdfsdfsdklfjnetwork",
		"Privategrity",
		"Privategrity_Corp",
		"exxplostion",
	}

	for _, username := range usernames {
		if !m.IsRestricted(validation.Canonicalize(username)) {
			t.Errorf("Username %q was not recognized as restricted when it "+
				"should have been.", username)
		}
	}
}

func TestManager_fileWatch(t *testing.T) {
	usernames := map[string]struct{}{
		"privategrity":      {},
		"privategrity_corp": {},
	}

	regexes := []*regexp.Regexp{
		regexp.MustCompile("xx"),
		regexp.MustCompile("xx.*?network"),
	}

	// Create manager and check it initialised the files correctly.
	m, deleteFunc := newTestManager(
		usernamesToList(usernames), regexesToList(regexes), t)
	defer deleteFunc()

	if !reflect.DeepEqual(m.usernames, usernames) {
		t.Errorf("Usernames in memory do not match usernames in file."+
			"\nexpected: %v\nreceived: %v", usernames, m.usernames)
	}
	if !reflect.DeepEqual(m.regexes, regexes) {
		t.Errorf("Regexes in memory do not match regexes in file."+
			"\nexpected: %v\nreceived: %v", regexes, m.regexes)
	}

	// Modify each file
	usernames = map[string]struct{}{
		"privategrity":      {},
		"privategrity_corp": {},
		"xxnetwork":         {},
	}

	regexes = []*regexp.Regexp{
		regexp.MustCompile("xx.*?network"),
		regexp.MustCompile("david.*?chaum"),
	}

	err := utils.WriteFile(m.usernamePath, []byte(usernamesToList(usernames)),
		utils.FilePerms, utils.DirPerms)
	if err != nil {
		t.Errorf("Failed to write username file: %+v", err)
	}
	err = utils.WriteFile(m.regexPath, []byte(regexesToList(regexes)),
		utils.FilePerms, utils.DirPerms)
	if err != nil {
		t.Errorf("Failed to write regex file: %+v", err)
	}

	time.Sleep(100 * time.Millisecond)

	// Check that the lists in memory match the new files
	if !reflect.DeepEqual(m.usernames, usernames) {
		t.Errorf("Usernames in memory do not match usernames in file."+
			"\nexpected: %v\nreceived: %v", usernames, m.usernames)
	}
	if !reflect.DeepEqual(m.regexes, regexes) {
		t.Errorf("Regexes in memory do not match regexes in file."+
			"\nexpected: %v\nreceived: %v", regexes, m.regexes)
	}

	// Modify each file
	usernames = map[string]struct{}{}

	regexes = []*regexp.Regexp{
		regexp.MustCompile("hi"),
	}

	err = utils.WriteFile(m.usernamePath, []byte(usernamesToList(usernames)),
		utils.FilePerms, utils.DirPerms)
	if err != nil {
		t.Errorf("Failed to write username file: %+v", err)
	}
	err = utils.WriteFile(m.regexPath, []byte(regexesToList(regexes)),
		utils.FilePerms, utils.DirPerms)
	if err != nil {
		t.Errorf("Failed to write regex file: %+v", err)
	}

	time.Sleep(100 * time.Millisecond)

	// Check that the lists in memory match the new files
	if !reflect.DeepEqual(m.usernames, usernames) {
		t.Errorf("Usernames in memory do not match usernames in file."+
			"\nexpected: %v\nreceived: %v", usernames, m.usernames)
	}
	if !reflect.DeepEqual(m.regexes, regexes) {
		t.Errorf("Regexes in memory do not match regexes in file."+
			"\nexpected: %v\nreceived: %v", regexes, m.regexes)
	}
}

// newTestManager creates two files for each list and loads them into a new
// manager. A function is returned to remove the files after the test.
func newTestManager(usernameList, regexList string, t *testing.T) (
	*Manager, func()) {
	timeNow := strconv.Itoa(int(time.Now().UnixNano()))
	usernamesPath := "restrictedUsernames-" + timeNow + ".txt"
	regexPath := "restrictedRegex-" + timeNow + ".txt"
	err := utils.WriteFile(
		usernamesPath, []byte(usernameList), utils.FilePerms, utils.DirPerms)
	if err != nil {
		t.Errorf("Failed to write file: %+v", err)
	}

	err = utils.WriteFile(
		regexPath, []byte(regexList), utils.FilePerms, utils.DirPerms)
	if err != nil {
		t.Errorf("Failed to write file: %+v", err)
	}

	quit := make(chan struct{})
	deleteFunc := func() {
		quit <- struct{}{}
		err = os.RemoveAll(usernamesPath)
		if err != nil {
			t.Errorf("Error deleting test file %q: %+v", usernamesPath, err)
		}
		err = os.RemoveAll(regexPath)
		if err != nil {
			t.Errorf("Error deleting test file %q: %+v", regexPath, err)
		}
	}

	m, err := NewManager(usernamesPath, regexPath, quit)
	if err != nil {
		t.Errorf("NewManager returned an error: %+v", err)
	}

	return m, deleteFunc
}

// usernamesToList converts the map of usernames to line-seperated list string.
func usernamesToList(usernames map[string]struct{}) string {
	usernameList := ""
	for username := range usernames {
		usernameList += username + "\n"
	}

	return usernameList
}

// regexesToList converts a list of regexp.Regexp to line-seperated list string.
func regexesToList(regexes []*regexp.Regexp) string {
	regexList := ""
	for _, regex := range regexes {
		regexList += regex.String() + "\n"
	}

	return regexList
}
