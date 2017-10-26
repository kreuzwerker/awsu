VERSION := "2.0"

BUILD := $(shell git rev-parse --short HEAD)
FLAGS	:= "-s -w -X=main.build=$(BUILD) -X=main.version=$(VERSION)"

.PHONY: build clean test

build:
	gox -parallel=8 -osarch="darwin/amd64 linux/amd64" -ldflags $(FLAGS) -output "build/{{.OS}}-{{.Arch}}/awsu" ./bin/ && find build -type f -print0 | xargs -0 -P 8 upx -1q

clean:
	rm -rf build

test:
	go list ./... | grep -v exp | xargs go test -cover
