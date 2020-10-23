# tickers

A beatiful data pipeline implemented in golang that concurrently collects live transaction data for anything traded against bitcoin on Binance 

### To run locally 
	1. `cp config.env.template config.env`
	2. insert env vars into `config.env`
	3. Either 
		a) `make build` and `./bin/candles`
		b) `make docker-build` and `make docker-run`

### Deployment 


### Testing 



### Links 
1. https://github.com/binance-exchange/binance-official-api-docs/blob/master/rest-api.md
2. https://binance-docs.github.io/apidocs/spot/en/#websocket-market-streams]
