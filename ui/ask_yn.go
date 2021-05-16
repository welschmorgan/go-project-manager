package ui

import (
	"regexp"

	"github.com/manifoldco/promptui"
)

func AskYN(label string, validators ...Validator) (bool, error) {
	prompt := promptui.Prompt{
		Label:     label,
		IsConfirm: true,
		Validate:  NewMultiValidator(validators...),
	}
	res, _ := prompt.Run()
	yesRule := regexp.MustCompile("(?i)y(?:es)?|1")
	return yesRule.MatchString(res), nil
}
