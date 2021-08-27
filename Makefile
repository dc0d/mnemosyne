.PHONY: test
test:
	clear
	go test -count=1 -timeout 10s -cover -p 1 ./...

lint:
	golangci-lint run ./...
