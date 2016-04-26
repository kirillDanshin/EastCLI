package main

import "github.com/gorilla/websocket"

type Message struct {
	Name string `json:"name"`
	Type string `json:"type"`
}

type helloMessage struct {
	Hello Message `json:"hello"`
}

type WSOptions struct {
}

type WSArgs struct {
	ModuleName  string    `json:"module_name"`
	UseListener bool      `json:"use_listener"`
	Options     WSOptions `json:"options"`
}

type WSMessage struct {
	Command string `json:"command"`
	Args    WSArgs `json:"args"`
}

type StartMessage struct {
	Hello   string
	GetInfo string
}

type Module struct {
	CVEName      string   `json:"CVE NAME"`
	Vendor       string   `json:"VENDOR"`
	Name         string   `json:"NAME"`
	Links        []string `json:"LINKS, omitempty"`
	Changelog    string   `json:"CHANGELOG"`
	Notes        string   `json:"NOTES"`
	DownloadLink string   `json:"DOWNLOAD_LINK, omitempty"`
	Path         string   `json:"PATH"`
	IsFile       bool     `json:"isFile"`
	Description  string   `json:"DESCRIPTION"`
}

type AllModules struct {
	Module []Module `json:"modules"`
}

type DirModules struct {
	Children    []Module `json:"children"`
	IsFile      bool     `json:"isFile"`
	Description string   `json:"DESCRIPTION"`
	Name        string   `json:"NAME"`
}
type ModulesArray struct {
	Modules []DirModules `json:"modules"`
}
type InfoServer struct {
}

type ArgsServer struct {
	Version string `json:"version"`
	Modules []byte `json:"modules"`
}

type AllData struct {
	Args    []byte `json:"args"`
	Command string `json:"command"`
}
type CommandHello struct {
	Command string `json:"command"`
}

func (M Module) CheckName(name string) bool {
	if M.Name == name {
		return true
	}
	return false
}

type Option struct {
	Name  string `json:"option"`
	Value struct {
		Type  string `json:"type"`
		Value string `json:"value"`
	} `json:"value"`
}
type ShowOptions struct {
	Args    []Option `json:"args"`
	Command string   `json:"command"`
}
type ModuleOptions struct {
	Options ShowOptions
	Name    string
}
type ModuleRunArgs struct {
	ModuleName string   `json:"module_name"`
	Listener   bool     `json:"use_listener"`
	Option     []Option `json:"options"`
}
type ModuleRun struct {
	Args    ModuleRunArgs `json:"args"`
	Command string        `json:"command"`
}
type ProgrammStatus struct {
	ws               *websocket.Conn
	Path             string
	ModuleName       string
	MessageEast      string
	MdlOptions       []Option
	State            string
	AvailableModules []Module
	AvailableCommand []string
	MdlNumb          int
	RunModule        ModuleRun
	UseListener      bool
}
