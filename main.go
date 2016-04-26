package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"log"
	"strings"

	"golang.org/x/net/websocket"
)

var (
	origin = "http://localhost/"
	url    = "ws://localhost:49999/"
	state  = &ProgrammStatus{}
	ws     *websocket.Conn
)

func main() {
	var (
		rcvMessage   []byte
		helloMessage = StartMessage{
			Hello:   `{"hello": {"name": "EastCLI", "type": "cli"}}`,
			GetInfo: `{"command": "get_all_server_data", "args": ""}`,
		}
		err error
	)
	ws, err = websocket.Dial(url, "", origin)
	if err != nil {
		log.Fatal(err)
	}
	send := []byte(helloMessage.Hello)
	rcvMessage = sendRcvMessage(ws, send)
	send = []byte(helloMessage.GetInfo)
	rcvMessage = sendRcvMessage(ws, send)
	version := parseAllData(rcvMessage)
	fmt.Printf("Version of East Server: %s\n", version)
	fmt.Println("Available modules on East server: ")
	for i := 0; i < len(state.AvailableModules); i++ {
		fmt.Println("\t" + state.AvailableModules[i].Name)
	}
	state.Path = "/"
	runCLI()
}

func moduleCommand(ws *websocket.Conn, command string) string {
	switch {
	case strings.HasPrefix(command, "show") && len(strings.Split(command, " ")) > 1:
		return mdlShow(ws, strings.SplitN(command, " ", 2)[1])
	case strings.HasPrefix(command, "set") && len(strings.Split(command, " ")) > 1:
		return setOption(strings.SplitN(command, " ", 2)[1])
	case strings.HasPrefix(command, "run"):
		return runModule(ws)
	default:
		return "Availaible commands: show, set, start\n"
	}
}

func sendRcvMessage(ws *websocket.Conn, mess []byte) []byte {
	var recieve []byte
	websocket.Message.Send(ws, mess)
	websocket.Message.Receive(ws, &recieve)
	return recieve
}

func parseAllData(rcv []byte) string {
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
	var AvailableModules []Module
	err := json.Unmarshal([]byte(modulesString), &AvailableModules)
	state.AvailableModules = AvailableModules
	if err != nil {
		log.Fatalf("Error while parsing data: %s", err)
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
