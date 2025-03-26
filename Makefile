# Select the preferred available shell, prioritizing ash > zsh > bash
SHELL_PATHS := /bin/ash /bin/zsh /bin/bash 
SHELL := $(firstword $(wildcard $(SHELL_PATHS)))

# Define images/dependencies
KIND := kindest/node:v1.29.14
POSTGRES        := postgres:17.2

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

run-help: 
	go run app/services/sales-api/main.go --help | go run app/tooling/logfmt/main.go

curl:
	curl -il http://localhost:3000/v1/hack

curl-auth:
	curl -il -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/v1/hackauth

curl-create:
	curl -il -X POST -H 'Content-Type: application/json' \
		-d '{"name":"Joe","email":"joe@foo.com","roles":["ADMIN"], \
		"department":"IT","password":"42","passwordConfirm":"42"}' \
		http://localhost:3000/v1/users
		
load:
	hey -m GET -c 100 -n 100000 "http://localhost:3000/v1/hack"

admin:
	go run app/tooling/sales-admin/main.go

ready:
	curl -il http://localhost:3000/v1/readiness

live:
	curl -il http://localhost:3000/v1/liveness

# ==============================================================================
# Running from within k8s/kind

dev-up:
	kind create cluster \
	--image $(KIND) \
	--name $(KIND_CLUSTER) \
	--config zarf/k8s/dev/kind-config.yaml

	kubectl wait --timeout=120s --namespace=local-path-storage --for=condition=Available deployment/local-path-provisioner

	kind load docker-image $(POSTGRES) --name $(KIND_CLUSTER)

dev-down:
	kind delete cluster --name $(KIND_CLUSTER)

# ==============================================================================

dev-load:
	kind load docker-image $(SERVICE_IMAGE) --name $(KIND_CLUSTER)

dev-apply:
	kustomize build zarf/k8s/dev/database | kubectl apply -f -
	kubectl rollout status --namespace=$(NAMESPACE) --watch --timeout=120s sts/database

	kustomize build zarf/k8s/dev/sales | kubectl apply -f -
	kubectl wait pods --namespace=$(NAMESPACE) --selector app=$(APP) --timeout=120s --for=condition=Ready

dev-restart:
	kubectl rollout restart deployment $(APP) --namespace=$(NAMESPACE)

dev-update: all dev-load dev-restart

dev-update-apply: all dev-load dev-apply

# ==============================================================================

dev-logs:
	kubectl logs --namespace=$(NAMESPACE) -l app=$(APP) --all-containers=true \
	-f --tail=100 --max-log-requests=6 | go run app/tooling/logfmt/main.go -service=$(SERVICE_NAME)

dev-describe-deployment:
	kubectl describe deployment --namespace=$(NAMESPACE) $(APP)

dev-describe-sales:
	kubectl describe pod --namespace=$(NAMESPACE) -l app=$(APP)

dev-logs-db:
	kubectl logs --namespace=$(NAMESPACE) -l app=database --all-containers=true -f --tail=100

dev-logs-init:
	kubectl logs --namespace=$(NAMESPACE) -l app=$(APP) -f --tail=100 -c init-migrate

pgcli:
	pgcli postgresql://postgres:postgres@localhost

# ==============================================================================

dev-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces
	
# ==============================================================================
# Metrics and tracing
metrics-view-sc:
	expvarmon -ports="localhost:4000" \
	-vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"

# ==============================================================================
# Modules support

tidy:
	go mod tidy
	go mod vendor

test-race:
	CGO_ENABLED=1 go test -race -count=1 ./...

test-only:
	CGO_ENABLED=0 go test -count=1 ./...

lint:
	CGO_ENABLED=0 go vet ./...
	staticcheck -checks=all ./...

vuln-check:
	govulncheck ./...

test: test-only lint vuln-check

test-race: test-race lint vuln-check
