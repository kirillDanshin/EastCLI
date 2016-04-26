package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strings"

	"golang.org/x/net/websocket"
)

func runModule(ws *websocket.Conn) string {
	var rcv []byte
	state.RunModule.Command = "exploit"
	state.RunModule.Args.ModuleName = state.ModuleName
	fmt.Printf("%s", "Use listener for this module???(y,n)\n")
	reader := bufio.NewReader(os.Stdin)
	comm, _ := reader.ReadString('\n')
	re := regexp.MustCompile(`\r?\n`)
	comm = strings.ToLower(re.ReplaceAllString(comm, ""))
	switch comm {
	case "y":
		state.RunModule.Args.Listener = true
	case "n":
		state.RunModule.Args.Listener = false
	default:
		fmt.Printf("%s", "Using default value false\n")
		state.RunModule.Args.Listener = false
	}
	state.RunModule.Args.Option = state.MdlOptions
	send, err := json.Marshal(state.RunModule)
	if err != nil {
		fmt.Printf("%s\n", err)
	}
	rcv = sendRcvMessage(ws, send)
	fmt.Printf("%s", string(rcv))
	websocket.Message.Receive(ws, &rcv)
	fmt.Printf("%s", string(rcv))
	return "hello world " + state.ModuleName
}
