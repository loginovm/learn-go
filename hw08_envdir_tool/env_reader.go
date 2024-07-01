package main

import (
	"bufio"
	"bytes"
	"errors"
	"os"
	"path"
	"strings"
)

type Environment map[string]EnvValue

// EnvValue helps to distinguish between empty files and files with the first empty line.
type EnvValue struct {
	Value      string
	NeedRemove bool
}

// ReadDir reads a specified directory and returns map of env variables.
// Variables represented as files where filename is name of variable, file first line is a value.
func ReadDir(dir string) (Environment, error) {
	entries, err := os.ReadDir(dir)
	if err != nil {
		return nil, err
	}

	env := Environment{}
	for _, e := range entries {
		if !e.IsDir() {
			fileName := e.Name()
			if err = validateFileName(fileName); err != nil {
				return nil, err
			}

			ev, err := readFile(path.Join(dir, fileName))
			if err != nil {
				return nil, err
			}

			env[fileName] = ev
		}
	}

	return env, nil
}

func validateFileName(fileName string) error {
	if strings.Contains(fileName, "=") {
		return errors.New(`"=" is not allowed in file fileName`)
	}

	return nil
}

func readFile(fileName string) (EnvValue, error) {
	file, err := os.Open(fileName)
	if err != nil {
		return EnvValue{}, err
	}
	defer func() {
		if er := file.Close(); er != nil && err == nil {
			err = er
		}
	}()

	fi, err := file.Stat()
	if err != nil {
		return EnvValue{}, err
	}
	if fi.Size() == 0 {
		return EnvValue{NeedRemove: true}, nil
	}

	scanner := bufio.NewScanner(file)
	scanner.Scan()
	firstLine := scanner.Text()
	firstLine = strings.TrimRight(firstLine, "\x20\t")
	b := bytes.ReplaceAll([]byte(firstLine), []byte("\x00"), []byte("\n"))

	return EnvValue{Value: string(b)}, err
}
