startdocker:
	@docker-compose down
	@docker-compose build
	@docker-compose up -d

stopdocker:
	@docker-compose down

build:
	@echo "---- Building Application ----"
	@go build -o bin/consumer cmd/consumer/main.go
	@go build -o bin/producer cmd/producer/main.go

consumer:
	@echo "---- Running Consumer ----"
	@go run cmd/consumer/main.go

producer:
	@echo "---- Running Producer ----"
	@go run cmd/producer/main.go
