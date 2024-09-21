SRC := ./src
BIN := aptgit

APTGIT_HOME := $(HOME)/.aptgit
INSTALL_DIR := $(HOME)/.local/bin

all:
	@cd $(SRC) && $(MAKE) $@

run:
	@cd $(SRC) && $(MAKE) $@

fmt:
	@cd $(SRC) && $(MAKE) $@

install: $(BIN)
	install -v -D -t $(INSTALL_DIR) -m 775 $(BIN)
	mkdir -p $(APTGIT_HOME)/bin
	mkdir -p $(APTGIT_HOME)/installed
	mkdir -p $(APTGIT_HOME)/downloads
	cp -r ./gpkgs $(APTGIT_HOME)
	touch $(APTGIT_HOME)/aptgit.lock

clean:
	@cd $(SRC) && $(MAKE) $@
