package main

import (
	"bytes"
	"os"
	"path"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const (
	testDir = "testdata"
)

func TestRunCmd(t *testing.T) {
	testScript := "executor_exit_test.sh"
	for expectedReturnCode := range []int{
		0,   // Success
		1,   // Error with exit code 1
		200, // Error with exit code 200
	} {
		t.Run("Return executed script error code", func(t *testing.T) {
			script := path.Join(testDir, testScript)
			cmd := []string{"/bin/sh", script, strconv.Itoa(expectedReturnCode)}
			env := make(Environment)
			actual := RunCmd(cmd, env)

			assert.Equal(t, expectedReturnCode, actual)
		})
	}

	testScript = "executor_test.sh"
	envVariableName := "OTUS_HW08_ENVDIR_TOOL_EXECUTOR_TEST"
	for _, test := range []struct {
		Name           string
		EnvVariableVal string
		NeedRemove     bool
		Expected       string
	}{
		{
			Name:           "Set value",
			EnvVariableVal: "Hello",
			NeedRemove:     false,
			Expected:       "Hello",
		},
		{
			Name:           "Set empty value",
			EnvVariableVal: "",
			NeedRemove:     false,
			Expected:       "",
		},
		{
			Name:           "Unset value",
			EnvVariableVal: "Hello",
			NeedRemove:     true,
			Expected:       "",
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			err := os.Setenv(envVariableName, "SHOULD_REPLACE")
			require.NoError(t, err)

			env := make(Environment)
			env[envVariableName] = EnvValue{
				Value:      test.EnvVariableVal,
				NeedRemove: test.NeedRemove,
			}

			script := path.Join(testDir, testScript)
			cmd := []string{"/bin/bash", script, envVariableName}

			b := &bytes.Buffer{}
			RunCmdWithWriters(cmd, env, b, b)

			actual := strings.TrimRight(b.String(), "\n")
			assert.Equal(t, test.Expected, actual)

			os.Unsetenv(envVariableName)
		})
	}
}
