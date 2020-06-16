REGISTRY = registry.hal.davecutting.uk/40178456/qse
NAME = ui
VERSION = latest

default:
	@echo "=============Building Docker image============="
	docker build -f Dockerfile -t $(REGISTRY)/$(NAME):$(VERSION) ../

lint:
	@echo "=============Linting============="
	golint -set_exit_status ./...

test:
	@echo "=============Testing============="
	go test ./...

run:
	@echo "=============Running============="
	go run cmd/qse-ui/main.go
