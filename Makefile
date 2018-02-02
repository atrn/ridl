.PHONY: all ridl clean doc clean-doc test realclean
all: ridl
ridl:; go build
clean:; rm -f ridl README.html
realclean: clean; @$(MAKE) --no-print-directory -C tests clean
doc:; markdown README.md > README.html
clean-doc: rm -f README.html
test: ridl; @$(MAKE) --no-print-directory -C tests
