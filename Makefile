SHELL := /bin/bash
BASEDIR = $(shell pwd)

# build with verison infos
versionDir = "github.com/lughong/gin-api-demo/global/version"
gitTag = $(shell if [ "`git describe --tags --abbrev=0 2>/dev/null`" != "" ];then git describe --tags --abbrev=0; else git log --pretty=format:'%h' -n 1; fi)
buildDate = $(shell TZ=Asia/Shanghai date +%FT%T%z)
gitCommit = $(shell git log --pretty=format:'%H' -n 1)
gitTreeState = $(shell if git status|grep -q 'clean';then echo clean; else echo dirty; fi)

ldflags="-w -X ${versionDir}.gitTag=${gitTag} -X ${versionDir}.buildDate=${buildDate} -X ${versionDir}.gitCommit=${gitCommit} -X ${versionDir}.gitTreeState=${gitTreeState}"

all: production
production: gotool
	@echo "building project..."
	@go build -v -ldflags ${ldflags} github.com/lughong/gin-api-demo/cmd/...
	@echo "build finished!"
clean:
	rm -rf gin-api-demo
	find . -name "[._]*.s[a-w][a-z]" | xargs -I {} rm -f {}
gotool:
	gofmt -w .
	go vet ./... 2>&1 | grep -v vendor;true
gotest:
	go test -v ./... -bench=.
run:
	./gin-api-demo &
help:
	@echo "make - compile the source code"
	@echo "make clean - remove binary file and vim swp files"
	@echo "make gotool - run go tool 'fmt' and 'vet'"

.PHONY: clean gotool help
