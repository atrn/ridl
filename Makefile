.PHONY: all build clean docs test realclean tag tarball
all: build
build:; go build
clean:; rm -f ridl README.html
realclean: clean; @$(MAKE) --no-print-directory -C tests clean
docs:; markdown README.md > README.html
test: ridl; @$(MAKE) --no-print-directory -C tests
tag:; git tag v`cat version.txt`
tarball:
	@d=$$(basename "`pwd`"); v=$$(cat version.txt); \
	tar -C .. -czf ../ridl-$$v.tar.gz "$$d" && ls -l ../ridl-$$v.tar.gz
