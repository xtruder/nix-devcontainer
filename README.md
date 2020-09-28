# debian-nix-devcontainer

## Summary

VSCode devcontainer that uses `debian` as base system and `nix-shell` for
management of development environment

| Metadata | Value |  
|----------|-------|
| *Contributors* | Jaka Hudoklin <jaka@x-truder.net> (github.com/offlinehacker) |
| *Definition type* | standalone or Docker Compose |
| *Works in Codespaces* | Yes |
| *Container host OS support* | Linux, macOS, Windows |
| *Languages, platforms* | All languages that nix supports |
| *Base image* | xtruder/debian-nix-devcontainer |
| *Image tags* | latest, flakes |

## Description

Debian nix devcontainer is [vscode devcontainer](https://code.visualstudio.com/docs/remote/containers)
that uses debian as a base system and nix for project development environment.
Idea is that system tools are managed by **apt** and development environment is
managed by [nix-shell](https://nixos.org/manual/nix/stable/#sec-nix-shell).

### Goals:

- **Keep it simple**

   Gool is to keep devcontainer simple, so it is easy to maintain. I really
   dislike how microsoft is [hacking up development containers](https://github.com/microsoft/vscode-dev-containers/blob/master/containers/go/.devcontainer/base.Dockerfile) by smashing together a bunch of shell scripts. At the same time i do not like the nix docker image building
   facility, which is again [bunch of bash scripts](https://github.com/NixOS/nixpkgs/blob/master/pkgs/build-support/docker/default.nix)
   producing a docker image.

   Less is more, less hacks means more fun.

- **Be more compatible**

   Be as compatible with different systems as possible. Whether using with
   devcontainers or with codespaces, or just using docker-compose locally.
   It should be easy to transition.

- **Use right tool for the job**

   I don't want to have multiple devcontainers for different programming
   languages. I think that project environment should be managed by
   independent tool, that you can run independent whether you develop your
   project in containers, vms or on physical machine.

   [nix-shell](https://nixos.org/manual/nix/stable/#sec-nix-shell) provides
   reproducible development environments and can be used across many
   languages.

- **Keep good productivity**

   While pure development environments provide much greater guarantees,
   they can provide burden for a developer.

### Features

- minimal debian as base system
- [nix](https://nixos.org/) for development environment with cached nix
  store even between container rebuilds (using docker volumes)
- [nix-shell](https://nixos.org/manual/nix/stable/#sec-nix-shell) for
  bootstraping of nix development environment
- [direnv](https://direnv.net/) for integration of nix-shell with your shell
- [nix-environment-selector](https://marketplace.visualstudio.com/items?itemName=arrterian.nix-env-selector)
  **vscode extension** for running vscode inside of nix shell environment.
- optional [docker-compose](https://docs.docker.com/compose/)
  (instead of using a single devcontainer)
- optional [nix-flakes](https://www.tweag.io/blog/2020-05-25-flakes/)
  support for better determinism

## Adding devcontainer definition into your project

If this is your first time using a development container, please follow the
[getting started steps](https://aka.ms/vscode-remote/containers/getting-started)
to set up your machine.

This project provides several [examples](examples/) of devcontainers for
different use cases and also a [template](template/) from which you can
generate required devcontainer definitions to use in your project.

### Examples

- **simple-project-with-flakes**
- **workspace-docker-compose-with-flakes**

### Template

This project contains [cookiecutter](https://github.com/cookiecutter/cookiecutter) template
with which you can quickly bootstrap devcontainer setup in your project.

Using template is simple as:

```shell
cookiecutter --directory template https://github.com/xtruder/debian-nix-devcontainer.git
```

Example usage:

```
cookiecutter --directory template https://github.com/xtruder/debian-nix-devcontainer.git
project_slug [project]: 
image []: 
workspace [n]: 
flakes [y]: 
compose [y]: 
niv [n]: 
nixpkgs_branch [nixos-20.09]
```

cookiecutter asks for several inputs required for generation:

- **project_slug**: name of the project and also name of output direcory where files are generated
- **image**: base image to use, by default uses `debian:latest` as base image
- **workspace**: wheter to generate workspaced project
- **flakes**: whether to use experimental [nix-flakes](https://www.tweag.io/blog/2020-05-25-flakes/)
- **compose**: whether to use `docker-compose` instead of single devcontainer
- **niv**: whether to use [niv](https://github.com/nmattia/niv) for nix dependency managment
- **nixpkgs_branch**: nixpkgs branch to use (to use unstable choose: `nixos-unstable`)

If you have existing project and would like to generate files into that project, select
`project-slug` with name of your project. You might need to force generate if you have existing
file and want to overwrite them.

## Running

Whether running an example or using a template, the setup should work out of the box. Open vscode press <kbd>F1</kbd>
and choose **Remote-Containers: Reopen Folder in Container**, vscode will build your devcontainer and you should be
ready to go. Also make sure **nix-env-selector** has set the correct environment set. In the status bar of vscode
you should see **Environment: shell.nix**.

When **opening terminal** in vscode you will have to use `direnv allow` command to allow direnv to
load nix shell. This is only required first time, since nix direnv allows are cached.

If opening a workspace, please make sure you open workspace in devcontainer and not a project. VSCode
does not currently provide a popup for opening a workspace in devcontainer, but will only provide you
with an option to open a project in devcontainer or to open workspace locally, but this is not what you want.

If `uid` or `gid` of user running `vscode` is not equal to `1000:1000`, you will have to set `USER_UID` and
`USER_GID` environment variables. This should be set globally or in a shell you are running `code` command from.
`vscode` inherits environment from where it is started, but make sure these environemnt variables are set or
files will have invalid permissions. If you change your local user `uid` or `gid` while already running a project
in devcontainer, you will have to remove all named volumes associated with devcontainer and also do rebuild of
devcontainer.

### Updating your nix shell environment

After updating `shell.nix` or `flake.nix` (whether you are using nix flakes or not) with changes that
affect your development environment, you will need to reload `nix-env-selector` environemnt by clicking
`Environment: shell.nix` in vscode status bar. This will reload nix shell environemnt for vscode.
After that you will need to reload vscode, but that should be fast operation, since vscode will reuse
already running container. Your open terminals should auto reload nix shell, due `direnv`.

### Adding another service

If you need to add another service, please make sure you are using `docker-compose` setup and not a simple
`devcontainer`, since it does not support running multiple services.

You can add other services to your `docker-compose.yml` file [as described in Docker's documentaiton](https://docs.docker.com/compose/compose-file/#service-configuration-reference). However, if you want anything running in this service to be available
in the container on localhost, or want to forward the service locally, be sure to add this line to the service config:

```
# Runs the service on the same network as the app container, allows "forwardPorts" in devcontainer.json function.
network_mode: service:app
```

This will make sure your devcontainer and service container are running in same network namespace.

## How it works

**[If the definition provides a pattern you think will be useful for others, describe the it here.]**

## License

Copyright (c) X-Truder. All rights reserved.

Licensed under the MIT License. See [LICENSE](https://github.com/xtruder/debian-nix-devcontainer/blob/master/LICENSE).
