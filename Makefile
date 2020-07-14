DOCKER_RUN=docker run -t -v $(CURDIR):/work:rw paasta-deb-builder-$*
CMDS=$(wildcard cmd/*)
UID:=$(shell id -u)
GID:=$(shell id -g)

GO_VERSION=1.12.7
VERSION=0.0.5

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

gen-paasta-api:
	rm -rf pkg/paastaapi
	mkdir -p pkg/paastaapi
	rm swagger.json
	curl -o swagger.json https://raw.githubusercontent.com/Yelp/paasta/master/paasta_tools/api/api_docs/swagger.json
	docker run \
		--rm -it \
		--user "$$(id -u):$$(id -g)" \
		-e GOPATH=$$HOME/go:/go \
		-v $$HOME:$$HOME \
		-w $$(pwd) quay.io/goswagger/swagger \
		generate client -f ./swagger.json -t pkg/paastaapi
	@echo "Due to bug in goswagger you may need to add an import for paastaapi/client/operations"
	@echo "in pkg/paastaapi/client/paasta_client.go, run 'go build ./...' to check."
	@echo
	@echo "Do not forget to 'git add' and 'git commit' updated swagger.json and paasta-api"

paasta_go:
ifeq ($(PAASTA_ENV),YELP)
	GOPROXY=http://athens.paasta-norcal-devc.yelp \
	GO111MODULE=on go build -tags yelp -modfile int.mod -v -o paasta_go ./cmd/paasta
else
	GO111MODULE=on go build -v -o paasta_go ./cmd/paasta
endif

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
