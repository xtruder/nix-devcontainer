{%- if cookiecutter.image -%}
FROM {{ cookiecutter.image}}
{%- else -%}
FROM xtruder/debian-nix-devcontainer: {%- if cookiecutter.flakes == "y" -%}flakes{%- else -%}latest{%- endif -%}
{% endif %}

# create volume for direnv
RUN sudo -u user mkdir -p /home/${USERNAME}/.config/direnv/allow
VOLUME /home/${USERNAME}/.config/direnv/allow
