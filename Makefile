SERVICE ?= account

gqlgen:
	cd ./graphql && \
	go run github.com/99designs/gqlgen generate

protoc:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    $(SERVICE)/pb/$(SERVICE).proto

run:
ifeq ($(SERVICE),graphql)
	cd ./$(SERVICE) && \
	go run .
else
	cd ./$(SERVICE)/cmd/$(SERVICE) && \
	go run .
endif

up:
	docker-compose up --build

down:
	docker-compose down -v