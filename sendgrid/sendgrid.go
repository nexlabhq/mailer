package sendgrid

import (
	"errors"

	"github.com/nexlabhq/mailer/common"
	goSendGrid "github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"github.com/sirupsen/logrus"
)

type SendGrid struct {
	logger *logrus.Entry
	*goSendGrid.Client
}

func New(apiKey string) *SendGrid {
	return &SendGrid{
		Client: goSendGrid.NewSendClient(apiKey),
	}
}

func (sg *SendGrid) SetLogger(logger *logrus.Entry) {
	sg.logger = logger
}

func (sg *SendGrid) Send(input *common.SendRequest) error {

	email := sg.buildEmailRequest(input)
	response, err := sg.Client.Send(email)
	if err != nil {
		return err
	}

	if sg.logger != nil {
		sg.logger.WithFields(logrus.Fields{
			"type":             "email",
			"vendor":           "sendgrid",
			"status_code":      response.StatusCode,
			"response_headers": response.Headers,
		}).Infof(response.Body)
	}

	if response.StatusCode > 299 {
		return errors.New(response.Body)
	}
	return nil
}

func (sg *SendGrid) buildEmailRequest(input *common.SendRequest) *mail.SGMailV3 {

	email := mail.NewV3Mail()
	email.Subject = input.Subject
	if input.PlainTextContent != "" {
		email.AddContent(mail.NewContent("text/plain", input.PlainTextContent))
	}
	if input.HTMLContent != "" {
		email.AddContent(mail.NewContent("text/html", input.HTMLContent))
	}
	email.SetFrom(mail.NewEmail(input.FromName, input.From))

	var addresses []string
	p := mail.NewPersonalization()
	for _, to := range input.To {
		addresses = append(addresses, to.Address)
		p.AddTos(mail.NewEmail(to.Name, to.Address))
	}
	for _, cc := range input.CC {
		if !common.HasString(addresses, cc.Address) {
			addresses = append(addresses, cc.Address)
			p.AddCCs(mail.NewEmail(cc.Name, cc.Address))
		}
	}
	for _, bcc := range input.BCC {
		if !common.HasString(addresses, bcc.Address) {
			addresses = append(addresses, bcc.Address)
			p.AddBCCs(mail.NewEmail(bcc.Name, bcc.Address))
		}
	}

	email.AddPersonalizations(p)

	return email
}
