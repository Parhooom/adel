DATABASE_URL = "host=localhost user=admin password=adminpassword dbname=adel port=5432 sslmode=disable"
MIGRATION_DIR = migrations


.PHONY: migrate-up
migrate-up:
	goose -dir $(MIGRATION_DIR) postgres $(DATABASE_URL) up

.PHONY: migrate-down
migrate-down:
	goose -dir $(MIGRATION_DIR) postgres $(DATABASE_URL) down

.PHONY: migrate-up-to
migrate-up-to:
	goose -dir $(MIGRATION_DIR) postgres $(DATABASE_URL) up-to $(VERSION)

.PHONY:	migrate-down-to
migrate-down-to:
	goose -dir $(MIGRATION_DIR) postgres $(DATABASE_URL) down-to $(VERSION)

.PHONY: migrate-reset
migrate-reset:
	goose -dir $(MIGRATION_DIR) postgres $(DATABASE_URL) reset 

.PHONY: migrate-status
migrate-status:
	goose -dir $(MIGRATION_DIR) postgres $(DATABASE_URL) status

.PHONY: swagger
swagger:
	swag init -g main.go --output docs
