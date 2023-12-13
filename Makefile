.PHONY: all binary man

PREFIX := /usr/local/
PWD := $(shell pwd)

ifeq ($(OS),Windows_NT)
	EXEEXT := .exe
endif


all: binary man

binary:
	mkdir -p ${PWD}/build/bin
	cd ${PWD}/src && go build -o ${PWD}/build/bin/ambunpack${EXEEXT} ./ambunpack/ambunpack.go
	cd ${PWD}/src && go build -o ${PWD}/build/bin/amb2html${EXEEXT} ./amb2html/amb2html.go

man:
	mkdir -p ${PWD}/build/man
	cp -v ${PWD}/man/*.1 ${PWD}/build/man
	gzip -f -9 ${PWD}/build/man/*.1

clean:
	rm -rf build

install: all
	mkdir -p ${PREFIX}/bin
	cp -v ${PWD}/build/bin/ambunpack ${PREFIX}/bin/ambunpack
	cp -v ${PWD}/build/bin/amb2html  ${PREFIX}/bin/amb2html
	mkdir -p ${PREFIX}/share/man/man1
	cp -v ${PWD}/build/man/*.1.gz ${PREFIX}/share/man/man1/

