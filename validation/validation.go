////////////////////////////////////////////////////////////////////////////////
// Copyright © 2022 xx foundation                                             //
//                                                                            //
// Use of this source code is governed by a license that can be found in the  //
// LICENSE file.                                                              //
////////////////////////////////////////////////////////////////////////////////

package validation

import (
	"github.com/pkg/errors"
	"regexp"
	"strings"
)

// Todo: move validation of username logic to the face package in primitives

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

// IsValidUsername determines whether the username is of an acceptable length and
// whether it contains allowed character. The allowed characters are defined
// by usernameRegex.
func IsValidUsername(username string) error {
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

// Canonicalize reduces the username to its canonical form. For the purposes
// of internal usage only.
func Canonicalize(username string) string {
	return strings.ToLower(username)
}
