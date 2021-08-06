SHELL := /bin/bash

# curl -il http://localhost:3000/testerror

# go install github.com/divan/expvarmon@latest
# expvarmon -ports=":4000" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"

run-local:
	go run app/sales-api/main.go | jq

# ==============================================================================

# $(shell git rev-parse --short HEAD)
VERSION := 1.0

sales-api:
	docker build \
		-f zarf/docker/dockerfile.sales-api \
		-t sales-api-amd64:$(VERSION) \
		--build-arg VCS_REF=$(VERSION) \
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.

# ==============================================================================
# Running from within k8s/kind

KIND_CLUSTER := ardan-starter-cluster

# Upgrade to latest Kind (>=v0.11): e.g. brew upgrade kind
# For full Kind v0.11 release notes: https://github.com/kubernetes-sigs/kind/releases/tag/v0.11.0
# Kind release used for our project: https://github.com/kubernetes-sigs/kind/releases/tag/v0.11.1
# The image used below was copied by the above link and supports both amd64 and arm64.

kind-up:
	kind create cluster \
		--image kindest/node:v1.21.1@sha256:69860bda5563ac81e3c0057d654b5253219618a22ec3a346306239bba8cfa1a6 \
		--name $(KIND_CLUSTER) \
		--config zarf/k8s/kind/kind-config.yaml
	kubectl config set-context --current --namespace=sales-system

kind-down:
	kind delete cluster --name $(KIND_CLUSTER)

kind-load:
	kind load docker-image sales-api-amd64:$(VERSION) --name $(KIND_CLUSTER)

kind-services:
	kustomize build zarf/k8s/kind/sales-pod | kubectl apply -f -

kind-update: sales-api kind-load
	kubectl rollout restart deployment sales-pod

kind-logs:
	kubectl logs -l app=sales --all-containers=true -f --tail=100 | go run app/logfmt/main.go

kind-logs-sales:
	kubectl logs -l app=sales --all-containers=true -f --tail=100 | go run app/logfmt/main.go -service=SALES-API

kind-status-all:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --all-namespaces

kind-status:
	kubectl get nodes -o wide
	kubectl get svc -o wide
	kubectl get pods -o wide --watch --namespace=sales-system

kind-describe:
	kubectl describe nodes
	kubectl describe svc
	kubectl describe pod -l app=sales

# ==============================================================================
# Modules support

tidy:
	go mod tidy
	go mod vendor