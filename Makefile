.PHONY: gen-mock
gen-mock:
	go generate ./...

.PHONY: test
test:
	go test -v ./...

.PHONY: run
run:
	go run main.go --port 8082

.PHONY: gen-chat-swagger
gen-docs:
	swag init -g api.go -d ./pkg/blockaction/api -o pkg/blockaction/api/swagger -pd

.PHONY: build-base-image
build-base-image:
	docker build -f deployment/dockerfile.yaml -t blockaction-api .

.PHONY:
deploy: build-base-image
	docker-compose -f deployment/docker-compose.yaml up -d