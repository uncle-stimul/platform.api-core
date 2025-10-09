all: add-healthcheck build-auto

add-healthcheck:
	@echo "Загрузка утилиты platform.healthcheck"
	wget https://raw.githubusercontent.com/uncle-stimul/platform.healthcheck/refs/heads/main/healthcheck.go \
		-O ./healthcheck.go

build-auto:
	@echo "Сборка docker образа platform.api-core"
	docker build -t platform.api-core:latest -f Dockerfile .