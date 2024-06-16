package main

import (
	"errors"
	"log"
	"os"
)

func main() {
	if err := validateArgs(os.Args); err != nil {
		log.Fatal(err)
	}

	env, err := ReadDir(os.Args[1])
	if err != nil {
		log.Fatal(err)
	}

	returnCode := RunCmd(os.Args[2:], env)
	os.Exit(returnCode)
}

func validateArgs(args []string) error {
	if len(args) < 3 {
		return errors.New("directory and command to execute should be specified")
	}

	return nil
}
