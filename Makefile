include .env
export

PROJECT_NAME := azeroth-digest
EXEC_NAME := azeroth-digest
SSH_USER := ansible
DEPLOY_TARGET_IP := 100.100.77.57
DB_URL := duckdb:data/azdg.duckdb

.PHONY: help ## print this
help:
	@echo ""
	@echo "$(PROJECT_NAME) Development CLI"
	@echo ""
	@echo "Usage:"
	@echo "  make <command>"
	@echo ""
	@echo "Commands:"
	@grep '^.PHONY: ' Makefile | sed 's/.PHONY: //' | awk '{split($$0,a," ## "); printf "  \033[34m%0-10s\033[0m %s\n", a[1], a[2]}'

.PHONY: init ## initialize the project
init:
	go run . init

.PHONY: run ## Run the project
run:
	@go run . bot serve --token=$$AZDG_DISC_TOKEN

.PHONY: sync ## Sync discord commands with gateway
sync:
	@go run . bot sync --token=$$AZDG_DISC_TOKEN

.PHONY: deps ## install dependencies used for development
deps:
	@go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest

.PHONY: clean ## delete generated code
clean:
	rm -rf generated

.PHONY: build ## builds the project
build:
	go build -ldflags='-s' -o=./bin/${PROJECT_NAME} .
	GOOS=linux GOARCH=amd64 go build -ldflags='-s' -o=./bin/linux_amd64/${PROJECT_NAME} .

.PHONY: lint ## run golangci-lint
lint:
	golangci-lint run

.PHONY: test ## run tests
test:
	go test -v ./...

.PHONY: fmt ## format the project
fmt:
	go fmt ./...


.PHONY: db/migration/status ## get the status of the db migrations
db/migration/status:
	goose duckdb $(DB_URL) -dir migrations status

.PHONY: db/migrate ## run database migrations
db/migrate:
	goose duckdb $(DB_URL) -dir migrations up

.PHONY: production/connect ## connects to production deployment server
production/connect:
	ssh ${SSH_USER}@${DEPLOY_TARGET_IP}

.PHONY: production/deploy ## deploys to production deployment server
production/deploy:
	$(MAKE) build
	rsync -P ./bin/linux_amd64/${PROJECT_NAME} ${SSH_USER}@${DEPLOY_TARGET_IP}:~
	rsync -P ./remote/${PROJECT_NAME}.service ${SSH_USER}@${DEPLOY_TARGET_IP}:~
	rsync -P ./remote/${PROJECT_NAME}.timer ${SSH_USER}@${DEPLOY_TARGET_IP}:~
	rsync -P .env ${SSH_USER}@${DEPLOY_TARGET_IP}:~
	ssh ${SSH_USER}@${DEPLOY_TARGET_IP} 'chmod 600 ~/.env'
	ssh -t ${SSH_USER}@${DEPLOY_TARGET_IP} '\
	  sudo mv ~/${PROJECT_NAME}.service /etc/systemd/system/ \
	  && sudo mv ~/${PROJECT_NAME}.timer /etc/systemd/system/ \
	  && sudo systemctl daemon-reload \
	  && sudo systemctl enable ${PROJECT_NAME}.timer \
	  && sudo systemctl restart ${PROJECT_NAME}.timer \
	'

.PHONY: production/logs ## gets the logs for the service
production/logs:
	ssh -t ${SSH_USER}@${DEPLOY_TARGET_IP} 'sudo journalctl -u ${PROJECT_NAME}.service'
