package mailer

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestNotificationTemplateParser(t *testing.T) {
	templates := []struct {
		Input     EmailTemplate
		Variables interface{}
		Output    EmailTemplate
	}{
		{
			EmailTemplate{
				ID: "test_template",
				Subjects: map[string]string{
					"en": "Test subject en",
					"vi": "Test subject vi",
				},
				Contents: map[string]string{
					"en": "Test contents en",
					"vi": "Test contents vi",
				},
				HTMLContents: map[string]string{
					"en": "<p>Test contents en<p>",
					"vi": "<p>Test contents vi<p>",
				},
			},
			nil,
			EmailTemplate{
				ID: "test_template",
				Subjects: map[string]string{
					"en": "Test subject en",
					"vi": "Test subject vi",
				},
				Contents: map[string]string{
					"en": "Test contents en",
					"vi": "Test contents vi",
				},
				HTMLContents: map[string]string{
					"en": "<p>Test contents en<p>",
					"vi": "<p>Test contents vi<p>",
				},
			},
		},
		{
			EmailTemplate{
				ID: "test_template_1",
				Subjects: map[string]string{
					"en": "Test subject {{.Foo}}",
					"vi": "Test subject {{.Bar}}",
				},
				Contents: map[string]string{
					"en": "Test contents {{.Foo}}",
					"vi": "Test contents {{.Bar}}",
				},
				HTMLContents: map[string]string{
					"en": "<p>Test contents {{.Foo}}<p>",
					"vi": "<p>Test contents {{.Bar}}<p>",
				},
			},
			struct {
				Foo string
				Bar string
			}{
				"foo",
				"bar",
			},
			EmailTemplate{
				ID: "test_template_1",
				Subjects: map[string]string{
					"en": "Test subject foo",
					"vi": "Test subject bar",
				},
				Contents: map[string]string{
					"en": "Test contents foo",
					"vi": "Test contents bar",
				},
				HTMLContents: map[string]string{
					"en": "<p>Test contents foo<p>",
					"vi": "<p>Test contents bar<p>",
				},
			},
		},
		{
			EmailTemplate{
				ID: "test_template_1",
				Subjects: map[string]string{
					"en": "Test subject {{.Foo}}",
					"vi": "Test subject {{.Bar}}",
				},
				Contents: map[string]string{
					"en": "Test contents {{.Foo}}",
					"vi": "Test contents {{.Bar}}",
				},
			},
			struct {
				Foo string
				Bar string
			}{
				"foo",
				"bar",
			},
			EmailTemplate{
				ID: "test_template_1",
				Subjects: map[string]string{
					"en": "Test subject foo",
					"vi": "Test subject bar",
				},
				Contents: map[string]string{
					"en": "Test contents foo",
					"vi": "Test contents bar",
				},
			},
		},
	}

	for _, template := range templates {
		parsedTemplates, err := ParseTemplate(&template.Input, template.Variables)
		assert.NoError(t, err)

		assert.Equal(t, template.Output.ID, parsedTemplates.ID)
		assert.Equal(t, template.Output.Subjects, parsedTemplates.Subjects)
		assert.Equal(t, template.Output.Contents, parsedTemplates.Contents)
	}

}
