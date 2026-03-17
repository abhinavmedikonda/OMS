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

dockup:
	docker-compose up --build

dockdown:
	docker-compose down -v

dockapp:
	docker build -t abhinavmedikonda/oms-$(SERVICE)-api:latest -f $(SERVICE)/app.dockerfile .

dockdb:
	docker build -t abhinavmedikonda/oms-$(SERVICE)-db:latest -f $(SERVICE)/db.dockerfile $(SERVICE)/

k6:
	k6 run k6.js
