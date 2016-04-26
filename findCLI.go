package main

// search-related functions goes here

import "strings"

func searchInCVE(searchQuery string) (results string) {
	for i := 0; i < len(state.AvailableModules); i++ {
		if strings.Contains(state.AvailableModules[i].CVEName, searchQuery) {
			results += state.AvailableModules[i].Name + "\n"
		}
	}
	return
}

func searchInDescr(searchQuery string) (results string) {
	for i := 0; i < len(state.AvailableModules); i++ {
		if strings.Contains(state.AvailableModules[i].Description, searchQuery) {
			results += state.AvailableModules[i].Name + "\n"
		}
	}
	return
}

func searchInNotes(searchQuery string) (results string) {
	for i := 0; i < len(state.AvailableModules); i++ {
		if strings.Contains(state.AvailableModules[i].Notes, searchQuery) {
			results += state.AvailableModules[i].Name + "\n"
		}
	}
	return
}

func searchInNames(searchQuery string) (results string) {
	for i := 0; i < len(state.AvailableModules); i++ {
		if strings.Contains(state.AvailableModules[i].Name, searchQuery) {
			results += state.AvailableModules[i].Name + "\n"
		}
	}
	return
}

func searchInVendors(searchQuery string) (results string) {
	for i := 0; i < len(state.AvailableModules); i++ {
		if strings.Contains(state.AvailableModules[i].Vendor, searchQuery) {
			results += state.AvailableModules[i].Name + "\n"
		}
	}
	return
}

func searchEverywhere(searchQuery string) (results string) {
	for i := 0; i < len(state.AvailableModules); i++ {
		if strings.Contains(state.AvailableModules[i].Name, searchQuery) {
			results += state.AvailableModules[i].Name + " finding in name\n"
		} else if strings.Contains(state.AvailableModules[i].CVEName, searchQuery) {
			results += state.AvailableModules[i].Name + " finding in cve\n"
		} else if strings.Contains(state.AvailableModules[i].Description, searchQuery) {
			results += state.AvailableModules[i].Name + " finding in description\n"
		} else if strings.Contains(state.AvailableModules[i].Notes, searchQuery) {
			results += state.AvailableModules[i].Name + " finding in notes\n"
		} else if strings.Contains(state.AvailableModules[i].Vendor, searchQuery) {
			results += state.AvailableModules[i].Name + " finding in vendor\n"
		}
	}
	return
}

const searchUsage = "Using search: 'search [all, cve, name, notes, description, vendor] value_for_search'"

func findCli(comm string) string {
	var resultString string
	searchQuery := strings.SplitN(comm, " ", 3)

	switch {
	case searchQuery[1] == "cve":
		resultString += searchInCVE(searchQuery[2])
	case searchQuery[1] == "description":
		resultString += searchInDescr(searchQuery[2])
	case searchQuery[1] == "notes":
		resultString += searchInNotes(searchQuery[2])
	case searchQuery[1] == "name":
		resultString += searchInNames(searchQuery[2])
	case searchQuery[1] == "vendor":
		resultString += searchInVendors(searchQuery[2])
	case searchQuery[1] == "all":
		resultString += searchEverywhere(searchQuery[2])
	case len(searchQuery) == 2:
		resultString += searchEverywhere(searchQuery[1])
	default:
		resultString += searchUsage
	}

	if len(resultString) == 0 {
		resultString += "Nothing find\n"
	}
	return resultString
}
