package ui

import (
	"fmt"

	"github.com/manifoldco/promptui"
)

func Ask(label, defaultValue interface{}, validators ...Validator) (string, error) {
	// templates := &promptui.PromptTemplates{
	// 	Prompt:  "{{ . }} ",
	// 	Valid:   "{{ . | green }} ",
	// 	Invalid: "{{ . | red }} ",
	// 	Success: "{{ . | bold }} ",
	// }

	prompt := promptui.Prompt{
		Label: label,
		// Templates: templates,
		AllowEdit: true,
		Default:   fmt.Sprintf("%v", defaultValue),
		Validate:  NewMultiValidator(validators...),
	}
	return prompt.Run()
}
