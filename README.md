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

