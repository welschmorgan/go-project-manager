TARGET=grlm

${TARGET}: main.go
	go build -o $@ $^

installdeps:
	go get ./...

clean:
	rm -f ${TARGET}
	go clean

re: clean ${TARGET}

.PHONY: installdeps clean