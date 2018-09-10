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
