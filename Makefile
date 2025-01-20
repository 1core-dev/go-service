# Select the preferred available shell, prioritizing ash > zsh > bash
SHELL_PATHS := /bin/ash /bin/zsh /bin/bash 
SHELL := $(firstword $(wildcard $(SHELL_PATHS)))

# Define images/dependencies
KIND := kindest/node:v1.29.12

KIND_CLUSTER := starter-cluster
BASE_IMAGE_NAME := baseimage/service
SERVICE_NAME    := sales-api
VERSION         := 0.0.1
SERVICE_IMAGE   := $(BASE_IMAGE_NAME)/$(SERVICE_NAME):$(VERSION)

# ==============================================================================
all: service

service:
	docker build \
	-f zarf/docker/Dockerfile.service \
	-t $(SERVICE_IMAGE) \
	--build-arg BUILD_REF=$(VERSION) \
	--build-arg BUILD_DATE=`date -u +"%d-%m-%YT%H:%M:%SZ"` \
	.

run: 
	go run app/services/sales-api/main.go | go run app/tooling/logfmt/main.go

dev-up:
	kind create cluster \
	--image $(KIND) \
	--name $(KIND_CLUSTER) \
	--config zarf/k8s/dev/kind-config.yaml

	kubectl wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner

dev-down:
	kind delete cluster --name $(KIND_CLUSTER)

dev-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces
	