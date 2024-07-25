.PHONY: run
run:
	go run main.go --port 8081

.PHONY: gen-chat-swagger
gen-docs:
	swag init -g api.go -d ./pkg/blockaction/api -o pkg/blockaction/api/swagger -pd