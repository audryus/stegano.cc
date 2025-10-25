.PHONY: build run test
.SILENT: build run
.DEFAULT_GOAL := run

build:
	@go build -o bin/app -buildvcs=false

run: build 
	@./bin/app

air: 
	@go build -o bin/app.exe

version:
	@lazyver semver

test:
	@go clean -testcache
	@go run gotest.tools/gotestsum@latest --packages="./tests/..." --format testdox

release:
	@goreleaser release --clean

docker:
	docker build -f ./build/Dockerfile -t audryus/stegano-site .