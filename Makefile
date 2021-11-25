VERSION := "2.3.3"

BUILD := $(shell git rev-parse --short HEAD)
FLAGS	:= "-s -w -X=main.build=$(BUILD) -X=main.time=`TZ=UTC date '+%FT%TZ'` -X=main.version=$(VERSION)"
REPO := awsu
USER := kreuzwerker

build/awsu-linux-amd64:
	@mkdir -p build
	nice docker container run --rm \
		-v $(PWD):/build/awsu \
		-w /build/awsu \
		golang:1.17-stretch bash -c \
		"apt-get update -q && apt-get install -qqy libpcsclite-dev && go mod download && go build -o $@ -ldflags $(FLAGS) awsu.go"
		
build/awsu-linux-amd64-ubuntu:
	@mkdir -p build
	nice docker container run --rm -e DEBIAN_FRONTEND=noninteractive \
		-v $(PWD):/build/awsu \
		-w /build/awsu \
		ubuntu:20.04 bash -c \
		"apt-get update -q && apt-get install -qqy build-essential software-properties-common pkg-config wget libpcsclite-dev && wget -c https://dl.google.com/go/go1.17.2.linux-amd64.tar.gz -O - | tar -xz -C /usr/local && export PATH=$$PATH:/usr/local/go/bin && go mod download && go build -o $@ -ldflags $(FLAGS) awsu.go"

# Test within the container
# curl -sL https://git.io/goreleaser | bash -s -- --rm-dist --skip-publish --snapshot --skip-sign --debug

build/awsu-darwin-amd64:
	@mkdir -p build
	nice go build -o $@ -ldflags $(FLAGS) awsu.go

build: build/awsu-darwin-amd64 build/awsu-linux-amd64;

clean:
	rm -rf build

test:
	go list ./... | grep -v exp | xargs go test -cover

.PHONY: build clean release retract test
