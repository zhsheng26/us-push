package udp

import (
	"os"
	"os/exec"
	"runtime"
)

/**
* @author zhangsheng
* @date 2019/7/5
 */

type MessageType int

const (
	// FUNC for functionnal messages ie technical messages from the client to the server
	FUNC MessageType = iota
	// CLASSIQUE message for messages sent by the end user
	CLASSIQUE
)

// ConnectionStatus is self explained
type ConnectionStatus int

const (
	JOINING ConnectionStatus = iota
	LEAVING
)

func CallClear() {
	clear := make(map[string]func()) //Initialize it
	clear["linux"] = func() {
		cmd := exec.Command("clear") //Linux example, its tested
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	}
	clear["windows"] = func() {
		cmd := exec.Command("cmd", "/c", "cls") //Windows example, its tested
		cmd.Stdout = os.Stdout
		_ = cmd.Run()
	}
	value, ok := clear[runtime.GOOS] //runtime.GOOS -> linux, windows, darwin etc.
	if ok {                          //if we defined a clear func for that platform:
		value() //we execute it
	} else { //unsupported platform
		panic("Your platform is unsupported! I can't clear terminal screen :(")
	}
}
