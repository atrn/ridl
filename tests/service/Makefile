ridl=		../../ridl
template1?=	types.template
template2?=	zmq.template
clang_format?=	clang-format --style=file:../.clang-format -i

test: types.hpp messages.hpp
	dmake

types.hpp : protocol.ridl $(template1)
	$(ridl) -t $(template1) -o $@ protocol.ridl
	$(clang_format) $@

messages.hpp : protocol.ridl $(template2)
	$(ridl) -t $(template2) -o $@ protocol.ridl
	$(clang_format) $@

clean:
	dmake clean
	rm -f *.hpp
