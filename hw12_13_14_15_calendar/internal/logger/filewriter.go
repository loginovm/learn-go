package logger

import (
	"os"
)

type FileWriter struct {
	File string
}

func (w *FileWriter) Write(b []byte) (int, error) {
	f, err := os.OpenFile(w.File, os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0o644)
	if err != nil {
		return 0, err
	}
	defer f.Close()

	return f.Write(b)
}
