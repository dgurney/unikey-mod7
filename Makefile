LDFLAGS=-ldflags "-w"
PROGRAMPATH=github.com/dgurney/unikey-mod7
PROGRAM=unikey-mod7

install:
	go install ${LDFLAGS} ${PROGRAMPATH}
windows:
	GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o build/windows/amd64/${PROGRAM}.exe ${PROGRAMPATH}
	GOOS=windows GOARCH=386 go build ${LDFLAGS} -o build/windows/386/${PROGRAM}.exe ${PROGRAMPATH}
	GOOS=windows GOARM=7 GOARCH=arm go build ${LDFLAGS} -o build/windows/arm/${PROGRAM}.exe ${PROGRAMPATH}
darwin:
	GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o build/darwin/amd64/${PROGRAM} ${PROGRAMPATH}
freebsd:
	GOOS=freebsd GOARCH=amd64 go build ${LDFLAGS} -o build/freebsd/amd64/${PROGRAM} ${PROGRAMPATH}
	GOOS=freebsd GOARCH=386 go build ${LDFLAGS} -o build/freebsd/386/${PROGRAM} ${PROGRAMPATH}
openbsd:
	GOOS=openbsd GOARCH=amd64 go build ${LDFLAGS} -o build/openbsd/amd64/${PROGRAM} ${PROGRAMPATH}
	GOOS=openbsd GOARCH=386 go build ${LDFLAGS} -o build/openbsd/386/${PROGRAM} ${PROGRAMPATH}
linux:
	GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o build/linux/amd64/${PROGRAM} ${PROGRAMPATH}
	GOOS=linux GOARCH=arm64 go build ${LDFLAGS} -o build/linux/arm64/${PROGRAM} ${PROGRAMPATH}
	GOOS=linux GOARCH=arm GOARM=7 go build ${LDFLAGS} -o build/linux/armv7/${PROGRAM} ${PROGRAMPATH}
	GOOS=linux GOARCH=386 go build ${LDFLAGS} -o build/linux/386/${PROGRAM} ${PROGRAMPATH}
clean:
	rm -rf build/
docker-image:
	docker build --build-arg UID=$(shell id -u) --build-arg GID=$(shell id -g) -t unikey-package .
docker-package:
	mkdir -p build
	docker run -it --mount type=bind,source=${CURDIR}/build,target=/go/src/unikey-mod7/build unikey-package make -j$(shell nproc) package
cross: windows darwin freebsd linux openbsd
package: windows darwin linux
	cd build; zip -r windows.zip windows && zip -r darwin.zip darwin && zip -r linux.zip linux
release: clean docker-image docker-package
all: install cross
