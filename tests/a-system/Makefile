ridl=		../../ridl
template1?=	gen-cpp.template
template2?=	gen-zmq.template
clang_format?=	clang-format
clang_format_args=--style=file:../.clang-format -i

test: a_system.h a_system_messages.h
	dmake

a_system.h : a_system.ridl $(template1)
	$(ridl) -t $(template1) -o $@ .
	$(clang_format) $(clang_format_args) $@

a_system_messages.h : a_system.ridl $(template2)
	$(ridl) -t $(template2) -o $@ .
	$(clang_format) $(clang_format_args) $@

clean:
	dmake clean; rm -f *.h
