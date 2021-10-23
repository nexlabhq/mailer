package sendgrid

import (
	"fmt"
	"testing"

	"github.com/nexlabhq/mailer/common"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/stretchr/testify/assert"
)

func TestSendGridClient(t *testing.T) {

	apiKey := "random_secret"
	client := New(apiKey)
	assert.Equal(t, fmt.Sprintf("Bearer %s", apiKey), client.Client.Headers["Authorization"])

	input := &common.SendRequest{
		From:     "toan.nguyen@nexlab.tech",
		FromName: "test",
		To: []*common.Email{
			{Address: "user@telehealth.nexlab", Name: "To user"},
		},
		CC: []*common.Email{
			{Address: "user@telehealth.nexlab", Name: "Cc user"},
		},
		BCC: []*common.Email{
			{Address: "user@telehealth.nexlab", Name: "Bcc user"},
		},
		Subject:          "Test email",
		PlainTextContent: "Hello world",
		HTMLContent:      "<h1>Hello</h1>",
	}
	email := client.buildEmailRequest(input)
	settings := mail.NewMailSettings()
	settings.SandboxMode = mail.NewSetting(true)

	email.SetMailSettings(settings)
	assert.Equal(t, "Test email", email.Subject)
	assert.Equal(t, input.Subject, email.Subject)
	for _, content := range email.Content {
		if content.Type == "text/html" {
			assert.Equal(t, input.HTMLContent, content.Value)
		} else {
			assert.Equal(t, input.PlainTextContent, content.Value)
		}
	}

	for i, to := range email.Personalizations[0].To {
		assert.Equal(t, input.To[i].Address, to.Address)
		assert.Equal(t, input.To[i].Name, to.Name)
	}
	for i, cc := range email.Personalizations[0].CC {
		assert.Equal(t, input.CC[i].Address, cc.Address)
		assert.Equal(t, input.CC[i].Name, cc.Name)
	}
	for i, bcc := range email.Personalizations[0].BCC {
		assert.Equal(t, input.BCC[i].Address, bcc.Address)
		assert.Equal(t, input.BCC[i].Name, bcc.Name)
	}

	response, err := client.Client.Send(email)
	assert.NoError(t, err)
	assert.Equal(t, 401, response.StatusCode)

}
