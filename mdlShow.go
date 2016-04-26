package main

import "golang.org/x/net/websocket"

func mdlShow(ws *websocket.Conn, comm string) string {
	switch {
	case comm == "options":
		showOptCli(ws)
		return ""
	case comm == "description":
		return state.AvailableModules[state.MdlNumb].Description
	case comm == "cve":
		return state.AvailableModules[state.MdlNumb].CVEName
	case comm == "vendor":
		return state.AvailableModules[state.MdlNumb].Vendor
	case comm == "notes":
		return state.AvailableModules[state.MdlNumb].Notes
	case comm == "links":
		links := ""
		for i := 0; i < len(state.AvailableModules[state.MdlNumb].Links); i++ {
			links += state.AvailableModules[state.MdlNumb].Links[i] + "\n"
		}
		return links
	default:
		return "Use show [options, description, cve, vendor, notes, links]\n"
	}
}
