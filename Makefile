SERVICE=rct
BINARY=rct

DOCKER_IMAGE_NAME=atuldaemon/rct

.DEFAULT_GOAL := help

check: test vet ## Runs all tests

test: ## Run the unit tests
	go test -race -v $(shell go list ./... | grep -v /vendor/)

#lint: ## Lint all files
#	go list ./... | grep -v /vendor/ | xargs -L1 golint -set_exit_status

vet: ## Run the vet tool
	go vet $(shell go list ./... | grep -v /vendor/)

clean: ## Clean up build artifacts
	go clean

dep: ## Ensure dependencies are pulled
	go get -u github.com/golang/dep/cmd/dep && dep ensure -v

local-build: dep
	go build

local-run: local-build
	./${SERVICE} -http.addr=:8080

docker-build: dep  ## Build docker image
	docker build -t ${DOCKER_IMAGE_NAME} .

docker-push: docker-build ## Push Docker image to registry
	docker tag ${DOCKER_IMAGE_NAME}:latest ${DOCKER_IMAGE_NAME}
	docker push ${DOCKER_IMAGE_NAME}

docker-run: ## Run docker image
	docker run -p 8080:8080 -it ${DOCKER_IMAGE_NAME}

docker-up: ## Spin up the project
	docker-compose up --build ${BINARY}

docker-stop: ## Stop running containers
	docker stop ${SERVICE}_${BINARY}_1

help: ## Display this help message
	@cat $(MAKEFILE_LIST) | grep -e "^[a-zA-Z_\-]*: *.*## *" | awk 'BEGIN {FS = ":.*?## "}; {printf "\033[36m%-30s\033[0m %s\n", $$1, $$2}'

.SILENT: build test lint vet clean docker-build docker-push help
