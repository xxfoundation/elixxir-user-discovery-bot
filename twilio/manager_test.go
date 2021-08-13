package twilio

import (
	"git.xx.network/elixxir/user-discovery-bot/interfaces/params"
	"testing"
)

func TestNewManager(t *testing.T) {
	m := NewManager(params.Twilio{}, nil)
	if m == nil {
		t.Error("This should not happen")
	}
}
