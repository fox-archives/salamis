package main

import (
	"fmt"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/pelletier/go-toml"
)

func doInit(opts Options) {
	// create config
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
	p(err)

	// write config file only if it doesn't already exist
	file, err := os.OpenFile(opts.ConfigFile, os.O_CREATE|os.O_EXCL|os.O_RDWR, 0644)
	defer func() {
		file.Close()
		p(err)
	}()

	if err != nil {
		if os.IsExist(err) {
			fmt.Printf("%s already exists. Remove it before continuing. Exiting\n", opts.ConfigFile)
			os.Exit(1)
			return
		}
		panic(err)
	}

	_, err = file.Write(configRaw)
	p(err)
}

func doUpdate(opts Options) {

}

func doCheck(opts Options) {
	extensions := getVscodeExtensions()
	config := readConfig()

	fmt.Println("Extensions that are installed globally, but could not be found local. Add them to your extensions.toml if you want to use them")
	for _, globalExtension := range extensions {
		if globalExtension == "" {
			continue
		}

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
	fmt.Println("Extensions that are installed locally, but not globally. Remove these from your extensions.toml, or install the extension globally, and re-clone your extensions")
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
			for _, group := range config.Workspaces {
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

}

func doLaunch(opts Options, workspaceName string) {
	extensionsDir := filepath.Join(opts.WorkspaceDir, workspaceName)

	cmd := exec.Command("code", "--extensions-dir", extensionsDir, ".")
	cmd.Stderr = os.Stderr
	stdout, err := cmd.Output()
	p(err)

	fmt.Println(stdout)
}
