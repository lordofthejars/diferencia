FROM golang:1.10.1-alpine3.7 as builder

RUN apk update && apk upgrade && \
    apk add --no-cache git

RUN go get -u github.com/golang/dep/cmd/dep

WORKDIR /go/src/github.com/lordofthejars/diferencia

COPY . .

RUN wget https://github.com/gobuffalo/packr/releases/download/v1.11.1/packr_1.11.1_linux_amd64.tar.gz
RUN tar -zxvf packr_1.11.1_linux_amd64.tar.gz 
RUN cp packr /usr/local/bin

RUN dep ensure
RUN GOOS=linux GOARCH=amd64 packr build -o binaries/diferencia

FROM alpine:3.7

RUN addgroup -S diferencia && adduser -S -G diferencia diferencia 
USER diferencia

EXPOSE 8080
EXPOSE 8081
EXPOSE 8082

WORKDIR /home/diferencia
COPY --from=builder /go/src/github.com/lordofthejars/diferencia/binaries/diferencia .

ENTRYPOINT ["./diferencia"]