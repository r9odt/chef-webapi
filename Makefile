GIT_BRANCH := "unknown"
GIT_HASH := $(shell git log --pretty=format:%H -n 1)
GIT_HASH_SHORT := $(shell echo "${GIT_HASH}" | cut -c1-7)
GIT_TAG := $(shell git describe --always --tags --abbrev=0 | tail -c+2)
GIT_COMMIT := $(shell git rev-list v${GIT_TAG}..HEAD --count)
GIT_COMMIT_DATE := $(shell git show -s --format=%ci | cut -d\  -f1)

VERSION_FEATURE := ${GIT_TAG}-${GIT_BRANCH}
VERSION_NIGHTLY := ${GIT_COMMIT_DATE}-${GIT_HASH_SHORT}
VERSION_RELEASE := ${GIT_TAG}.${GIT_COMMIT}

GO_VERSION := $(shell go version | cut -d' ' -f3)
GO111MODULE := on

GOLANGCI_CERSION := "1.41.1"

.PHONY: lint
lint:
	go get github.com/golangci/golangci-lint/cmd/golangci-lint@v${GOLANGCI_CERSION}
	GOGC=30 golangci-lint run
	golint ./...
	revive -formatter friendly ./...

.PHONY: test
test:
	echo 'mode: atomic' > coverage.txt && go list ./... | \
	xargs -n1 -I{} sh -c 'go test -v -bench=. -covermode=atomic \
	-coverprofile=coverage.tmp {} && tail -n +2 coverage.tmp >> coverage.txt' && \
	rm coverage.tmp

.PHONY: build
build:
	CGO_ENABLED=0 GOOS=linux go build \
	-a -installsuffix cgo \
	-ldflags "-X main.Version=${VERSION_RELEASE} \
	-X main.GoVersion=${GO_VERSION} \
	-X main.GitCommit=${GIT_HASH}" \
	-o bin/web github.com/JIexa24/chef-webapi/cmd/web; 	