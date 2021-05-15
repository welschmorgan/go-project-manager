package ui

import (
	"github.com/manifoldco/promptui"
)

func Ask(label, defaultValue string, validator func(string) error) (string, error) {
	templates := &promptui.PromptTemplates{
		Prompt:  "{{ . }} ",
		Valid:   "{{ . | green }} ",
		Invalid: "{{ . | red }} ",
		Success: "{{ . | bold }} ",
	}

	prompt := promptui.Prompt{
		Label:     label,
		Templates: templates,
		Default:   defaultValue,
	}
	if validator != nil {
		prompt.Validate = validator
	}

	return prompt.Run()
}

func Select(label string, items []string, validator func(string) error) (string, error) {

	prompt := promptui.Select{
		Label: label,
		Items: items,
	}

	_, result, err := prompt.Run()
	return result, err
}
