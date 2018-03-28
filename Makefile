VERSION := "2.1.0"

BUILD := $(shell git rev-parse --short HEAD)
FLAGS	:= "-s -w -X=main.build=$(BUILD) -X=main.time=`TZ=UTC date '+%FT%TZ'` -X=main.version=$(VERSION)"
REPO := awsu
TOKEN = $(shell cat .token)
USER := kreuzwerker

build:
	gox -parallel=8 -osarch="darwin/amd64 linux/amd64" -ldflags $(FLAGS) -output "build/awsu-{{.OS}}-{{.Arch}}" ./bin/
	parallel upx --best --ultra-brute --quiet {} ::: build/awsu-*-*

clean:
	rm -rf build

check:
	ifeq ($(strip $(shell git status --porcelain 2>/dev/null)),)
		$(error git state is not clean)
	endif

release: check clean build
	git tag $(VERSION) -f && git push --tags -f
	github-release release --user $(USER) --repo $(REPO) --tag $(VERSION) -s $(TOKEN)
	parallel github-release upload --user $(USER) --repo $(REPO) --tag $(VERSION) -s $(TOKEN) --name {/} --file {} ::: build/*

retract:
	github-release delete --tag $(VERSION) -s $(TOKEN)

test:
	go list ./... | grep -v exp | xargs go test -cover

.PHONY: build clean check release retract test
