# debian-nix-devcontainer

## Summary

VSCode devcontainer that uses `debian` as base system and `nix-shell` for
management of development environment

| Metadata | Value |  
|----------|-------|
| *Contributors* | Jaka Hudoklin <jaka@x-truder.net> (github.com/offlinehacker) |
| *Definition type* | Docker Compose |
| *Works in Codespaces* | Yes |
| *Container host OS support* | Linux, macOS, Windows |
| *Languages, platforms* | All languages that nix supports |

## Description

Debian nix devcontainer is [vscode devcontainer](https://code.visualstudio.com/docs/remote/containers)
that uses debian as a base system and nix for project development environment.
Idea is that system tools are managed by apt and development environment is
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

## Using this definition with an existing folder

**[Optional] Include any special setup requirements here. For example:**

This devcontainer should work out of the box. If you have different `uid`
or `gid` of your user on system, you will have to set `USER_UID` and
`USER_GID` global environment variables in environment that vscode is running
from.

### Adding another service

You can add other services to your `docker-compose.yml` file [as described in Docker's documentaiton](https://docs.docker.com/compose/compose-file/#service-configuration-reference). However, if you want anything running in this service to be available in the container on localhost, or want to forward the service locally, be sure to add this line to the service config:

# Runs the service on the same network as the app container, allows "forwardPorts" in devcontainer.json function.
network_mode: service:app

### Adding the definition to your project

1. If this is your first time using a development container, please follow the [getting started steps](https://aka.ms/vscode-remote/containers/getting-started) to set up your machine.

2. To use VS Code's copy of this definition:
   1. Start VS Code and open your project folder.
   2. Press <kbd>F1</kbd> select and **Remote-Containers: Add Development Container Configuration Files...** from the command palette.
   3. Select the **YOUR NAME HERE** definition.

3. To use latest-and-greatest copy of this definition from the repository:
   1. Clone this repository.
   2. Copy the contents of this folder in the cloned repository to the root of your project folder.
   3. Start VS Code and open your project folder.

4. After following step 2 or 3, the contents of the `.devcontainer` folder in your project can be adapted to meet your needs.

5. Finally, press <kbd>F1</kbd> and run **Remote-Containers: Reopen Folder in Container** to start using the definition.

## Testing the definition

This definition includes some test code that will help you verify it is working as expected on your system. Follow these steps:

1. If this is your first time using a development container, please follow the [getting started steps](https://aka.ms/vscode-remote/containers/getting-started) to set up your machine.
2. Clone this repository.
3. Start VS Code, press <kbd>F1</kbd>, and select **Remote-Containers: Open Folder in Container...**
4. Select this folder from the cloned repository.
5. Wait for devcontainer to get started.
6. Check if nix environment selector has used `shell.nix` for its environment.
7. Check if user id and group id inside container is same as user that vscode is running from.
8. Check if `nix-shell` works from terminal.

## [Optional] How it works

**[If the definition provides a pattern you think will be useful for others, describe the it here.]**

## License

Copyright (c) X-Truder. All rights reserved.

Licensed under the MIT License. See [LICENSE](https://github.com/Microsoft/vscode-dev-containers/blob/master/LICENSE).
