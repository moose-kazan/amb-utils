.PHONY: all

all:
	mkdir -p ${PWD}/build
	cd ${PWD}/src && go build -o ${PWD}/build/ambunpack ./ambunpack/ambunpack.go

clean:
	rm -rf build