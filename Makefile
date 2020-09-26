examples = simple-project-with-flakes workspace-docker-compose-with-flakes

.PHONY: cookiecutter
cookiecutter: clean $(examples)

clean:
	-rm -r examples/*

$(examples):
	cookiecutter --no-input --config-file .cookiecutter/$@.yaml -f -o examples template
	cd examples/$@; nix flake update

.PHONY: build
build:
	docker build -t xtruder/debian-nix-devcontainer:latest .

.PHONY: build-flakes
build-flakes:
	docker build \
		--build-arg NIX_INSTALL_SCRIPT=https://github.com/numtide/nix-flakes-installer/releases/download/nix-3.0pre20200820_4d77513/install \
		--build-arg EXTRA_NIX_CONFIG="experimental-features = nix-command flakes" \
		-t xtruder/debian-nix-devcontainer:flakes .

.PHONY: test
test:
	docker build --build-arg USER_UID= --build-arg USER_GID= -f test.Dockerfile .