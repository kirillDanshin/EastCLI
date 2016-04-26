package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"regexp"
	"strings"
)

func runCLI() {
	reader := bufio.NewReader(os.Stdin)
	re := regexp.MustCompile(`\r?\n`)

	for true {
		fmt.Printf("East cli@" + state.Path + ">>> ")
		comm, err := reader.ReadString('\n')
		if err != nil {
			log.Fatalf("Error occurred: %s", err)
		}
		comm = re.ReplaceAllString(comm, "")

		close := false

		if strings.HasPrefix(comm, "help") {
			comm, close = cliCommand(ws, "help")
		}
		// @kirilldanshin
		// we want to quit by quit, exit and \q
		if strings.HasPrefix(comm, "quit") || strings.HasPrefix(comm, "exit") {
			comm, close = "Close Programm...", true
		}

		switch {
		case close:
			break
		case strings.HasPrefix(comm, "use"):
			comm, close = cliCommand(ws, comm)
		case strings.HasPrefix(comm, "search") && len(strings.Split(comm, " ")) > 0:
			comm = findCli(comm)
		case len(state.ModuleName) > 0:
			comm = moduleCommand(ws, comm)
		default:
			comm = "Unknown command."
		}
		if strings.HasSuffix(comm, "\n") {
			fmt.Printf("%s", comm)
		} else {
			fmt.Printf("%s\n", comm)
		}
		if close {
			break
		}
	}
}
