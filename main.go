package main

import (
	"fmt"
	"os"

	"github.com/teemuteemu/batman/cmd"
)

var version = "unset"

func main() {
	err := cmd.Execute(version)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
