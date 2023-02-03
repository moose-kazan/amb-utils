.PHONY: all

PREFIX=/usr/local/

all:
	mkdir -p ${PWD}/build
	cd ${PWD}/src && go build -o ${PWD}/build/ambunpack ./ambunpack/ambunpack.go
	cd ${PWD}/src && go build -o ${PWD}/build/amb2html ./amb2html/amb2html.go

clean:
	rm -rf build

install: all
	mkdir -p ${PREFIX}/bin
	cp ${PWD}/build/ambunpack -v ${PREFIX}/bin/ambunpack
	cp ${PWD}/build/amb2html  -v ${PREFIX}/bin/amb2html

