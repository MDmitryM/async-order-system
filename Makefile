.PHONY: build-api push-api clean-api

API_IMAGE_NAME = mdmitrym/order-system-api
#TAG по умолчанию всегда latest
TAG ?= latest

#сборка api сервера
build-api:
	docker build -f services/api/Dockerfile -t $(API_IMAGE_NAME):$(TAG) .
#пуш образа
push-api:
	docker push $(API_IMAGE_NAME):$(TAG)
#удалить образ api сервера по тегу
clean-api:
	docker rmi $(API_IMAGE_NAME):$(TAG)