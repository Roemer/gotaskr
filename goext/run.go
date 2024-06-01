package goext

import (
	"errors"
	"fmt"
)

// An option that can be used to run which is applied before the run and reverted after.
type RunOption interface {
	// The method that is applied before the run.
	apply() error
	// The method that is applied after the run.
	revert() error
}

// Runs a given method with additional options.
func RunWithOptions(f func() error, options ...RunOption) (err error) {
	// Apply the options
	for _, option := range options {
		err = errors.Join(err, option.apply())
	}
	// Make sure to revert all options, in reverse order
	defer func() {
		for index := len(options) - 1; index >= 0; index-- {
			option := options[index]
			err = errors.Join(err, option.revert())
		}
	}()
	// Execute the function
	methodErr := f()
	if methodErr != nil {
		err = errors.Join(err, fmt.Errorf("inner method failed: %v", methodErr))
	}
	return
}

// Runs a given method with returns one parameter with additional options.
func RunWithOptions1P[P1 any](f func() (P1, error), options ...RunOption) (P1, error) {
	var p1 P1
	return p1, RunWithOptions(func() error {
		var err error
		p1, err = f()
		return err
	}, options...)
}

// Runs a given method with returns two parameters with additional options.
func RunWithOptions2P[P1 any, P2 any](f func() (P1, P2, error), options ...RunOption) (P1, P2, error) {
	var p1 P1
	var p2 P2
	return p1, p2, RunWithOptions(func() error {
		var err error
		p1, p2, err = f()
		return err
	}, options...)
}

// Runs a given method with returns three parameters with additional options.
func RunWithOptions3P[P1 any, P2 any, P3 any](f func() (P1, P2, P3, error), options ...RunOption) (P1, P2, P3, error) {
	var p1 P1
	var p2 P2
	var p3 P3
	return p1, p2, p3, RunWithOptions(func() error {
		var err error
		p1, p2, p3, err = f()
		return err
	}, options...)
}
