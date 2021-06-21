package main

import (
	"fmt"
	"os"

	"github.com/welschmorgan/go-release-manager/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "\033[1;31merror\033[0m: %s\n", err.Error())
	}
}
