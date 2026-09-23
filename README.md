# Marketplace

Микросервисный backend интернет-магазина на Go.

## Архитектура

Проект состоит из 4 сервисов:

* **Gateway** — внешний HTTP API
* **Order** — создание и управление заказами
* **Product** — товары и остатки
* **Analytics** — статистика и аналитика

Внутри сервисы общаются через **gRPC**.

Каждый сервис имеет свою PostgreSQL базу данных.

Redis используется в Analytics для кеширования статистики.

## Технологии

* Go
* PostgreSQL 15
* Redis 7
* gRPC
* Protocol Buffers
* Docker 
* pgx
* go-redis
* testify
* miniredis

## API

Gateway:

| Метод  | Endpoint                   | Назначение         |
| ------ | -------------------------- | ------------------ |
| GET    | `/ping`                    | Проверка сервиса   |
| GET    | `/products`                | Список товаров     |
| GET    | `/products/{id}`           | Получить товар     |
| POST   | `/products`                | Создать товар      |
| POST   | `/orders`                  | Создать заказ      |
| GET    | `/orders/{id}`             | Получить заказ     |
| DELETE | `/orders/{id}`             | Отменить заказ     |
| GET    | `/analytics/revenue`       | Выручка            |
| GET    | `/analytics/orders`        | Количество заказов |
| GET    | `/analytics/average-check` | Средний чек        |
| GET    | `/analytics/top-products`  | Топ товаров        |

Gateway доступен на:

```text
http://localhost:8080
```

## Создание заказа

При создании заказа:

1. Gateway передаёт запрос в Order.
2. Order проверяет товар через Product.
3. Product уменьшает остаток.
4. Order сохраняется в PostgreSQL.
5. Order передаёт данные в Analytics.
6. Analytics сохраняет данные в свою PostgreSQL.
7. Статистика кешируется в Redis.

## Базы данных

У каждого сервиса отдельная база:

| Сервис    | База       | Порт |
| --------- | ---------- | ---: |
| Order     | PostgreSQL | 5433 |
| Product   | PostgreSQL | 5434 |
| Analytics | PostgreSQL | 5435 |

Сервисы не обращаются напрямую к базе данных друг друга.

## Redis

Analytics кеширует:

* выручку
* количество заказов
* средний чек
* топ товаров

При создании нового заказа соответствующий кеш инвалидируется.

## Тестирование

Для сервисов написаны:

* unit-тесты
* тесты сервисного слоя
* тесты HTTP handlers
* тесты gRPC
* интеграционные-тесты PostgreSQL
* тесты Redis через miniredis

Запуск тестов:

```powershell
cd gateway
go test ./...

cd ../order
go test ./...

cd ../product
go test ./...

cd ../analytics
go test ./...
```

## Docker

Запуск всего проекта:

```powershell
docker compose up -d --build
```

Проверить контейнеры:

```powershell
docker compose ps
```

Остановить:

```powershell
docker compose down
```

Основные порты:

```text
Gateway      8080
Order        8081
Product      8082
Analytics    8083
Redis        6379
```

## Структура

```text
marketplace/
├── gateway/
├── order/
├── product/
├── analytics/
├── docker-compose.yml
└── README.md
```
