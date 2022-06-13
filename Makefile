setup:
	docker-compose up -d db zookeeper kafka kafkaui

down:
	docker-compose down

setup-logs:
	docker-compose logs -f db zookeeper kafka kafkaui

server:
	docker-compose up --remove-orphans -d validate-balance-worker validate-user-status-worker orchestrator web

server-build:
	docker-compose up --build --remove-orphans -d validate-balance-worker validate-user-status-worker orchestrator web

server-restart:
	docker-compose restart web validate-balance-worker validate-user-status-worker orchestrator

server-down:
	docker-compose stop web validate-balance-worker validate-user-status-worker orchestrator

server-logs:
	docker-compose logs -f web validate-balance-worker validate-user-status-worker orchestrator

swag:
	swag init