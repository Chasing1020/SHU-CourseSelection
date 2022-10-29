BINARY_NAME:=course-selection
GO:=$(shell which go)
GOFMT:=$(shell which gofmt)
GO_FILES:=$(shell find . -name "*.go" -type f)

export GO111MODULE:=on

.PHONY: all
all: fmt build

.PHONY: fmt
fmt:
	$(GO) mod tidy -go=1.17
	$(GOFMT) -s -w $(GO_FILES)

.PHONY: build
build:
	$(GO) build -o $(BINARY_NAME)

.PHONY:
run: fmt
	$(GO) run $(GO_FILES)

.PHONY: clean
clean:
	$(GO) clean
	rm -f $(BINARY_NAME)


.PHONY: build-release
build-release: build-windows-amd64 build-linux-amd64 build-macos-arm64 build-macos-amd64

.PHONY: clean-release
clean-release:
	rm windows-amd64-$(BINARY_NAME).zip linux-amd64-$(BINARY_NAME).zip darwin-arm64-$(BINARY_NAME).zip darwin-amd64-$(BINARY_NAME).zip

.PHONY: build-windows-amd64
build-windows-amd64:
	GOOS=windows GOARCH=amd64 $(GO) build -o windows-amd64-$(BINARY_NAME).exe
	zip windows-amd64-$(BINARY_NAME).zip windows-amd64-$(BINARY_NAME).exe config.yaml
	rm windows-amd64-$(BINARY_NAME).exe

.PHONY: build-linux-amd64
build-linux-amd64:
	GOOS=linux GOARCH=amd64 $(GO) build -o linux-amd64-$(BINARY_NAME)
	zip linux-amd64-$(BINARY_NAME).zip linux-amd64-$(BINARY_NAME) config.yaml
	rm linux-amd64-$(BINARY_NAME)

.PHONY: build-macos-arm64
build-macos-arm64:
	GOOS=darwin GOARCH=arm64 $(GO) build -o darwin-arm64-$(BINARY_NAME)
	zip darwin-arm64-$(BINARY_NAME).zip darwin-arm64-$(BINARY_NAME) config.yaml
	rm -f darwin-arm64-$(BINARY_NAME)

.PHONY: build-macos-amd64
build-macos-amd64:
	GOOS=darwin GOARCH=amd64 $(GO) build -o darwin-amd64-$(BINARY_NAME)
	zip darwin-amd64-$(BINARY_NAME).zip darwin-amd64-$(BINARY_NAME) config.yaml config
	rm -f darwin-amd64-$(BINARY_NAME)