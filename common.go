package main

import "fmt"

func quitCli() (string, bool) {
	return "Close programm...", true
}

func helpCli() (string, bool) {
	return "Help page for working with EaST CLI", false
}

func useCli() string {
	return fmt.Sprintf(`Using %s type "show options" for show available options`, state.ModuleName)
}
