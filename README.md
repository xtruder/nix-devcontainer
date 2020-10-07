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

   Goal is to keep devcontainer simple, so it is easy to maintain. I really
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
  bootstrapping of nix development environment
- [direnv](https://direnv.net/) for integration of nix-shell with your shell
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

- [simple-project](examples/simple-project)

   A simple devcontainer project using nix shell and direnv for bootstrapping of
   dev environment.

- [simple-project-with-niv](examples/simple-project-with-niv)

   A simple devcontainer project using nix shell and direnv for bootstrapping of
   dev environment. Additionally it uses [niv](https://github.com/nmattia/niv) for
   nix dependency management. 

- [simple-project-with-flakes](examples/simple-project-with-flakes)

   A simple devcontainer project using nix shell and direnv for bootstrapping of
   dev environment. Additionally it uses nix flakes for nix dependency management. 

- [workspace-docker-compose-with-flakes](examples/workspace-docker-compose-with-flakes)

   Workspace devcontainer using docker-compose and nix flakes.

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

- **project_slug**: name of the project and also name of output directory where files are generated
- **image**: base image to use, by default uses `debian:latest` as base image
- **workspace**: whether to generate workspace project
- **flakes**: whether to use experimental [nix-flakes](https://www.tweag.io/blog/2020-05-25-flakes/)
- **compose**: whether to use `docker-compose` instead of single devcontainer
- **niv**: whether to use [niv](https://github.com/nmattia/niv) for nix dependency management
- **nixpkgs_branch**: nixpkgs branch to use (to use unstable choose: `nixos-unstable`)

If you have existing project and would like to generate files into that project, select
`project-slug` with name of your project. You might need to force generate if you have existing
file and want to overwrite them.

## Running

Whether running an example or using a template, the setup should work out of the box. Open vscode press <kbd>F1</kbd>
and choose **Remote-Containers: Reopen Folder in Container**, vscode will build your devcontainer and you should be
ready to go.

When **opening terminal** in vscode you will have to use `direnv allow` command to allow direnv to
load nix shell. This is only required first time, since nix direnv allows are cached.

If opening a workspace, please make sure you open workspace in devcontainer and not a project. VSCode
does not currently provide a popup for opening a workspace in devcontainer, but will only provide you
with an option to open a project in devcontainer or to open workspace locally, but this is not what you want.

If `uid` or `gid` of user running `vscode` is not equal to `1000:1000`, you will have to set `USER_UID` and
`USER_GID` environment variables. This should be set globally or in a shell you are running `code` command from.
`vscode` inherits environment from where it is started, but make sure these environment variables are set or
files will have invalid permissions. If you change your local user `uid` or `gid` while already running a project
in devcontainer, you will have to remove all named volumes associated with devcontainer and also do rebuild of
devcontainer.

### Reloading your nix shell environment

After updating `shell.nix` or `flake.nix` (whether you are using nix flakes or not) with changes that
affect your development environment, you will need to reload vscode by pressing <kbd>Ctrl+R</kbd> or by selecting
`Developer: Reload Window` in command menu.
This will restart vscode server running in container, which will probe for new environment by running
`direnv`.

### Adding personalized dotfiles

VSCode has internal for cloning dotfiles repo and running custom install script when devcontainer is opened.
Check how to personalize your environment [here](https://code.visualstudio.com/docs/remote/containers#_personalizing-with-dotfile-repositories).
An example of dotfiles repository is available [here](https://github.com/offlinehacker/dotfiles).

### Adding another service

If you need to add another service, please make sure you are using `docker-compose` setup and not a simple
`devcontainer`, since it does not support running multiple services.

You can add other services to your `docker-compose.yml` file [as described in Dockers documentation](https://docs.docker.com/compose/compose-file/#service-configuration-reference). However, if you want anything running in this service to be available
in the container on localhost, or want to forward the service locally, be sure to add this line to the service config:

```yaml
# Runs the service on the same network as the app container, allows "forwardPorts" in devcontainer.json function.
network_mode: service:app
```

This will make sure your devcontainer and service container are running in same network namespace.

### Caching additional directories using docker volumes

You can cache additional directories using docker volumes. You will need to make sure volumes
have right permissions. To do that you need to first create directory in your `Dockerfile` and
set it as a volume and later put named volume in your `docker-compose.yml` or in `devcontainer.json`.

1. Add `VOLUME` definition to `Dockerfile`

```Dockerfile
RUN sudo -u user mkdir -p /home/${USERNAME}/.cache
VOLUME /home/${USERNAME}/.cache
```

2. Add named volume to `docker-compose`:

If you are using `docker-compose` you can add the following to your `docker-compose.yml` in
`.devcontainer` directory. This will mount named volume to desired location.

```yaml
version: '3'
services:
   dev:
      ...
      volumes:
         - home-cache:/home/user/.cache
volumes:
   home-cache
```

3. Add named volume to `devcontainer.json` run arguments:

If not using `docker-compose`, you can achieve the same result by adding named volume to
your `devcontainer.json` mounts:

```json
{
   "name": "project-name",
   ...,
   "mounts: [
      "source=project-name_devcontainer_home-cache,target=/home/user/.cache,type=volume"
   ]
}
```

### Using `docker` inside devcontainer

There are several ways how to use docker inside devcontainer. The easiest would be just
to mount `docker.sock` in devcontainer:

1. Add mount in `docker-compose.yml` or `devcontainer.json`:

```yaml
version: '3'
services:
   dev:
      ...
      volumes:
         - type: bind
           source: /var/run/docker.sock
           target: /var/run/docker.sock
```

```json
{
   ...,
   "mounts": [
     "source=/var/run/docker.sock,target=/var/run/docker.sock,type=bind"
   ]
}
```

2. (Optional) Add user to docker group in your `Dockerfile`

```Dockerfile
ARG DOCKER_GID=966
RUN groupadd -g ${DOCKER_GID} docker && usermod -a -G docker ${USERNAME} 
```

**Exposing docker socket to your development environment is a security risk, as it
exposes your system to potentially malicious development environment. Make sure
you never run untrusted code in such environment.**

Better alternative is to run rootless docker in docker as separate privileged service
via `docker-compose`. That should mitigate most common security risks, but please be
aware 

## License

Copyright (c) X-Truder. All rights reserved.

Licensed under the MIT License. See [LICENSE](https://github.com/xtruder/debian-nix-devcontainer/blob/master/LICENSE).
