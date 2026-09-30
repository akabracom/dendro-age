# Dendro Age

Web application for analyzing anomalous dendrochronological patterns and estimating the age of findings.

## Technology

- Go
- Gin
- HTML templates
- CSS
- MinIO
- Docker Compose

## Lab 1

Три GET-страницы паттернов дендрошкалы (`dendro_scale_patterns`): лента, черновик, плитка.

```bash
docker compose up -d
go run ./cmd/dendro-age
```

- Лента: http://localhost:8088/dendro_scale_patterns/feed
- Черновик: http://localhost:8088/dendro_scale_patterns/draft
- Плитка: http://localhost:8088/dendro_scale_patterns
- MinIO: http://localhost:9001 (логин `minioadmin`)
- Adminer: http://localhost:8081 — System `PostgreSQL`, Server `postgres`, User `dendro`, Password `dendro1234`, Database `dendro_age` (Postgres с хоста: порт `5433`)

Лаба 2: таблицы `dendro_researchers`, `dendro_scale_patterns`, `dendro_pattern_likes`. Создание и публикация через GORM, удаление — `SQL UPDATE` без ORM. Новые фото/видео на SSR не сохраняются; пустые ссылки заменяются `/static/defaults/`. Фильтр плитки — календарь `date_from`/`date_to` (год события берётся из даты).

Лаба 3 REST API (`/api`):

| Метод | URL |
|-------|-----|
| GET | `/api/dendro_scale_patterns?date_from=&date_to=` |
| GET | `/api/dendro_scale_patterns/feed` |
| GET | `/api/dendro_scale_patterns/feed/:id?next=true` |
| GET | `/api/dendro_scale_patterns/draft` |
| POST | `/api/dendro_scale_patterns` (multipart: title, image, video → MinIO) |
| PUT | `/api/dendro_scale_patterns/publish` |
| DELETE | `/api/dendro_scale_patterns/:id` |
| POST | `/api/dendro_scale_patterns/:id/like` body `{"like":1}` |
| POST | `/api/dendro_researchers/register` |
| POST | `/api/dendro_researchers/login` (заглушка) |
| POST | `/api/dendro_researchers/logout` (заглушка) |

Создатель во всех методах — singleton `currentresearcher.Login()` = `shebe`.

