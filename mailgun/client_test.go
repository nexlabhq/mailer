package mailgun

import (
	"testing"
	"time"

	"github.com/nexlabhq/mailer/common"
	"github.com/stretchr/testify/assert"
)

func TestMailgunClient(t *testing.T) {

	client := New(&Config{
		ApiKey:  "random_secret",
		Domain:  "example.com",
		Timeout: 10 * time.Second,
	})

	input := &common.SendRequest{
		From:     "sender@nexlab.local",
		FromName: "test",
		To: []*common.Email{
			{Address: "user@nexlab.local", Name: "To user"},
		},
		CC: []*common.Email{
			{Address: "cc@nexlab.local", Name: "Cc user"},
		},
		BCC: []*common.Email{
			{Address: "bcc@nexlab.local", Name: "Bcc user"},
		},
		Subject:          "Test email",
		PlainTextContent: "Hello world",
		HTMLContent:      "<h1>Hello</h1>",
	}

	err := client.Send(input)
	assert.Error(t, err)

}
