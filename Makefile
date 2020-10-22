mkfile_path := $(abspath $(lastword $(MAKEFILE_LIST)))
base_dir := $(notdir $(patsubst %/,%,$(dir $(mkfile_path))))

SERVICE ?= $(base_dir)

DOCKER_REGISTRY=registry.uw.systems
DOCKER_REPOSITORY_NAMESPACE=telecom
DOCKER_ID=telco
DOCKER_REPOSITORY_IMAGE=$(SERVICE)
DOCKER_REPOSITORY=registry.uw.systems/$(DOCKER_REPOSITORY_NAMESPACE)/$(DOCKER_REPOSITORY_IMAGE)
LINT_EXCLUDE=internal

K8S_NAMESPACE=$(DOCKER_REPOSITORY_NAMESPACE)
K8S_DEPLOYMENT_NAME=$(DOCKER_REPOSITORY_IMAGE)
K8S_CONTAINER_NAME=$(K8S_DEPLOYMENT_NAME)

BUILDENV :=
BUILDENV += CGO_ENABLED=0
GIT_HASH := $(CIRCLE_SHA1)
ifeq ($(GIT_HASH),)
  GIT_HASH := $(shell git rev-parse HEAD)
endif
LINKFLAGS :=-s -X main.gitHash=$(GIT_HASH) -extldflags "-static"
TESTFLAGS := -v -cover
LINT_FLAGS :=--disable-all --enable=vet --enable=vetshadow --enable=golint --enable=ineffassign --enable=goconst
LINTER_EXE := gometalinter.v1
LINTER := $(GOPATH)/bin/$(LINTER_EXE)

EMPTY :=
SPACE := $(EMPTY) $(EMPTY)
join-with = $(subst $(SPACE),$1,$(strip $2))

LEXC :=
ifdef LINT_EXCLUDE
	LEXC := $(call join-with,|,$(LINT_EXCLUDE))
endif

.PHONY: install
install:
	go get -v -t -d ./...

$(LINTER):
	go get -u gopkg.in/alecthomas/$(LINTER_EXE)
	$(LINTER) --install

.PHONY: lint
lint: $(LINTER)
ifdef LEXC
	$(LINTER) --exclude '$(LEXC)' $(LINT_FLAGS) ./...
else
	$(LINTER) $(LINT_FLAGS) ./...
endif

.PHONY: clean
clean:
	rm -f $(SERVICE)

# builds our binary
$(SERVICE):
	$(BUILDENV) go build -o $(SERVICE) -a -ldflags '$(LINKFLAGS)' ./cmd/$(SERVICE)

build: $(SERVICE)

.PHONY: test
test:
	@cp swagger.json /tmp/swagger.json
	$(BUILDENV) go test $(TESTFLAGS) ./...
	@rm /tmp/swagger.json

docker-image:
	docker build -t $(DOCKER_REPOSITORY):local . --build-arg SERVICE=$(SERVICE) --build-arg GITHUB_TOKEN=$(GITHUB_TOKEN)

ci-docker-auth:
	@echo "Logging in to $(DOCKER_REGISTRY) as $(DOCKER_ID)"
	@docker login -u $(DOCKER_ID) -p $(DOCKER_PASSWORD) $(DOCKER_REGISTRY)

ci-docker-build: ci-docker-auth
	docker build -t $(DOCKER_REPOSITORY):$(CIRCLE_SHA1) . --build-arg SERVICE=$(SERVICE) --build-arg GITHUB_TOKEN=$(GITHUB_TOKEN)
	docker tag $(DOCKER_REPOSITORY):$(CIRCLE_SHA1) $(DOCKER_REPOSITORY):latest
	docker push $(DOCKER_REPOSITORY)

K8S_URL=https://elb.master.k8s.dev.uw.systems/apis/extensions/v1beta1/namespaces/$(K8S_NAMESPACE)/deployments/$(K8S_DEPLOYMENT_NAME)
K8S_PAYLOAD={"spec":{"template":{"spec":{"containers":[{"name":"$(K8S_CONTAINER_NAME)","image":"$(DOCKER_REPOSITORY):$(CIRCLE_SHA1)"}]}}}}

ci-kubernetes-push:
	test "$(shell curl -o /dev/null -w '%{http_code}' -s -X PATCH -k -d '$(K8S_PAYLOAD)' -H 'Content-Type: application/strategic-merge-patch+json' -H 'Authorization: Bearer $(K8S_DEV_TOKEN)' '$(K8S_URL)')" -eq "200"

