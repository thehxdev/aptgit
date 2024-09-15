GO := go
SRC := ./src

all:
	@cd $(SRC) && $(MAKE) $@

run:
	@cd $(SRC) && $(MAKE) $@

fmt:
	@cd $(SRC) && $(MAKE) $@

clean:
	@cd $(SRC) && $(MAKE) $@
