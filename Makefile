.PHONY: all ridl clean doc
all: ridl
ridl:; @go build
clean:; rm -f ridl2 README.html
doc:; markdown README.md > README.html
clean-doc: rm -f README.html
