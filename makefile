up:
	docker compose up -d

down:
	docker compose down --remove-orphans --volumes

materialize:
	docker compose exec pg psql -U materialize -d materialize -p 6875 -h materialize

up_pg:
	docker compose up -d pg

up_redis:
	docker compose up -d redis

up_analytics:
	docker compose up -d es kibana

psql:
	docker compose exec pg psql -U postgres -d postgres

kibana:
	open http://localhost:5601

es:
	open http://localhost:9200

redis_cli:
	docker compose exec redis redis-cli