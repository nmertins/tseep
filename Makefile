ARCH := $(shell uname -m)

build:
	@GOOS=linux GOARCH=arm64 go build -o bin/tseep-linux-arm64 cmd/tseep/main.go
	@GOOS=linux GOARCH=amd64 go build -o bin/tseep-linux-amd64 cmd/tseep/main.go

clean:
	@rm -rf bin

install:
ifeq ($(ARCH), x86_64)
	@cp bin/tseep-linux-amd64 /usr/local/bin/tseep
endif
ifeq ($(ARCH), aarch64)
	@cp bin/tseep-linux-arm64 /usr/local/bin/tseep
endif

uninstall:
	@rm -f /usr/local/bin/tseep

test:
	@go test -v

all: build