SHELL := /bin/bash

# $(shell git rev-parse --short HEAD)
VERSION := 1.0

sales-api:
	docker build \
		-f zarf/docker/dockerfile.sales-api \
		-t sales-api-amd64:$(VERSION) \
		--build-arg VCS_REF=$(VERSION)\
		--build-arg BUILD_DATE=`date -u +"%Y-%m-%dT%H:%M:%SZ"` \
		.
