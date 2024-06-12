up:
	docker compose up -d

down:
	docker compose down --remove-orphans --volumes

materialize:
	docker compose exec pg psql -U materialize -d materialize -p 6875 -h materialize

up_pg:
	docker compose up -d pg

up_analytics:
	docker compose up -d es kibana

psql:
	docker compose exec pg psql -U postgres -d postgres

kibana:
	open http://localhost:5601

es:
	open http://localhost:9200