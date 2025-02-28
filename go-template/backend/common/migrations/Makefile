.PHONY: migrationdir
MIGRATION_NAME := $(shell date  +%Y%m%d%H%M%S)_$(name)
migrationdir: migrationfile
	mkdir -p backend/common/migrations/$(MIGRATION_NAME)

.PHONY: migrationfile
migrationfile:
	cp backend/common/migrations/template.txt backend/common/migrations/$(MIGRATION_NAME).go

.PHONY: migrate
migrate:
	cd backend && go run pkg/cmd/main.go init
	cd backend && go run pkg/cmd/main.go migrate
	cd backend && go run pkg/cmd/main.go status

.PHONY: rollback
rollback:
	cd backend && go run pkg/cmd/main.go rollback

.PHONY: seed
seed:
	cd backend && go run pkg/cmd/main.go seed

.PHONY: resetdb
resetdb: rollback migrate seed
