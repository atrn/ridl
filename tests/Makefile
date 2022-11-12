.PHONY: all test clean

ridl=../ridl
templates=../templates
clang_format_args?=--style=Microsoft -i

all:

demo.hpp : demo.ridl ../templates/c++-header.template $(ridl)
	$(ridl) --debug -T $(templates) -t c++-header -o demo.hpp demo.ridl
	clang-format $(clang_format_args) demo.hpp

demo.cpp : demo.ridl ../templates/c++-source.template $(ridl)
	$(ridl) --debug -T $(templates) -t c++-source.template -o demo.cpp demo.ridl
	clang-format $(clang_format_args) demo.cpp

test: demo.hpp demo.cpp
	dmake

clean:
	dmake clean; rm -f *.cpp *.hpp; $(MAKE) --no-print-directory -C a-system clean

