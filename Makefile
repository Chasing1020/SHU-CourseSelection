export GO111MODULE:=on

GO_FILES:=$(shell find . -name "*.go" -type f)

.PHONY: fmt
fmt:
	go mod tidy
	gofmt -s -w $(GO_FILES)

.PHONY: build
build:
	go build -race

.PHONY: run
run:
	go run -race *.go