# Check if Go is installed
ifneq ($(shell command -v go 1>/dev/null 2>&1; echo $$?), 0)
	GOLANG := golang
endif

.PHONY: build-deps
build-deps: golang

.PHONY: golang
golang:
	curl -Lo /tmp/go1.24.0.linux-amd64.tar.gz \
		https://golang.org/dl/go1.24.0.linux-amd64.tar.gz
	tar -C /usr/local -xvzf /tmp/go1.24.0.linux-amd64.tar.gz
	ln -sf /usr/local/go/bin/go /usr/local/bin/go
	ln -sf /usr/local/go/bin/gofmt /usr/local/bin/gofmt
	rm -f /tmp/go1.24.0.linux-amd64.tar.gz

.PHONY: build
build:
	CGO_ENABLED=1 go build -o sqlvalid .

.PHONY: test
test:
	CGO_ENABLED=1 go test -v ./...

.PHONY: clean
clean:
	rm -f sqlvalid

.PHONY: vet
vet:
	go vet ./...
