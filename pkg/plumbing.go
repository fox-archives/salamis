package pkg

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

// downloads currently installed extensions.
func doExtensionsInstall(opts Options) {
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
		handle(err)

		fmt.Println(string(stdout))
	}

	// now, we have to rename all the files to remove the version number
	dirs, err := ioutil.ReadDir(opts.ExtensionsDir)
	handle(err)

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
		handle(err)

		fmt.Printf("Renaming file: '%s'\n", new)
	}
}

func doExtensionsRemove(opts Options) {
	err := os.RemoveAll(opts.ExtensionsDir)
	handle(err)

	err = os.MkdirAll(opts.ExtensionsDir, 0755)
	handle(err)
}

func doExtensionsSymlink(opts Options) {
	config := readConfig(opts)

	err := os.RemoveAll(opts.WorkspaceDir)
	handle(err)

	err = os.MkdirAll(opts.WorkspaceDir, 0755)
	handle(err)

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

func doExtensionsUnsymlink(opts Options) {
	err := os.RemoveAll(opts.WorkspaceDir)
	handle(err)

	err = os.MkdirAll(opts.WorkspaceDir, 0755)
	handle(err)
}

func doXdgInstall(opts Options) {
	type DesktopEntry struct {
		Name          string
		ExtensionsDir string
		Exec          string
		Icon          string
	}

	for _, workspace := range readConfig(opts).Workspaces {
		entry := &DesktopEntry{
			Name:          workspace.Name,
			ExtensionsDir: opts.ExtensionsDir,
			Exec:          "path",
			Icon:          "icon",
		}

		tmpl, err := template.New("entry").Parse(`[Desktop Entry]
Name={{.Name}} VSCode
Comment=Code Editing. Redefined.
GenericName=Text Editor
Exec=/usr/share/code/code --no-sandbox --unity-launch {{.ExtensionsDir}}
Icon=com.visualstudio.code
Type=Application
StartupNotify=false
StartupWMClass=Code
Categories=Utility;TextEditor;Development;IDE;
MimeType=text/plain;inode/directory;application/x-code-workspace;
Actions=new-empty-window;
Keywords=vscode;{{.Name}};`)
		handle(err)

		var buf bytes.Buffer
		err = tmpl.Execute(&buf, entry)
		handle(err)

		ioutil.WriteFile(filepath.Join(opts.ApplicationsDir, "salamis."+strings.ToLower(entry.Name)+".desktop"), []byte(buf.String()), 0o644)
	}
}

func doXdgRemove(opts Options) {
	matches, err := filepath.Glob(filepath.Join(opts.ApplicationsDir, "salamis.*"))
	handle(err)

	for _, file := range matches {
		err := os.Remove(file)
		handle(err)
	}
}
