package main

import (
	"math"
	"os"
	"path"
	"testing"

	"github.com/stretchr/testify/require"
)

const (
	testFilesDir = "testdata"
	inputFile    = "input.txt"
	outFile      = "out.txt"
)

func TestCopy(t *testing.T) {
	for _, test := range []struct {
		Name             string
		Offset           int64
		Limit            int64
		ExpectedFileName string
	}{
		{
			Name:             "Copy file with offset 0 limit 0",
			Offset:           0,
			Limit:            0,
			ExpectedFileName: "out_offset0_limit0.txt",
		},
		{
			Name:             "Copy file with offset 0 limit 1000",
			Offset:           0,
			Limit:            1000,
			ExpectedFileName: "out_offset0_limit1000.txt",
		},
		{
			Name:             "Copy file with offset 100 limit 1000",
			Offset:           100,
			Limit:            1000,
			ExpectedFileName: "out_offset100_limit1000.txt",
		},
		{
			Name:             "Limit exceeds file size",
			Offset:           0,
			Limit:            math.MaxInt64,
			ExpectedFileName: "out_offset0_limit0.txt",
		},
		{
			Name:             "Copy file with offset 0 limit 1",
			Offset:           0,
			Limit:            1,
			ExpectedFileName: "out_offset0_limit1.txt",
		},
		{
			Name:             "Copy file with offset to EOF limit 0",
			Offset:           -1,
			Limit:            0,
			ExpectedFileName: "out_offsetEOF_limit0.txt",
		},
	} {
		t.Run(test.Name, func(t *testing.T) {
			fromPath := path.Join(testFilesDir, inputFile)
			outputFile, err := os.CreateTemp("", "hw07_out*.txt")
			require.NoError(t, err)
			toPath := outputFile.Name()

			if test.Offset == -1 {
				fi, _ := os.Stat(fromPath)
				test.Offset = fi.Size()
			}
			err = Copy(fromPath, toPath, test.Offset, test.Limit)
			require.NoError(t, err)

			expectedFilePath := path.Join(testFilesDir, test.ExpectedFileName)
			expectedFileData, err := os.ReadFile(expectedFilePath)
			require.NoError(t, err)
			expected := string(expectedFileData)

			actualFileData, err := os.ReadFile(toPath)
			require.NoError(t, err)
			actual := string(actualFileData)

			require.Equal(t, expected, actual)

			os.Remove(toPath)
			require.NoError(t, err)
		})
	}

	t.Run("Offset exceeds file size", func(t *testing.T) {
		fromPath := path.Join(testFilesDir, inputFile)
		toPath := path.Join(testFilesDir, outFile)

		var offset int64 = math.MaxInt64
		err := Copy(fromPath, toPath, offset, 0)
		require.ErrorIs(t, err, ErrOffsetExceedsFileSize)

		_, e := os.Stat(toPath)
		require.ErrorIs(t, e, os.ErrNotExist)
	})
}
