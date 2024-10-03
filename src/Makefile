GO := go
GOFLAGS := -ldflags='-s -w -buildid=' -trimpath
BIN := aptgit

OS := $(shell go env GOOS)
ARCH := $(shell go env GOARCH)

all: $(BIN)

$(BIN):
	CGO_ENABLED=0 GOOS=$(OS) ARCH=$(ARCH) $(GO) build $(GOFLAGS) -o ../$(BIN)

run:
	@$(GO) run .

fmt:
	@$(GO) fmt ./...

clean:
	rm -rf ../$(BIN)
