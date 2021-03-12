BIN_NAME=conduit
GOMOCK := $(shell command -v mockgen 2> /dev/null)

.PHONY: build 
build:
	@echo building binary...
	@GOPRIVATE=github.com/volatrade CGO_ENABLED=0 go build -a -tags netgo -o bin/${BIN_NAME}

.PHONY: deps 
deps:
	git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/volatrade/".insteadOf "https://github.com/volatrade/" && go mod download

.PHONY: test 
test:
	cd internal && go test -cover ./...

test-integration: docker-up
	cd integration_tests && go test -cover ./... 
	docker-compose down 

.PHONY: docker-build 
docker-build:
	docker build -t ${BIN_NAME} . --build-arg GITHUB_TOKEN=${GITHUB_TOKEN}

.PHONY: docker-run
docker-run:
	docker run --name conduit --network=conduit-compose --log-opt max-size=10m --log-opt max-file=5 ${BIN_NAME}

ecr-push-image:
	docker push ${ECR_URI}/${BIN_NAME}

ecr-login:
	aws ecr get-login-password | docker login --username AWS --password-stdin ${ECR_URI}

tag:
	git tag ${NEW_VERSION} && echo ${NEW_VERSION} >> version

.PHONY : gen-mocks
gen-mocks : setup/gomock go-gen-mocks

.PHONY: setup/gomock
setup/gomock:
ifeq ('$(GOMOCK)','')
	@echo "Installing gomock"
	@GO111MODULE=off go get github.com/golang/mock/gomock >/dev/null
	@GO111MODULE=off go install github.com/golang/mock/mockgen >/dev/null
endif

.PHONY: go-gen-mocks
go-gen-mocks:
	@echo "generating go mocks..."
	@GO111MODULE=on go generate --run "mockgen*" ./...

.PHONY: build-linux
build-linux:
	@echo "\033[0;34m» Building Conduit Linux Binary\033[0;39m"
	@CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -a -tags netgo -o bin/$(BIN_NAM)
	@echo "\033[0;32m» Successfully Built Binary :) \033[0;39m"

.PHONY: docker-up
docker-up:
	@echo "\033[0;34m» Creating Conduit Service Dependencies \033[0;39m"
	docker-compose up -d

.PHONY: docker-dev-build
docker-dev-build: build-linux
	@echo "\033[0;34m» Building Conduit Image \033[0;39m"
	@docker build -t ${BIN_NAME} -f Dockerfile.dev .  
	@echo "\033[0;32m» Successfully Built Test Image :) \033[0;39m"

docker-dev-down:
	docker stop conduit
	docker rm conduit

