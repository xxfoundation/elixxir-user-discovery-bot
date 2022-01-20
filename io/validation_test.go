////////////////////////////////////////////////////////////////////////////////
// Copyright © 2020 Privategrity Corporation                                   /
//                                                                             /
// All rights reserved.                                                        /
////////////////////////////////////////////////////////////////////////////////

package io

import "testing"

// Tests whether good usernames are considered valid.
func TestIsValidUsername_GoodUsernames(t *testing.T) {
	// Construct a list of good username
	goodUsernames := []string{
		"abcdefghijklmnopqrstuvwxyzABCDEF",
		"GHIJKLMNOPQRSTUVWXYZ0123456789",
		"john_doe",
		"daMan",
		"Mr.George",
		"josh-@+#b",
		"A........5",
	}

	// Test whether every good username is valid
	for _, goodUsername := range goodUsernames {
		err := isValidUsername(goodUsername)
		if err != nil { // If invalid, fail test
			t.Errorf("isValidUsername failed with username %q: %v", goodUsername, err)
		}
	}

}

// Tests whether invalid usernames are considered invalid.
func TestIsValidUsername_BadUsernames(t *testing.T) {
	// Construct a list of bad usernames
	badUsernames := []string{
		"",
		"  ",
		"pie",
		"123456789012345678901234567890123",
		"аdМіnіѕТrаТоr",
		"ÁdmïNIstrátör",
		"𝔞𝔡𝔪𝔦𝔫",
		"á̵͖͔͇̟͙̜͙̀͑̒̀̕ď̶̦̣̲m̴̡̬̺̯̩͂i̶͚͍̝̋ͅn̶̤̙̩͎̠͍̙̱̹̏",
		"ﬁnished",
		"GHIJKLMNOPQRSTUVWXYZ0123456789_-",
		"!@#$%^*?",
		"josh!!!!!",
	}

	// Test if every bad username is invalid
	for _, badUsername := range badUsernames {
		err := isValidUsername(badUsername)
		if err == nil { // If considered valid, fail test
			t.Errorf("isValidUsername did not fail with username %q", badUsername)
		}
	}

}

// Consistency test for the Canonicalize function.
func TestCanonicalize(t *testing.T) {
	inputList := []string{
		"John_Doe",
		"HELLO",
		"hello",
		"abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789_-!@#$%^*?",
	}

	expected := []string{
		"john_doe",
		"hello",
		"hello",
		"abcdefghijklmnopqrstuvwxyzabcdefghijklmnopqrstuvwxyz0123456789_-!@#$%^*?",
	}

	for i, input := range inputList {
		received := Canonicalize(input)
		if received != expected[i] {
			t.Errorf("Canonicalize did not produce expeccted result with %q"+
				"\nExpected: %s"+
				"\nReceived: %s", input, expected[i], received)
		}
	}
}
