# store-service

Сервис для управления каталогом, заказами и отчетами (Go + chi + pgx + zap).

## Стек
- Go 1.22, chi, pgx, zap
- envconfig для конфигурации
- Docker / docker-compose для локального окружения
- Postgres 16, pgAdmin для UI
- migrate/migrate для миграций, seed через psql

## Быстрый старт без Docker
```bash
export POSTGRES_DSN="postgres://user:pass@localhost:5432/store?sslmode=disable"
export LOG_LEVEL=info
export HTTP_ADDR=:8080
go run ./cmd/store-service
```

## Локально через Docker Compose
Из корня репозитория:
```bash
cd deploy/debug
# подготовить .env (можно скопировать пример)
cp .env.example .env

# запуск
docker compose --env-file .env up --build          # в форграунде
# или
docker compose --env-file .env up --build -d       # в фоне
```

Сервисы:
- `postgres` — база
- `migrate` — применяет миграции из `deploy/debug/migrations`
- `seed` — наполняет тестовыми данными из `deploy/debug/seed/seed.sql`
- `store-service` — приложение
- `pgadmin` — UI на http://localhost:5050 (логин/пароль в `.env`)
- `docs` (внутри бинаря) — Swagger UI по /docs, openapi.yaml по /docs/openapi.yaml

Порты по умолчанию:
- API: `http://localhost:8080`
- Postgres: `5432`
- pgAdmin: `5050`

Остановка и очистка:
```bash
docker compose --env-file .env down
docker compose --env-file .env down -v   # с удалением данных/pgadmin volume
```

## Основные ручки
- `GET /healthz`
- Категории: `GET/POST /categories`, `GET/PUT/DELETE /categories/{id}`
- Клиенты: `GET/POST /customers`, `GET/PUT/DELETE /customers/{id}`
- Товары: `GET/POST /products`, `GET/PUT/DELETE /products/{id}`
- Заказы: `GET/POST /orders`, `GET/PUT/DELETE /orders/{id}`, `POST /orders/{id}/items`
- Отчеты:
  - `GET /reports/customer-totals`
  - `GET /reports/category-children`
  - `GET /reports/top-products-last-month`

## Миграции и сиды вручную
```bash
# миграции
migrate -path deploy/debug/migrations -database "$POSTGRES_DSN" up
# сиды
psql "$POSTGRES_DSN" -f deploy/debug/seed/seed.sql
```

## Документация
- Swagger UI: `http://localhost:8080/docs`
- OpenAPI: `http://localhost:8080/docs/openapi.yaml`
- Postman: `docs/postman_collection.json`

## Конфигурация (.env)
См. `deploy/debug/.env`. Ключевые переменные:
- `POSTGRES_DSN` — строка подключения
- `HTTP_ADDR` — адрес HTTP сервера
- `LOG_LEVEL` — уровень логирования
- `PGADMIN_DEFAULT_EMAIL` / `PGADMIN_DEFAULT_PASSWORD` — доступ в pgAdmin

