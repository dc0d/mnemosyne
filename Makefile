.PHONY: test
test:
	clear
	go test -count=1 -timeout 10s -cover ./...

race:
	clear
	go test -count=1 -timeout 10s -race ./...

lint:
	golangci-lint run ./...

cover:
	go test -count=1 -timeout 10s -coverprofile=cover-profile.out -covermode=set -coverpkg=./... ./...

cover-html: cover
	go tool cover -html=cover-profile.out -o cover-coverage.html
