BIN_NAME=tickers

build:
	@echo building wire....
	@wire 
	@echo building binary...
	@GOPRIVATE=github.com/volatrade CGO_ENABLED=0 go build -a -tags netgo -o bin/${BIN_NAME};

docker-build:
	docker build -t ${BIN_NAME} . --build-arg GITHUB_TOKEN=$(GITHUB_TOKEN)


docker-run:
	docker run --restart=always -d ${BIN_NAME}

integration-test:
	docker-compose up --remove-orphans

ecr-push-image:
	docker push ${ECR_URI}/${BIN_NAME}

ecr-login:
	aws ecr get-login-password | docker login --username AWS --password-stdin ${ECR_URI}

run:
	python3 control_panel/driver.py

