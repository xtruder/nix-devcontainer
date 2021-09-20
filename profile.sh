#!/bin/bash

PROJECT_DIR=${PROJECT_DIR:-/workspace}

# load nix into profile
if [ -e $HOME/.nix-profile/etc/profile.d/nix.sh ] ; then
    . $HOME/.nix-profile/etc/profile.d/nix.sh
fi

# whether running via vscode env probe
if shopt -q login_shell && [[ "$BASH_EXECUTION_STRING" =~ $HOME/.vscode-server/bin/.*/node ]]; then
    # load direnv if avalible
    if [ -f "$PROJECT_DIR/.envrc" ]; then
        direnv allow $PROJECT_DIR
        eval "$(direnv exec $PROJECT_DIR direnv dump bash)"
    fi
fi
