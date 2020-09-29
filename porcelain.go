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
		err := file.Close()
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

func doUpdate(opts Options) {}

func doCheck(opts Options) {
	extensions := getVscodeExtensions()
	config := readConfig(opts)

	fmt.Println(`Extensions saved in Sparta, but not used in the config`)
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
			fmt.Printf("- %s\n", globalExtension)
		}
	}

	fmt.Println()
	fmt.Println("Extensions that are used in the config, but not saved in Sparta")
	for _, spartaExtension := range config.Extensions {
		isGlobal := false

		for _, globalExtension := range extensions {

			if spartaExtension.Name == globalExtension {

				isGlobal = true
				continue
			}
		}

		if !isGlobal {
			fmt.Printf("- %s\n", spartaExtension.Name)
		}
	}

	fmt.Println()
	fmt.Println("Extensions with missing tags")
	for _, extension := range config.Extensions {
		if len(extension.Tags) == 0 {
			fmt.Printf("- %s\n", extension.Name)
		}
	}

	fmt.Println()
	fmt.Println("Tags assigned to extensions that aren't referenced by any workspace")
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
				fmt.Printf("- %s\n", tag)
			}
		}
	}

}

func doLaunch(opts Options, workspaceName string) {
	extensionsDir := filepath.Join(opts.WorkspaceDir, workspaceName)

	_, err := os.Stat(extensionsDir)
	if err != nil {
		if os.IsPermission(err) {
			fmt.Printf("Could not access extension folder '%s'\n", extensionsDir)
			os.Exit(1)
		} else if os.IsNotExist(err) {
			fmt.Printf("Folder '%s' does not exist\n", extensionsDir)
			os.Exit(1)
		}
		panic(err)

	}

	cmd := exec.Command("code", "--extensions-dir", extensionsDir, ".")
	cmd.Stderr = os.Stderr
	stdout, err := cmd.Output()
	p(err)

	fmt.Println(stdout)
}
