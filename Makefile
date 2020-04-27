DOCKER_RUN=docker run -t -v $(CURDIR):/work:rw paasta-deb-builder-$*
CMDS=$(wildcard cmd/*)
UID:=$(shell id -u)
GID:=$(shell id -g)

GO_VERSION=1.12.7
VERSION=0.0.2

.PHONY: cmd $(CMDS)

all: build test
test:
	GO111MODULE=on go test -failfast -v ./...

build:
	GO111MODULE=on go build -v ./...

clean:
	rm -rf bin
	rm -rf dist/*

cmd: cmd/*

$(CMDS): 
	[ -d bin ] || mkdir -p bin
	GO111MODULE=on go build -o bin/paasta-tools-$(subst cmd/,,$@) $@/*.go

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
