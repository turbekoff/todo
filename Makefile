.PHONY: local development
.SILENT:


run:
	export CONFIG_PATH=/configs/api/local.env && go run ./cmd/api/


local:
	echo "Starting local environment"
	docker-compose -f docker-compose.local.yml up --build

development:
	echo "Starting docker environment"
	docker-compose -f docker-compose.development.yml up --build

clean:
	docker system prune -f
