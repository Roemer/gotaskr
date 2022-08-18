# gotaskr
[![Go Reference](https://pkg.go.dev/badge/github.com/roemer/gotaskr.svg)](https://pkg.go.dev/github.com/roemer/gotaskr)
[![Docs github.io](https://img.shields.io/badge/Docs-github.io-blue.svg)](https://roemer.github.io/gotaskr/)
![GitHub](https://img.shields.io/github/license/roemer/gotaskr)
![GitHub tag (latest SemVer)](https://img.shields.io/github/v/tag/roemer/gotaskr)

## Introduction
gotaskr (Go-Task-Runner) is a generic task runner which is invoked via CLI.

The tasks are written in plan Go and can easily be called from the CLI.
This is especially usefull for tasks in the CI.

The tasks can be chained and in the end, there is a small statistic about
the tasks that were executed and their runtime.

There are some inbuilt helpers for often used tasks for various DevOps tasks.

## Features
- Compileable or directly runnable with go run
- Fanzy statistics after the execution
- Output from subprocesses directly visible
- Tasks are written in plain Go
- Each task can have custom arguments
- Chainable tasks
- Inbuilt helpers for various DevOps tasks
- Small footprint and easily extendable

https://proxy.golang.org/github.com/roemer/gotaskr/@v/v0.0.1.info

## Quick-Start
Create a new go project:
```
go mod init my-project
```

Add gotaskr:
```
go get github.com/roemer/gotaskr
```

Create a `build.go` file:
```go
package main

import (
	"os"

	"github.com/roemer/gotaskr"
)

func main() {
	os.Exit(gotaskr.Execute())
}
```
Now you need to register your tasks you want to be able to run:
```go
func init() {
	gotaskr.Task("My-Task", func() error {
		fmt.Println("Hello from My-Task")
		return nil
	})
}
```
Now invoke this task by running:
```
go run . --target My-Task
```

## Examples
Have a look at the [examples](examples) from this repository.
