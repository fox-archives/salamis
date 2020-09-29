function __fish_sparta_needs_command
    set cmd (commandline -opc)

    if [ (count $cmd) -eq 1 ]
        return 0
    end

    return 1
end

function __fish_sparta_no_subcommand -d 'Test if apt has yet to be given the subcommand'
    for i in (commandline -opc)
        if contains -- $i update upgrade dselect-upgrade dist-upgrade install remove purge source build-dep check clean autoclean changelog
            return 1
        end
    end
    return 0
end

complete -f -c sparta -a h -l help -d "Help"

complete -c sparta -n '__fish_sparta_no_subcommand' -a 'init' -d "Initiate"
complete -c sparta -n '__fish_sparta_no_subcommand' -a 'update' -d 'Update'
complete -c sparta -n '__fish_sparta_no_subcommand' -a 'check' -d 'Check'
complete -c sparta -n '__fish_sparta_no_subcommand' -a 'launch' -d 'Launch'
complete -r -c sparta -n '__fish_sparta_needs_command' -a 'plumbing' -d 'More plumbing commands'

complete -c sparta -n '__fish_sparta_no_subcommand' -a 'download-extensions' -d 'Download Extensions'
complete -c sparta -n '__fish_sparta_no_subcommand' -a 'remove-extensions' -d 'Remove Extensions'
complete -c sparta -n '__fish_sparta_no_subcommand' -a 'symlink-extensions' -d 'Symlink Extensions'
complete -c sparta -n '__fish_sparta_no_subcommand' -a 'remove-symlinks' -d 'Remove Symlinks'
