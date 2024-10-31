package helpers

import (
	"bytes"
	"errors"
	"fmt"
	"math"
	"os"
	"path"
	"path/filepath"
	"runtime/debug"
	"time"
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

// Converts duration to minutes and seconds with milliseconds mm:ss.sss
func ConvertTimeToHuman(totalTime time.Duration) string {
	lapTimeTotalMinutes := int(math.Trunc(totalTime.Minutes()))
	lapTimeTotalSeconds := totalTime.Seconds()
	lapTimeSeconds := lapTimeTotalSeconds - (float64(lapTimeTotalMinutes) * 60)
	return fmt.Sprintf("%d:%06.3f", lapTimeTotalMinutes, lapTimeSeconds)
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

func FindInFolder(root, pattern string) ([]string, error) {
	var matches []string
	err := filepath.Walk(root, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if info.IsDir() {
			return nil
		}
		if matched, err := filepath.Match(pattern, filepath.Base(path)); err != nil {
			return err
		} else if matched {
			matches = append(matches, path)
		}
		return nil
	})
	if err != nil {
		return nil, err
	}
	return matches, nil
}

func CheckFileExists(filePath string) bool {
	_, error := os.Stat(filePath)
	//return !os.IsNotExist(err)
	return !errors.Is(error, os.ErrNotExist)
}
