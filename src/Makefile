GO := go
GOFLAGS := -ldflags='-s -w -buildid=' -trimpath
BIN := aptgit

all: $(BIN)

$(BIN):
	CGO_ENABLED=0 $(GO) build $(GOFLAGS) -o ../$(BIN)

run:
	@$(GO) run .

fmt:
	@$(GO) fmt ./...

clean:
	rm -rf ../$(BIN)
