run-db:
	docker run --name db -d -p 3306:3306 -e MYSQL_USER=mysql -e MYSQL_PASSWORD=password -e MYSQL_ROOT_PASSWORD=password -e MYSQL_DATABASE=currencies mysql

create-migrations:
	migrate create -dir migrations -seq -ext sql $(MIGRATION_NAME)

up-migrations:
	migrate --source file://migrations/ -database $(DSN) up

down-migrations:
	migrate --source file://migrations/ -database $(DSN) down

.PHONY: run-db create-migrations up-migrations down-migrations
