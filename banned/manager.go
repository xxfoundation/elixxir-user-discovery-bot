///////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 xx network SEZC                                          //
//                                                                           //
// Use of this source code is governed by a license that can be found in the //
// LICENSE file                                                              //
///////////////////////////////////////////////////////////////////////////////

package banned

import (
	"github.com/pkg/errors"
	jww "github.com/spf13/jwalterweatherman"
	"gitlab.com/elixxir/user-discovery-bot/validation"
	"regexp"
	"strings"
	"testing"
)

// Manager contains two lists of banned/reserved usernames. It handles
// checking for banned usernames on user registration.
type Manager struct {
	// A simple banned user lookup. Any user that exactly matches something
	// in this map will be considered banned/reserved.
	bannedUserList map[string]struct{}

	// A more complex banned user lookup. Contains a list of regular expressions,
	// and any username which matches any of these regular expressions will be
	// considered banned/reserved.
	bannedRegexList []*regexp.Regexp
}

// NewManager constructs the banned.Manager object. NewManager is passed in
// two text files containing a list where values are separated by the
// Linux newline ("\n"). NewManager will parse these two lists separately to create a
// Manager.bannedUserList and a Manager.bannedRegexList respectively.
func NewManager(bannedUserFile, bannedRegexFile string) (*Manager, error) {
	// Construct a map of banned/reserved usernames
	bannedUsers := make(map[string]struct{})
	if bannedUserFile != "" {
		bannedUserList := strings.Split(bannedUserFile, "\n")
		for _, bannedUser := range bannedUserList {
			if bannedUser == "" { // Skip any empty lines
				continue
			}
			bannedUsers[validation.Canonicalize(bannedUser)] = struct{}{}
		}
	}

	// Construct a regex list for banned/reserved usernames
	bannedRegexList := make([]*regexp.Regexp, 0)
	if bannedRegexFile != "" {
		regexList := strings.Split(bannedRegexFile, "\n")
		for _, bannedRegex := range regexList {
			if bannedRegex == "" { // Skip any empty lines
				continue
			}

			// Compile regex expression
			regex, err := regexp.Compile(bannedRegex)
			if err != nil {
				return nil, errors.Errorf("Failed to compile banned user regex %q: %v", bannedRegex, err)
			}

			bannedRegexList = append(bannedRegexList, regex)

		}
	}

	return &Manager{
		bannedUserList:  bannedUsers,
		bannedRegexList: bannedRegexList,
	}, nil
}

// IsBanned checks if the username is in Manager's bannedUserList or
// matched to any banned regular expression.
func (m *Manager) IsBanned(username string) bool {
	_, exists := m.bannedUserList[username]
	if exists {
		return exists
	}

	return m.isRegexBanned(username)
}

// isRegexBanned checks is the username matches any banned regular expression.
func (m *Manager) isRegexBanned(username string) bool {
	for _, regex := range m.bannedRegexList {
		if regex.MatchString(username) {
			return true
		}
	}

	return false
}

// SetBannedTest is a testing only helper function which sets a username
// in Manager's bannedUserList.
func (m *Manager) SetBannedTest(username string, t *testing.T) {
	if t == nil {
		jww.FATAL.Panic("Cannot use this outside of testing")
	}

	m.bannedUserList[username] = struct{}{}
}
