tidy:
	go mod tidy

generate: tidy
	templ generate

run: generate
	go run cmd/main.go

build: generate
	go build -o bin/blogx cmd/main.go 

buildtemp: generate
	go build -o ./tmp/main cmd/main.go