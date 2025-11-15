.PHONY: run swagger

run:
	docker compose up --build

swagger:
	docker compose -f swagger/compose.yml up -d