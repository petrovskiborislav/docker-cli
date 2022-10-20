package prompt

import (
	"gopkg.in/AlecAivazis/survey.v1"
)

// Prompt is an interface for a prompt.
type Prompt interface {
	SelectPrompt(label string, items []string, opts ...survey.AskOpt) ([]string, error)
}

// NewPrompt creates a new prompt.
func NewPrompt() Prompt {
	return &prompt{}
}

type prompt struct{}

// SelectPrompt creates a prompt which allows the user to select multiple options.
func (p prompt) SelectPrompt(label string, items []string, opts ...survey.AskOpt) ([]string, error) {
	var result []string
	prompt := &survey.MultiSelect{
		Message:  label,
		Options:  items,
		PageSize: len(items),
	}

	err := survey.AskOne(prompt, &result, nil, opts...)
	if err != nil {
		return nil, err
	}

	return result, nil
}
