package log

import "fmt"

var Newline string = fmt.Sprintln()

var Information = logWrite
var Informationf = logWritef

var Debug = logNoop
var Debugf = logNoopf

func Initialize(verbose bool) {
	if verbose {
		Debug = func(a ...any) (n int, err error) {
			return logWrite(a...)
		}
		Debugf = func(format string, a ...any) (n int, err error) {
			return logWritef(format, a...)
		}
	}
}

func logNoop(a ...any) (n int, err error) {
	return 0, nil
}

func logNoopf(format string, a ...any) (n int, err error) {
	return 0, nil
}

func logWrite(a ...any) (n int, err error) {
	return fmt.Println(a...)
}

func logWritef(format string, a ...any) (n int, err error) {
	text := fmt.Sprintf(format, a...)
	return logWrite(text)
}
