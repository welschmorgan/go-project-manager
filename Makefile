NAME=grlm
DIST_DIR?=dist
INSTALL_DIR?=/usr/local

LD_FLAGS="-H windowsgui"

TARGET=${DIST_DIR}/${NAME}

all: ${TARGET}

platforms:
	GOPATH=/home/darkboss/development/go xgo -branch develop github.com/welschmorgan/go-release-manager

${TARGET}: assets main.go
	env GOOS=linux GOARCH=amd64 go build -a -o $@ main.go

installdeps:
	go get ./...
	go get -u github.com/go-bindata/go-bindata/...

clean:
	rm -f ${TARGET}
	go clean -x

re: clean all

assets:
	go-bindata -fs -pkg gui -prefix cmd/gui/web-app -o cmd/gui/assets.go  cmd/gui/web-app/...

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

.PHONY: installdeps clean install devinst uninstall re all assets phony