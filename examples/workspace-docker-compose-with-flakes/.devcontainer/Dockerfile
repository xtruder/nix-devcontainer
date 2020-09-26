FROM xtruder/debian-nix-devcontainer:flakes

# create volume for direnv
RUN sudo -u user mkdir -p /home/${USERNAME}/.config/direnv/allow
VOLUME /home/${USERNAME}/.config/direnv/allow
