package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pelletier/go-toml"

	"github.com/otiai10/copy"
)

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

	first := args[0]

	if contains(args, "--help") || contains(args, "-h") {
		fmt.Println(`sparta

Description:
    Contextual management of vscode extensions

Commands:
    generate
    Generates all extension folders to be used by vscode

    lauch [workspace]
    launches a particular workspace in vscode`)
		os.Exit(0)
	}

	if first == "generate" {
		type Extension struct {
			Name string   `json:"name"`
			Tags []string `json:"stags"`
		}

		type Extensions struct {
			Version    string      `json:"version"`
			Extensions []Extension `json:"extension"`
		}

		var extensionsToml Extensions

		t, err := ioutil.ReadFile(filepath.Join("extensions.toml"))
		if err != nil {
			panic(err)
		}
		if err = toml.Unmarshal(t, &extensionsToml); err != nil {
			panic(err)
		}

		for _, extension := range extensionsToml.Extensions {
			fmt.Printf("processing %s\n", extension.Name)
			for _, tag := range extension.Tags {
				fmt.Printf("foo %s\n", tag)
				if err := copy.Copy(filepath.Join("extensions", extension.Name), filepath.Join("workspaces", tag, extension.Name)); err != nil {
					panic(err)
				}
			}
		}

		fmt.Println(string(t))
	} else if first == "workspace" {
		second := args[1]
		name := args[2]

		if name == "" {
			fmt.Println("need to pass in a workspace name")
			os.Exit(1)
		}

		if second == "create" {
			if err := os.MkdirAll(filepath.Join("workspaces", name), os.ModePerm); err != nil {
				panic(err)
			}
		} else if second == "delete" {
			if name == "" {
				fmt.Println("Must pass in a workspace name")
				os.Exit(1)
			}

			if err := os.RemoveAll(filepath.Join("workspaces", name)); err != nil {
				panic(err)
			}
		}
	} else if first == "launch" {
		name := args[1]

		if name == "" {
			fmt.Println("need to pass in a workspace name")
			os.Exit(1)
		}

		extensionsDir := filepath.Join("workspaces", name)
		cmd := exec.Command("code", "--extensions-dir", extensionsDir, ".")
		stdout, err := cmd.Output()

		if err != nil {
			panic(err)
		}
		fmt.Println(stdout)
	} else {
		fmt.Println("unknown option")
		os.Exit(1)
	}

}
