version ?= latest

.PHONY: install
install:
	dep ensure
	go build -o diferencia

.PHONY: build
build:
	go build -o diferencia