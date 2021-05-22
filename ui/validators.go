package ui

import (
	"errors"
	"os"
	"strings"
)

func StrMustBeNonEmpty(s string) error {
	if len(s) == 0 {
		return errors.New("value must be non-empty")
	}
	return nil
}

func StrMustNotContainOnlySpaces(s string) error {
	if len(strings.TrimSpace(s)) == 0 {
		return errors.New("value must not contain only spaces")
	}
	return nil
}

func PathMustExist(p string) error {
	if _, err := os.Stat(p); err != nil && os.IsNotExist(err) {
		return errors.New("path does not exist")
	}
	return nil
}

func PathMustBeDir(p string) error {
	if fi, err := os.Stat(p); err != nil {
		return err
	} else if !fi.IsDir() {
		return errors.New("path is not a directory")
	}
	return nil
}
