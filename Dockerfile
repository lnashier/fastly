FROM golang:1.16.6-alpine3.14 as build

ADD . /go/src/github.com/lnashier/fastly
WORKDIR /go/src/github.com/lnashier/fastly

ARG VERSION
ARG GOOS="linux"
ARG GOARCH="amd64"

RUN cd web \
    && GOOS=${GOOS} GOARCH=${GOARCH} CGO_ENABLED=0 \
      go build -a -o fastly \
      -ldflags "-s -w -extldflags \"-static\"" *.go \
    && chmod 0755 fastly

FROM alpine:3.14

COPY --from=build /go/src/github.com/lnashier/fastly/web/fastly /fastly
COPY --from=build /go/src/github.com/lnashier/fastly/web/configs /configs

WORKDIR /

EXPOSE 8080
ENTRYPOINT ["/fastly"]
CMD [""]
