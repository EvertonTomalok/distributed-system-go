up:
	docker-compose up -d

down:
	docker-compose down

logs:
	docker-compose logs -f

server:
	go run . server