BINPATH ?= build

BUILD_TIME=$(shell date +%s)
GIT_COMMIT=$(shell git rev-parse HEAD)
VERSION ?= $(shell git tag --points-at HEAD | grep ^v | head -n 1)
CORE_ASSETS_PATH=$(shell go list -f '{{.Dir}}' -m github.com/rav-pradhan/test-modules/render)

.PHONY: audit
audit:
	go list -m all | nancy sleuth

.PHONY: build
build:
	go build -tags 'production' -o $(BINPATH)/dp-frontend-cookie-controller -ldflags "-X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT) -X main.Version=$(VERSION)"

.PHONY: debug
debug: generate-debug
	go build -tags 'debug' -o $(BINPATH)/dp-frontend-cookie-controller -ldflags "-X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT) -X main.Version=$(VERSION)"
	HUMAN_LOG=1 DEBUG=1 $(BINPATH)/dp-frontend-cookie-controller

.PHONY: test
test: generate-debug
	go test -race -cover ./...

.PHONY: generate-debug
generate-debug:
	# build the dev version
	go get github.com/rav-pradhan/test-modules/render
	cd assets; go run github.com/kevinburke/go-bindata/go-bindata -prefix $(CORE_ASSETS_PATH)/assets -debug -o data.go -pkg assets locales/... templates/... $(CORE_ASSETS_PATH)/../...
	{ echo "// +build debug\n"; cat assets/data.go; } > assets/debug.go.new
	mv assets/debug.go.new assets/data.go

.PHONY: test-local
test-local:
	echo $(CORE_ASSETS_PATH)