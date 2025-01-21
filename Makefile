# Select the preferred available shell, prioritizing ash > zsh > bash
SHELL_PATHS := /bin/ash /bin/zsh /bin/bash 
SHELL := $(firstword $(wildcard $(SHELL_PATHS)))

# Define images/dependencies
KIND := kindest/node:v1.29.12

KIND_CLUSTER    := 1core-starter-cluster
NAMESPACE       := sales-system
APP             := sales
BASE_IMAGE_NAME := 1core/service
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

# ==============================================================================
# Running from within k8s/kind

dev-up:
	kind create cluster \
	--image $(KIND) \
	--name $(KIND_CLUSTER) \
	--config zarf/k8s/dev/kind-config.yaml

	kubectl wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner

dev-down:
	kind delete cluster --name $(KIND_CLUSTER)

# ==============================================================================

dev-load:
	kind load docker-image $(SERVICE_IMAGE) --name $(KIND_CLUSTER)

dev-apply:
	kustomize build zarf/k8s/dev/sales | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=$(APP) --timeout=120s --for=condition=Ready

# ==============================================================================

dev-logs:
	kubectl logs --namespace=$(NAMESPACE) -l app=$(APP) --all-containers=true \
	-f --tail=100 --max-log-requests=6 | go run app/tooling/logfmt/main.go -service=$(SERVICE_NAME)

dev-describe-deployment:
	kubectl describe deployment --namespace=$(NAMESPACE) $(APP)

dev-describe-sales:
	kubectl describe pod --namespace=$(NAMESPACE) -l app=$(APP)

# ==============================================================================

dev-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces
	