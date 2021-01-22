#
# shellcheck shell=bash
#

# based on https://tylerthrailkill.com/2019-01-19/writing-bash-completion-script-with-subcommands/

_log() {
	local -r file=".salamis-debug"
	[[ -v DEBUG ]] && {
		"$@" >> "$file"
	} 1>&2 2>/dev/null
}

_salamis_launch() {
	local cur="${COMP_WORDS[COMP_CWORD]}"
	local dirs

	# TODO
	dirs="$(\ls --color=never -c1 ~/.cache/salamis/workspaces)"
	# shellcheck disable=SC2207
	COMPREPLY=($(compgen -W "$dirs" -- "$cur"))
}

_salamis_plumbing() {
	local i=1 subcommand_index

	while [[ $i -lt $COMP_CWORD ]]; do
		local s="${COMP_WORDS[i]}"
		case "$s" in
			subcommand)
				# we set the current subcommand
				subcommand_index=$i
				break
				;;
		esac

		(( i++ ))
	done

	while [[ $subcommand_index -lt $COMP_CWORD ]]; do
		local s="${COMP_WORDS[subcommand_index]}"
		case "$s" in
		download-extensions)
			COMPREPLY=()
			return
			;;
		remove-extensions)
			COMPREPLY=()
			return
			;;
		symlink-extensions)
			COMPREPLY=()
			return
			;;
		remove-symlinks)
			COMPREPLY=()
			return
			;;
		esac
		(( subcommand_index++ ))
	done

	local cur="${COMP_WORDS[COMP_CWORD]}"
	# shellcheck disable=SC2207
	COMPREPLY=($(compgen -W "download-extensions remove-extensions symlink-extensions remove-symlinks" -- "$cur"))
}

_salamis() {
	local i=1 cmd

	# iterate over COMP_WORDS (ending at currently completed word)
	# this ensures we get command completion even after passing flags
	while [[ "$i" -lt "$COMP_CWORD" ]]; do
		local s="${COMP_WORDS[i]}"
		case "$s" in
		# if our current word starts with a '-', it is not a subcommand
		-*) ;;
		# we are completing a subcommand, set cmd
		*)
			cmd="$s"
			break
			;;
		esac
		(( i++ ))
	done

	# check if we're completing 'salamis'
	if [[ "$i" -eq "$COMP_CWORD" ]]; then
		local cur="${COMP_WORDS[COMP_CWORD]}"
		# shellcheck disable=SC2207
		COMPREPLY=($(compgen -W "init update check launch plumbing -h --help" -- "$cur"))
		return
	fi

	# if we're not completing 'salamis', then we're completing a subcommand
	case "$cmd" in
		init)
			COMPREPLY=()
			;;
		update)
			COMPREPLY=()
			;;
		check)
			COMPREPLY=()
			;;
		launch)
			_salamis_launch
			;;
		plumbing)
			_salamis_plumbing
			;;
		*)
			;;
	esac

	return 0
}

complete -F _salamis salamis
