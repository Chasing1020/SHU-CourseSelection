BINARY_NAME:=CourseSelection
GO:=$(shell which go)
GOFMT:=$(shell which gofmt)
GO_FILES:=$(shell find . -name "*.go" -type f)

export GO111MODULE:=on

.PHONY: fmt
fmt:
	$(GO) mod tidy -go=1.17
	$(GO) -s -w $(GO_FILES)

.PHONY: build
build:
	$(GO) build -o $(BINARY_NAME)

.PHONY: fmt
run:
	$(GO) run $(GO_FILES)

.PHONY: clean
clean:
	$(GO) clean
	rm -f $(BINARY_NAME)
