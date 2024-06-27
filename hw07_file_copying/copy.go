package main

import (
	"errors"
	"fmt"
	"io"
	"os"
	"strings"

	progressbar "github.com/schollz/progressbar/v3"
)

var (
	ErrParamValueIsEmpty     = errors.New("value is empty")
	ErrUnsupportedFile       = errors.New("unsupported file")
	ErrOffsetExceedsFileSize = errors.New("offset exceeds file size")
)

func Copy(fromPath, toPath string, offset, limit int64) error {
	fromPath = strings.TrimSpace(fromPath)
	toPath = strings.TrimSpace(toPath)

	if err := validateFilePathParam(fromPath, "fromPath"); err != nil {
		return err
	}
	if err := validateFilePathParam(toPath, "toPath"); err != nil {
		return err
	}

	fi, err := os.Stat(fromPath)
	if err != nil {
		return err
	}

	if offset > fi.Size() {
		return ErrOffsetExceedsFileSize
	}

	srcFile, err := os.Open(fromPath)
	if err != nil {
		return err
	}
	defer func() {
		if e := srcFile.Close(); e != nil && err == nil {
			err = e
		}
	}()

	_, err = srcFile.Seek(offset, 0)
	if err != nil {
		return err
	}

	dstFile, err := os.Create(toPath)
	if err != nil {
		return err
	}
	defer func() {
		if e := dstFile.Close(); e != nil && err == nil {
			err = e
		}
	}()

	pbar := initProgressBar(fi.Size(), offset, limit)
	var chunkSize int64 = 512
	var totalCopied int64
	for {
		if limit > 0 {
			if totalCopied+chunkSize > limit {
				chunkSize = limit - totalCopied
			}
		}

		copied, ce := io.CopyN(dstFile, srcFile, chunkSize)
		if ce != nil && !errors.Is(ce, io.EOF) {
			return ce
		}

		if copied > 0 {
			totalCopied += copied
			pbar.Add(int(copied))
		}

		if limit > 0 && totalCopied >= limit || errors.Is(ce, io.EOF) {
			break
		}
	}

	return err
}

func initProgressBar(fileSize, offset, limit int64) *progressbar.ProgressBar {
	bytesToCopy := fileSize - offset
	if limit > 0 && limit < bytesToCopy {
		bytesToCopy = limit
	}

	return progressbar.Default(bytesToCopy)
}

func validateFilePathParam(paramValue string, paramName string) error {
	if paramValue == "" {
		return fmt.Errorf("%s : %w", paramName, ErrParamValueIsEmpty)
	}

	if paramValue[0] != '/' {
		paramValue = "/" + paramValue
	}

	if strings.HasPrefix(paramValue, "/dev") {
		return ErrUnsupportedFile
	}

	if strings.HasPrefix(paramValue, "/proc") {
		return ErrUnsupportedFile
	}

	return nil
}
