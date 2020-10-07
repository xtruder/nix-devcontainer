FROM xtruder/debian-nix-devcontainer:flakes

ARG DOCKER_GID=966
RUN groupadd -g ${DOCKER_GID} docker && usermod -a -G docker ${USERNAME} 
