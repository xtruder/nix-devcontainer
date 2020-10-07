{%- if cookiecutter.image -%}
FROM {{ cookiecutter.image}}
{%- else -%}
FROM xtruder/debian-nix-devcontainer: {%- if cookiecutter.flakes == "y" -%}flakes{%- else -%}latest{%- endif -%}
{% endif %}