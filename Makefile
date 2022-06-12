up:
	docker-compose up -d db zookeeper kafka kafkaui

down:
	docker-compose down

logs:
	docker-compose logs -f db zookeeper kafka kafkaui

server:
	docker-compose up --remove-orphans -d web validate-balance-worker validate-user-status-worker

server-restart:
	docker-compose restart web validate-balance-worker validate-user-status-worker

server-down:
	docker-compose stop web validate-balance-worker validate-user-status-worker

server-logs:
	docker-compose logs -f web validate-balance-worker validate-user-status-worker