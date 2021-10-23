package common

import (
	"errors"
	"strings"
)

type SendRequest struct {
	TemplateID       string
	From             string
	FromName         string
	To               []*Email
	CC               []*Email
	BCC              []*Email
	Subject          string
	PlainTextContent string
	HTMLContent      string
}

type Email struct {
	Address string
	Name    string
}

func ParseEmail(input string) (*Email, error) {
	parts := strings.Split(input, " ")
	if len(parts) == 0 || parts[0] == "" {
		return nil, errors.New("email input is empty")
	}

	var address, name string

	if len(parts) == 1 {
		address = parts[0]
	} else {
		address = parts[len(parts)-1]
		name = strings.Join(parts[:len(parts)-1], " ")
	}

	return &Email{
		Address: strings.TrimSuffix(strings.TrimPrefix(address, "<"), ">"),
		Name:    name,
	}, nil
}
