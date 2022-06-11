up:
	docker-compose up -d db zookeeper kafka kafkaui

down:
	docker-compose down

logs:
	docker-compose logs -f db zookeeper kafka kafkaui

server:
	docker-compose up -d web

server-restart:
	docker-compose restart web

server-down:
	docker-compose stop web

server-logs:
	docker-compose logs -f web