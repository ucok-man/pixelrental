include .env

# ------------------------------------------------------------------ #
#                                 app                                #
# ------------------------------------------------------------------ #
.PHONY: go/build
go/build:
	@echo "building app..."
	@go clean -cache
	@go build -o ./bin/snippetbox main.go
	@echo "done!"

.PHONY: go/run
go/run:
	@echo "starting app..."
	@go run main.go

.PHONY: go/start
go/start: go/build
	@echo "running app..."
	@./bin/snippetbox


# ------------------------------------------------------------------ #
#                             migrations                             #
# ------------------------------------------------------------------ #
DSN=${DB_DSN}

.PHONY: migrate/create
migrate/create:
	@migrate -path=migrations -database $(DSN) create -dir=migrations -ext=.sql -seq $(name)

.PHONY: migrate/goto
migrate/goto:
	@migrate -path=migrations -database $(DSN) goto $(version)

.PHONY: migrate/up
migrate/up:
	@migrate -path=migrations -database $(DSN) up

.PHONY: migrate/down
migrate/down:
	@migrate -path=migrations -database $(DSN) down
	
.PHONY: migrate/force
migrate/force:
	@migrate -path=migrations -database $(DSN) force $(version)

.PHONY: migrate/drop
migrate/drop:
	@migrate -path=migrations -database $(DSN) drop

.PHONY: migrate/version
migrate/version:
	@migrate -path=migrations -database $(DSN) version