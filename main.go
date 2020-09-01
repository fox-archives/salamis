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

// Group is a group of extensions
type Group struct {
	Name string   `toml:"name"`
	Use  []string `toml:"use"`
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
	Groups     []Group     `toml:"groups"`
}

func readConfig() Config {
	var config Config

	configRaw, err := ioutil.ReadFile(filepath.Join("extensions.toml"))
	if err != nil {
		panic(err)
	}
	if err = toml.Unmarshal(configRaw, &config); err != nil {
		panic(err)
	}

	return config
}

func getVscodeExtensions() []string {
	cmd := exec.Command("code", "--list-extensions")

	cmd.Stderr = os.Stderr
	fmt.Println("h")
	stdout, err := cmd.Output()
	fmt.Println("thing")

	if err != nil {
		panic(err)
	}
	return strings.Split(string(stdout), "\n")
}
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

func main() {
	args := os.Args[1:]

	if contains(args, "--help") || contains(args, "-h") {
		fmt.Println(`sparta

Description:
  Contextual vscode extension management

Commands:
  init
    Initiates an 'extensions.toml' folder that contains all extensions for tagging

  generate
    Generates all extension folders to be used by vscode

  clear
    Removes all downloaded extensions from their workspaces

  check
    Prints all extensions mismatches between default globally installed and ones defined in extensions.toml

  launch [workspace]
    Launches a particular workspace in vscode`)
		os.Exit(0)
	}

	ensureLength(args, 1, "Must specify command. Exiting")

	command := args[0]
	if command == "generate" {
		config := readConfig()

		if err := os.MkdirAll("workspaces", 0755); err != nil && !os.IsExist(err) {
			panic(err)
		}
		if err := os.MkdirAll("aggregations", 0755); err != nil && !os.IsExist(err) {
			panic(err)
		}

		// generate workspaces
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

		// generate every combination of tags

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

		extensions := getVscodeExtensions()
		for _, extension := range extensions {
			if extension == "" {
				continue
			}

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

	} else if command == "check" {

		extensions := getVscodeExtensions()
		config := readConfig()

		fmt.Println("Extensions that are installed globally, but could not be found local")
		for _, globalExtension := range extensions {
			isHere := false
			for _, spartaExtension := range config.Extensions {
				if globalExtension == spartaExtension.Name {
					isHere = true
					continue
				}
			}

			if !isHere {
				fmt.Printf("NOT LOCAL: %s\n", globalExtension)
			}
		}

		fmt.Println()
		fmt.Println("Extensions that are installed locally, but not globally")
		for _, spartaExtension := range config.Extensions {
			isGlobal := false

			for _, globalExtension := range extensions {

				if spartaExtension.Name == globalExtension {

					isGlobal = true
					continue
				}
			}

			if !isGlobal {
				fmt.Printf("NOT GLOBAL: %s\n", spartaExtension.Name)
			}
		}

		fmt.Println()
		fmt.Println("Extensions that don't have any tags")
		for _, extension := range config.Extensions {
			if len(extension.Tags) == 0 {
				fmt.Printf("NO TAGS: %s\n", extension.Name)
			}
		}

		fmt.Println()
		fmt.Println("Extensions tags that aren't in a group")
		for _, extension := range config.Extensions {
			for _, tag := range extension.Tags {
				inGroup := false
			g:
				for _, group := range config.Groups {
					for _, usedTag := range group.Use {
						if usedTag == tag {
							inGroup = true
							continue g
						}
					}

				}
				if !inGroup {
					fmt.Printf("TAG NOT USED: %s\n", tag)
				}
			}
		}

	} else {
		log.Fatalln("Unknown Command. Exiting")
	}

}
