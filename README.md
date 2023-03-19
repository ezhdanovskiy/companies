# companies
Microservice to handle companies.

## Integration tests
Run `make test/int/docker-compose`.  
It runs postgres and kafka in docker-compose; creates topic in kafka; run integration tests; stops containers.

## Local run
1. Run `make run/local`.  
It runs postgres and kafka in docker-compose; creates topic in kafka; build app; run app.
2. Run `make company/livecycle` in other console.
It creates, patches, get and remove a company using curl.
3. Run `make kafka/topic/consume` in third console to see events.

## Linter
Run `make lint`.