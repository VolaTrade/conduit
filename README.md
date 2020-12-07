# conduit

A beatiful data pipeline implemented in golang that concurrently collects live transaction data for anything traded against bitcoin on Binance 

### To run locally 
	1. `cp config.env.template config.env`
	2. insert env vars into `config.env`
	3. Either 
		a) `make build` and `./bin/candles`
		b) `make docker-build` and `make docker-run`

### Deployment 
**Before deploying to prod:**
- Ensure all unit tests pass
- Ensure all integration tests pass
- Create a temp ec2 and deploy to there and monitor the service

Follow these steps to deploy:
1. Ensure you're on master branch
2. Double check all of the above steps are complete and at least 2 other team members have reviewed the changes
3. Tag the branch via `NEW_VERSION=<vx.x> make tag"`
	- This will create a git tag for the new version and update the version file
4. Push the tags via `git push origin ${NEW_VERSION}

### Testing 



### Links 
1. https://github.com/binance-exchange/binance-official-api-docs/blob/master/rest-api.md
2. https://binance-docs.github.io/apidocs/spot/en/#websocket-market-streams]
