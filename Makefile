build:
	go build -o bin/myapp main.go
run:build
	./bin/myapp

proto:
	protoc --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    proto/blog.proto


.PHONY:proto			