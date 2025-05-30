.PHONY: casino transactions

test-coverage: 
	go test -coverprofile=coverage.out $(shell go list ./... | grep -v /pkg | grep -v /cmd)
	go tool cover -func=coverage.out
	go tool cover -html=coverage.out

test:
	go clean -testcache
	go test $(shell go list ./... | grep -v /pkg | grep -v /cmd)

api-docs:
	swag init -g src/main.go --output src/docs