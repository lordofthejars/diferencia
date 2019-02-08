version ?= latest

.PHONY: install
install:
	dep ensure
	packr build -o diferencia

.PHONY: test
test:
	go test -v -race $(go list ./... | grep -v "/vendor/")

.PHONY: build
build:
	packr build -o diferencia

.PHONY: crossbuild
crossbuild:
	docker run -it --rm -v "$$PWD":/go/src/github.com/lordofthejars/diferencia -w /go/src/github.com/lordofthejars/diferencia -e "version=${version}" lordofthejars/diferenciarelease:0.3.0 crossbuild.sh
