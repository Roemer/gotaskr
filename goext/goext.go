// Package goext adds various extensions for the go language.
package goext

import (
	"fmt"
	"os"
	"strconv"
)

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

func RunInDirectory1P[P1 any](path string, f func() (P1, error)) (P1, error) {
	var p1 P1
	return p1, RunInDirectory(path, func() error {
		var err error
		p1, err = f()
		return err
	})
}

func RunInDirectory2P[P1 any, P2 any](path string, f func() (P1, P2, error)) (P1, P2, error) {
	var p1 P1
	var p2 P2
	return p1, p2, RunInDirectory(path, func() error {
		var err error
		p1, p2, err = f()
		return err
	})
}

func RunInDirectory3P[P1 any, P2 any, P3 any](path string, f func() (P1, P2, P3, error)) (P1, P2, P3, error) {
	var p1 P1
	var p2 P2
	var p3 P3
	return p1, p2, p3, RunInDirectory(path, func() error {
		var err error
		p1, p2, p3, err = f()
		return err
	})
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

// Pass is a no-op function that can be used to set variables to used.
// Usefull during development but must be removed afterwards!
func Pass(i ...interface{}) {
	// No-Op
}
