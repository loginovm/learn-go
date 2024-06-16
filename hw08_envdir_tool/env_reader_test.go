package main

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestReadDir(t *testing.T) {
	envFileNameTemplate := "evar*"
	for _, test := range []struct {
		Name               string
		FileContents       string
		ExpectedVal        string
		NeedRemoveVariable bool
	}{
		{
			Name:               "Simple value",
			FileContents:       "Hello",
			ExpectedVal:        "Hello",
			NeedRemoveVariable: false,
		},
		{
			Name: "Multiline value",
			FileContents: `bar
PLEASE IGNORE SECOND LINE`,
			ExpectedVal:        "bar",
			NeedRemoveVariable: false,
		},
		{
			Name:         "Value with 0x00",
			FileContents: "foo\u0000with new line",
			ExpectedVal: `foo
with new line`,
			NeedRemoveVariable: false,
		},
		{
			Name:               "Empty value",
			FileContents:       "",
			ExpectedVal:        "",
			NeedRemoveVariable: true,
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			dir, err := os.MkdirTemp("", "Env*")
			require.NoError(t, err)

			f, err := os.CreateTemp(dir, envFileNameTemplate)
			require.NoError(t, err)
			envVarFilePath := f.Name()

			err = os.WriteFile(envVarFilePath, []byte(test.FileContents), 0o644)
			require.NoError(t, err)

			env, err := ReadDir(dir)
			require.NoError(t, err)
			assert.Len(t, env, 1)

			varName := filepath.Base(envVarFilePath)
			assert.Contains(t, env, varName)

			v := env[varName]
			assert.Equal(t, test.ExpectedVal, v.Value)
			assert.Equal(t, test.NeedRemoveVariable, v.NeedRemove)

			err = os.RemoveAll(dir)
			require.NoError(t, err)
		})
	}
}
