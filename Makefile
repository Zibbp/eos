dev-setup:
	go install github.com/sqlc-dev/sqlc/cmd/sqlc@latest
	go install github.com/air-verse/air@latest
	go install github.com/joho/godotenv/cmd/godotenv@latest
	go install github.com/a-h/templ/cmd/templ@latest
	npm i

dev-server:
	rm -f ./tmp/server
	air -c ./.server.air.toml

dev-worker:
	rm -f ./tmp/worker
	air -c ./.worker.air.toml

migration-sql: ##@Development Create an empty SQL migration file
	@if [ -z "${name}" ]; then echo "Usage: make migration-sql name=name_of_migration_file"; exit 1; fi
	go run github.com/pressly/goose/v3/cmd/goose@latest -dir internal/db/migrations create ${name} sql

sqlc-generate:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest generate

sqlc-lint:
	go run github.com/sqlc-dev/sqlc/cmd/sqlc@latest vet

# get goose status using env vars from .env
goose-status:
	@export $(shell grep -v '^#' .env | xargs) && \
	go run github.com/pressly/goose/v3/cmd/goose@latest -dir internal/db/migrations postgres "host=$$DB_HOST user=$$DB_USER password=$$DB_PASS dbname=$$DB_NAME sslmode=disable" status

goose-reset:
	@export $(shell grep -v '^#' .env | xargs) && \
	go run github.com/pressly/goose/v3/cmd/goose@latest -dir internal/db/migrations postgres "host=$$DB_HOST user=$$DB_USER password=$$DB_PASS dbname=$$DB_NAME sslmode=disable" reset

db-reset:
	@export $(shell grep -v '^#' .env | xargs) && \
	PGPASSWORD=$$DB_PASS psql -h localhost -p $$DB_PORT -U $$DB_USER -c "DROP SCHEMA public CASCADE; CREATE SCHEMA public; GRANT ALL ON SCHEMA public TO postgres; GRANT ALL ON SCHEMA public TO public;"

river-webui:
	curl -L https://github.com/riverqueue/riverui/releases/latest/download/riverui_linux_amd64.gz | gzip -d > /tmp/riverui
	chmod +x /tmp/riverui
	@export $(shell grep -v '^#' .env | xargs) && \
	VITE_RIVER_API_BASE_URL=http://localhost:8080/api DATABASE_URL=postgres://$$DB_USER:$$DB_PASS@localhost:$$DB_PORT/$$DB_NAME /tmp/riverui

templ-generate:
	templ generate

templ-watch:
	templ generate --watch --proxy="http://localhost:3000" --proxybind="0.0.0.0" --open-browser=false

watch-assets-css:
	npx -y postcss-cli internal/assets/app.css -o public/assets/styles.css --watch
build-assets-css:
	npx -y postcss-cli internal/assets/app.css -o public/assets/styles.css

watch-assets-js:
	npx -y esbuild internal/assets/index.js --format=esm --bundle --outdir=public/assets --watch=forever
build-assets-js:
	npx -y esbuild internal/assets/index.js --format=esm --bundle --outdir=public/assets --minify

dev-web:
	make -j4 templ-watch dev-server watch-assets-css watch-assets-js

nfs-mount:
	sudo mkdir -p /mnt/videos
	@export $(shell grep -v '^#' .env | xargs) && \
	sudo mount -t nfs -o nolock $$NFS_SHARE /mnt/videos -vv