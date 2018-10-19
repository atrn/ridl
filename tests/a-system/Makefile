.PHONY: all test clean

ridl=../../ridl
template1?=c++-header
template2?=zmq-header
templates=../../templates
clang_format?=clang-format-mp-6.0
clang_format_args=--style=webkit -i

all test:
	$(ridl) -t $(templates)/$(template1) -o a_system.h
	$(ridl) -t $(templates)/$(template2) -o a_system_messages.h
	$(clang_format) $(clang_format_args) *.h
	dmake

clean:
	dmake clean; rm -f *.h
