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

Лаба 2: таблицы `dendro_researchers`, `dendro_scale_patterns`, `dendro_pattern_likes`. Создание и публикация через GORM, удаление на страницах — `SQL UPDATE` без ORM. Новые фото и видео на SSR не сохраняются; пустые ссылки заменяются `/static/defaults/`. Фильтр плитки — календарь `date_from` / `date_to`.

## Таблицы

`dendro_researchers` — исследователи

| столбец | тип | обязательно |
|---|---|---|
| id | BIGINT, PK | да |
| login | VARCHAR(64), уникальный | да |
| password | VARCHAR(128) | да |
| full_name | VARCHAR(128) | да |

`dendro_scale_patterns` — паттерны дендрошкалы

| столбец | тип | обязательно |
|---|---|---|
| id | BIGINT, PK | да |
| title | VARCHAR(160) | да |
| status | VARCHAR(20): draft, published, deleted | да |
| created_at | TIMESTAMPTZ | да |
| creator_id | BIGINT, FK на исследователя | да |
| image_url | VARCHAR(255) | нет, NULL |
| video_url | VARCHAR(255) | нет, NULL |
| short_description | VARCHAR(500) | нет, NULL до публикации |
| full_description | VARCHAR(2000) | нет, NULL |
| event_year | BIGINT, год начала | нет, NULL до публикации |
| ring_count | BIGINT, число колец | нет, NULL до публикации |
| formed_at | TIMESTAMPTZ | нет, NULL до публикации |

При публикации метод требует краткое описание, год начала и число колец. В таблице эти колонки остаются необязательными, потому что черновик создаётся раньше.

`dendro_pattern_likes` — лайки, связь многие-ко-многим

| столбец | тип | обязательно |
|---|---|---|
| id | BIGINT, PK | да |
| researcher_id | BIGINT, FK | да |
| pattern_id | BIGINT, FK | да |

Пара `(researcher_id, pattern_id)` уникальна. Удаление связей каскадом запрещено.

## Лаба 3–4, десять методов

Коллекция: `docs/dendro-age-10-methods.postman_collection.json`. Услуга = паттерн дендрошкалы, пользователь = исследователь.

### Домен паттерна дендрошкалы

- **GET** список с фильтрацией (только опубликованные). Признак 0/1, если создатель паттерна совпадает с текущим исследователем. Поле `mine`. `/api/dendro_scale_patterns?date_from=&date_to=`
- **GET** лента (только опубликованные). Признак 0/1, если текущий исследователь лайкнул этот паттерн. Поле `liked`. `/api/dendro_scale_patterns/feed`
- **GET** получения черновика (не более 1 записи для исследователя), ид не указывается. `/api/dendro_scale_patterns/draft`
- **POST** добавление + добавление файлов изображения и короткого видео (как файлы, а не просто URL). Названия файлов сохраняются в полях БД, сами файлы в MinIO. Названия генерируются на латинице. `/api/dendro_scale_patterns`
- **PUT** публикации (смена статуса). `/api/dendro_scale_patterns/publish`
- **DELETE** удаление (только soft delete). Только паттерны этого исследователя. `/api/dendro_scale_patterns/:id`
- **POST** like (от текущего исследователя). Поле 0/1 (`0` отменяет лайк, `1` ставит). `/api/dendro_scale_patterns/:id/like`

### Домен исследователя

- **POST** регистрация. `/api/dendro_researchers/register`
- **POST** аутентификация. `/api/dendro_researchers/login`
- **POST** деавторизация. `/api/dendro_researchers/logout`

## Лаба 4

JWT и чёрный список в Redis. Swagger: http://localhost:8088/swagger

Без токена доступны только GET списка и ленты, регистрация и вход.
Создание, черновик, публикация, удаление, лайк и выход требуют заголовок `Authorization: Bearer <access_token>`.
Создатель в этих методах берётся из JWT, не из константы.

```bash
docker compose up -d
go run ./cmd/dendro-age
```

Вход сида: логин `shebe`, пароль `shebe123`.
Redis проекта слушает на хосте порт 6380. Ключ после выхода: `docker exec dendro-age-redis redis-cli KEYS "dendro:jwt:blacklist:*"`

