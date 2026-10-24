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

## Лаба 3, десять методов

Коллекция: `docs/dendro-age-10-methods.postman_collection.json`. Создатель в API берётся из JWT. HTML-страницы по-прежнему открываются от `shebe`.

| № | Метод | URL | Что делает |
|---|---|---|---|
| 1 | GET | `/api/dendro_scale_patterns?date_from=&date_to=` | Список опубликованных. Фильтр по году из календарных дат. Удалённые не отдаются |
| 2 | POST | `/api/dendro_scale_patterns` | Черновик. Multipart: `title`, `image`, `video`. Файлы в MinIO, имена латиницей |
| 3 | GET | `/api/dendro_scale_patterns/draft` | Черновик текущего создателя, id в адресе нет |
| 4 | PUT | `/api/dendro_scale_patterns/publish` | Публикация: краткое описание, год начала, число колец. Статус ставит сервер |
| 5 | GET | `/api/dendro_scale_patterns/feed` | Лента без id: первый опубликованный |
| 6 | GET | `/api/dendro_scale_patterns/feed/:id?next=true` | Лента с этого id или следующий при `next=true` |
| 7 | POST | `/api/dendro_scale_patterns/:id/like` | Тело `{"like":1}` ставит лайк, `{"like":0}` снимает |
| 8 | DELETE | `/api/dendro_scale_patterns/:id` | Логическое удаление, статус `deleted` |
| 9 | POST | `/api/dendro_researchers/register` | Регистрация: `login`, `password`, `full_name` |
| 10 | POST | `/api/dendro_researchers/login` | Вход: `login` и `password`. Ответ — JWT |

В ленте и в списке два разных числа: `like_count` — сколько лайков всего, `liked` — `0` или `1`, лайкнул ли исследователь из токена. Без токена `liked` равен `0`.

Деавторизация: `POST /api/dendro_researchers/logout` — токен попадает в Redis blacklist.

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

