tidy:
	go mod tidy

generate: tidy
	templ generate

tailwind : generate
	npx tailwindcss -i internal/assets/css/input.css -o internal/assets/css/input.css

run: tailwind
	go run cmd/main.go

build: tailwind
	go build -o bin/blogx cmd/main.go 

buildtemp: tailwind
	go build -o ./tmp/main cmd/main.go