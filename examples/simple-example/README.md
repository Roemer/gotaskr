# simple-example
This is a simple introduction sample.

The following commands can be used to execute the sample:

Without any target, a list of all targets is listed.
```
go run .
```

The target `Hello` is called without any arguments, so the default is used.
```
go run . --target Hello
```

The target `Hello` is called with the name argument, so this value is used.
```
go run . --target Hello --name Dietmar
```
