//go:build integration
// +build integration

package mailer

import (
	"context"
	"net/http"
	"os"
	"testing"
	"time"

	"github.com/hasura/go-graphql-client"
)

func cleanup(t *testing.T, client *Client) {

	_, err := client.CancelEmails(map[string]interface{}{})
	if err != nil {
		t.Fatal(err)
	}
}

// hasuraTransport transport for Hasura GraphQL Client
type hasuraTransport struct {
	adminSecret string
	headers     map[string]string
	// keep a reference to the client's original transport
	rt http.RoundTripper
}

// RoundTrip set header data before executing http request
func (t *hasuraTransport) RoundTrip(r *http.Request) (*http.Response, error) {
	if t.adminSecret != "" {
		r.Header.Set("X-Hasura-Admin-Secret", t.adminSecret)
	}
	for k, v := range t.headers {
		r.Header.Set(k, v)
	}
	return t.rt.RoundTrip(r)
}

func newGqlClient() *graphql.Client {
	adminSecret := os.Getenv("HASURA_GRAPHQL_ADMIN_SECRET")
	httpClient := &http.Client{
		Transport: &hasuraTransport{
			rt:          http.DefaultTransport,
			adminSecret: adminSecret,
		},
		Timeout: 30 * time.Second,
	}
	return graphql.NewClient(os.Getenv("DATA_URL"), httpClient)
}

func TestSendEmails(t *testing.T) {

	client := New(newGqlClient())
	defer cleanup(t, client)

	contents := "Test contents"
	results, err := client.Send([]SendEmailInput{
		{
			PlainTextContent: contents,
			HTMLContent:      "<p>test content</p>",
			Subject:          "test subject",
			To:               NewEmails("0123456789", "User"),
			Save:             true,
		},
	}, nil)
	if err != nil {
		t.Fatal(err)
	}

	var getQuery struct {
		EmailRequests []struct {
			ID string `json:"id"`
		} `graphql:"email_request(where: $where)"`
	}

	getVariables := map[string]interface{}{
		"where": email_request_bool_exp{
			"id": map[string]interface{}{
				"_eq": results.Responses[0].RequestID,
			},
		},
	}
	err = client.client.Query(context.TODO(), &getQuery, getVariables)
	if err != nil {
		t.Fatal(err)
	}
	if len(getQuery.EmailRequests) != 1 {
		t.Fatalf("expected 1 request, got: %d", len(getQuery.EmailRequests))
	}
}
