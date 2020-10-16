package main

import (
	"bytes"
	"fmt"
	"html/template"
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

func doInstallXdgDesktopEntries(opts Options) {
	home, _ := os.UserHomeDir()
	dataDir := os.Getenv("XDG_DATA_HOME")
	if dataDir == "" {
		dataDir = filepath.Join(home, ".local", "share")
	}
	destDir := filepath.Join(dataDir, "applications")

	fmt.Printf("Installing to '%s'\n", destDir)

	type DesktopEntry struct {
		Name                 string
		ExtensionCacheFolder string
		Exec                 string
		Icon                 string
	}

	for _, workspace := range readConfig(opts).Workspaces {
		entry := &DesktopEntry{
			Name:                 workspace.Name,
			ExtensionCacheFolder: filepath.Join(home, ".cache", "salamis", "workspaces", workspace.Name),
			Exec:                 "path",
			Icon:                 "icon",
		}

		tmpl, err := template.New("entry").Parse(`[Desktop Entry]
Name={{.Name}} VSCode
Comment=Code Editing. Redefined.
GenericName=Text Editor
Exec=/usr/share/code/code --no-sandbox --unity-launch {{.ExtensionCacheFolder}}
Icon=com.visualstudio.code
Type=Application
StartupNotify=false
StartupWMClass=Code
Categories=Utility;TextEditor;Development;IDE;
MimeType=text/plain;inode/directory;application/x-code-workspace;
Actions=new-empty-window;
Keywords=vscode;{{.Name}};`)
		if err != nil {
			panic(err)
		}

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, entry)
		if err != nil {
			panic(err)
		}

		ioutil.WriteFile(filepath.Join(destDir, "salamis."+strings.ToLower(entry.Name)+".desktop"), []byte(buf.String()), 0o644)
	}
}

func doRemoveXdgDesktopEntries(opts Options) {
	home, _ := os.UserHomeDir()
	dataDir := os.Getenv("XDG_DATA_HOME")
	if dataDir == "" {
		dataDir = filepath.Join(home, ".local", "share")
	}
	destDir := filepath.Join(dataDir, "applications")

	matches, err := filepath.Glob(filepath.Join(destDir, "salamis.*"))
	if err != nil {
		panic(err)
	}
	for _, file := range matches {
		err := os.Remove(file)
		if err != nil {
			panic(err)
		}
	}
}
