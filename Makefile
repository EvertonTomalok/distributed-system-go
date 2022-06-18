setup:
	docker-compose up -d db mongodb zookeeper kafka kafkaui user-api

setup-down:
	docker-compose stop db mongodb zookeeper kafka kafkaui user-api

down:
	docker-compose down

setup-logs:
	docker-compose logs -f --tail 10 db mongodb zookeeper kafka kafkaui user-api

server:
	docker-compose up --remove-orphans -d validate-balance-worker validate-user-status-worker orchestrator web

server-build:
	docker-compose up --build --remove-orphans -d validate-balance-worker validate-user-status-worker orchestrator web

server-restart:
	docker-compose restart validate-balance-worker validate-user-status-worker orchestrator web

server-down:
	docker-compose stop web validate-balance-worker validate-user-status-worker orchestrator

server-logs:
	docker-compose logs -f --tail 10 validate-balance-worker validate-user-status-worker orchestrator web

swag:
	swag init

build:
	go build .