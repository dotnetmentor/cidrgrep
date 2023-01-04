PKG = github.com/cli_template_owner/cli_template_name
COMMIT = $$(git describe --tags --always)

export GO111MODULE=on

BUILD_LDFLAGS = -X $(PKG)/version.Commit=$(COMMIT)

default: test

test:
	go test ./... -coverprofile=coverage.out -covermode=count

build:
	go build -ldflags="$(BUILD_LDFLAGS)"

prerelease:
	@test $${VER?Environment variable VER is required}
	git pull origin main --tag
	go mod tidy
	@echo "${VER}">> versions
	git add go.mod go.sum versions
	git commit -m "Bumped version number to ${VER}"
	git tag ${VER}

release:
	@test $${GITHUB_TOKEN?Environment variable GITHUB_TOKEN is required}
	git push origin main --tag
	goreleaser --rm-dist

.PHONY: default test
