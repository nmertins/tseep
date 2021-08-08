# t👀p

`tseep` monitors TCP connections on the host and prints new connections every 10 seconds. When `tseep` starts, all existing connections will be reported as new as part of the initialization process.

## Building

Download and install [Go](https://golang.org/doc/install). Run `go build` from the root of the repository.

```shell
git clone https://github.com/nmertins/tseep.git
cd tseep
make build
```

## Running

After building, run the executable `tseep`. Press Ctrl + C to stop execution.

```shell
./bin/tseep-linux-amd64
```