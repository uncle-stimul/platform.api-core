all: add-healthcheck build

add-healthcheck:
	@echo "Загрузка утилиты platform.healthcheck"
	wget https://raw.githubusercontent.com/uncle-stimul/platform.healthcheck/refs/heads/main/healthcheck.go \
		-O ./healthcheck.go

build:
	@echo "Сборка docker образа platform.api-core"
	docker build -t platform.api-core:latest -f Dockerfile .

cleanup:
	@echo "Очистка сборочных данных platform.api-core"
	docker system prune --force