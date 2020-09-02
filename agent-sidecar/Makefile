.EXPORT_ALL_VARIABLES:
OUT_DIR := ./_output
BIN_DIR := ./bin

$(shell mkdir -p $(OUT_DIR) $(BIN_DIR))

APP_NAME=gsecret-agent
IMAGE_NAME=imrenagi/$(APP_NAME)
PACKAGE=github.com/imrenagi/google-secret-k8s/agent-sidecar

CURRENT_DIR=$(shell pwd)
BUILD_DATE=$(shell date -u +'%Y-%m-%dT%H:%M:%SZ')
GIT_COMMIT=$(shell git rev-parse --short HEAD)
GIT_TAG=$(shell git describe --exact-match --tags HEAD 2>/dev/null)
GIT_TREE_STATE=$(shell if [ -z "`git status --porcelain`" ]; then echo "clean" ; else echo "dirty"; fi)

# perform static compilation
STATIC_BUILD?=true

override LDFLAGS += \
  -X ${PACKAGE}/version.buildDate=${BUILD_DATE} \
  -X ${PACKAGE}/version.gitCommit=${GIT_COMMIT} \
  -X ${PACKAGE}/version.gitTreeState=${GIT_TREE_STATE} \
	-X ${PACKAGE}/version.gitTag=${GIT_TAG}

ifeq (${STATIC_BUILD}, true)
override LDFLAGS += -extldflags "-static"
endif

ifneq (${GIT_TAG},)
IMAGE_TAG=${GIT_TAG}
else
IMAGE_TAG?=$(GIT_COMMIT)
endif

# Main Test Targets (without docker)
.PHONY: test
test:
	go test ./...

.PHONY: build.binaries
build.binaries:
	CGO_ENABLED=0 go build -a -ldflags '${LDFLAGS}' -o ${BIN_DIR}/${APP_NAME} ./main.go

.PHONY: build.image
build.image: 	
	DOCKER_BUILDKIT=1 docker build --cache-from $(IMAGE_NAME):latest -t $(IMAGE_NAME):$(IMAGE_TAG) -t $(IMAGE_NAME):latest --build-arg BUILDKIT_INLINE_CACHE=1 .
	
.PHONY: release.image
release.image: 	
	docker push $(IMAGE_NAME):latest
	docker push $(IMAGE_NAME):$(IMAGE_TAG)
	