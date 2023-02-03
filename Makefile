.PHONY: all

all:
	mkdir -p ${PWD}/build
	cd ${PWD}/src && go build -o ${PWD}/build/ambunpack ./ambunpack/ambunpack.go
	cd ${PWD}/src && go build -o ${PWD}/build/amb2html ./amb2html/amb2html.go

clean:
	rm -rf build