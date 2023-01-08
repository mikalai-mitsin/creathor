.PHONY: build
.DEFAULT_GOAL := build
SEMVER_REGEX := \([0-9]*\.[0-9]*\.[0-9]*\)

build:
	CGO_ENABLED=0 go build -v -o ./dist/creathor

lint:
	golangci-lint run ./... --timeout 5m0s

clean:
	golangci-lint run ./... --fix

log:
	git-chglog --config docs/.chglog/config.yml --output docs/CHANGELOG.md --next-tag $(tag)


release:
	git flow release start ${tag}
	sed -i "" 's/const version = "${SEMVER_REGEX}"/const version = "${tag}"/' ./main.go
	git-chglog --config docs/.chglog/config.yml --output docs/CHANGELOG.md --next-tag $(tag)
	golangci-lint run ./... --timeout 5m0s
	go test ./...
	git add .
	git commit -m "bumped the version number"
	git flow release finish ${tag}