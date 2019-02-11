version ?= latest

CUR_DIR:=$(shell pwd)
BINARY_DIR:=$(CUR_DIR)/binaries

.PHONY: all
all: tools install format lint build ## (default) Runs 'tools deps format lint compile' targets

.PHONY: install
install:
	dep ensure

.PHONY: tools
tools: ## Installs required go tools
	go get -u github.com/golangci/golangci-lint/cmd/golangci-lint
	go get -u golang.org/x/tools/cmd/goimports
	go get -u github.com/onsi/ginkgo/ginkgo
	go get -u github.com/onsi/gomega
	go get -u github.com/gobuffalo/packr/packr

.PHONY: format
format: ## Removes unneeded imports and formats source code
	goimports -l -w ./core/ ./difference/ ./exporter/ ./log/ ./metrics/

.PHONY: lint
lint: install ## Concurrently runs a whole bunch of static analysis tools
	golangci-lint run

.PHONY: test
test:
	ginkgo -r

.PHONY: build
build: install test
	packr build -o $(BINARY_DIR)/diferencia

.PHONY: crossbuild
crossbuild:
	docker run -it --rm -v "$$PWD":/go/src/github.com/lordofthejars/diferencia -w /go/src/github.com/lordofthejars/diferencia -e "version=${version}" lordofthejars/diferenciarelease:0.3.0 crossbuild.sh
