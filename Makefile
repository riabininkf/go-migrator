.PHONY: build lint integration-tests

build:
	docker build -f Dockerfile -t gomigrator .
lint:
	golangci-lint run --config=./.golangci.yml
test:
	go test -race -count 100 ./pkg/...

integration-tests:
	docker-compose -f docker-compose.test.yml up --build

