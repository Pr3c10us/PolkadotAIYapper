include ./.env

define db_up
	    migrate -path ./stack/migrations -database "postgres://${PG_DB_USERNAME}:${PG_DB_PASSWORD}@${PG_DB_HOST}:${PG_DB_PORT}/${PG_DB_NAME}?sslmode=${PG_SSL_MODE}"  up
endef

define db_down
	    migrate -path ./stack/migrations -database "postgres://${PG_DB_USERNAME}:${PG_DB_PASSWORD}@${PG_DB_HOST}:${PG_DB_PORT}/${PG_DB_NAME}?sslmode=${PG_SSL_MODE}"  down
endef

define db_force
	    migrate -database "postgres://${PG_DB_USERNAME}:${PG_DB_PASSWORD}@${PG_DB_HOST}:${PG_DB_PORT}/${PG_DB_NAME}?sslmode=${PG_SSL_MODE}" -path ./stack/migrations force $(version)
endef

define db_up_test
	    migrate -path ./stack/migrations -database "postgres://${PG_DB_USERNAME_TEST}:${PG_DB_PASSWORD_TEST}@${PG_DB_HOST_TEST}:${PG_DB_PORT_TEST}/${PG_DB_NAME_TEST}?sslmode=${PG_SSL_MODE_TEST}"  up
endef

define db_down_test
	    migrate -path ./stack/migrations -database "postgres://${PG_DB_USERNAME_TEST}:${PG_DB_PASSWORD_TEST}@${PG_DB_HOST_TEST}:${PG_DB_PORT_TEST}/${PG_DB_NAME_TEST}?sslmode=${PG_SSL_MODE_TEST}"  down
endef

define db_force_test
	    migrate -database "postgres://${PG_DB_USERNAME_TEST}:${PG_DB_PASSWORD_TEST}@${PG_DB_HOST_TEST}:${PG_DB_PORT_TEST}/${PG_DB_NAME_TEST}?sslmode=${PG_SSL_MODE_TEST}" -path ./stack/migrations force $(version)
endef

define db_postgresql_create
	migrate create -ext sql -dir ./stack/migrations -seq $(migration)
endef

db_force:
	$(call db_force)

db_down:
	$(call db_down)

db_up:
	$(call db_up)

db_create:
	$(call db_postgresql_create)

generate_queries: pg_db_up
	sqlc generate -f stack/sqlc/sqlc.yml

dev:
	@refresh

build:
	@go build -o bin/main ./cmd/

run: build
	@bin/main

host:
	@ngrok http --domain=sponge-guided-minnow.ngrok-free.app 5000
