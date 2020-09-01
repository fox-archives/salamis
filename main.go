package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pelletier/go-toml"
)

func ensureLength(arr []string, minLength int, message string) {
	if len(arr) < minLength {
		log.Fatalln(message)
	}
}

func contains(arr []string, query string) bool {
	for _, el := range arr {
		if el == query {
			return true
		}
	}
	return false
}

// Extension gives information about a particular vscode extension
type Extension struct {
	Name string   `json:"name"`
	Tags []string `json:"stags"`
}

// Config gives information about the whole configuration file
type Config struct {
	Version    string      `json:"version"`
	Extensions []Extension `json:"extension"`
}

func main() {
	args := os.Args[1:]

	if contains(args, "--help") || contains(args, "-h") {
		fmt.Println(`sparta

Description:
    Contextual management of vscode extensions

Commands:
    generate
    Generates all extension folders to be used by vscode

    launch [workspace]
    launches a particular workspace in vscode`)
		os.Exit(0)
	}

	ensureLength(args, 1, "Must specify command. Exiting")

	command := args[0]
	if command == "generate" {
		var config Config

		configRaw, err := ioutil.ReadFile(filepath.Join("extensions.toml"))
		if err != nil {
			panic(err)
		}
		if err = toml.Unmarshal(configRaw, &config); err != nil {
			panic(err)
		}

		for _, extension := range config.Extensions {
			fmt.Printf("EXTENSION: %s\n", extension.Name)

			for _, tag := range extension.Tags {
				fmt.Printf("tag: %s\n", tag)

				// install extension
				extensionsDir := filepath.Join("workspaces", tag)

				cmd := exec.Command("code", "--extensions-dir", extensionsDir, "--install-extension", extension.Name, "--force")
				cmd.Stderr = os.Stderr
				stdout, err := cmd.Output()
				if err != nil {
					panic(err)
				}
				fmt.Println(string(stdout))
			}
		}
	} else if command == "launch" {
		ensureLength(args, 2, "Must pass in a workspace name")
		workspaceName := args[1]
		extensionsDir := filepath.Join("workspaces", workspaceName)

		cmd := exec.Command("code", "--extensions-dir", extensionsDir, ".")
		stdout, err := cmd.Output()
		if err != nil {
			panic(err)
		}
		fmt.Println(stdout)
	} else {
		log.Fatalln("Unknown Command. Exiting")
	}

}
