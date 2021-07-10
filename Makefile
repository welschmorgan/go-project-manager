NAME=grlm
DIST_DIR?=dist
INSTALL_DIR?=/usr/local

TARGET=${DIST_DIR}/${NAME}

${TARGET}: main.go
	go build -a -o $@ $^

installdeps:
	go get ./...

clean:
	rm -f ${TARGET}
	go clean -x

re: clean ${TARGET}

test-wksp:
	cd /tmp && rm -rf /tmp/test-wksp && 7z x $$OLDPWD/test-wksp.7z

install:
	[ -e "${INSTALL_DIR}" ] || mkdir -p ${INSTALL_DIR}
	[ -e "${INSTALL_DIR}/bin" ] || mkdir -p ${INSTALL_DIR}/bin
	cp ${TARGET} ${INSTALL_DIR}/bin/${NAME}

uninstall:
	[ ! -e "${INSTALL_DIR}/bin/${NAME}" ] || rm -f ${INSTALL_DIR}/bin/${NAME}

.PHONY: installdeps clean