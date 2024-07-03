build: 
	@go build -o bin/GoPay

run: build 
	@./bin/GoPay

test:
	@go test -v ./...
