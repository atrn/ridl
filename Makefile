.PHONY: all ridl clean doc clean-doc test
all: ridl
ridl:; @go build
clean:; rm -f ridl2 README.html
doc:; markdown README.md > README.html
clean-doc: rm -f README.html

test:ridl
	./ridl -t templates/c++-header.template test1.ridl > test1.hpp && \
	./ridl -t templates/c++-source.template test1.ridl > test1.cpp && \
	clang-format-mp-5.0 --style=webkit -i test.cpp test.hpp
