TARGET=grlm

${TARGET}: main.go
	go build -a -o $@ $^

installdeps:
	go get ./...

clean:
	rm -f ${TARGET}
	go clean -x

re: clean ${TARGET}

test-wksp:
	cd /tmp && 7z x $$OLDPWD/test-wksp.7z

.PHONY: installdeps clean