# Project Setup

## Development

- Install [Go](https://golang.org/doc/install)
- (optional) Install [Docker](https://docs.docker.com/get-docker/)
- Install preferred IDE

## Build & Run

### Standalone
```
$ cd $WORKSPACE/fastly/cmd/web
$ CGO_ENABLED=0 go build -a -ldflags '-w -extldflags=-static' -o fastly
$ ./fastly
```

### Docker
```
$ cd $WORKSPACE/fastly
$ make
$ docker run -p 8080:8080 fastly:<version>
```
