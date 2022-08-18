package main

import (
	"fmt"
	"os"

	"github.com/roemer/gotaskr"
	"github.com/roemer/gotaskr/goext"
)

var varFunc = func() error {
	fmt.Println("Hi from varFunc")
	return nil
}

func RegularFunction() error {
	fmt.Println("Hi from RegularFunction")
	return nil
}

func init() {
	gotaskr.Task("Inline-Register", func() error {
		fmt.Println("Hi from inline")
		return nil
	})
	gotaskr.Task("Var-Register", varFunc)
	gotaskr.Task("Regular-Register", RegularFunction)
	gotaskr.Task("All", goext.Noop).DependsOn("Inline-Register", "Var-Register", "Regular-Register")
}

func main() {
	os.Exit(gotaskr.Execute())
}
