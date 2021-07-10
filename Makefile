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

install: ${TARGET}
	[ -e "${INSTALL_DIR}" ] || mkdir -p ${INSTALL_DIR}
	[ -e "${INSTALL_DIR}/bin" ] || mkdir -p ${INSTALL_DIR}/bin
	cp ${TARGET} ${INSTALL_DIR}/bin/${NAME}

uninstall:
	[ ! -e "${INSTALL_DIR}/bin/${NAME}" ] || rm -f ${INSTALL_DIR}/bin/${NAME}

devinst: ${TARGET}
	@mkdir -p ${DIST_DIR}
	@rm -rf ${DIST_DIR}/test-wksp
	@cd ${DIST_DIR}; 7z x $$OLDPWD/test-wksp.7z >/dev/null || (echo failed to extract base workspace; exit 1)
	@echo "export PATH=${DIST_DIR}:$$PATH"

.PHONY: installdeps clean install devinst uninstall re