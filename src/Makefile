GO := go
GOFLAGS := -ldflags='-s -w -buildid=' -trimpath
BIN := aptgit

all: $(BIN)

$(BIN):
	$(GO) build $(GOFLAGS) -o ../$(BIN)

run:
	@$(GO) run .

fmt:
	@$(GO) fmt ./...

clean:
	rm -rf ../$(BIN)
