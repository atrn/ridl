ridl	= ../../ridl
tests	?= *.ridl

test:
	@set -e; \
	for file in $(tests) ; do \
	    echo '----------------------------------------------------------------'; \
	    echo $${file} ; \
	    echo '----------------------------------------------------------------'; \
	    template=`echo $${file} | sed 's/-.*//'`; \
	    $(ridl) -t $${template} $${file}; \
	    echo; \
	done

parse-module-test:
	@$(ridl) .
