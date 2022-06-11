setup:
	docker-compose up -d

setup-down:
	docker-compose down

logs:
	docker-compose logs -f

server:
	go run . server