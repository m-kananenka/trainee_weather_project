.PHONY: user_service

build:
	go build main.go

run:
	go run main.go

test:
	go test ./... -coverprofile=coverage.out

mock:
	go generate ./...

run-migration-up:
	migrate -database postgres://postgres:secretpassword@localhost:5432/postgres?sslmode=disable -path db/migrations up

run-migration-down:
	migrate -database postgres://postgres:secretpassword@localhost:5432/postgres?sslmode=disable -path db/migrations down

proto:
	protoc --go_out=. --go-grpc_out=. api/user.proto

docker-up:
	docker-compose up

docker-down:
	docker-compose down

lint:
	golangci-lint run
