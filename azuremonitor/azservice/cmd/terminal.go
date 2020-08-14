package cmd

import (
	"os"
	"os/exec"
)

func init() {

	setTerminalEnv()
}

var clear map[string]func() //create a map for storing clear funcs

func setTerminalEnv() {
	clear = make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		cmd.Run()
	}
}
