package ui

import "github.com/manifoldco/promptui"

func Select(label string, items []string) (string, error) {

	prompt := promptui.Select{
		Label: label,
		Items: items,
	}
	_, result, err := prompt.Run()
	return result, err
}
