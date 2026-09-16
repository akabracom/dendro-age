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

Лаба 2: таблицы `dendro_researchers`, `dendro_scale_patterns`, `dendro_pattern_likes`. Создание и публикация через GORM, удаление — `SQL UPDATE` без ORM. Новые фото/видео на сервер не сохраняются; пустые ссылки заменяются `/static/defaults/`.

