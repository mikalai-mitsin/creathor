.PHONY: build
.DEFAULT_GOAL := build

build:
	CGO_ENABLED=0 go build -v -o ./dist/creathor

lint:
	golangci-lint run ./... --timeout 5m0s

clean:
	golangci-lint run ./... --fix

log:
	git-chglog --config docs/.chglog/config.yml --output docs/CHANGELOG.md --next-tag $(tag)
