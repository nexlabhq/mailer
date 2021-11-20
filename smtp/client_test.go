package smtp

import (
	"errors"
	"io"
	"io/ioutil"
	"log"
	"testing"
	"time"

	"github.com/emersion/go-smtp"
	"github.com/nexlabhq/mailer/common"
	"github.com/stretchr/testify/assert"
)

// The mockBackend implements SMTP server methods.
type mockBackend struct{}

// Authenticate a user. Return smtp.ErrAuthUnsupported if you don't want to
// support this.
func (mb *mockBackend) Login(state *smtp.ConnectionState, username, password string) (smtp.Session, error) {
	ss := &mockSession{}
	return ss, ss.AuthPlain(username, password)
}

// Called if the client attempts to send mail without logging in first.
// Return smtp.ErrAuthRequired if you don't want to support this.
func (mb *mockBackend) AnonymousLogin(state *smtp.ConnectionState) (smtp.Session, error) {
	return mb.NewSession(state, "")
}

func (mb *mockBackend) NewSession(_ *smtp.ConnectionState, _ string) (smtp.Session, error) {
	return &mockSession{}, nil
}

// A mockSession is returned after EHLO.
type mockSession struct{}

func (s *mockSession) AuthPlain(username, password string) error {
	if username != "user@example.com" || password != "password" {
		return errors.New("Invalid username or password")
	}
	return nil
}

func (s *mockSession) Mail(from string, opts smtp.MailOptions) error {
	log.Println("Mail from:", from)
	return nil
}

func (s *mockSession) Rcpt(to string) error {
	log.Println("Rcpt to:", to)
	return nil
}

func (s *mockSession) Data(r io.Reader) error {
	if b, err := ioutil.ReadAll(r); err != nil {
		return err
	} else {
		log.Println("Data:", string(b))
	}
	return nil
}

func (s *mockSession) Reset() {}

func (s *mockSession) Logout() error {
	return nil
}

func createMockServer() *smtp.Server {

	be := &mockBackend{}

	s := smtp.NewServer(be)

	s.Addr = ":1025"
	s.Domain = "localhost"
	s.ReadTimeout = 10 * time.Second
	s.WriteTimeout = 10 * time.Second
	s.MaxMessageBytes = 1024 * 1024
	s.MaxRecipients = 50
	s.AllowInsecureAuth = true

	return s
}

func TestSendGridClient(t *testing.T) {

	server := createMockServer()

	go func() {
		if err := server.ListenAndServe(); err != nil {
			log.Fatal(err)
		}
	}()

	defer server.Close()

	client := New(&Config{
		User:     "user@example.com",
		Password: "password",
		Host:     "localhost",
		Port:     1025,
	})

	input := &common.SendRequest{
		From:     "user@example.com",
		FromName: "test",
		To: []*common.Email{
			{Address: "to@nexlab.local", Name: "To user"},
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
	assert.NoError(t, err)

}
