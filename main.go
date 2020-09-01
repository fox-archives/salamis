package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pelletier/go-toml"

	"github.com/otiai10/copy"
)

func main() {
	args := os.Args[1:]

	first := args[0]

	if first == "copy" {
		fmt.Println("copying current extensions")

		home, err := os.UserHomeDir()
		if err != nil {
			log.Println("Home directory not found")
			panic(err)
		}
		extensionsDir := filepath.Join(home, ".vscode/extensions")
		wd, err := os.Getwd()
		if err != nil {
			log.Println("Could not get current working directory")
			panic(err)
		}
		programDir := filepath.Join(wd, "extensions")
		if err = copy.Copy(extensionsDir, programDir); err != nil {
			log.Println("Could not copy extensions directory")
			panic(err)
		}
	} else if first == "copy2" {
		type Extension struct {
			Name string   `json:'name"`
			Tags []string `json:"tags"`
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
			wd, err := os.Getwd()
			if err != nil {
				panic(err)
			}

			if err := os.MkdirAll(filepath.Join(wd, "workspaces", name), os.ModePerm); err != nil {
				panic(err)
			}
		} else if second == "delete" {
			if name == "" {
				fmt.Println("Must pass in a workspace name")
				os.Exit(1)
			}

			wd, err := os.Getwd()
			if err != nil {
				panic(err)
			}

			if err := os.RemoveAll(filepath.Join(wd, "workspaces", name)); err != nil {
				panic(err)
			}
		}
	} else if first == "open" {
		name := args[1]

		if name == "" {
			fmt.Println("need to pass in a workspace name")
			os.Exit(1)
		}
		wd, err := os.Getwd()
		if err != nil {
			panic(err)
		}
		extensionsDir := filepath.Join(wd, "workspaces", name)
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
