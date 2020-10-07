{%- if cookiecutter.flakes != "y" -%}
{ {% if cookiecutter.niv == "y" -%}
  sources ? import ./nix/sources.nix,
  pkgs ? import sources.nixpkgs {}
{%- else -%}
  src ? builtins.fetchTarball "https://github.com/NixOS/nixpkgs/archive/{{cookiecutter.nixpkgs_branch}}.tar.gz",
  pkgs ? import src {}
{%- endif -%} }:

pkgs.mkShell {
  buildInputs = with pkgs; [
  {%- if cookiecutter.niv == "y" %}
    niv
  {% endif -%}
  ];
}
{%- else -%}
(import (fetchTarball https://github.com/edolstra/flake-compat/archive/master.tar.gz) {
  src = ./.;
}).shellNix
{%- endif %}