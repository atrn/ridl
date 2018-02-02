.PHONY: all test clean

ridl=../ridl
templates=../templates
format=clang-format-mp-5.0 --style=webkit -i

all test:
	$(ridl) -t $(templates)/c++-header.template -o test.hpp
	$(ridl) -t $(templates)/c++-source.template -o test.cpp
	$(format) test.cpp test.hpp

clean:
	rm -f *.cpp *.hpp
