BIN := trewoga

PKG := github.com/binarycode/trewoga

ARCH := amd64

VERSION := $(shell git describe --tags --always --dirty)

SRC_DIRS := cmd internal

IMAGE := golang:1.11-stretch

all: build

build: bin/$(ARCH)/$(BIN)

bin/$(ARCH)/$(BIN): build-dirs
	@echo "building: $@"
	@docker run                                                           \
	  -ti                                                                 \
	  --rm                                                                \
	  -u $$(id -u):$$(id -g)                                              \
	  -v "$$(pwd)/.go:/go"                                                \
	  -v "$$(pwd):/go/src/$(PKG)"                                         \
	  -v "$$(pwd)/bin/$(ARCH):/go/bin"                                    \
	  -v "$$(pwd)/.go/std/$(ARCH):/usr/local/go/pkg/linux_$(ARCH)_static" \
	  -v "$$(pwd)/.go/cache:/root/.cache"                                 \
	  -w /go/src/$(PKG)                                                   \
	  $(IMAGE)                                                            \
	  bash -c "                                                           \
	      ARCH=$(ARCH)                                                    \
	      VERSION=$(VERSION)                                              \
	      PKG=$(PKG)                                                      \
	      ./build/build.sh                                                \
	  "

shell: build-dirs
	@echo "launching a shell in the containerized build environment"
	@docker run                                                           \
	  -ti                                                                 \
	  --rm                                                                \
	  -u $$(id -u):$$(id -g)                                              \
	  -v "$$(pwd)/.go:/go"                                                \
	  -v "$$(pwd):/go/src/$(PKG)"                                         \
	  -v "$$(pwd)/bin/$(ARCH):/go/bin"                                    \
	  -v "$$(pwd)/.go/std/$(ARCH):/usr/local/go/pkg/linux_$(ARCH)_static" \
	  -v "$$(pwd)/.go/cache:/root/.cache"                                 \
	  -w /go/src/$(PKG)                                                   \
	  $(IMAGE)                                                            \
	  bash $(CMD)

server: build
	@docker-compose run --rm -p 8080:80 server

version:
	@echo $(VERSION)

build-dirs:
	@mkdir -p bin/$(ARCH)
	@mkdir -p .go .go/cache .go/src/$(PKG) .go/pkg .go/bin .go/std/$(ARCH)

clean: container-clean bin-clean

bin-clean:
	rm -rf .go bin
