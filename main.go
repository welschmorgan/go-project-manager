package main

import (
	"github.com/welschmorgan/go-release-manager/cmd"
)

func main() {
	if err := cmd.Execute(); err != nil {
		panic(err)
	}
}
