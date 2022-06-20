run:
	go run cmd/api/main.go

migrate-create:
	@echo "---Creating migration files---"
	# another - migrate create -ext sql -dir $(MPATH) -seq -digits 5 $(NAME)
	go run cmd/migrate/main.go create $(NAME) sql

migrate-up:
	go run cmd/migrate/main.go up

migrate-down:
	go run cmd/migrate/main.go down

migrate-force:
	go run cmd/migrate/main.go force $(VERSION)
