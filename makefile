SHELL := /bin/bash

# curl -il http://localhost:3000/testerror
# curl -il -H "Authorization: Bearer ${TOKEN}" http://localhost:3000/testauth
# export TOKEN=eyJhbGciOiJSUzI1NiIsImtpZCI6IjU0YmIyMTY1LTcxZTEtNDFhNi1hZjNlLTdkYTRhMGUxZTJjMSIsInR5cCI6IkpXVCJ9.eyJleHAiOjE2NTk4MDE4NjIsImlhdCI6MTYyODI2NTg2MiwiaXNzIjoic2VydmljZSBwcm9qZWN0Iiwic3ViIjoiMTIzNDU2Nzg5IiwiUm9sZXMiOlsiQURNSU4iXX0.Llefo3A8QPbihMDv3gb6x36qvOMm2GejFtmu0t_PCDviV7vYqsUJvx9SS_kAFYZjDieEF4AcIXGCcMkJaayezC2mK4KnBUaQjjBdo-2lEzoJyzvuJAU9lfjCPFx5_Y7DD6S17m-wMADRV0VbZuyGac1qBYl6msXIRe_ZtHIGQI5Kt6LlOyx9uBpc1BFCaJZ1Q3k1glTFpl6DK-_ksTGJXvN8909QHrosGfqwU2eNvnjQ80vNx3rzPVQjYbHriNumpwJuIJCNGWsQeIH9raJwW8xhxBqDh5sUuCYW8TY8ejcsKiEwwES50xVNX5NE33IHBi7FErAfXZmzMsQndeHglA
#
# go install github.com/divan/expvarmon@latest
# expvarmon -ports=":4000" -vars="build,requests,goroutines,errors,panics,mem:memstats.Alloc"
#
# To generate a private/public key PEM file.
# openssl genpkey -algorithm RSA -out private.pem -pkeyopt rsa_keygen_bits:2048
# openssl rsa -pubout -in private.pem -out public.pem
# ./sales-admin genkey

run-admin:
	go run app/admin/main.go

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

kind-services: kind-apply

kind-apply:
	kustomize build zarf/k8s/kind/database-pod | kubectl apply -f -
	kubectl wait --namespace=database-system --timeout=120s --for=condition=Available deployment/database-pod
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

kind-status-db:
	kubectl get pods -o wide --watch --namespace=sales-system

# ==============================================================================
# Modules support

tidy:
	go mod tidy
	go mod vendor