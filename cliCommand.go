package main

import (
	"strings"

	"golang.org/x/net/websocket"
)

func cliCommand(ws *websocket.Conn, comm string) (string, bool) {
	if ((len(strings.Split(comm, " ")) == 1) && (len(state.ModuleName) == 0)) || (comm == "quit") {
		availableCommand := map[string]interface{}{
			"quit": quitCli,
			"help": helpCli,
		}
		switch comm {
		case "quit":
			return availableCommand["quit"].(func() (string, bool))()
		default:
			state.Path = "/help"
			return availableCommand["help"].(func() (string, bool))()
		}
	} else if (state.Path == "/") || (state.Path == "/help") {
		commOptions := strings.SplitN(comm, " ", 2)
		availableCommand := map[string]interface{}{
			"use":  useCli,
			"help": helpCli,
		}
		switch commOptions[0] {
		case "use":
			for i := 0; i < len(state.AvailableModules); i++ {
				if state.AvailableModules[i].Name == commOptions[1] {
					state.Path = "/" + commOptions[1]
					state.MdlNumb = i
					state.ModuleName = commOptions[1]
					return availableCommand["use"].(func(*ProgrammStatus) string)(state), false
				}
			}
			return availableCommand["help"].(func() (string, bool))()
		default:
			return availableCommand["help"].(func() (string, bool))()
		}
	} else {
		// commands for module
		return "here", false

	}
}
