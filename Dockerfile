FROM golang:1.10.1-alpine3.7 as builder

RUN apk update && apk upgrade && \
    apk add --no-cache git

RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/lordofthejars/diferencia

COPY . .

RUN dep ensure
RUN GOOS=linux GOARCH=amd64 go build -o binaries/diferencia

FROM alpine:3.7

RUN addgroup -S diferencia && adduser -S -G diferencia diferencia 
USER diferencia

EXPOSE 8080

WORKDIR /home/diferencia
COPY --from=builder /go/src/github.com/lordofthejars/diferencia/binaries/diferencia .

ENTRYPOINT ["./diferencia"]