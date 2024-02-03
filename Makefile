.PHONY: run es

run:
	go run main.go

es:
	docker-compose up -d