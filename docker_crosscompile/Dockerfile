FROM golang:1.9.2-alpine3.7

MAINTAINER Alex Soto <asotobu@gmail.com>

RUN apk update && apk upgrade && \
    apk add --no-cache git

RUN go get -u github.com/golang/dep/cmd/dep

ADD crossbuild.sh /usr/local/bin/crossbuild.sh
RUN chmod 755 /usr/local/bin/crossbuild.sh