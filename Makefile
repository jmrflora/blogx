tidy:
	go mod tidy

generate: tidy
	templ generate

run: generate
	go run cmd/main.go