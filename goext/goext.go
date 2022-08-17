package goext

import (
	"fmt"
	"os"
	"strconv"
)

func Ternary[T any](cond bool, vtrue, vfalse T) T {
	if cond {
		return vtrue
	}
	return vfalse
}

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

func RemoveEmpty(s []string) []string {
	var r []string
	for _, str := range s {
		if str != "" {
			r = append(r, str)
		}
	}
	return r
}

func RunInDirectory(path string, f func() error) error {
	pwd, err := os.Getwd()
	if err != nil {
		return fmt.Errorf("Cannot get current directory: %v", err)
	}
	err = os.Chdir(path)
	if err != nil {
		return fmt.Errorf("Cannot change to directory %s: %v", strconv.Quote(path), err)
	}
	err = f()
	if err != nil {
		return fmt.Errorf("Inner method failed: %v", err)
	}
	err = os.Chdir(pwd)
	if err != nil {
		return fmt.Errorf("Cannot change back to directory %s: %v", strconv.Quote(pwd), err)
	}
	return nil
}
