# Release version
RELEASE ?= latest
# Docker image repo
REPO ?= b4fun

# Get the currently used golang install path (in GOPATH/bin, unless GOBIN is set)
ifeq (,$(shell go env GOBIN))
GOBIN=$(shell go env GOPATH)/bin
else
GOBIN=$(shell go env GOBIN)
endif

docker-build-azure-usage-prom:
	GOOS=linux GOARCH=amd64 go build -o docker_bin/azure-usage-prom ./cmd/azure-usage-prom
	docker build docker_bin \
		-f ./cmd/azure-usage-prom/Dockerfile \
		-t ${REPO}/azure-usage-prom:${RELEASE}

docker-push-azure-usage-prom:
	docker push ${REPO}/azure-usage-prom:${RELEASE}

docker-build: docker-build-azure-usage-prom

docker-push: docker-push-azure-usage-prom