build:
	@echo 'Building the project...'
	go build -o songs/cmd/service/main ./cmd/service

run: build
	@echo 'Running the project...'
	./songs/cmd/service/main

clean:
	@echo 'Cleaning...'
	go clean
	rm -f songs/cmd/service/main

lint:
	@echo 'Linting the project...'
	gofumpt -w .
	go mod tidy
	golangci-lint run --config .golangci.yaml

test: up
	go test -v ./...

up:
	docker compose up -d

down:
	docker compose down
