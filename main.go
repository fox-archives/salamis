package main

import (
	"log"
	"os"
	"path/filepath"
)

// Options User-customizable
type Options struct {
	ConfigFile      string
	ExtensionsDir   string
	WorkspaceDir    string
	ApplicationsDir string
}

func main() {
	configDir, err := os.UserConfigDir()
	handle(err)

	cacheDir, err := os.UserCacheDir()
	handle(err)

	dataDir := os.Getenv("XDG_DATA_HOME")
	if dataDir == "" {
		home, err := os.UserHomeDir()
		handle(err)

		dataDir = filepath.Join(home, ".local", "share", "applications")
	} else {
		dataDir = filepath.Join(dataDir, "applications")
	}

	opts := Options{
		ConfigFile:      filepath.Join(configDir, "salamis", "extensions.toml"),
		ExtensionsDir:   filepath.Join(cacheDir, "salamis", "extensions"),
		WorkspaceDir:    filepath.Join(cacheDir, "salamis", "workspaces"),
		ApplicationsDir: filepath.Join(dataDir),
	}

	args := os.Args[1:]

	if contains(args, "--help") || contains(args, "-h") {
		printHelp()
		os.Exit(0)
	}

	ensureLength(args, 1, "Must specify command. Exiting")

	command := args[0]
	switch command {
	case "init":
		doInit(opts)
		break

	case "update":
		doExtensionsRemove(opts)
		doExtensionsInstall(opts)
		doExtensionsUnsymlink(opts)
		doExtensionsSymlink(opts)
		doXdgRemove(opts)
		doXdgInstall(opts)
		break

	case "list":
		doList(opts)
		break

	case "edit":
		doEdit(opts)
		break

	case "check":
		doCheck(opts)
		break

	case "launch":
		ensureLength(args, 2, "Error: Must pass in a workspace name")
		workspaceName := args[1]
		doLaunch(opts, workspaceName)
		break

	case "plumbing":
		ensureLength(args, 2, "Error: Must pass in a 'plumbing' subcommand")
		command = args[1]

		switch command {
		case "extensions-install":
			doExtensionsInstall(opts)
			break

		case "extensions-remove":
			doExtensionsRemove(opts)
			break

		case "extensions-symlink":
			doExtensionsSymlink(opts)
			break

		case "extensions-unsymlink":
			doExtensionsUnsymlink(opts)
			break

		case "xdg-install":
			doXdgInstall(opts)
			break

		case "xdg-remove":
			doXdgRemove(opts)
			break

		default:
			log.Fatalln("Unknown subcommand. Exiting")
			break
		}
		break

	default:
		log.Fatalln("Unknown command. Exiting")
		break
	}
}
