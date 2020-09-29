package main

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
)

// downloads currently installed extensions
func doDownloadExtensions(opts Options) {
	currentExtensions := getVscodeExtensions()
	for _, extension := range currentExtensions {
		if extension == "" {
			continue
		}

		fmt.Printf("EXTENSION: %s\n", extension)

		// install extensions, automatically creates opts.ExtensionsDir if it doesn't already exist
		cmd := exec.Command("code", "--extensions-dir", opts.ExtensionsDir, "--install-extension", extension, "--force")
		cmd.Stderr = os.Stderr
		stdout, err := cmd.Output()
		p(err)

		fmt.Println(string(stdout))
	}

	// now, we have to rename all the files to remove the version number
	dirs, err := ioutil.ReadDir(opts.ExtensionsDir)
	p(err)

	for _, dir := range dirs {
		// if extension doesn't have version, it has already
		// been renamed
		if !extensionHasVersion(dir.Name()) {
			continue
		}

		// take off the version number and make it lowercase
		// vscode lowercases the extension name, but not the author, so this normalizes it, giving a stable link target
		parts := strings.Split(dir.Name(), "-")
		newName := strings.Join(parts[:len(parts)-1], "-")
		newName = strings.ToLower(newName)

		old := filepath.Join(opts.ExtensionsDir, dir.Name())
		new := filepath.Join(opts.ExtensionsDir, newName)

		err := os.Rename(old, new)
		p(err)

		fmt.Printf("Renaming file: '%s'\n", new)
	}
}

func doRemoveExtensions(opts Options) {
	err := os.RemoveAll(opts.ExtensionsDir)
	p(err)

	err = os.MkdirAll(opts.ExtensionsDir, 0755)
	p(err)
}

func doSymlinkExtensions(opts Options) {
	config := readConfig(opts)

	err := os.RemoveAll(opts.WorkspaceDir)
	p(err)

	err = os.MkdirAll(opts.WorkspaceDir, 0755)
	p(err)

	if len(config.Workspaces) == 0 {
		fmt.Println("You have no workspaces specified")
		os.Exit(1)
	}

	// generate workspaces
	for _, workspace := range config.Workspaces {
		fmt.Printf("WORKSPACE: %s\n", workspace.Name)
		err := os.MkdirAll(filepath.Join(opts.WorkspaceDir, workspace.Name), 0755)
		if err != nil && !os.IsExist(err) {
			panic(err)
		}

		for _, extension := range config.Extensions {
			fmt.Printf("EXTENSION: %s\n", extension.Name)

			for _, tag := range extension.Tags {
				// if any tag in current extension is used in the workspace
				if contains(workspace.Use, tag) {
					src := filepath.Join("../..", opts.ExtensionsDir, strings.ToLower(extension.Name))
					dest := filepath.Join(opts.WorkspaceDir, workspace.Name, extension.Name)

					err := os.Symlink(src, dest)
					if err != nil && !os.IsExist(err) {
						panic(err)
					}

					// go to next extension of a workspace
					continue
				}
			}
		}
	}
}

func doSymlinkRemove(opts Options) {
	err := os.RemoveAll(opts.WorkspaceDir)
	p(err)

	err = os.MkdirAll(opts.WorkspaceDir, 0755)
	p(err)
}
