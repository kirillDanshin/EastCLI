package main

import (
	"encoding/json"
	"fmt"
	"strings"

	"golang.org/x/net/websocket"
)

func setOption(command string) string {
	switch {
	case len(strings.Split(command, " ")) == 2:
		opt := strings.Split(command, " ")
		badOption := "Option: " + opt[0] + " not found\nAvailable options for " + state.ModuleName
		for i := 0; i < len(state.MdlOptions); i++ {
			if opt[0] == state.MdlOptions[i].Name {
				state.MdlOptions[i].Value.Value = opt[1]
				return state.MdlOptions[i].Name + "=>" + opt[1]
			}
			badOption += state.MdlOptions[i].Name + "\n"
		}
		return badOption
	default:
		return "Using set command 'set ModuleOption OptionValue\n"
	}
}

func showOptCli(ws *websocket.Conn) {
	if len(state.MdlOptions) == 0 {
		var rcv []byte
		var ParseOptions ShowOptions
		message := `{"command": "options", "args": {"module_name": "` + state.ModuleName + `"}}`
		websocket.Message.Send(ws, []byte(message))
		websocket.Message.Receive(ws, &rcv)
		state.MdlOptions = parseOpt(rcv, &ParseOptions)
	}
	for i := 0; i < len(state.MdlOptions); i++ {
		fmt.Printf("Name: %s, Type: %s, Value: %s\n", state.MdlOptions[i].Name,
			state.MdlOptions[i].Value.Type, state.MdlOptions[i].Value.Value)
	}
}

func parseOpt(rcv []byte, ParseOptions *ShowOptions) []Option {
	rcvString := string(rcv)
	json.Unmarshal(rcv, &ParseOptions)
	for i := 0; i < len(ParseOptions.Args); i++ {
		if ParseOptions.Args[i].Value.Type == "list" {
			ParseOptions.Args[i].Value.Value = strings.Split(strings.Split(rcvString, `"options": [`)[1], "]")[0]
		}
	}
	return ParseOptions.Args
}
