.PHONY: clean vet build test

GOCMD=GO111MODULE=on go

clean:
	rm -f plexmatch.go y.output

vet:
	${GOCMD} vet ${GOARGS} ./...

build:
	$(MAKE) clean
	goyacc -l -o plexmatch.go plexmatch.y

test:
	$(MAKE) build
	go test
