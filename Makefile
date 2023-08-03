test:
	go test -count=1 -cover ./...
	
cover:
	go test -coverprofile=coverage.out ./...
	go tool cover -html=coverage.out
	rm coverage.out