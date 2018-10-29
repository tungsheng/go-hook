DIST := dist
EXECUTABLE := go-hook

DEPLOY_ACCOUNT := tonka
DEPLOY_IMAGE := $(EXECUTABLE)
GOFMT ?= gofmt "-s"

TARGETS ?= linux darwin windows
ARCHS ?= amd64 386
PACKAGES ?= $(shell go list ./... | grep -v /vendor/)
GOFILES := $(shell find . -name "*.go" -type f -not -path "./vendor/*")
SOURCES ?= $(shell find . -name "*.go" -type f)
TEMPLATES ?= $(shell find template -name "*.html" -type f ! -name "master.html")
TAGS ?=
LDFLAGS ?= -X 'main.Version=$(VERSION)'
TMPDIR := $(shell mktemp -d 2>/dev/null || mktemp -d -t 'tempdir')
STYLESHEETS := $(wildcard assets/dist/less/innhp.less  assets/dist/less/_*.less)

ifneq ($(shell uname), Darwin)
	EXTLDFLAGS = -extldflags "-static" $(null)
else
	EXTLDFLAGS =
endif

ifneq ($(DRONE_TAG),)
	VERSION ?= $(DRONE_TAG)
else
	VERSION ?= $(shell git describe --tags --always || git rev-parse --short HEAD)
endif

all: build

.PHONY: tar
tar:
	tar -zcvf release.tar.gz bin env Dockerfile Makefile

.PHONY: check_image
check_image:
	if [ "$(shell docker ps -aq -f name=$(EXECUTABLE))" ]; then \
		docker rm -f $(EXECUTABLE); \
	fi

.PHONY: dev
dev: build_image check_image
	docker run -d --name $(DEPLOY_IMAGE) --env-file env/env.$@ --net host -p 3003:3003 --log-opt max-size=10m --log-opt max-file=10 --restart always $(DEPLOY_ACCOUNT)/$(DEPLOY_IMAGE)

.PHONY: prod
prod: build_image check_image
	docker run -d --name $(DEPLOY_IMAGE) --env-file env/env.$@ --net host -p 3003:3003 --log-opt max-size=10m --log-opt max-file=10 --restart always $(DEPLOY_ACCOUNT)/$(DEPLOY_IMAGE)

.PHONY: fmt
fmt:
	$(GOFMT) -w $(GOFILES)

.PHONY: fmt-check
fmt-check:
	# get all go files and run go fmt on them
	@diff=$$($(GOFMT) -d $(GOFILES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make fmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

embedmd:
	@hash embedmd > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/campoy/embedmd; \
	fi
	embedmd -d *.md

gtfmt-check:
	@hash gtfmt > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/gotpl/gtfmt; \
	fi
	# get all go files and run gtfmt on them
	@diff=$$(gtfmt -l $(TEMPLATES)); \
	if [ -n "$$diff" ]; then \
		echo "Please run 'make gtfmt' and commit the result:"; \
		echo "$${diff}"; \
		exit 1; \
	fi;

gtfmt:
	@hash gtfmt > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/gotpl/gtfmt; \
	fi
	gtfmt $(TEMPLATES)

.PHONY: test-vendor
test-vendor:
	@hash govendor > /dev/null 2>&1; if [ $$? -ne 0 ]; then \
		go get -u github.com/kardianos/govendor; \
	fi
	govendor list +unused | tee "$(TMPDIR)/wc-gitea-unused"
	[ $$(cat "$(TMPDIR)/wc-gitea-unused" | wc -l) -eq 0 ] || echo "Warning: /!\\ Some vendor are not used /!\\"

	govendor list +outside | tee "$(TMPDIR)/wc-gitea-outside"
	[ $$(cat "$(TMPDIR)/wc-gitea-outside" | wc -l) -eq 0 ] || exit 1

	govendor status || exit 1
