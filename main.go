package main

import (
	"github.com/welschmorgan/go-project-manager/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
