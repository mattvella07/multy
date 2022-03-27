SHELL=/bin/bash -O globstar
PORT=8000

build:
	protoc --go_out=. --go_opt=paths=source_relative     --go-grpc_out=. --go-grpc_opt=paths=source_relative   api/proto/**/*.proto
	go build -v ./...

run: build
	go run main.go serve --port=$(PORT)

clean:
	find api/proto -name '*.pb.go' -delete
	find . -name '*.lock.hcl' -delete
	find . -name '*.tfstate*' -delete
	rm -rf ./test/**/.terraform