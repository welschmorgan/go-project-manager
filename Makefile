TARGET=grlm

${TARGET}: main.go
	go build -a -o $@ $^

installdeps:
	go get ./...

clean:
	rm -f ${TARGET}
	go clean -x

re: clean ${TARGET}

.PHONY: installdeps clean