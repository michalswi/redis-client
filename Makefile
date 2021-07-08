GOLANG_VERSION := 1.15.6
ALPINE_VERSION := 3.13

GIT_REPO := github.com/michalswi/redis-client
DOCKER_REPO := michalsw
APPNAME := redis-client

VERSION ?= $(shell git describe --tags --always)
# if '.git/' doesn't exist
# VERSION ?= $(shell git describe --tags --always 2>/dev/null || echo testdev)
BUILD_TIME ?= $(shell date -u '+%Y-%m-%d %H:%M:%S')
LAST_COMMIT_USER ?= $(shell git log -1 --format='%cn <%ce>')
LAST_COMMIT_HASH ?= $(shell git log -1 --format=%H)
LAST_COMMIT_TIME ?= $(shell git log -1 --format=%cd --date=format:'%Y-%m-%d %H:%M:%S')

REDIS_HOST ?= localhost
REDIS_PORT ?= 6379
REDIS_PASS ?= ""
REDIS_TLS ?= ""

SERVICE_ADDR ?= 8080
DNS_NAME ?= localhost

AZ_RG ?= redisrg
AZ_LOCATION ?= westeurope
AZ_RANDOM ?=$(shell head /dev/urandom | tr -dc a-z0-9 | head -c 7)
AZ_DNS_LABEL ?= $(APPNAME)-$(AZ_RANDOM)

.DEFAULT_GOAL := help
.PHONY: test go-run go-build all docker-build docker-run docker-stop azure-rg azure-rg-del azure-aci azure-aci-logs azure-vnet-create azure-redis-vnet azure-aci-delete

help:
	@awk 'BEGIN {FS = ":.*##"; printf "\nUsage:\n  make \033[36m<target>\033[0m\n\nTargets:\n"} /^[a-zA-Z_-]+:.*?##/ \
	{ printf "  \033[36m%-15s\033[0m %s\n", $$1, $$2 }' $(MAKEFILE_LIST)

test:
	go test -v ./...

go-run: 		## Run redis client - no binary
	$(info -run redis-client (no binary)-)
	REDIS_HOST=$(REDIS_HOST) \
	REDIS_PORT=$(REDIS_PORT) \
	REDIS_PASS=$(REDIS_PASS) \
	REDIS_TLS=$(REDIS_TLS) \
	SERVICE_ADDR=$(SERVICE_ADDR) \
	DNS_NAME=$(DNS_NAME) \
	go run .	

go-build: 		## Build binary
	$(info -build redis-client binary-)
	CGO_ENABLED=0 \
	go build \
	-v \
	-ldflags "-s -w -X '$(GIT_REPO)/version.AppVersion=$(VERSION)' \
	-X '$(GIT_REPO)/version.BuildTime=$(BUILD_TIME)' \
	-X '$(GIT_REPO)/version.LastCommitUser=$(LAST_COMMIT_USER)' \
	-X '$(GIT_REPO)/version.LastCommitHash=$(LAST_COMMIT_HASH)' \
	-X '$(GIT_REPO)/version.LastCommitTime=$(LAST_COMMIT_TIME)'" \
	-o $(APPNAME)-$(VERSION) .

all: test go-build

docker-build:	## Build docker image
	$(info -build redis-client docker image-)
	docker build \
	--pull \
	--build-arg GOLANG_VERSION="$(GOLANG_VERSION)" \
	--build-arg ALPINE_VERSION="$(ALPINE_VERSION)" \
	--build-arg APPNAME="$(APPNAME)" \
	--build-arg VERSION="$(VERSION)" \
	--build-arg BUILD_TIME="$(BUILD_TIME)" \
	--build-arg LAST_COMMIT_USER="$(LAST_COMMIT_USER)" \
	--build-arg LAST_COMMIT_HASH="$(LAST_COMMIT_HASH)" \
	--build-arg LAST_COMMIT_TIME="$(LAST_COMMIT_TIME)" \
	--label="build.version=$(VERSION)" \
	--tag="$(DOCKER_REPO)/$(APPNAME):latest" \
	--tag="$(DOCKER_REPO)/$(APPNAME):$(VERSION)" \
	.

docker-run:		## Run docker image with default parameters (or overwrite)
	$(info -run redis-client in docker-)

	docker run -d --rm \
	--name $(APPNAME) \
	-e REDIS_HOST=$(REDIS_HOST) \
	-e REDIS_PORT=$(REDIS_PORT) \
	-e REDIS_PASS=$(REDIS_PASS) \
	-e REDIS_TLS=$(REDIS_TLS) \
	-p $(SERVICE_ADDR):$(SERVICE_ADDR) \
	$(DOCKER_REPO)/$(APPNAME):latest

docker-stop:	## Stop running docker
	$(info -stop redis-client in docker-)
	docker stop $(APPNAME)

azure-rg:	## Create the Azure Resource Group
	az group create --name $(AZ_RG) --location $(AZ_LOCATION)

azure-rg-del:	## Delete the Azure Resource Group
	az group delete --name $(AZ_RG)

azure-aci:	## Run redis-client app (Azure Container Instance)
	az container create \
	--resource-group $(AZ_RG) \
	--name $(APPNAME) \
	--image michalsw/$(APPNAME) \
	--restart-policy Always \
	--ports $(SERVICE_ADDR) \
	--dns-name-label $(AZ_DNS_LABEL) \
	--location $(AZ_LOCATION) \
	--environment-variables \
		SERVICE_ADDR=$(SERVICE_ADDR) \
		REDIS_HOST=$(REDIS_HOST) \
		REDIS_PORT=$(REDIS_PORT) \
		REDIS_PASS=$(REDIS_PASS) \
		REDIS_TLS=$(REDIS_TLS) \
		DNS_NAME=$(AZ_DNS_LABEL).$(AZ_LOCATION).azurecontainer.io

azure-aci-logs:	## Get redis-client app logs (Azure Container Instance)
	az container logs \
	--resource-group $(AZ_RG) \
	--name $(APPNAME)

azure-vnet-create:	## Create the private VNet
	az network vnet create \
	--name "$(APPNAME)-vnet" \
	--resource-group $(AZ_RG) \
	--address-prefixes "172.20.0.0/23" \
	--subnet-name "$(APPNAME)-subnet" \
	--subnet-prefixes "172.20.0.0/24"

azure-redis-vnet:	## Run redis (Azure Container Instance)
	az container create \
	--resource-group $(AZ_RG) \
	--name redis \
	--image redis \
	--restart-policy Always \
	--ports 6379 \
	--location $(AZ_LOCATION) \
	--vnet "$(APPNAME)-vnet" \
	--subnet "$(APPNAME)-subnet"

azure-aci-delete:	## Delete redis-client app (Azure Container Instance)
	az container delete \
	--resource-group $(AZ_RG) \
	--name $(APPNAME)