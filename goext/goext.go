// Package goext adds various extensions for the go language.
package goext

import (
	"fmt"
	"os"
	"strconv"
)

// AddIf adds the given values to an array if the condition is fulfilled.
func AddIf[T any](args []T, cond bool, values ...T) []T {
	if cond {
		args = append(args, values...)
	}
	return args
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
	n1, err1 := fmt.Printf(format, a...)
	n2, err2 := fmt.Println()
	n = n1 + n2
	err = nil
	if err1 != nil {
		err = err1
	}
	if err2 != nil {
		err = err2
	}
	return
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
		return fmt.Errorf("Cannot get current directory: %v", err)
	}
	// Make sure to reset to the previous folder
	defer func() {
		err = os.Chdir(pwd)
		if err != nil {
			err = fmt.Errorf("Cannot change back to directory %s: %v", strconv.Quote(pwd), err)
		}
	}()
	// Change the path
	err = os.Chdir(path)
	if err != nil {
		return fmt.Errorf("Cannot change to directory %s: %v", strconv.Quote(path), err)
	}
	// Execute the function
	err = f()
	if err != nil {
		return fmt.Errorf("Inner method failed: %v", err)
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
