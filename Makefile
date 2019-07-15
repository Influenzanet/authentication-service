.PHONY: build test test-short docker

DOCKER_OPTS ?= --rm

help:
	@echo "Service building targets"
	@echo "	 build : build service command"
	@echo "  test  : run test suites"
	@echo "  test-short : run short version of test suites"
	@echo "  docker: build docker image"
	@echo "Env:"
	@echo "  DOCKER_OPTS : default docker build options (default : $(DOCKER_OPTS))"
	@echo "  TEST_ARGS : Arguments to pass to go test call"

build:
	go build .

test:
	JWT_TOKEN_KEY="Izit4dhcUXUJN3Shy/qUQgq/44/DyyWVrVcUIVK/QMo=" TOKEN_EXPIRATION_MIN="10" TOKEN_MINIMUM_AGE_MIN="2" go test $(TEST_ARGS)

test-short: TEST_ARGS += -short
test-short: test

docker:
	docker build $(DOCKER_OPTS) .