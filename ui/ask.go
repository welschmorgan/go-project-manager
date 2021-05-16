package ui

import "github.com/manifoldco/promptui"

func Ask(label, defaultValue string, validators ...Validator) (string, error) {
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
		Default:   defaultValue,
		Validate:  NewMultiValidator(validators...),
	}
	return prompt.Run()
}
