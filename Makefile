PROJECTNAME := $(shell basename "$(PWD)")
ARCH := $(shell uname -m)

INSTALL_PATH=/usr/local/bin/tseep
GOMAIN=cmd/tseep/main.go

all: help

## build: Build project binaries
build:
	@go get -d -v
	@go install -v
	@GOOS=linux GOARCH=arm64 go build -o bin/tseep-linux-arm64 $(GOMAIN)
	@GOOS=linux GOARCH=amd64 go build -o bin/tseep-linux-amd64 $(GOMAIN)

## clean: Delete binaries generated by build
clean:
	@rm -rf bin

## install: Install binary for your platform in /usr/local/bin
install:
ifeq ($(ARCH), x86_64)
	@cp bin/tseep-linux-amd64 $(INSTALL_PATH)
endif
ifeq ($(ARCH), aarch64)
	@cp bin/tseep-linux-arm64 $(INSTALL_PATH)
endif

## uninstall: Delete installed binary
uninstall:
	@rm -f $(INSTALL_PATH)

## test: Run automated tests
test:
	@go test -test.coverprofile "" -v

.PHONY: help
help: Makefile
	@echo
	@echo " Choose a command run in "$(PROJECTNAME)":"
	@echo
	@sed -n 's/^##//p' $< | column -t -s ':' |  sed -e 's/^/ /'
	@echo