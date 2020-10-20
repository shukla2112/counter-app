APP?=counter
REGISTRY?=gcr.io/images
COMMIT_SHA=$(shell git rev-parse --short HEAD)

.PHONY: build
## build: build the application
build: clean
	@echo "Building..."
	@go build -o ${APP} main.go

.PHONY: run
## run: runs go run main.go
run:
	go run -race main.go

.PHONY: clean
## clean: cleans the binary
clean:
	@echo "Cleaning"
	@go clean

.PHONY: docker-build
## docker-build: build the docker image 
#docker-build: build
docker-build:
	docker build -t ${APP} .
	docker tag ${APP} ${APP}:${COMMIT_SHA}

.PHONY: docker-push
## docker-push: push the docker image to repository
docker-push: check-environment docker-build
	docker push ${REGISTRY}/${ENV}/${APP}:${COMMIT_SHA}

check-environment:
ifndef ENV
    $(error ENV not set, allowed values - `dev` or `staging` or `production`)
endif

.PHONY: setup
## setup: setup go modules
setup:
	@go mod init

.PHONY: help
## help: prints this help message
help:
	@echo "Usage: \n"
	@sed -n 's/^##//p' ${MAKEFILE_LIST} | column -t -s ':' |  sed -e 's/^/ /'

