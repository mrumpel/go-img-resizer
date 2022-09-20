run:
	docker-compose up

build:
	go build -o go-img-resizer cmd/*.go

test:
	go test -race ./internal/...