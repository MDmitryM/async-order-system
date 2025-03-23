.PHONY: build-api push-api clean-api

API_IMAGE_NAME = mdmitrym/order-system-api
BILLING_IMAGE_NAME = mdmitrym/order-system-bill
#TAG по умолчанию всегда latest
TAG ?= latest

#сборка api сервера
build-api:
	docker build -f services/api/Dockerfile -t $(API_IMAGE_NAME):$(TAG) .
#пуш api образа
push-api:
	docker push $(API_IMAGE_NAME):$(TAG)
#удалить образ api сервера по тегу
clean-api:
	docker rmi $(API_IMAGE_NAME):$(TAG)

#сборка биллинга
build-bill:
	docker build -f services/billing/Dockerfile -t $(BILLING_IMAGE_NAME):$(TAG) .
#пуш биллинг сервиса
push-bill:
	docker push $(BILLING_IMAGE_NAME):$(TAG)
#удалить образ биллинга по тегу
clean-bill:
	docker rmi $(BILLING_IMAGE_NAME):$(TAG)

