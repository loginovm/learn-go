
.PHONY: run-img-db
run-img-db:
	docker run -d --name pg \
	-e POSTGRES_USER=otus_user \
	-e POSTGRES_PASSWORD=dev_pass \
	-e POSTGRES_DB=calendar \
	-e PGDATA=/var/lib/postgresql/data/pgdata \
	-v pg_data_mloginov:/var/lib/postgresql/data \
	-p 5432:5432 \
	postgres:14