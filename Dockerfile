FROM golang:1.14-alpine3.12 AS build-env

RUN apk update && apk upgrade && \
   apk add --no-cache bash git pkgconfig gcc g++ libc-dev ca-certificates

RUN update-ca-certificates

ENV GO111MODULE=on
ENV GOPROXY=https://proxy.golang.org

ENV TZ Europe/Amsterdam

WORKDIR /go/src/app

ADD . /go/src/app

# Because of how the layer caching system works in Docker, the go mod download
# command will _ only_ be re-run when the go.mod or go.sum file change
# (or when we add another docker instruction this line)
RUN go mod download

# set crosscompiling fla 0/1 => disabled/enabled
ENV CGO_ENABLED=0
# compile linux only
ENV GOOS=linux
# run tests
RUN go test ./... -covermode=atomic

RUN go build -v -ldflags='-s -w -linkmode auto' -a -installsuffix cgo -o /atom atom.go

FROM scratch

# important for time conversion
ENV TZ Europe/Amsterdam

WORKDIR /

# Import from builder.
COPY --from=build-env /atom /
COPY --from=build-env /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/

CMD ["atom"]