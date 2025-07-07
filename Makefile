wire:
	wire gen ./cmd/api/
test-db:
	docker compose -f compose.test.yml up -d
test-down:
	docker compose -f compose.test.yml down
test:
	source .env.test && go test ./... -v
