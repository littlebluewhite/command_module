run:
	go run cmd/api/main.go

migrate-create:
	@echo "---Creating migration files---"
	# another - migration create -ext sql -dir $(MPATH) -seq -digits 5 $(NAME)
	go run cmd/migrate/main.go create $(NAME) sql

migrate-up:
	go run cmd/migrate/main.go up

migrate-down:
	go run cmd/migrate/main.go down

migrate-force:
	go run cmd/migrate/main.go force $(VERSION)

swagger-docs-generate:
	swag init --parseDependency --parseInternal --parseDepth 1 -d cmd/api

docker-api-image-create-windows:
	docker build -t schedule:latest -f deploy/api/windows/Dockerfile .

docker-api-image-create-linux:
	docker build -t schedule:latest -f deploy/api/linux/Dockerfile .

docker-migrate-image-create-windows:
	docker build -t schedule-migrate:latest -f deploy/migrate/windows/Dockerfile .

docker-migrate-image-create-linux:
	docker build -t schedule-migrate:latest -f deploy/migrate/linux/Dockerfile .

docker-api-run-windows:
	docker run --name schedule -p 5487:5487 -e DB_HOST=host.docker.internal -v ${PWD}/deploy/api/dockerLog:/app/log schedule:latest

docker-api-run-linux:
	docker run --name schedule -p 5487:5487 --network="host" -v ${PWD}/deploy/api/dockerLog:/app/log schedule:latest

docker-migrate-up-windows:
	docker run --name schedule-migrate --rm -e DB_HOST=host.docker.internal schedule-migrate:latest

docker-migrate-up-linux:
	docker run --name schedule-migrate --rm --network="host" schedule-migrate:latest

docker-compose-up:
	docker-compose up

docker-compose-down:
	docker-compose down

# docker exec -it <container-id> sh