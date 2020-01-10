package testutil

import (
	"gitlab.com/elixxir/client/cmixproto"
	"gitlab.com/elixxir/primitives/id"
	"strings"
)

type MockSender struct {
	Val     bool
	find    string
	exclude string
}

func (d *MockSender) Send(recipientID *id.User, msg string, msgType cmixproto.Type) {
	if strings.Contains(msg, d.find) && (d.exclude == "" || !strings.Contains(msg, d.exclude)) {
		d.Val = true
	}
}

func NewMockSender(find, exclude string) *MockSender {
	d := MockSender{
		Val:     false,
		find:    find,
		exclude: exclude,
	}
	return &d
}
