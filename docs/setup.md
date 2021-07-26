# Project Setup

## Development

- Install [Go](https://golang.org/doc/install)
- Install [memcached](https://memcached.org/downloads)
- (optional) Install [Docker](https://docs.docker.com/get-docker/)

## Build & Run

### Standalone

```shell
$ cd $WORKSPACE/fastly/cmd/web
$ CGO_ENABLED=0 go build -a -ldflags '-w -extldflags=-static' -o fastly
$ ./fastly
```

### Docker

```shell
$ cd $WORKSPACE/fastly
$ make
$ docker run -p 8080:8080 fastly:<version>
```

## Samples

### Post content

```shell
curl POST 'http://localhost:8080/' \
  --form 'myfile=@"/dummy.bin"'

curl POST 'http://localhost:8080/' \
  --header 'Content-Type: text/plain' \
  --data-raw 'Hello, World!'
  
curl POST 'http://localhost:8080/' \
  --header 'Content-Type: application/octet-stream' \
  --data-binary '@/dummy.bin'
```

### Get content

```shell
curl --request GET 'http://localhost:8080/YOUR_KEY'
```