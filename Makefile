.PHONY: build
build:
	docker build -t xtruder/nix-devcontainer:latest -f src/Dockerfile src

.PHONY: test
test: build
	docker build -t nix-devcontainer-test \
		--build-arg USERNAME=user \
		--build-arg USER_UID=1001 \
		--build-arg USER_GID=1001 \
		test -f test/Dockerfile

	docker run -ti -e PRELOAD_EXTENSIONS="test" nix-devcontainer-test \
		bash -ic '$${PROMPT_COMMAND}; . /workspace/test.sh; pkill -TERM ext-preloader'
