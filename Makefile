BIN_NAME = tickers

build:
	@echo building wire....
	@wire 
	@echo building binary...
	@GOPRIVATE=github.com/volatrade CGO_ENABLED=0 go build -a -tags netgo -o bin/$(BIN_NAME);

docker-build:
	docker build -t candles . --build-arg GITHUB_TOKEN=$(GITHUB_TOKEN)


docker-run:
	docker run -e DB_PORT=5432 -e DB_HOST=docker.for.mac.host.internal -d candles 

integration-test:
	docker-compose up --remove-orphans

ecr-push-image:
	docker push 752939442315.dkr.ecr.us-west-2.amazonaws.com/candles

ecr-login:
	aws ecr get-login-password --profile volatrade | docker login --username AWS --password-stdin 752939442315.dkr.ecr.us-west-2.amazonaws.com


run:
	python3 control_panel/driver.py

