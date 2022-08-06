package mailer

import "time"

type json map[string]string
type email_request_bool_exp map[string]interface{}

// Email represents email address and name information
type Email struct {
	Address string `graphql:"address" json:"address"`
	Name    string `graphql:"name" json:"name,omitempty"`
}

// SendEmailInput represents the email request payload
type SendEmailInput struct {
	TemplateID       string    `json:"template_id,omitempty" graphql:"template_id"`
	From             string    `json:"from,omitempty" graphql:"from"`
	FromName         string    `json:"from_name,omitempty" graphql:"from_name"`
	To               []*Email  `json:"to,omitempty" graphql:"to" scalar:"true"`
	CC               []*Email  `json:"cc,omitempty" graphql:"cc" scalar:"true"`
	BCC              []*Email  `json:"bcc,omitempty" graphql:"bcc" scalar:"true"`
	Subject          string    `json:"subject,omitempty" graphql:"subject"`
	PlainTextContent string    `json:"content,omitempty" graphql:"content"`
	HTMLContent      string    `json:"html_content,omitempty" graphql:"html_content"`
	SendAfter        time.Time `json:"send_after,omitempty" graphql:"send_after"`
	Save             bool      `json:"save"`
	Locale           string    `json:"locale"`
}

// SendEmailResponse represents email response from external service
type SendEmailResponse struct {
	Success   bool        `json:"success" graphql:"success"`
	RequestID *string     `json:"request_id,omitempty" graphql:"request_id"`
	MessageID string      `json:"message_id,omitempty" graphql:"message_id"`
	Error     interface{} `json:"error,omitempty" graphql:"error"`
}

// SendEmailOutput represents the summary result of sending emails
type SendEmailOutput struct {
	Responses    []SendEmailResponse `json:"responses" graphql:"responses"`
	SuccessCount int                 `json:"success_count" graphql:"success_count"`
	FailureCount int                 `json:"failure_count" graphql:"failure_count"`
}

// NewEmails a shortcut for creating Email array
func NewEmails(address string, name string) []*Email {
	return []*Email{
		{
			Address: address,
			Name:    name,
		},
	}
}

// NewEmails a shortcut for creating Email instance
func NewEmail(address string, name string) *Email {
	return &Email{
		Address: address,
		Name:    name,
	}
}
