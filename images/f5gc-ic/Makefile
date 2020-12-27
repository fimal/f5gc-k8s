## add pipefail
SHELL=/bin/bash -o pipefail

GOCMD=GOARCH=amd64 GOOS=linux CGO_ENABLED=0 go

## version information
GIT_COMMIT ?= $(shell git describe --dirty --always --match '')
BRANCH_NAME ?= $(shell git rev-parse --abbrev-ref HEAD)

BINARY_OUT ?= build/ic

TEST_OPTS := -timeout 30s

## disable caching
ifneq ($(NO_CACHE),)
	TEST_OPTS+= -count=1
	REBUILD=-a
endif

## verify that GOPRIVATE is set
ifeq ($(GOPRIVATE),)
 $(error go private not set)
endif

## targets

ic: build

.PHONY: build
build: nofmt
	$(GOCMD) build ${REBUILD} -v -o ${BINARY_OUT}


# make sure no files need formatting
.PHONY: nofmt
nofmt:
	go env
	echo "gofmt - making sure no files need formatting"
	test -z $(shell gofmt -l .)

.PHONY: build-container
build-container:
	docker build -t ic .
