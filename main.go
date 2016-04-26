package main

import (
	"bufio"
	"bytes"
	"encoding/json"
	"fmt"
	"golang.org/x/net/websocket"
	"log"
	"os"
	"regexp"
	"strings"
)

var origin = "http://localhost/"
var url = "ws://localhost:49999/"

func main() {
	var state ProgrammStatus
	var rcvMessage []byte
	var helloMessage StartMessage
	helloMessage.Hello = `{"hello": {"name": "EastCLI", "type": "cli"}}`
	helloMessage.GetInfo = `{"command": "get_all_server_data", "args": ""}`
	ws, err := websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	send := []byte(helloMessage.Hello)
	rcvMessage = sendRcvMessage(ws, send)
	send = []byte(helloMessage.GetInfo)
	rcvMessage = sendRcvMessage(ws, send)
	version := parseAllData(rcvMessage, &state)
	fmt.Printf("Version of East Server: %s\n", version)
	fmt.Println("Available modules on East server: ")
	for i := 0; i < len(state.AvailableModules); i++ {
		fmt.Println("\t" + state.AvailableModules[i].Name)
	}
	state.Path = "/"
	reader := bufio.NewReader(os.Stdin)
	for i := 0; i < 1; {
		fmt.Printf("East cli@" + state.Path + ">>>")
		comm, _ := reader.ReadString('\n')
		re := regexp.MustCompile(`\r?\n`)
		comm = re.ReplaceAllString(comm, "")
		switch {
		case strings.HasPrefix(comm, "help"):
			comm, i = cliCommand(ws, "help", &state)
		case strings.HasPrefix(comm, "quit"):
			comm, i = "Close Programm...", 1
		case strings.HasPrefix(comm, "use"):
			comm, i = cliCommand(ws, comm, &state)
		case strings.HasPrefix(comm, "search") && len(strings.Split(comm, " ")) > 0:
			comm, i = findCli(comm, &state), 0
		case len(state.ModuleName) > 0:
			comm, i = moduleCommand(ws, comm, &state), 0
		}
		if strings.HasSuffix(comm, "\n") {
			fmt.Printf("%s", comm)
		} else {
			fmt.Printf("%s\n", comm)
		}
	}
}

func findCli(comm string, state *ProgrammStatus) string {
	var resultString string
	searchString := strings.SplitN(comm, " ", 3)
	if len(searchString) > 2 {
		switch searchString[1] {
		case "cve":
			for i := 0; i < len(state.AvailableModules); i++ {
				if strings.Contains(state.AvailableModules[i].CVEName, searchString[2]) {
					resultString += state.AvailableModules[i].Name + "\n"
				}
			}
			return resultString
		case "description":
			for i := 0; i < len(state.AvailableModules); i++ {
				if strings.Contains(state.AvailableModules[i].Description, searchString[2]) {
					resultString += state.AvailableModules[i].Name + "\n"
				}
			}
			return resultString
		case "notes":
			for i := 0; i < len(state.AvailableModules); i++ {
				if strings.Contains(state.AvailableModules[i].Notes, searchString[2]) {
					resultString += state.AvailableModules[i].Name + "\n"
				}
			}
			return resultString
		case "name":
			for i := 0; i < len(state.AvailableModules); i++ {
				if strings.Contains(state.AvailableModules[i].Name, searchString[2]) {
					resultString += state.AvailableModules[i].Name + "\n"
				}
			}
			return resultString
		case "vendor":
			for i := 0; i < len(state.AvailableModules); i++ {
				if strings.Contains(state.AvailableModules[i].Vendor, searchString[2]) {
					resultString += state.AvailableModules[i].Name + "\n"
				}
			}
			return resultString
		default:
			resultString += "Using search: 'search [cve, name, notes, description, vendor] value_for_search'"
		}
	} else if len(searchString) == 2 {
		for i := 0; i < len(state.AvailableModules); i++ {
			if strings.Contains(state.AvailableModules[i].Name, searchString[1]) {
				resultString += state.AvailableModules[i].Name + " finding in name\n"
			} else if strings.Contains(state.AvailableModules[i].CVEName, searchString[1]) {
				resultString += state.AvailableModules[i].Name + " finding in cve\n"
			} else if strings.Contains(state.AvailableModules[i].Description, searchString[1]) {
				resultString += state.AvailableModules[i].Name + " finding in description\n"
			} else if strings.Contains(state.AvailableModules[i].Notes, searchString[1]) {
				resultString += state.AvailableModules[i].Name + " finding in notes\n"
			} else if strings.Contains(state.AvailableModules[i].Vendor, searchString[1]) {
				resultString += state.AvailableModules[i].Name + " finding in vendor\n"
			}
		}
	} else {
		resultString += "Using search: 'search [cve, name, notes, description, vendor] value_for_search\n"
	}
	if len(resultString) == 0 {
		resultString += "Nothing find\n"
	}
	return resultString
}
func moduleCommand(ws *websocket.Conn, command string, state *ProgrammStatus) string {
	switch {
	case strings.HasPrefix(command, "show") && len(strings.Split(command, " ")) > 1:
		return mdlShow(ws, strings.SplitN(command, " ", 2)[1], state)
	case strings.HasPrefix(command, "set") && len(strings.Split(command, " ")) > 1:
		return setOption(strings.SplitN(command, " ", 2)[1], state)
	case strings.HasPrefix(command, "run"):
		return runModule(ws, state)
	default:
		return "Availaible commands: show, set, start\n"
	}
}

