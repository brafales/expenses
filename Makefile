GOPACKAGES = $(shell go list ./...  | grep -v /vendor/)

build:
	dep ensure -v
	env GOOS=linux go build -ldflags="-s -w" -o bin/monzo monzo/main.go
	env GOOS=linux go build -ldflags="-s -w" -o bin/toshl toshl/main.go

.PHONY: clean
clean:
	rm -rf ./bin ./vendor Gopkg.lock

.PHONY: deploy
deploy: clean build
	sls deploy --verbose

silent_deploy: clean build
	sls deploy > /dev/null 2>&1

test:
	@go test -v $(GOPACKAGES)
