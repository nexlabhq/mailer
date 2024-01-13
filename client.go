package mailer

import (
	"context"
	"time"

	"github.com/hasura/go-graphql-client"
	"github.com/hgiasac/graphql-utils/client"
)

// Client a high level mail client that communicates the backend through GraphQL client
type Client struct {
	client client.Client
}

// New constructs a mail client
func New(client client.Client) *Client {
	return &Client{
		client: client,
	}
}

// Send a mail request
func (c *Client) Send(inputs []SendEmailInput, variables map[string]interface{}) (*SendEmailOutput, error) {
	if len(inputs) == 0 {
		return &SendEmailOutput{}, nil
	}

	for i, input := range inputs {
		if input.SendAfter.IsZero() {
			inputs[i].SendAfter = time.Now()
		}
	}
	var mutation struct {
		SendEmails SendEmailOutput `graphql:"sendEmails(data: $data, variables: $variables)"`
	}

	inputVariables := map[string]interface{}{
		"data":      inputs,
		"variables": json(variables),
	}

	err := c.client.Mutate(context.Background(), &mutation, inputVariables, graphql.OperationName("SendEmails"))
	if err != nil {
		return nil, err
	}

	return &mutation.SendEmails, nil
}

// CancelEmails cancel and delete email requests
func (c *Client) CancelEmails(where map[string]interface{}) (int, error) {

	var mutation struct {
		DeleteEmails struct {
			AffectedRows int `graphql:"affected_rows"`
		} `graphql:"delete_email_request(where: $where)"`
	}

	variables := map[string]interface{}{
		"where": email_request_bool_exp(where),
	}

	err := c.client.Mutate(context.Background(), &mutation, variables, graphql.OperationName("DeleteEmailRequests"))
	if err != nil {
		return 0, err
	}

	return mutation.DeleteEmails.AffectedRows, nil
}
