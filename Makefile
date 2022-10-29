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
