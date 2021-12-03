FROM xtruder/nix-devcontainer

VOLUME /nix

WORKDIR /workspace
COPY --chown=${USERNAME}:${USERNAME} test.sh shell.nix .envrc /workspace/
