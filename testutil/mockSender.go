package testutil

import (
	"gitlab.com/elixxir/client/cmixproto"
	"gitlab.com/elixxir/primitives/id"
)

type MockSender struct {
	ch chan string
}

func (d *MockSender) Send(recipientID *id.User, msg string, msgType cmixproto.Type) {
	d.ch <- msg
}

func NewMockSender(ch chan string) *MockSender {
	d := MockSender{
		ch: ch,
	}
	return &d
}
