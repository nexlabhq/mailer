package common

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestParseEmail(t *testing.T) {
	for i, s := range []struct {
		Input  string
		Output Email
	}{
		{
			"test@example.com",
			Email{Address: "test@example.com", Name: ""},
		},
		{
			"<test@example.com>",
			Email{Address: "test@example.com", Name: ""},
		},
		{
			"Test user <user@example.com>",
			Email{Address: "user@example.com", Name: "Test user"},
		},
		{
			"Test user 2 user2@example.com",
			Email{Address: "user2@example.com", Name: "Test user 2"},
		},
	} {
		output, err := ParseEmail(s.Input)
		assert.NoError(t, err, "%d: no error", i)
		assert.Equal(t, s.Output.Address, output.Address, "%d: address", i)
		assert.Equal(t, s.Output.Name, output.Name, "%d: name", i)
	}
}
