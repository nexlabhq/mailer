package email

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"

	"github.com/hasura/go-graphql-client"
	"github.com/nexlabhq/mailer/common"
	"github.com/sirupsen/logrus"
)

type Emailer struct {
	config     EmailConfig
	emailer    IEmailer
	locale     string
	dataClient *graphql.Client
	logger     *logrus.Entry
}

func New(config Config) (*Emailer, error) {
	var emailer IEmailer
	var err error
	// don't require emailer vendor if we only save to queue
	if config.EmailVendor != "" {
		emailer, err = config.getIEmail()
		if err != nil {
			return nil, err
		}
	}

	return &Emailer{
		config:     config.EmailConfig,
		emailer:    emailer,
		dataClient: config.DataClient,
		logger:     config.Logger,
		locale:     config.EmailLocale,
	}, nil
}

func (em *Emailer) Send(inputs []*SendRequest) error {
	if em.emailer == nil {
		return errors.New("required emailer vendor")
	}

	for _, input := range inputs {
		if input.From == "" {
			input.From = em.config.EmailFrom
		}
		if input.FromName == "" {
			input.FromName = em.config.EmailFromName
		}

		err := em.emailer.Send(input.Model())
		if err != nil {
			return err
		}
	}

	return nil
}

func (em *Emailer) SendQueue(inputs []*SendRequest) error {
	if em.dataClient == nil {
		return errors.New("DATA_CLIENT is required")
	}

	for _, input := range inputs {
		err := em.emailer.Send(input.Model())
		if err != nil {
			return err
		}
	}

	return nil
}

func (em *Emailer) SendWithTemplates(inputs []*SendRequest, variables interface{}) error {
	if len(inputs) == 0 {
		return nil
	}

	if em.dataClient == nil {
		return errors.New("DATA_CLIENT is required")
	}

	templateIDs := uniqueStrings{}
	newInputs := make([]*SendRequest, 0, len(inputs))

	for _, input := range inputs {
		if input.TemplateID != "" {
			templateIDs.Add(input.TemplateID)
		}
	}
	if !templateIDs.IsEmpty() {
		templates, err := em.GetTemplateByIDs(templateIDs.Value()...)
		if err != nil {
			return err
		}

		for _, item := range inputs {
			if item.TemplateID != "" {
				template, ok := templates[item.TemplateID]
				if !ok {
					return fmt.Errorf("email template not found: %s", item.TemplateID)
				}
				newItem, err := ParseTemplate(template, variables)
				if err != nil {
					return err
				}

				item.Subject = em.getLocaleContent(newItem.Subjects)
				item.PlainTextContent = em.getLocaleContent(newItem.Contents)
				item.HTMLContent = em.getLocaleContent(newItem.HTMLContents)
				newInputs = append(newInputs, item)
			} else {
				newInputs = append(newInputs, item)
			}
		}
	}

	return em.Send(newInputs)
}

func (em *Emailer) GetTemplateByIDs(ids ...string) (map[string]*EmailTemplate, error) {
	results := make(map[string]*EmailTemplate)
	if len(ids) == 0 {
		return results, nil
	}

	var query struct {
		NotificationTemplates []EmailTemplateRaw `graphql:"email_template(where: $where)" json:"email_template"`
	}

	variables := map[string]interface{}{
		"where": email_template_bool_exp{
			"id": map[string]interface{}{
				"_in": ids,
			},
		},
	}

	bytes, err := em.dataClient.NamedQueryRaw(context.Background(), "GetEmailTemplatesByIds", &query, variables)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(*bytes, &query)
	if err != nil {
		return nil, err
	}

	for _, ntr := range query.NotificationTemplates {
		nt, err := ntr.Parse()
		if err != nil {
			return nil, err
		}
		results[nt.ID] = nt
	}
	return results, nil
}

func (em *Emailer) UpsertTemplates(inputs []*EmailTemplate) ([]*EmailTemplate, error) {
	if len(inputs) == 0 {
		return []*EmailTemplate{}, nil
	}

	var mutation struct {
		UpsertEmailTemplates struct {
			Returning []EmailTemplateRaw `graphql:"returning" json:"returning"`
		} `graphql:"insert_email_template(objects: $objects, on_conflict: { constraint: notification_template_pkey, update_columns: [contents, headings] })" json:"insert_notification_template"`
	}

	objects := make([]email_template_insert_input, len(inputs))
	for i, nt := range inputs {
		objects[i] = email_template_insert_input(*nt)
	}

	variables := map[string]interface{}{
		"objects": objects,
	}

	bytes, err := em.dataClient.NamedMutateRaw(context.Background(), "UpsertEmailTemplates", &mutation, variables)
	if err != nil {
		return nil, err
	}

	err = json.Unmarshal(*bytes, &mutation)
	if err != nil {
		return nil, err
	}

	results := make([]*EmailTemplate, 0, len(mutation.UpsertEmailTemplates.Returning))

	for _, ntr := range mutation.UpsertEmailTemplates.Returning {
		nt, err := ntr.Parse()
		if err != nil {
			return nil, err
		}
		results = append(results, nt)
	}
	return results, nil
}

func (em *Emailer) getLocaleContent(contents map[string]string) string {
	if em.config.EmailLocale == "" {
		return common.GetFirstStringInMap(contents)
	}
	if s, ok := contents[em.locale]; ok {
		return s
	}
	return common.GetFirstStringInMap(contents)
}
