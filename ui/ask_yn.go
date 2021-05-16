package ui

import (
	"regexp"

	"github.com/manifoldco/promptui"
)

func AskYN(label string, validators ...Validator) (bool, error) {
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
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
	res, _ := prompt.Run()
	yesRule := regexp.MustCompile("(?i)y(?:es)?|1")
	return yesRule.MatchString(res), nil
}
