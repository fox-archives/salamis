package pkg

import (
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"strings"

	"github.com/pelletier/go-toml"
)

// Extension contains information about each extension. This may be autogenerated.
type Extension struct {
	Name string   `toml:"name"`
	Tags []string `toml:"tags"`
}

// Workspace is the aggregation of extensions which is meant to be used for a specific programming language or other subdomain.
type Workspace struct {
	Name string   `toml:"name"`
	Use  []string `toml:"use"`
}

// Config gives information about the whole configuration file.
type Config struct {
	Extensions []Extension `toml:"extensions"`
	Workspaces []Workspace `toml:"workspaces"`
}

func handle(err error) {
	if err != nil {
		panic(err)
	}
}

// tests if an extension has a version.
func extensionHasVersion(str string) bool {
	str = str[len(str)-1:]

	if strings.Contains("1234567890", str) {
		return true
	}

	return false
}

func isFolderEmpty(path string) bool {
	dirs, err := ioutil.ReadDir(path)
	handle(err)

	if len(dirs) == 0 {
		return true
	}

	return false
}

func readConfig(opts Options) Config {
	var config Config

	configRaw, err := ioutil.ReadFile(filepath.Join(opts.ConfigFile))
	if os.IsNotExist(err) {
		fmt.Println("Error: extensions.toml not found. Did you forget to init?")
		os.Exit(1)
	}
	handle(err)

	err = toml.Unmarshal(configRaw, &config)
	handle(err)

	return config
}

// returns array of extensions
// example: ["yzhang.markdown-all-in-one@3.3.0"]
func getVscodeExtensions() []string {
	cmd := exec.Command("code", "--list-extensions")

	cmd.Stderr = os.Stderr
	stdout, err := cmd.Output()
	handle(err)

	return strings.Split(string(stdout), "\n")
}

func ensureLength(arr []string, minLength int, message string) {
	if len(arr) < minLength {
		fmt.Println(message)
		os.Exit(1)
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

func printHelp() {
	fmt.Println(`salamis

Description:
  Contextual vscode extension management

Commands:
  init
    Initiates an 'extensions.toml' folder that contains all extensions for tagging

  update
    Removes, re-downloads, and unsymlinks, resymlinks all specified extensions. Then, it removes and reinstalls XDG Desktop Entries

  list
    List all workspaces and their specifying tags

  edit
    Opens the 'workspaces.toml' file in your default editor

  check
    Prints all workspace/extensions/tag relationships that may need notice / be of interest. See the documentation for more information about each type of check performed.

  launch [workspace]
    Launches a the specified workspace in vscode

  plumbing extensions-install
    Installs (downloads) all current extensions into the salamis cache folder

  plumbing extensions-remove
    Removes all extensions found in the salamis cache folder

  plumbing extensions-symlink
    For each tag of each workspace, symlink all extensions that match that tag to a subfolder of that workspace

  plumbing extensions-unsymlink
     Remove all created symlinks for each workspace

  plumbing xdg-install
    Install (write) all solamis.*.desktop entries

  plumbing xdg-remove
    Remove all solamis.desktop entries`)
}
