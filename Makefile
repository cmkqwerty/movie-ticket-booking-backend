build:
	@go build -o bin/api

run: build
	@./bin/api

seed:
	@go run scripts/seed.go

docker:
	@docker build -t api .
	echo "running"
	@docker run -p 3000:3000 api

test:
	@go test -v ./...