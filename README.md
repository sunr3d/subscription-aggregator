## Subscription Aggregator

Мини‑сервис учёта подписок (REST + Postgres). Сборка/запуск через Docker Compose.

### Быстрый старт


Запуск:
```bash
- make up
``` 
Требуются: `Docker Engine` + `Make`

Остановить: 
```bash
make down
```  

Сбросить БД/переиграть миграции: 
```bash
make clean && make up
```

### Переменные окружения (.env)

- HTTP_PORT=8080
- HTTP_TIMEOUT=30s
- LOG_LEVEL=info
- POSTGRES_HOST=db
- POSTGRES_PORT=5432
- POSTGRES_USER=postgres
- POSTGRES_PASSWORD=postgres
- POSTGRES_DB_NAME=postgres
- POSTGRES_SSL_MODE=disable
- POSTGRES_MIN_CONNS=1
- POSTGRES_MAX_CONNS=10
- POSTGRES_MAX_CONN_TTL=1h
- POSTGRES_HEALTH_CHECK_PERIOD=30s
- POSTGRES_PING_TIMEOUT=5s

### Команды Make

- `make up` — собрать и поднять
- `make down` — остановить
- `make clean` — остановить и удалить volume БД
- `make restart` - перезапустить
- `make logs` — логи приложения
- `make build` — собрать docker‑образ
- `make test` — юнит‑тесты домена

### API (коротко)

- POST /subscriptions — создать запись о подписке
- GET /subscriptions — список подписок по фильтру (?user_id, ?service_name, ?limit, ?offset)
- GET /subscriptions/{id} — получить запись по id
- PATCH /subscriptions/{id} — частичное обновление записи
- DELETE /subscriptions/{id} — удалить запись
- GET /subscriptions/total — сумма за период (?period_start, ?period_end, +фильтры по имени и сервису)
  
  
### ПОДРОБНАЯ SWAGGER ДОКУМЕНТАЦИЯ — `http://localhost:8081`.

### Архитектура

- `models/` — домен
- `internal/interfaces/*` — интерфейсы services/infra
- `internal/services/` — бизнес‑логика
- `internal/infra/` — Postgres адаптер
- `internal/api/` — HTTP‑хендлеры и DTO
- `internal/middleware/`, `internal/server/`, `internal/config/`, `internal/logger/`, `internal/entrypoint/`
