package smtp

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"github.com/emersion/go-message"
	"github.com/nexlabhq/mailer/common"
)

func buildEmailRequest(input *common.SendRequest) ([]string, string, error) {

	var b bytes.Buffer
	var h message.Header
	var emails []string
	charset := map[string]string{"charset": "utf-8"}
	h.Set("Subject", input.Subject)
	h.Set("From", fmt.Sprintf("%s <%s>", input.FromName, input.From))
	h.SetContentType("multipart/alternative", charset)

	to := make([]string, 0)
	for _, t := range input.To {
		to = append(to, t.String())
		emails = append(emails, t.Address)
	}
	h.Set("To", strings.Join(to, ", "))

	if len(input.CC) > 0 {
		cc := make([]string, 0)

		for _, t := range input.CC {
			cc = append(cc, t.String())
			emails = append(emails, t.Address)
		}
		h.Set("Cc", strings.Join(cc, ", "))
	}

	// Sending "Bcc" messages is accomplished by including an email address
	// in the to parameter but not including it in the msg headers.
	if len(input.BCC) > 0 {
		for _, t := range input.BCC {
			emails = append(emails, t.Address)
		}
	}

	w, err := message.CreateWriter(&b, h)
	if err != nil {
		return nil, "", err
	}

	if input.PlainTextContent != "" {
		var plainTextHeader message.Header
		plainTextHeader.SetContentType("text/plain", charset)
		w1, err := w.CreatePart(plainTextHeader)
		if err != nil {
			return nil, "", err
		}

		io.WriteString(w1, input.PlainTextContent)

		w1.Close()
	}

	if input.HTMLContent != "" {
		var htmlHeader message.Header
		htmlHeader.SetContentType("text/html", charset)
		w2, err := w.CreatePart(htmlHeader)
		if err != nil {
			return nil, "", err
		}

		io.WriteString(w2, input.HTMLContent)

		w2.Close()
	}

	w.Close()

	return emails, b.String(), err
}
