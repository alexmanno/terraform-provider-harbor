TEST?=$$(go list ./... | grep -v 'vendor')
HOSTNAME=hashicorp.com
NAMESPACE=alexmanno
NAME=harbor
BINARY=terraform-provider-${NAME}
VERSION=1.0.0
OS_ARCH=linux_amd64

default: install

build:
	go build -o ${BINARY}

release:
	docker run --rm \
      -v $(PWD):/go/src/github.com/$(NAMESPACE)/$(BINARY) \
	  -w /go/src/github.com/$(NAMESPACE)/$(BINARY) \
	  -e GITHUB_TOKEN \
	  -e GPG_FINGERPRINT \
      goreleaser/goreleaser release --rm-dist


pre-release:
	docker run --rm \
      -v $(PWD):/go/src/github.com/$(NAMESPACE)/$(BINARY) \
      -w /go/src/github.com/$(NAMESPACE)/$(BINARY) \
      goreleaser/goreleaser release --snapshot --rm-dist

install: build
	mkdir -p ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}
	mv ${BINARY} ~/.terraform.d/plugins/${HOSTNAME}/${NAMESPACE}/${NAME}/${VERSION}/${OS_ARCH}

test:
	go test -i $(TEST) || exit 1
	echo $(TEST) | xargs -t -n4 go test $(TESTARGS) -timeout=30s -parallel=4

testacc:
	TF_ACC=1 go test $(TEST) -v $(TESTARGS) -timeout 120m

tf-init:
	terraform -chdir=example init

tf-plan:
	terraform -chdir=example plan
