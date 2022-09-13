package log

import "fmt"

var Newline string = fmt.Sprintln()

var Information = logWrite
var Informationf = logWritef

var Debug = logNoop
var Debugf = logNoopf

func Initialize(verbose bool) {
	if verbose {
		Debug = func(a ...any) int {
			return logWrite(a...)
		}
		Debugf = func(format string, a ...any) int {
			return logWritef(format, a...)
		}
	}
}

func logNoop(a ...any) int {
	return 0
}

func logNoopf(format string, a ...any) int {
	return 0
}

func logWrite(a ...any) int {
	n, _ := fmt.Println(a...)
	return n
}

func logWritef(format string, a ...any) int {
	text := fmt.Sprintf(format, a...)
	return logWrite(text)
}
