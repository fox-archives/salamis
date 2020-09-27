package main

import (
	"log"
	"os"
)

// Options User-customizable
type Options struct {
	ConfigFile    string
	ExtensionsDir string
	WorkspaceDir  string
}

func main() {
	opts := Options{
		ConfigFile:    "extensions.toml",
		ExtensionsDir: "extensions",
		WorkspaceDir:  "workspaces",
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
		doUpdate(opts)
		break

	case "check":
		doCheck(opts)
		break

	case "launch":
		ensureLength(args, 2, "Must pass in a workspace name")
		workspaceName := args[1]
		doLaunch(opts, workspaceName)
		break

	case "plumbing":
		ensureLength(args, 2, "Must pass in a 'plumbing' subcommand")
		command = args[1]

		switch command {
		case "download-extensions":
			doDownloadExtensions(opts)
			break

		case "remove-extensions":
			doRemoveExtensions(opts)
			break

		case "symlink-extensions":
			doSymlinkExtensions(opts)
			break

		case "remove-symlinks":
			doSymlinkRemove(opts)
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
