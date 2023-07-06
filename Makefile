all: build

.PHONY: all

BINARY_NAME=reverse-proxy
TARGET_SERVER_EXAMPLE_NAME=targetserverexample

export GOPATH?=${HOME}/go

build:
	go build -o ${BINARY_NAME} -v cmd/reverseproxy/*.go

tests:
	go test -v ./...

example: build
	chmod +x ${BINARY_NAME}
	./${BINARY_NAME} &
	go run -v test/*.go

clean:
	go clean
	rm -f ${BINARY_NAME}