import os
import sys
import shutil

REMOVE_PATHS = [
    '{% if cookiecutter.compose != "y" %}.devcontainer/docker-compose.yml{% endif %}',
    '{% if cookiecutter.workspace != "y" %}{{cookiecutter.project_slug}}.code-workspace{% endif %}',
    '{% if cookiecutter.workspace == "y" %}.vscode/settings.json{% endif %}',
    '{% if cookiecutter.flakes != "y" %}flake.nix{% endif %}',
    '{% if cookiecutter.niv != "y" %}nix{% endif %}'
]

for path in REMOVE_PATHS:
    path = path.strip()
    if path and os.path.exists(path):
        if os.path.isdir(path):
            shutil.rmtree(path)
        else:
            os.unlink(path)