package nlog

import (
	"fmt"

	"github.com/TwiN/go-color"
)

//func Debug(msg string, args ...interface{}) {
//	print(msg, args)
//}

func Debugln(msg string) {
	ConsoleC(color.Gray, msg+"\n")
}

func Console(msg string) {
	fmt.Printf(msg)
}

func ConsoleC(clr string, msg string) {
	Console(color.Colorize(clr, msg))
}
