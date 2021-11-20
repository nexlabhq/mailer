package smtp

import (
	"errors"
	"fmt"
	"strings"

	"github.com/emersion/go-sasl"
	"github.com/emersion/go-smtp"
	"github.com/nexlabhq/mailer/common"
	"github.com/sirupsen/logrus"
)

type Config struct {
	Host     string `envconfig:"SMTP_HOST" default:"smtp.gmail.com"`
	Port     int    `envconfig:"SMTP_PORT" default:"587"`
	User     string `envconfig:"SMTP_USER"`
	Password string `envconfig:"SMTP_PASSWORD"`
}

func (c Config) Validate() error {
	if c.User == "" {
		return errors.New("SMTP_USER is required")
	}
	if c.Password == "" {
		return errors.New("SMTP_PASSWORD is required")
	}

	return nil
}

type Client struct {
	config *Config
	logger *logrus.Entry
}

func New(config *Config) *Client {
	return &Client{config: config}
}

func (c *Client) SetLogger(logger *logrus.Entry) {
	c.logger = logger
}

func (c *Client) Send(input *common.SendRequest) error {

	logError := func(err error) {
		if err != nil && c.logger != nil {
			c.logger.WithFields(logrus.Fields{
				"type":   "email",
				"vendor": "smtp",
				"error":  err,
			}).Error(err)
		}
	}
	// Setup authentication information.
	auth := sasl.NewPlainClient("", c.config.User, c.config.Password)

	// Connect to the server, authenticate, set the sender and recipient,
	// and send the email all in one step.
	mailAddresses, msg, err := buildEmailRequest(input)
	if err != nil {
		logError(err)
		return err
	}

	err = smtp.SendMail(
		fmt.Sprintf("%s:%d", c.config.Host, c.config.Port),
		auth,
		input.From, mailAddresses, strings.NewReader(msg),
	)
	logError(err)

	return err
}
