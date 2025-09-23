install:
	go install ./cmd/war

cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out
