// Package goext adds various extensions for the go language.
package goext

import (
	"fmt"
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
