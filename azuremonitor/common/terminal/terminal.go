package terminal

import (
	"os"
	"os/exec"
	"runtime"
)

func init() {

	setTerminalEnv()
}

var m map[string]func() //create a map for storing clear funcs

func Clear() {
	value, ok := m[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("wrong platform")
	}
}

func setTerminalEnv() {
	m = make(map[string]func()) //Initialize it
	m["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	m["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}
