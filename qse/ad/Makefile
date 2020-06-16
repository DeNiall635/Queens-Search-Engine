REGISTRY = registry.hal.davecutting.uk/40173800/qse
NAME = ad
VERSION = latest

# Local DB settings
LOCAL_DB_PORT = 5432
LOCAL_DB_NETWORK = qse
LOCAL_DB_CONTAINER = qse-postgres
LOCAL_DB_USER = postgres
LOCAL_DB_PASS = Rkciymov8
LOCAL_NETWORK_DB = postgres://$(LOCAL_DB_USER):$(LOCAL_DB_PASS)@$(LOCAL_DB_CONTAINER):5432/postgres?sslmode=disable
LOCAL_DB = postgres://$(LOCAL_DB_USER):$(LOCAL_DB_PASS)@localhost:5432/postgres?sslmode=disable

default:
	@echo "=============Building Docker image============="
	docker build -f Dockerfile -t $(REGISTRY)/$(NAME):$(VERSION) ../

lint:
	@echo "=============Linting============="
	golint -set_exit_status ./...

db:
	@echo "=============Setting up database and migrating it============="
	# Remove any existing DBs
	docker stop $(LOCAL_DB_CONTAINER) || true && docker rm $(LOCAL_DB_CONTAINER) || true
	# Remove any existing network
	docker network rm $(LOCAL_DB_NETWORK) || true
	# Set up docker network
	docker network create -d bridge $(LOCAL_DB_NETWORK)
	# Set up new DB
	docker run --network=$(LOCAL_DB_NETWORK) -p $(LOCAL_DB_PORT):5432 --name $(LOCAL_DB_CONTAINER) -e POSTGRES_PASSWORD=$(LOCAL_DB_PASS) -d postgres
	# Wait for DB to start
	sleep 5
	# Migrate DB
	docker run --network=$(LOCAL_DB_NETWORK) -v $(shell dirname `pwd`)/data/ad:/migrations migrate/migrate -path=/migrations/ -database $(LOCAL_NETWORK_DB) up

migrate:
	@echo "=============Running migrations============="
	docker run --network=$(LOCAL_DB_NETWORK) -v $(shell dirname `pwd`)/data/search:/migrations migrate/migrate -path=/migrations/ -database $(LOCAL_NETWORK_DB) up

test:
	@echo "=============Testing============="
	go test ./...

run:
	@echo "=============Running============="
	AD_DB=$(LOCAL_DB) go run cmd/qse-ad/main.go
