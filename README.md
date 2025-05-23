# Асинхронная Система Заказов

REST API сервис для управления заказами, написанный на Go с использованием микросервисной архитектуры. Система включает API сервис, сервис биллинга и сервис доставки, взаимодействующие через Kafka.

## Технологии

- Go
- Fiber (веб-фреймворк)
- PostgreSQL (pgxpool, sqlc)
- Apache Kafka
- Docker & Docker Compose
- Swagger (документация API)

## Предварительные требования

- Docker
- Docker Compose
- Go 1.23 или выше

## Архитектура

Система состоит из трех основных микросервисов:
- **API сервис** - основной REST API для управления заказами
- **Billing сервис** - обработка платежей
- **Shipping сервис** - управление доставкой

Микросервисы взаимодействуют через Kafka, что обеспечивает асинхронную обработку и отказоустойчивость.

## API Documentation

Полная API документация доступна через Swagger UI по адресу: `http://localhost:8080/swagger/`

### Доступные эндпоинты

#### Заказы (Orders)

| Метод | Эндпоинт      | Описание                    |
|-------|---------------|----------------------------|
| POST  | /order/create/ | Создание нового заказа     |
| GET   | /order/details/ | Получение информации о заказе по ID |
| DELETE| /order/delete/ | Удаление заказа по ID      |
| GET   | /order/list/   | Получение списка заказов с пагинацией |

## Установка и запуск

### Использование Docker Compose

1. Клонируйте репозиторий:
```bash
git clone https://github.com/MDmitryM/async-order-system.git
cd async-order-system
```

2. Создайте файл `.env` в директории `services/api/`:
```env
# Database
PORT=8080
POSTGRES_USER=your_pg_user
POSTGRES_PASSWORD=your_pg_pwd
POSTGRES_DB=your_pg_db
API_DB_PORT=5432
API_DB_SSL_MODE=disable
```

3. Запустите приложение:
```bash
docker-compose up -d
```

Приложение будет доступно по адресу `http://localhost:8080`

### Переменные окружения

| Переменная    | Описание                           |
|---------------|------------------------------------|
| POSTGRES_USER | Имя пользователя PostgreSQL        |
| POSTGRES_PASSWORD | Пароль для базы данных         |
| POSTGRES_DB   | Название базы данных               |
| API_DB_HOST   | Хост базы данных                   |
| ENV           | Окружение (development/production) |

### Порты

- 8080: API сервер
- 5432: PostgreSQL
- 9091-9093: Kafka брокеры
- 9020: Kafka UI
- 2181: Zookeeper

## Структура проекта

- `/services` - микросервисы
  - `/api` - основной API сервис
    - `/handler` - обработчики HTTP запросов
    - `/repository` - работа с базой данных
    - `/kafka` - взаимодействие с Kafka
  - `/billing` - сервис биллинга
  - `/shipping` - сервис доставки

## База данных

PostgreSQL база данных автоматически инициализируется при первом запуске. Данные сохраняются в Docker volume.

## Kafka

Система использует кластер Kafka с тремя брокерами для обеспечения отказоустойчивости. Создаются следующие топики:
- `orders` - для сообщений о заказах
- `payments` - для сообщений о платежах
- `shipping` - для сообщений о доставке

Kafka UI доступен по адресу `http://localhost:9020` для мониторинга сообщений и топиков.
