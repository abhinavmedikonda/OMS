SERVICE ?= account

gqlgen:
	cd ./graphql && \
	go run github.com/99designs/gqlgen generate

pb:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    $(SERVICE)/pb/$(SERVICE).proto

gorun:
	cd ./$(SERVICE)/cmd/$(SERVICE) && \
	go run .