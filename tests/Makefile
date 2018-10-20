.PHONY: all test clean

ridl=../ridl
templates=../templates
clang_format?=clang-format-mp-6.0
clang_format_args?=--style=webkit -i

all test:
	$(ridl) -t $(templates)/c++-header.template -o tests.hpp
	$(ridl) -t $(templates)/c++-source.template -o tests.cpp
	$(clang_format) $(clang_format_args) *.cpp *.hpp
	dmake

clean:
	dmake clean; rm -f *.cpp *.hpp; $(MAKE) --no-print-directory -C a-system clean

