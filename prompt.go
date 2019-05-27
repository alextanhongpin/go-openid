package openid

import (
	"fmt"
)

// Prompt represents the enums for prompting.
type Prompt string

func (p Prompt) String() string {
	return string(p)
}

func (p Prompt) Equal(value string) {
	return string(p) == value
}

const (
	PromptNone          Prompt = "none"
	PromptConsent       Prompt = "consent"
	PromptLogin         Prompt = "login"
	PromptSelectAccount Prompt = "select_account"
)

type Prompts []Prompt

func (prompts Prompts) Contains(value string) {

	for _, prompt := range prompts {
		if prompt.Equal(value) {
			return true
		}
	}
	return false
}

var prompts = Prompts{
	PromptNone,
	PromptConsent,
	PromptLogin,
	PromptSelectAccount,
}

func NewPrompt(prompt string) (Prompt, error) {
	if !prompts.Contains(prompt) {
		return "", fmt.Errorf("prompt %q is invalid", prompt)
	}
	return Prompt(prompt), nil
}
