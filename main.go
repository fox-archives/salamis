package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

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
	Name string   `toml:"name"`
	Tags []string `toml:"tags"`
}

// Config gives information about the whole configuration file
type Config struct {
	Version    string      `toml:"version"`
	Extensions []Extension `toml:"extensions"`
}

func main() {
	args := os.Args[1:]

	if contains(args, "--help") || contains(args, "-h") {
		fmt.Println(`sparta

Description:
    Contextual management of vscode extensions

Commands:
	 init
    initiates an 'extensions.toml' folder that contains all extensions for tagging
    generate
    Generates all extension folders to be used by vscode

	 Clear
	 Clears all workspaces

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
		cmd.Stderr = os.Stderr
		stdout, err := cmd.Output()
		if err != nil {
			panic(err)
		}
		fmt.Println(stdout)
	} else if command == "clear" {
		if err := os.RemoveAll("workspaces"); err != nil {
			panic(err)
		}
	} else if command == "init" {
		var config Config
		config.Version = "1"

		cmd := exec.Command("code", "--list-extensions")
		cmd.Stderr = os.Stderr
		stdout, err := cmd.Output()
		if err != nil {
			panic(err)
		}
		extensions := strings.Split(string(stdout), "\n")
		for _, extension := range extensions {
			config.Extensions = append(config.Extensions, Extension{
				Name: extension,
			})
		}

		configRaw, err := toml.Marshal(config)
		if err != nil {
			panic(err)
		}

		file, err := os.OpenFile("extensions.toml", os.O_CREATE|os.O_EXCL|os.O_RDWR, 0644)
		if err != nil {
			if os.IsExist(err) {
				fmt.Println("extensions.toml already exists. Remove it before continuing. Exiting")
				os.Exit(1)
				return
			}
			panic(err)
		}

		if _, err = file.Write(configRaw); err != nil {
			panic(err)
		}
		if err := file.Close(); err != nil {
			panic(err)
		}

	} else {
		log.Fatalln("Unknown Command. Exiting")
	}

}
