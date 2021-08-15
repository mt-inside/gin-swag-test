install-tools:
	go install github.com/swaggo/swag/cmd/swag@v1.7.1

gen-openapi:
	swag init

run:
	go run main.go
