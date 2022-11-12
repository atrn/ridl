.PHONY: all test clean

ridl=../ridl
templates=../templates
clang_format_args?=--style=file -i

all: test

demo.hpp : demo.ridl ../templates/c++-header.template $(ridl)
	$(ridl) --debug -T $(templates) -t c++-header -o demo.hpp demo.ridl
	clang-format $(clang_format_args) demo.hpp

test: demo.hpp
	echo '#include "demo.hpp"' > demo.cpp; echo 'int main() {}' >> demo.cpp
	dmake

clean:
	dmake clean; rm -f *.cpp *.hpp; $(MAKE) --no-print-directory -C a-system clean

