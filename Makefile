up:
	docker compose up -d
build:
	docker compose build --no-cache --force-rm
stop:
	docker compose stop
down:
	docker compose down --remove-orphans
restart:
	@make down
	@make up
rebuild:
	@make down
	@make build
	@make up
destroy:
	docker compose down --rmi all --volumes --remove-orphans
destroy-volumes:
	docker compose down --volumes --remove-orphans
ps:
	docker compose ps
logs:
	docker compose logs -f
app:
	docker compose exec app ash
amend:
	git commit --amend
fmt:
	docker compose exec app ash -c 'go fmt ./...'
run:
	docker compose exec app ash -c 'go run cmd/main.go'