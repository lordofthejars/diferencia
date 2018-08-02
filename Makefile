version ?= latest

.PHONY: install
install:
	dep ensure
	packr build -o diferencia

.PHONY: build
build:
	packr build -o diferencia

.PHONY: crossbuild
crossbuild:
	docker run -it --rm -v "$$PWD":/go/src/github.com/lordofthejars/diferencia -w /go/src/github.com/lordofthejars/diferencia -e "version=${version}" lordofthejars/diferenciarelease:0.0.2 crossbuild.sh
