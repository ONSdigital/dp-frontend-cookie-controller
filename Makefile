BINPATH ?= build

BUILD_TIME=$(shell date +%s)
GIT_COMMIT=$(shell git rev-parse HEAD)
VERSION ?= $(shell git tag --points-at HEAD | grep ^v | head -n 1)
LOCAL_DP_RENDERER_IN_USE = $(shell grep -c "\"github.com/ONSdigital/dp-renderer/v2\" =" go.mod)

.PHONY: audit
audit:
	go list -m all | nancy sleuth

.PHONY: build
build: generate-prod
	go build -tags 'production' -o $(BINPATH)/dp-frontend-cookie-controller -ldflags "-X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT) -X main.Version=$(VERSION)"

.PHONY: prepare-lint-go
prepare-lint-go:
	cp assets/assets.go assets/assets.go.bak
	echo 'func Asset(_ string) ([]byte, error) { return nil, nil }' >> assets/assets.go
	echo 'func AssetNames() []string { return []string{} }' >> assets/assets.go
	gofmt -w assets/assets.go

.PHONY: post-lint-go
post-lint-go:
	mv assets/assets.go.bak assets/assets.go

.PHONY: lint
lint: prepare-lint-go
	npx mega-linter-runner || { make post-lint-go; exit 1; }

.PHONY: lint-go
lint-go: prepare-lint-go
	golangci-lint run ./... || { make post-lint-go; exit 1; }

.PHONY: debug
debug: generate-debug
	go build -tags 'debug' -o $(BINPATH)/dp-frontend-cookie-controller -ldflags "-X main.BuildTime=$(BUILD_TIME) -X main.GitCommit=$(GIT_COMMIT) -X main.Version=$(VERSION)"
	HUMAN_LOG=1 DEBUG=1 $(BINPATH)/dp-frontend-cookie-controller

.PHONY: test
test: generate-prod
	go test -race -cover -tags 'production' ./...

.PHONY: generate-debug
generate-debug: fetch-dp-renderer
	# fetch the renderer library and build the dev version
	cd assets; go run github.com/kevinburke/go-bindata/go-bindata -prefix $(CORE_ASSETS_PATH)/assets -debug -o data.go -pkg assets locales/... templates/... $(CORE_ASSETS_PATH)/assets/locales/... $(CORE_ASSETS_PATH)/assets/templates/...
	{ printf "// +build debug\n"; cat assets/data.go; } > assets/debug.go.new
	mv assets/debug.go.new assets/data.go

.PHONY: generate-prod
generate-prod: fetch-dp-renderer
	# fetch the renderer library and build the prod version
	cd assets; go run github.com/kevinburke/go-bindata/go-bindata -prefix $(CORE_ASSETS_PATH)/assets -o data.go -pkg assets locales/... templates/... $(CORE_ASSETS_PATH)/assets/locales/... $(CORE_ASSETS_PATH)/assets/templates/...
	{ printf "// +build production\n"; cat assets/data.go; } > assets/data.go.new
	mv assets/data.go.new assets/data.go

.PHONY: fetch-dp-renderer
fetch-dp-renderer:
ifeq ($(LOCAL_DP_RENDERER_IN_USE), 1)
	$(eval CORE_ASSETS_PATH = $(shell grep -w "\"github.com/ONSdigital/dp-renderer/v2\" =>" go.mod | awk -F '=> ' '{print $$2}' | tr -d '"'))
else
	$(eval APP_RENDERER_VERSION=$(shell grep "github.com/ONSdigital/dp-renderer/v2" go.mod | cut -d ' ' -f2 ))
	$(eval CORE_ASSETS_PATH = $(shell go get github.com/ONSdigital/dp-renderer/v2@$(APP_RENDERER_VERSION) && go list -f '{{.Dir}}' -m github.com/ONSdigital/dp-renderer/v2))
endif

.PHONY: test-component
test-component:
	exit 0
