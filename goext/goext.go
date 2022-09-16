// Package goext adds various extensions for the go language.
package goext

import (
	"encoding/json"
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

// AddIf adds the given values to a slice if the condition is fulfilled.
func AddIf[T any](slice []T, cond bool, values ...T) []T {
	if cond {
		return append(slice, values...)
	}
	return slice
}

// AppendIfMissing appends the given element if it is missing in the slice.
func AppendIfMissing[T comparable](slice []T, value T) []T {
	for _, ele := range slice {
		if ele == value {
			return slice
		}
	}
	return append(slice, value)
}

// Ternary adds the missing ternary operator.
func Ternary[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}

// Printfln allows to use Printf and Println in one call.
func Printfln(format string, a ...any) (n int, err error) {
	text := fmt.Sprintf(format, a...)
	return fmt.Println(text)
}

// RemoveEmpty removes all empty strings from an array.
func RemoveEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

// RunInDirectory runs a given function inside the passed directory as working directory.
// It resets to the previous directory when finished (or an error occured).
func RunInDirectory(path string, f func() error) (err error) {
	// Get the current directory
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("cannot get current directory: %v", err)
	}
	// Make sure to reset to the previous folder
	defer func() {
		// Set the error only if the main error is not yet set
		if tempErr := os.Chdir(pwd); tempErr != nil && err == nil {
			err = fmt.Errorf("cannot change back to directory %s: %v", strconv.Quote(pwd), tempErr)
		}
	}()
	// Change the path
	err = os.Chdir(path)
	if err != nil {
		return fmt.Errorf("cannot change to directory %s: %v", strconv.Quote(path), err)
	}
	// Execute the function
	err = f()
	if err != nil {
		return fmt.Errorf("inner method failed: %v", err)
	}
	return
}

// PanicOnError panics if the given error is set.
func PanicOnError(err error) {
	if err != nil {
		panic(err)
	}
}

// Noop is a no operation function that can be used for tasks.
// For example if they are just used for multiple dependencies.
func Noop() error {
	return nil
}

// WriteJsonToFile writes the given object into a file.
func WriteJsonToFile(object any, outputFilePath string, indented bool) error {
	var data []byte
	var err error
	if indented {
		data, err = json.MarshalIndent(object, "", "  ")
	} else {
		data, err = json.Marshal(object)
	}
	if err != nil {
		return err
	}
	if err := os.WriteFile(outputFilePath, data, os.ModePerm); err != nil {
		return err
	}
	return nil
}

// EnvExists checks if the given environment variable exists or not.
func EnvExists(key string) bool {
	if _, ok := os.LookupEnv(key); ok {
		return true
	}
	return false
}

// GetEnvOrDefault returns the value if the environment variable exists or the default otherwise.
func GetEnvOrDefault(key string, defaultValue string) (string, bool) {
	if _, ok := os.LookupEnv(key); ok {
		return key, true
	}
	return defaultValue, false
}

// FileExists checks if a file exists (and it is not a directory).
func FileExists(filePath string) (bool, error) {
	info, err := os.Stat(filePath)
	if err == nil {
		return !info.IsDir(), nil
	}
	if errors.Is(err, os.ErrNotExist) {
		return false, nil
	}
	return false, err
}

// Pass is a no-op function that can be used to set variables to used.
// Usefull during development but must be removed afterwards!
func Pass(i ...interface{}) {
	// No-Op
}

func TrimAllPrefix(s, prefix string) string {
	for strings.HasPrefix(s, prefix) {
		s = s[len(prefix):]
	}
	return s
}

func TrimAllSuffix(s, suffix string) string {
	for strings.HasSuffix(s, suffix) {
		s = s[:len(s)-len(suffix)]
	}
	return s
}

func TrimNewlineSuffix(v string) string {
	return TrimAllSuffix(TrimAllSuffix(v, "\r\n"), "\n")
}

// SplitByNewLine splits the given value by newlines.
func SplitByNewLine(value string) []string {
	return strings.Split(strings.ReplaceAll(value, "\r\n", "\n"), "\n")
}
