.PHONY: all

ridl=../ridl
templates=../templates
format=clang-format-mp-5.0 --style=webkit -i

all:
	@for t in test-*.ridl ; do \
	    b=`basename $$t`; \
	    $(ridl) -t $(templates)/c++-header.template -o $${b}.hpp $$t; \
	    $(ridl) -t $(templates)/c++-source.template -o $${b}.cpp $$t; \
	    $(format) $${b}.cpp $${b}.hpp; \
	done


clean:
	rm -f *.cpp *.hpp
