.PHONY: all test clean

ridl		?=	../../ridl
templates=	../../templates
clang_format?=	clang-format --style=file:../.clang-format -i

test: demo.hpp
	echo '#include "demo.hpp"' > demo.cpp
	echo 'int main() {}' >> demo.cpp
	dmake

demo.hpp : demo.ridl $(templates)/c++-header.template $(ridl)
	$(ridl) -t $(templates)/c++-header.template -o $@ demo.ridl
	$(clang_format) demo.hpp

clean:
	dmake clean
	rm -f *.cpp *.hpp
