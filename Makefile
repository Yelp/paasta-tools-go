ifeq ($(findstring .yelpcorp.com,$(shell hostname -f)), .yelpcorp.com)
	PAASTA_ENV ?= YELP
else
	PAASTA_ENV ?= $(shell hostname --fqdn)
endif

DOCKER_RUN=docker run -t -v $(CURDIR):/work:rw paasta-deb-builder-$*
CMDS=$(wildcard cmd/*)
UID:=$(shell id -u)
GID:=$(shell id -g)

GO_VERSION=1.12.7
VERSION=0.0.22

ifeq ($(PAASTA_ENV),YELP)
	GO_TAGS=-tags yelp
	GO_MODFILE=-modfile int.mod
	GO_ENV=GONOSUMDB=*.yelpcorp.com GOPROXY=http://athens.paasta-norcal-devc.yelp GOPRIVATE=*github.yelpcorp.com
else
	GO_TAGS=
	GO_MODFILE=
	GO_ENV=
endif

GOBUILD=$(GO_ENV) CGO_ENABLED=0 GO111MODULE=on go build $(GO_TAGS) $(GO_MODFILE) -ldflags="\
	-X github.com/Yelp/paasta-tools-go/pkg/version.Version=$(VERSION) \
	-X github.com/Yelp/paasta-tools-go/pkg/version.PaastaVersion=$(PAASTA_VERSION)"

GOTEST=$(GO_ENV) GO111MODULE=on go test $(GO_TAGS) $(GO_MODFILE)

.PHONY: cmd $(CMDS)

all: build test
test:

	$(GOTEST) -failfast -v ./...

build:
	$(GOBUILD) -v ./...

clean:
	rm -rf bin
	rm -rf dist/*
	rm -f paasta_go
	go clean -testcache

cmd: cmd/*

$(CMDS):
	[ -d bin ] || mkdir -p bin
	$(GOBUILD) -o bin/paasta-tools-$(subst cmd/,,$@) $@/*.go

docker_build_%:
	@echo "Building build docker image for $*"
	[ -d dist/$* ] || mkdir -p dist/$*
	cd ./yelp_package/$* && docker build --build-arg GO_VERSION=$(GO_VERSION) -t paasta-deb-builder-$* .

deb_%: clean docker_build_%
	$(DOCKER_RUN) /bin/bash -c ' \
		$(MAKE) cmd && \
		mv bin/paasta{-tools-paasta,_go} && \
		fpm --output-type deb --input-type dir --version $(VERSION) \
			--deb-dist $* --deb-priority optional \
			--name paasta-tools-go --package dist \
			--description "CLI tools for PaaSTA in Go" \
			--package dist/$* \
			bin=/usr/ && \
		chown -R $(UID):$(GID) bin dist \
	'

itest_%: deb_%
	@echo "Built package for $*"

openapi-codegen:
	rm -rf pkg/paastaapi
	mkdir -p pkg/paastaapi
	rm oapi.yaml
	curl -o oapi.yaml https://raw.githubusercontent.com/Yelp/paasta/master/paasta_tools/api/api_docs/oapi.yaml
	docker run --rm -i --user `id -u`:`id -g` -v `pwd`:/src -w /src \
	        yelp/openapi-generator-cli:20201026 \
	        generate -i oapi.yaml -g go --package-name paastaapi -o pkg/paastaapi
	# Remove all files except *.go
	find `pwd`/pkg/paastaapi -mindepth 1 ! -name \*.go -delete
	@echo "Do not forget to 'git add' and 'git commit' updated oapi.yaml and paasta-api"

paasta_go:
	$(GOBUILD) -v -o paasta_go ./cmd/paasta

# Steps to release
# 1. Bump version in Makefile
# 2. `make release`
release:
	# docker run -it --rm -v "$(pwd)":/usr/local/src/paasta-tools-go \
	# 	ferrarimarco/github-changelog-generator \
	# 	-u Yelp \
	# 	-p paasta-tools-go \
	# 	--max-issues=100 \
	# 	--future-release $(VERSION) \
	# 	--output ../CHANGELOG.md
	@git diff
	@echo "Now Run:"
	@echo 'git commit -a -m "Released $(VERSION) via make release"'
	@echo 'git tag -a -m "Released $(VERSION) via make release" v$(VERSION)'
	@echo 'git push --tags origin master'
