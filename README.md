# nix-devcontainer

![workflow status](https://github.com/xtruder/nix-devcontainer/actions//workflows/ci.yml/badge.svg)


## Summary

**Swiss army knife container for vscode development environments**

| Metadata                    | Value                                                                         |
| --------------------------- | ----------------------------------------------------------------------------- |
| *Image*                     | [ghcr.io/xtruder/nix-devcontainer](https://github.com/xtruder/nix-devcontainer/pkgs/container/nix-devcontainer) |
| *Image tags*                | v1,latest,edge                                                                |
| *Definition type*           | standalone or Docker Compose                                                  |
| *Works in Codespaces*       | Yes                                                                           |
| *Container host OS support* | Linux, macOS, Windows                                                         |
| *Languages, platforms*      | All languages that nix supports                                               |
| *Contributors*              | [@offlinehacker](github.com/offlinehacker), [@Rizary](github.com/rizary)      |
| *Maintainer*                | Jaka Hudoklin <jaka@x-truder.net> [@offlinehacker](github.com/offlinehacker)  |

## Description

Nix devcontainer is an opinionated [vscode devcontainer](https://code.visualstudio.com/docs/remote/containers)
that uses debian image for a base system and [nix package manager](https://nixos.org/)
for management of your development environments. Combination
of a good base image and a best in class package manager, gives
you versatile, reproduible and deterministic development environment
that you can use everywhere.

### Components

- **[Debian slim](https://hub.docker.com/_/debian) docker image**

   Docker base image, which provides minimalistic environment in which both nix and
   vscode remote extension can run without any issues.

- **[Nix package manager](https://nixos.org/)**

  Used for providing declarative, deterministic and reporoducible development environment that
  you can run anywhere. Imagine better [conda](https://docs.conda.io/en/latest/) alternative.
  You write a single `shell.nix` file that describes your environment and all your tools and
  vscode extensions will be running with same exact versions of binaries, same environment
  variables, same libraries forever.

- **[Nix Environment Selector](https://marketplace.visualstudio.com/items?itemName=arrterian.nix-env-selector) vscode extension**

   Used to automatically load nix development environment into your vscode and provides
   capabilities to reload it later when environment changes, without having to rebuild
   docker image from scratch on every change.

- **[Direnv](https://direnv.net/) shell environment loader**

   While nix environment loader extension loads environment for vscode, you
   want `direnv` to manage you shell environment. `Direnv` loads nix environment
   (defined by `shell.nix` file) into your shell and reloads it automatically when it changes,
   keeping your environment fresh.

## Example templates

There are sevaral example templates you can use to quickly bootstrap your project:

- [nix-devcontainer-golang](https://github.com/xtruder/nix-devcontainer-golang/)
  
  Example project using `nix-devcontainer` for golang development, with docker-compose
  running docker-in-docker service for building docker images.

- [nix-devcontainer-python-jupyter](https://github.com/xtruder/nix-devcontainer-python-jupyter/)
  
  Example project using `nix-devcontainer` for python and jupyter notebooks,
  with python packages managed by nix.

## Adding devcontainer definition into your project

If this is your first time using a development container, please follow the
[getting started steps](https://aka.ms/vscode-remote/containers/getting-started)
to set up your machine, install docker and vscode remote extensions.

### Project setup

Make sure that your project has `shell.nix` that describes your development
environment. Internally nix environment selector vscode extension runs
`nix-shell` to configure vscode's development environment.

Here is minimal example of `shell.nix` to get you started:

```nix
{ pkgs ? import <nixpkgs> { } }:

pkgs.mkShell {
  # nativeBuildInputs is usually what you want -- tools you need to run
  nativeBuildInputs = with pkgs; [
     #hello
  ];
}
```

If you want `nix-shell` to automatically run for your shell environemnts
running in your development container, create `.envrc` file.

A minimal example of `.envrc` file:

```shell
use_nix
```

For more informattion on how to develop with `nix-shell` you can take a
look here: https://nixos.wiki/wiki/Development_environment_with_nix-shell

### Devcontainer integration

Integrating `devcontainer` into your project is as simple as creating `devcontainer.json` file and a
`Dockerfile` in `.devcontainer` directory.

Example `.devcontainer/devcontainer.json`:

```jsonc
// For format details, see https://aka.ms/vscode-remote/devcontainer.json or the definition README at
// https://github.com/microsoft/vscode-dev-containers/tree/master/containers/docker-existing-dockerfile
{
  "name": "devcontainer-project",
  "dockerFile": "Dockerfile",
  "context": "${localWorkspaceFolder}",
  "build": {
    "args": {
      "USER_UID": "${localEnv:USER_UID}",
      "USER_GID": "${localEnv:USER_GID}"
    },
  },

  // run arguments passed to docker
  "runArgs": [
    "--security-opt", "label=disable"
  ],

  "containerEnv": {
     // extensions to preload before other extensions
    "PRELOAD_EXTENSIONS": "arrterian.nix-env-selector"
  },

   // disable command overriding and updating remote user ID
  "overrideCommand": false,
  "userEnvProbe": "loginShell",
  "updateRemoteUserUID": false,

  // build development environment on creation, make sure you already have shell.nix
  "onCreateCommand": "nix-shell --command 'echo done building nix dev environment'",

  // Add the IDs of extensions you want installed when the container is created.
  "extensions": [
    // select nix environment
    "arrterian.nix-env-selector",

    // extra extensions
    //"fsevenm.run-it-on",
    //"jnoortheen.nix-ide",
    //"ms-python.python"
  ],

  // Use 'forwardPorts' to make a list of ports inside the container available locally.
  "forwardPorts": [],

  // Use 'postCreateCommand' to run commands after the container is created.
  // "postCreateCommand": "go version",
}
```

Example `.devcontainer/Dockerfile`:

```dockerfile
FROM ghcr.io/xtruder/nix-devcontainer:v1
```

**Dockerfile is needed for build triggers to run.** Build triggers will change
user `uid` and `gid` to one provided by `USER_UID` and `USER_GID` env variables
and change ownersip of `/nix` and `/home` folders. This is required, as
docker currently does not provide a way to map filesystem uids/gids.
If you don't need this and your host user always has `1000:1000` uid/gid,
you can also specify image directly by setting `image` parameter
in `devcontainer.json`

If you already have your `shell.nix`, you can also set to use it in your
project `.vscode/settings.json` file:

```json
{
    "nixEnvSelector.nixFile": "${workspaceRoot}/shell.nix",
}
```

### Using docker-compose instead

Alternatively you can use `docker-compose` instead. This allows you to run multiple
services and have more control over development environment. In this case you need
to specify path to compose file in your `devcontainer.json` file. Compose file
can also be in your project root if you prefer.

Example `.devcontainer/devcontainer.json`:

```jsonc
// For format details, see https://aka.ms/vscode-remote/devcontainer.json or the definition README at
// https://github.com/microsoft/vscode-dev-containers/tree/master/containers/docker-existing-dockerfile
{
  "name": "devcontainer-project",
  "dockerComposeFile": "docker-compose.yml",
  "service": "dev",
  "workspaceFolder": "/workspace",
  
  "userEnvProbe": "loginShell",
  "updateRemoteUserUID": false,

  // build development environment on creation
  "onCreateCommand": "nix-shell --command 'echo done building nix dev environment'",

  // Add the IDs of extensions you want installed when the container is created.
  "extensions": [
    // select nix environment
    "arrterian.nix-env-selector",

    // extra extensions
    //"fsevenm.run-it-on",
    //"jnoortheen.nix-ide",
    //"ms-python.python"
  ],
}
```

Example `.devcontainer/docker-compose.yml`:

```yaml
version: '3'
services:
  dev:
    build:
      context: ../
      dockerfile: .devcontainer/Dockerfile
      args:
        USER_UID: ${USER_UID:-1000}
        USER_GID: ${USER_GID:-1000}
    environment:
      # list of docker extensions to load before other extensions
      PRELOAD_EXTENSIONS: "arrterian.nix-env-selector"
    volumes:
      - ..:/workspace:cached
    security_opt:
      - label:disable
    network_mode: "bridge"
```

## Running

When you open a project vscode should ask you to open project in remote container.
If you click to open in remove container, dev environment should be built
automatically and after you should be ready to start coding.

Alternativelly you can choose **Remote-Containers: Reopen Folder in Container** from command menu.

If opening a workspace, please make sure you open workspace in devcontainer by selecting
**Remote-Containers: Open workspace in Container". Remote containers extension does not
currently provide a popup for opening a workspace, but will only provide you with an option
to open a project in devcontainer or to open workspace locally, but this is not what you want.

If running under linux and `uid` or `gid` of user running `vscode` are not equal to `1000:1000`,
you will have to set `USER_UID` and `USER_GID` environment variables.
These should be set globally or in a shell you are running `code` command from. `vscode`
inherits environment from where it is started, but make sure these environment variables
are set or files will have incorrect permissions.

This is a list of `nix-devcontainer` image build arguments and its defaults:

| Name     | Description                       | Default |
| -------- | --------------------------------- | ------- |
| USERNAME | Username of the user in container | code    |
| USER_UID | ID of the user in container       | 1000    |
| USER_GID | Group ID of the user in container | 1000    |

### Rebuilding your development environment

After updating `shell.nix` with changes that affect your development environment.
If you are using nix environment selector extension you can choose `Nix-Env: Hit environment`
from command menu and it will rebuild your environment and after that it will ask you
to reload your editor.

Alternatively you can also reload your editor by choosing `Developer: Rebuild Container`
and it will rebuild your environment on editor start.

### Atomatically rebuilding environment on environment changes

If you want to automatomatically rebuild your environment when `shell.nix` changes
you can use `fsevenm.run-it-on` extension and choose to automatically run
`nixEnvSelector.hitEnv` command. Here is an example of settings you can put in your
workspace settings (`.vscode/settings.json` file in your project.)

```jsonc
{
    "runItOn": {
        "commands": [
            {
                "match": "flake\\.nix",
                "isShellCommand": false,
                "cmd": "nixEnvSelector.hitEnv"
            }
        ]
    }
}
```

### Adding personalized dotfiles

VSCode remote containers have internal support for cloning dotfiles repo and running
custom install script when devcontainer is opened.
Check how to personalize your environment [here](https://code.visualstudio.com/docs/remote/containers#_personalizing-with-dotfile-repositories).
An example of dotfiles repository is available [here](https://github.com/offlinehacker/dotfiles).

### Adding another service

If you need to add another service, please make sure you are using `docker-compose` and not a plain
`devcontainer.json`, since it does not support running multiple services.

You can add other services to your `docker-compose.yml` file [as described in Dockers documentation](https://docs.docker.com/compose/compose-file/#service-configuration-reference). However, if you want anything running in these service to be available
in the dev container on localhost, or want to forward the service locally,
be sure to add this line to the `docker-compose` service config:

```yaml
# Runs the service in the same network namespace as the dev container,
# keeping services on localhost makes life easier
network_mode: service:app
```

This will make sure your dev container and service containers are running in same network namespace.

### Caching nix store

Caching nix store is as simple as adding named docker volume on `/nix`.

```dockerfile
FROM ghcr.io/xtruder/nix-devcontainer:v1
VOLUME /nix
```

Then add named docker volume to `devcontainer.json` or `docker-compose.yml` as explained
in next paragraph.

### Caching additional directories using docker volumes

You can cache additional directories using docker volumes. You will need to make sure volumes
have right permissions. To do that you need to first create directory in your `Dockerfile` and
set it as a volume and later put named volume in your `docker-compose.yml` or in `devcontainer.json`.

1. Create directory and add volume to your `Dockerfile`

   ```dockerfile
   RUN mkdir -p /home/${USERNAME}/.cache
   VOLUME /home/${USERNAME}/.cache
   ```

1. Add named volume to `devcontainer.json` or `docker-compose.yml`:

   ```jsonc
   {
      "name": "project-name",
      //...,
      "mounts: [
         "source=project-name_devcontainer_home-cache,target=/home/user/.cache,type=volume"
      ]
   }
   ```

   Alternatively If you are using `docker-compose` you can add the following to
   your `docker-compose.yml` in `.devcontainer` directory. This will mount named
   volume to your desired location.

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

### docker-in-docker or how to use `docker` inside devcontainer

There are several ways how to use docker inside devcontainer. The easiest is just
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

   ```jsonc
   {
     //...,
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

**Be aware that exposing docker socket to your development environment is a security risk, as it
exposes your system to potentially malicious development environment. Make sure you never run untrusted code in such environment.**

Better alternative is to run rootless docker in docker as separate privileged service via `docker-compose`. While this service runs
in privileged container, it mittigates the risk by running docker
without root.

```yaml
version: '3'
services:
  # your development container
  dev:
    ...

    # it's advised to use shared network namespace, as it simplifies
    # development and can allow connections to docker via localhost
    network_mode: "service:docker"

  docker:
    image: docker:dind-rootless
    environment:
      DOCKER_TLS_CERTDIR: ""
      DOCKER_DRIVER: overlay2
    privileged: true
    volumes:
      - ..:/workspace:cached
      - docker:/var/lib/docker
    security_opt:
      - label:disable
    network_mode: bridge

volumes:
  docker:
```

**Only linux kernels of version 5.11+ support overlay2 storage driver in rootless containers. You can use default vfs or fuse-overlayfs drivers, but both are a bit slow, so they are not recommended**

### Preload vscode extensions (run nix env selector before other extensions)

Some vscode extensions have issue that development environment is loaded too
late. Currently vscode does not support an option to make extensions load after
some other extension, but only supports extension dependencies, where one extension
can wait for other extension to load, see also https://github.com/microsoft/vscode/issues/57481.

To workaround this issue we have implemented a hack, a vscode extension preloader
that modifies extensions `package.json` on the fly to make extensions depend on
other set of extensions while they are installed, but before they are loaded.

You can enable this feature/hack by setting `PRELOAD_EXTENSIONS` environment
variable in your `devcontainer.json` or `docker-compose.yml`.

Example `devcontainer.json` settings:

```jsonc
{
  //...
  "containerEnv": {
    "PRELOAD_EXTENSIONS": "arrterian.nix-env-selector"
  },
  //...

  "extensions": [
    "arrterian.nix-env-selector",
    //...
  ]
}
```

### Using with podman

Running this devcontainer with Podman is currently not supported, due to vscode remote
containers not supporting passing build flags to `podman` and `podman-compose`.
See also https://github.com/microsoft/vscode-remote-release/issues/3545, there is also
a possible workaround, but i haven't tried it yet.

## Technical details

### How environment is loaded?

- Vscode starts devcontainer from `devcontainer.json` definition
- `.devcontainer/Dockerfile` (that inherits from this image) is
  used to build development container, several `ONBUILD` docker triggers
  are run to change user `uid` and `gid` of non-root user (`code` by default)
  in container and fix ownership of files.
- Upon start vscode installs extensions defined in `devcontainer.json`, including
  `arrterian.nix-env-selector`.
- `arrterian.nix-env-selector` extension evaluates development shell defined in
  `shell.nix` file and sets vscode environment variables based on that environment.
- All other extensions are loaded into vscode.
- When you start vscode terminal, environment variables are not automatically inherited
  from `vscode`, that's why this devcontainer sets `direnv` shell hook to automatically
  load environment into your shell.

## Development

It's recommended to open this project in vscode devcontainer or via github
codespaces. This will automatically prepare development environment with
all required dependencies.

### Project structure

- `examples` - devcontainer examples
- `src` - image source
- `test` - image smoke test

### Testing

For basic validity testing you should run `make test`, which will build test
image and runs tests in image. Sanity checks are run by first running
direnv hook which loads nix environment and then it sources `test.sh` script,
which does a few basic sanity checks.

You should also check if image works with example templates.

## License

Copyright (c) X-Truder. All rights reserved.

Licensed under the MIT License. See [LICENSE](https://github.com/xtruder/nix-devcontainer/blob/master/LICENSE).
