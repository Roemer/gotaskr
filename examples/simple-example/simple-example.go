package main

import (
	"fmt"
	"os"

	"github.com/roemer/gotaskr"
)

func init() {
	gotaskr.Task("Hello", func() error {
		name := gotaskr.GetArgument("name", "Wulfgang")
		fmt.Println("Hi ", name)
		return nil
	}).Description("Just a hello task")
}

func main() {
	os.Exit(gotaskr.Execute())
}
