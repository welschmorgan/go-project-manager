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
	}
	prompt.Validate = func(s string) error {
		var err error
		for _, v := range validators {
			if v != nil {
				if err = v(s); err != nil {
					return err
				}
			}
		}
		return nil
	}
	return prompt.Run()
}