func setOption(command string, state *ProgrammStatus) string {
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

func mdlShow(ws *websocket.Conn, comm string, state *ProgrammStatus) string {
	switch {
	case comm == "options":
		showOptCli(ws, state)
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
func quitCli() (string, int) {
	return "Close programm...", 1
}

func helpCli() (string, int) {
	return "Help page for working with EaST CLI", 0
}

func useCli(state *ProgrammStatus) string {
	fmt.Printf("Using %s", state.ModuleName)
	return " type \"show options\" for show available options"
}

func showOptCli(ws *websocket.Conn, state *ProgrammStatus) {
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

func cliCommand(ws *websocket.Conn, comm string, state *ProgrammStatus) (string, int) {
	if ((len(strings.Split(comm, " ")) == 1) && (len(state.ModuleName) == 0)) || (comm == "quit") {
		availableCommand := map[string]interface{}{
			"quit": quitCli,
			"help": helpCli,
		}
		switch comm {
		case "quit":
			return availableCommand["quit"].(func() (string, int))()
		default:
			state.Path = "/help"
			return availableCommand["help"].(func() (string, int))()
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
					return availableCommand["use"].(func(*ProgrammStatus) string)(state), 0
				}
			}
			return availableCommand["help"].(func() (string, int))()
		default:
			return availableCommand["help"].(func() (string, int))()
		}
	} else {
		// commands for module
		return "here", 0

	}
}

func sendRcvMessage(ws *websocket.Conn, mess []byte) []byte {
	var recieve []byte
	websocket.Message.Send(ws, mess)
	websocket.Message.Receive(ws, &recieve)
	return recieve
}

func parseAllData(rcv []byte, state *ProgrammStatus) string {
	var version, module string
	version = strings.Split(string(rcv), `version": `)[1]
	version = strings.SplitN(string(version), ",", 2)[0]
	version = strings.Replace(version, `"`, "", -1)
	modules := strings.Split(string(rcv), `"CVE Name"`)
	modulesString := "["
	for i := 1; i < len(modules); i++ {
		module = strings.SplitN(string(modules[i]), "}", 2)[0]
		module = fixLinksArr(module)
		modulesString += `{"CVE Name"` + module
		if i != (len(modules) - 1) {
			modulesString += "}, "
		} else {
			modulesString += "}"
		}
	}
	modulesString += "]"
	err := json.Unmarshal([]byte(modulesString), &state.AvailableModules)
	if err != nil {
		log.Fatal(err)
	}
	return version
}

func fixLinksArr(module string) string {
	var moduleString string
	if len(strings.Split(module, `"LINKS": `)) > 1 {
		moduleString = strings.Split(module, `"LINKS": `)[1]
		moduleInfo := strings.SplitN(moduleString, ",", 2)
		if bytes.HasPrefix([]byte(moduleString), []byte("[")) == false {
			moduleString = strings.Split(module, `"LINKS": `)[0] + `"LINKS": [` + moduleInfo[0] + `], ` + moduleInfo[1]
		} else {
			moduleString = strings.Split(module, `"LINKS": `)[0] + `"LINKS": ` + moduleString
		}
	} else {
		moduleString = module
	}
	return moduleString
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
func runModule(ws *websocket.Conn, state *ProgrammStatus) string {
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
