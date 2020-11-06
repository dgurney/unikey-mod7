LDFLAGS=-ldflags "-w"
COMMAND=github.com/dgurney/unikey-mod7
PROGRAMSHORT=unikey-mod7

install:
	go install ${LDFLAGS} ${COMMAND}
windows:
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o build/windows/amd64/${PROGRAMSHORT}.exe ${COMMAND}
	GOOS=windows GOARCH=386 go build ${LDFLAGS} -o build/windows/386/${PROGRAMSHORT}.exe ${COMMAND}
	GOOS=windows GOARM=7 GOARCH=arm go build ${LDFLAGS} -o build/windows/arm/${PROGRAMSHORT}.exe ${COMMAND}
darwin:
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o build/darwin/amd64/${PROGRAMSHORT} ${COMMAND}
freebsd:
	GOOS=freebsd GOARCH=amd64 go build ${LDFLAGS} -o build/freebsd/amd64/${PROGRAMSHORT} ${COMMAND}
	GOOS=freebsd GOARCH=386 go build ${LDFLAGS} -o build/freebsd/386/${PROGRAMSHORT} ${COMMAND}
openbsd:
	GOOS=openbsd GOARCH=amd64 go build ${LDFLAGS} -o build/openbsd/amd64/${PROGRAMSHORT} ${COMMAND}
	GOOS=openbsd GOARCH=386 go build ${LDFLAGS} -o build/openbsd/386/${PROGRAMSHORT} ${COMMAND}
linux:
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o build/linux/amd64/${PROGRAMSHORT} ${COMMAND}
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o build/linux/arm64/${PROGRAMSHORT} ${COMMAND}
	GOOS=linux GOARCH=arm GOARM=7 go build ${LDFLAGS} -o build/linux/armv7/${PROGRAMSHORT} ${COMMAND}
	GOOS=linux GOARCH=386 go build ${LDFLAGS} -o build/linux/386/${PROGRAMSHORT} ${COMMAND}
clean:
	rm -rf build/
docker-image:
	docker build --build-arg UID=$(shell id -u) --build-arg GID=$(shell id -g) -t unikey-package .
docker-package:
	mkdir build
	docker run -it --mount type=bind,source=${CURDIR}/build,target=/go/src/unikey-mod7/build unikey-package make -j$(shell nproc) package
cross: windows darwin freebsd linux openbsd
package: windows darwin linux
	cd build; zip -r windows.zip windows && zip -r darwin.zip darwin && zip -r linux.zip linux
release: clean docker-image docker-package
all: install cross
