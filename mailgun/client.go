package mailgun

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/mailgun/mailgun-go/v4"
	"github.com/nexlabhq/mailer/common"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Domain  string        `envconfig:"MAILGUN_DOMAIN"`
	ApiKey  string        `envconfig:"MAILGUN_API_KEY"`
	BaseURL string        `envconfig:"MAILGUN_BASE_URL"`
	Timeout time.Duration `envconfig:"MAILGUN_TIMEOUT" default:"30s"`
}

func (c Config) Validate() error {
	if c.Domain == "" {
		return errors.New("MAILGUN_DOMAIN is required")
	}
	if c.ApiKey == "" {
		return errors.New("MAILGUN_API_KEY is required")
	}

	return nil
}

type Client struct {
	client  *mailgun.MailgunImpl
	logger  *logrus.Entry
	timeout time.Duration
}

func New(config *Config) *Client {
	client := mailgun.NewMailgun(config.Domain, config.ApiKey)
	if config.BaseURL != "" {
		client.SetAPIBase(config.BaseURL)
	}

	return &Client{
		client:  client,
		timeout: config.Timeout,
	}
}

func (c *Client) SetLogger(logger *logrus.Entry) {
	c.logger = logger
}

func (c *Client) Send(input *common.SendRequest) error {

	logError := func(err error) {
		if err != nil && c.logger != nil {
			c.logger.WithFields(logrus.Fields{
				"type":   "email",
				"vendor": "mailgun",
				"error":  err,
			}).Error(err)
		}
	}

	message := c.buildEmailRequest(input)
	ctx, cancel := context.WithTimeout(context.Background(), c.timeout)
	defer cancel()

	resp, id, err := c.client.Send(ctx, message)

	if err != nil {
		logError(err)
		return err
	}

	if c.logger != nil {
		c.logger.WithFields(logrus.Fields{
			"type":       "email",
			"vendor":     "mailgun",
			"response":   resp,
			"message_id": id,
		}).Error(err)
	}
	return nil
}

func (c *Client) buildEmailRequest(input *common.SendRequest) *mailgun.Message {

	message := c.client.NewMessage(
		fmt.Sprintf("%s <%s>", input.FromName, input.From),
		input.Subject,
		input.PlainTextContent,
	)
	if input.HTMLContent != "" {
		message.SetHtml(input.HTMLContent)
	}

	for _, to := range input.To {
		message.AddRecipient(to.String())
	}
	for _, cc := range input.CC {
		message.AddCC(cc.String())
	}
	for _, bcc := range input.BCC {
		message.AddBCC(bcc.String())
	}

	return message
}
