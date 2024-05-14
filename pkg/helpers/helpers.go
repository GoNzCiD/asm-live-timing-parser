package helpers

import (
	"bytes"
	"os"
	"path"
	"path/filepath"
	"runtime/debug"
)

func SaveToFile(content string, filepath string) error {
	f, err := os.Create(filepath)
	if err != nil {
		return err
	}
	defer f.Close()

	_, err = f.WriteString(content)
	if err != nil {
		return err
	}

	return nil
}

func ReadFromFile(filepath string) (string, error) {
	b, err := os.ReadFile(filepath)
	if err != nil {
		return "", err
	}

	return string(b), nil
}

func GetFullPath(relPath string) (string, error) {
	if filepath.IsAbs(relPath) {
		return relPath, nil
	} else {
		basePath, err := filepath.Abs(relPath)
		if err != nil {
			return "", err
		}

		absPath, err := filepath.Abs(path.Join(path.Dir(basePath), relPath))
		if err != nil {
			return "", err
		}
		return absPath, nil
	}
}

// StackTraceAsList returns the stack trace of the calling goroutine
// as list of string. At most 50 lines are returned.
func StackTraceAsList() []string {
	var rv []string
	st := bytes.SplitN(debug.Stack(), []byte{'\n'}, 50)
	if len(st) == 50 {
		st = st[:50]
	}
	for _, line := range st {
		rv = append(rv, string(line))
	}
	return rv
}
