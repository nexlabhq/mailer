package email

import "github.com/nexlabhq/mailer/common"

func NewEmails(address string, name string) []*common.Email {
	return []*common.Email{
		{
			Address: address,
			Name:    name,
		},
	}
}

func NewEmail(address string, name string) *common.Email {
	return &common.Email{
		Address: address,
		Name:    name,
	}
}
