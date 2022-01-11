{ pkgs ? import (fetchTarball "https://github.com/NixOS/nixpkgs/archive/nixos-21.11.tar.gz") { } }:

pkgs.mkShell {
  # nativeBuildInputs is usually what you want -- tools you need to run
  nativeBuildInputs = with pkgs; [
    # needed by nix
    nixpkgs-fmt
    rnix-lsp

    # for building docker images
    docker-client
    gnumake
    git

    # go dev
    go
  ];
}