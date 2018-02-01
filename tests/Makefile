.PHONY: all

ridl=../ridl
templates=../templates

all:
	@for t in test-*.ridl ; do \
	    b=`basename $$t`; \
	    $(ridl) -t $(templates)/c++-header.template -o $${b}.hpp $$t; \
	    $(ridl) -t $(templates)/c++-source.template -o $${b}.cpp $$t; \
	done


clean:
	rm -f *.cpp *.hpp
