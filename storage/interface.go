////////////////////////////////////////////////////////////////////////////////
// Copyright © 2018 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

// Interface and enums for UDB storage systems
package storage

// The ValueType constant stores the allowable types we search on
// (e-mail, group, nickname, etc).
type ValueType int

// Note: because DB backends vary, and this list could have
//       items added and removed, we are not using an iota on purpose.
const (
	Email ValueType = 0 // An e-mail address
	Nick  ValueType = 1 // The user's nickname
	// TODO: Add more as necessary
)

// Print strings for ValueType
func (v ValueType) String() string {
	names := [...]string{
		"Email",
		"Nick",
	}
	if v < Email || v > Nick {
		return "Unknown"
	}
	return names[v]
}
