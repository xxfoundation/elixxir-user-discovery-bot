////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package io

import (
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

// Contains logic for username validation.

// Character limits for usernames.
const (
	minimumUsernameLength = 4
	maximumUsernameLength = 32
)

// usernameRegex is the regular expression for the enforcing the following characters only:
//  abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-+.@#
// Furthermore, the regex enforces usernames to begin and end with an alphanumeric character.
var usernameRegex = regexp.MustCompile("^[a-zA-Z0-9][a-zA-Z0-9_\\-+@.#]*[a-zA-Z0-9]$")

// isValidUsername determines whether the username is of an acceptable length and
// whether it contains allowed character. The allowed characters are defined
// by usernameRegex.
func isValidUsername(username string) error {
	// Check for acceptable length
	if len(username) < minimumUsernameLength || len(username) > maximumUsernameLength {
		return errors.Errorf("username length %d is not between %d and %d",
			len(username), minimumUsernameLength, maximumUsernameLength)
	}

	// Check is username contains allowed characters only
	if !usernameRegex.MatchString(username) {
		return errors.Errorf("username can only contain alphanumberics " +
			"and the symbols _, -, +, ., @, # and must start and end with an alphanumeric character")
	}

	return nil
}

// canonicalize reduces the username to its canonical form. For the purposes
// of internal usage only.
func canonicalize(username string) string {
	return strings.ToLower(username)
}
