BIN_NAME=conduit

build:
	@echo building binary...
	@GOPRIVATE=github.com/volatrade CGO_ENABLED=0 go build -a -tags netgo -o bin/${BIN_NAME}

deps:
	git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/volatrade/".insteadOf "https://github.com/volatrade/" && go mod download

test:
	go test -cover ./...

docker-build:
	docker build -t ${BIN_NAME} . --build-arg GITHUB_TOKEN=${GITHUB_TOKEN}

start-dev:
	docker compose -f docker-compose-dev.yaml up

start-prod:
	docker-compose -f docker-compose-prod.yaml up

docker-run:
	docker run --network="host" --log-opt max-size=10m --log-opt max-file=5 ${BIN_NAME}

ecr-push-image:
	docker push ${ECR_URI}/${BIN_NAME}

ecr-login:
	aws ecr get-login-password | docker login --username AWS --password-stdin ${ECR_URI}

run:
	python3 control_panel/driver.py

tag:
	git tag ${NEW_VERSION} && echo ${NEW_VERSION} >> version

