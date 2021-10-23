package email

import (
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strings"

	"github.com/hasura/go-graphql-client"
	"github.com/nexlabhq/mailer/common"
	"github.com/nexlabhq/mailer/sendgrid"
	"github.com/sirupsen/logrus"
)

type EmailVendor string

const (
	VendorSendGrid = "sendgrid"
)

type SendRequest common.SendRequest

func (sr SendRequest) Model() *common.SendRequest {
	srr := common.SendRequest(sr)
	return &srr
}

type IEmailer interface {
	Send(input *common.SendRequest) error
}

type EmailConfig struct {
	EmailVendor    string `envconfig:"EMAIL_VENDOR" required:"true"`
	EmailFrom      string `envconfig:"EMAIL_FROM"`
	EmailFromName  string `envconfig:"EMAIL_FROM_NAME"`
	SendGridAPIKey string `envconfig:"SENDGRID_API_KEY"`
	EmailLocale    string `envconfig:"EMAIL_LOCALE"`
}

type Config struct {
	EmailConfig
	DataClient *graphql.Client
	Logger     *logrus.Entry
}

func (c Config) getIEmail() (IEmailer, error) {
	switch c.EmailVendor {
	case VendorSendGrid:
		if c.SendGridAPIKey == "" {
			return nil, errors.New("SendGrid API Key is required")
		}

		client := sendgrid.New(c.EmailConfig.SendGridAPIKey)
		if c.Logger != nil {
			client.SetLogger(c.Logger)
		}
		return client, nil
	}

	return nil, fmt.Errorf("invalid vendor %s", c.EmailVendor)
}

type EmailTemplate struct {
	ID           string            `json:"id"`
	Subjects     map[string]string `json:"subjects"`
	Contents     map[string]string `json:"contents"`
	HTMLContents map[string]string `json:"html_contents"`
}

type EmailTemplateRaw struct {
	ID           string          `graphql:"id" json:"id"`
	Subjects     json.RawMessage `graphql:"subjects" json:"subjects"`
	Contents     json.RawMessage `graphql:"contents" json:"contents"`
	HTMLContents json.RawMessage `graphql:"html_contents" json:"html_contents"`
}

func (etr *EmailTemplateRaw) Parse() (*EmailTemplate, error) {

	var subjects map[string]string
	var contents map[string]string
	var htmlContents map[string]string
	if len(etr.Contents) > 0 {
		err := json.Unmarshal(etr.Contents, &contents)
		if err != nil {
			return nil, err
		}
	}
	if len(etr.Subjects) > 0 {
		err := json.Unmarshal(etr.Subjects, &subjects)
		if err != nil {
			return nil, err
		}
	}
	if len(etr.HTMLContents) > 0 {
		err := json.Unmarshal(etr.HTMLContents, &htmlContents)
		if err != nil {
			return nil, err
		}
	}

	return &EmailTemplate{
		ID:           etr.ID,
		Subjects:     subjects,
		Contents:     contents,
		HTMLContents: htmlContents,
	}, nil
}

type email_template_bool_exp map[string]interface{}
type email_template_insert_input EmailTemplate

// uniqueStrings is the special array string that only store unique values
type uniqueStrings map[string]bool

// Add append new value or skip if it's existing
func (us uniqueStrings) Add(values ...string) {
	for _, s := range values {
		if _, ok := us[s]; !ok {
			us[s] = true
		}
	}
}

// IsEmpty check if the array is empty
func (us uniqueStrings) IsEmpty() bool {
	return len(us) == 0
}

// Value return
func (us uniqueStrings) Value() []string {
	results := make([]string, 0, len(us))
	for k := range us {
		results = append(results, k)
	}
	return results
}

// String implement string interface
func (us uniqueStrings) String() string {
	results := us.Value()
	sort.Strings(results)
	return strings.Join(results, ",")
}
