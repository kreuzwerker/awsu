VERSION := "2.3.2"

BUILD := $(shell git rev-parse --short HEAD)
FLAGS	:= "-s -w -X=main.build=$(BUILD) -X=main.time=`TZ=UTC date '+%FT%TZ'` -X=main.version=$(VERSION)"
REPO := awsu
TOKEN = $(shell cat .token)
USER := kreuzwerker

build/awsu-linux-amd64:
	@mkdir -p build
	nice docker container run -it --rm -e "GO111MODULE=on" \
		-v $(PWD):/go/src/github.com/gesellix/awsu \
		golang:1.11-stretch bash -c \
		"apt-get update -q && apt-get install -qqy libpcsclite-dev && cd /go/src/github.com/gesellix/awsu && go mod download && go build -o $@ -ldflags $(FLAGS) awsu.go"

build/awsu-darwin-amd64:
	@mkdir -p build
	GO111MODULES=on nice go build -o $@ -ldflags $(FLAGS) awsu.go

build: build/awsu-darwin-amd64 build/awsu-linux-amd64;

clean:
	rm -rf build

release: clean build
	git tag $(VERSION) -f && git push --tags -f
	github-release release --user $(USER) --repo $(REPO) --tag $(VERSION) -s $(TOKEN)
	find build/* -type f -print0 | xargs -P8 -0J {} github-release upload --user $(USER) --repo $(REPO) --tag $(VERSION) -s $(TOKEN) --name {} --file {}

retract:
	github-release delete --tag $(VERSION) -s $(TOKEN)

test:
	go list ./... | grep -v exp | xargs go test -cover

.PHONY: build clean release retract test
