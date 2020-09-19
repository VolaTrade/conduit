BIN_NAME = candles 


build:
	@echo building....
	@GOPRIVATE=github.com/volatrade CGO_ENABLED=0 go build -a -tags netgo -o bin/$(BIN_NAME);
