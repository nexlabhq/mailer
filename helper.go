package email

import (
	"bytes"
	"fmt"
	htmlTemplate "html/template"
	"text/template"
)

func ParseTemplate(nt *EmailTemplate, variables interface{}) (*EmailTemplate, error) {
	subjects := make(map[string]string)
	contents := make(map[string]string)
	htmlContents := make(map[string]string)

	for k, v := range nt.Subjects {
		tName := fmt.Sprintf("%s:%s:%s", nt.ID, "subject", k)
		t, err := template.New(tName).Parse(v)
		if err != nil {
			return nil, err
		}
		var b bytes.Buffer
		if err = t.Execute(&b, variables); err != nil {
			return nil, err
		}
		subjects[k] = b.String()
	}

	for k, v := range nt.Contents {
		tName := fmt.Sprintf("%s:%s:%s", nt.ID, "content", k)
		t, err := template.New(tName).Parse(v)
		if err != nil {
			return nil, err
		}
		var b bytes.Buffer
		if err = t.Execute(&b, variables); err != nil {
			return nil, err
		}
		contents[k] = b.String()
	}

	for k, v := range nt.HTMLContents {
		tName := fmt.Sprintf("%s:%s:%s", nt.ID, "html", k)
		t, err := htmlTemplate.New(tName).Parse(v)
		if err != nil {
			return nil, err
		}
		var b bytes.Buffer
		if err = t.Execute(&b, variables); err != nil {
			return nil, err
		}
		htmlContents[k] = b.String()
	}

	return &EmailTemplate{
		ID:           nt.ID,
		Subjects:     subjects,
		Contents:     contents,
		HTMLContents: htmlContents,
	}, nil
}
