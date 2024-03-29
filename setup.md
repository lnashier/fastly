# Project Setup

## Development

- Install [Go](https://golang.org/doc/install)

## Build & Run

### Standalone

```shell
brew install memcached
memcached -vv
```

```shell
cd web
CGO_ENABLED=0 go build -trimpath -a -ldflags '-s -w -extldflags=-static' -o ./bin/fastly
./bin/fastly
```

### Docker

- Install [Docker](https://docs.docker.com/get-docker/)

```shell
docker network create fastnet
```

```shell
docker pull memcached
docker run -d --net=fastnet --name=memcached -p 11211:11211 memcached '-vv'
```

```shell
cd $WORKSPACE/fastly
make
docker run -d --net=fastnet --name=fastly -e ENV=docker -p 8080:8080 fastly:1.0.0
```

## Samples

### Post content

```shell
curl --request POST 'http://localhost:8080/' \
  --form 'myfile=@"/dummy.bin"'

curl --request POST 'http://localhost:8080/' \
  --header 'Content-Type: text/plain' \
  --data-raw 'Hello, World!'
  
curl --request POST 'http://localhost:8080/' \
  --header 'Content-Type: application/octet-stream' \
  --data-binary '@/dummy.bin'
```

### Get content

```shell
curl --request GET 'http://localhost:8080/YOUR_KEY'
```
