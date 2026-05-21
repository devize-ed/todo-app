# ToDo List REST API — полное техническое задание

## 1. Назначение проекта

Необходимо разработать учебное backend-приложение **ToDo List REST API** на Go для управления пользователями, задачами и получения статистики по выполнению задач.

Приложение должно предоставлять HTTP REST API, хранить данные в PostgreSQL, запускаться локально через Docker Compose, поддерживать миграции базы данных, конфигурацию через переменные окружения, логирование и удобные команды запуска через Makefile.

Проект является учебным, но должен быть реализован в стиле, близком к production-подходу: разделение слоёв, централизованная обработка ошибок, валидация входных данных, транзакции там, где они нужны, понятная структура проекта и покрытие ключевой бизнес-логики тестами.

## 2. Цели проекта

### Основные цели

1. Реализовать REST API для CRUD-операций над пользователями.
2. Реализовать REST API для CRUD-операций над задачами.
3. Связать задачи с пользователями-авторами.
4. Реализовать аналитический endpoint для статистики по задачам.
5. Настроить PostgreSQL и миграции.
6. Настроить Docker Compose для локального окружения.
7. Настроить Makefile для типовых команд разработки.
8. Реализовать логирование HTTP-запросов и ошибок.
9. Добавить базовую валидацию входных данных.
10. Добавить тесты для сервисного слоя и/или repository-слоя.

### Не входит в первую версию

1. Регистрация и логин.
2. JWT, sessions, OAuth или другая авторизация.
3. Роли пользователей.
4. Frontend.
5. Email-уведомления.
6. Очереди сообщений.
7. Redis-кеширование.
8. Полнотекстовый поиск.
9. Микросервисная архитектура.

## 3. Технологический стек

| Категория           | Технология                                             |
| ------------------- | ------------------------------------------------------ |
| Язык                | Go                                                     |
| API                 | REST API over HTTP                                     |
| База данных         | PostgreSQL                                             |
| Миграции            | golang-migrate или аналогичный инструмент              |
| Конфигурация        | Переменные окружения + `.env` для локальной разработки |
| Логирование         | `log/slog`, `zap`, `zerolog` или аналог                |
| Роутинг             | `net/http`, `chi`, `gorilla/mux` или аналог            |
| Работа с PostgreSQL | `pgx`, `database/sql` + `pgx` driver или аналог        |
| Контейнеризация     | Docker, Docker Compose                                 |
| Сборка и команды    | Makefile                                               |
| Тесты               | `testing`, `testify`, при необходимости testcontainers |

## 4. Общая архитектура

Рекомендуемая структура проекта:

```text
.
├── cmd/
│   └── todoapp/
│       └── main.go
├── internal/
│   ├── app/
│   │   └── app.go
│   ├── config/
│   │   └── config.go
│   ├── domain/
│   │   ├── user.go
│   │   ├── task.go
│   │   └── statistics.go
│   ├── service/
│   │   ├── user_service.go
│   │   ├── task_service.go
│   │   └── statistics_service.go
│   ├── repository/
│   │   ├── postgres/
│   │   │   ├── user_repository.go
│   │   │   ├── task_repository.go
│   │   │   └── statistics_repository.go
│   │   └── errors.go
│   ├── transport/
│   │   └── http/
│   │       ├── router.go
│   │       ├── user_handler.go
│   │       ├── task_handler.go
│   │       ├── statistics_handler.go
│   │       ├── middleware.go
│   │       └── response.go
│   └── validator/
│       └── validator.go
├── migrations/
│   ├── 000001_init.up.sql
│   └── 000001_init.down.sql
├── docker-compose.yml
├── Dockerfile
├── Makefile
├── .env.example
├── go.mod
└── README.md
```

### Слои приложения

| Слой             | Ответственность                                       |
| ---------------- | ----------------------------------------------------- |
| `transport/http` | HTTP handlers, парсинг request, формирование response |
| `service`        | Бизнес-логика, бизнес-правила, координация repository |
| `repository`     | Работа с PostgreSQL                                   |
| `domain`         | Доменные модели и типы                                |
| `config`         | Загрузка конфигурации                                 |
| `validator`      | Проверка входных данных                               |

Handlers не должны содержать SQL. Repository не должен знать об HTTP. Service не должен зависеть от конкретного HTTP-фреймворка.

## 5. Доменные сущности

## 5.1 User

Пользователь — автор задач.

### Поля

| Поле         | Тип                      | Обязательное | Описание                              |
| ------------ | ------------------------ | -----------: | ------------------------------------- |
| `id`         | UUID                     |           Да | Уникальный идентификатор пользователя |
| `name`       | string                   |           Да | Имя пользователя                      |
| `email`      | string                   |           Да | Email пользователя                    |
| `created_at` | timestamp with time zone |           Да | Дата создания                         |
| `updated_at` | timestamp with time zone |           Да | Дата последнего обновления            |

### Валидация `name`

1. Поле обязательно.
2. Минимальная длина: 3 символа.
3. Максимальная длина: 100 символов.
4. Строка не может состоять только из пробелов.
5. Перед сохранением пробелы в начале и конце должны быть удалены.

### Валидация `email`

1. Поле обязательно.
2. Email должен быть уникальным среди пользователей.
3. Минимальная длина: 6 символов.
4. Максимальная длина: 255 символов.
5. Строка не может состоять только из пробелов.
6. Перед сохранением пробелы в начале и конце должны быть удалены.
7. Email должен соответствовать базовому формату `local-part@domain`, например `user@example.com`.
8. Email должен содержать ровно один символ `@`.
9. После `@` должен быть домен с точкой.
10. В учебной версии достаточно простой проверки формата; полноценную RFC-валидацию делать не требуется.

Примеры валидных email:

```text
ivan@example.com
daniil.chirkov@example.de
user+test@gmail.com
```

Примеры невалидных email:

```text
ivan
ivan@
@example.com
ivan@example
ivan example@example.com
```

## 5.2 Task

Задача — элемент ToDo-листа, привязанный к пользователю.

### Поля

| Поле           | Тип                             | Обязательное | Описание                          |
| -------------- | ------------------------------- | -----------: | --------------------------------- |
| `id`           | UUID                            |           Да | Уникальный идентификатор задачи   |
| `user_id`      | UUID                            |           Да | Идентификатор пользователя-автора |
| `title`        | string                          |           Да | Заголовок задачи                  |
| `description`  | string / null                   |          Нет | Описание задачи                   |
| `completed`    | boolean                         |           Да | Статус выполнения                 |
| `created_at`   | timestamp with time zone        |           Да | Дата создания                     |
| `updated_at`   | timestamp with time zone        |           Да | Дата последнего обновления        |
| `completed_at` | timestamp with time zone / null |          Нет | Дата завершения задачи            |

### Валидация `title`

1. Поле обязательно.
2. Минимальная длина: 1 символ.
3. Максимальная длина: 100 символов.
4. Строка не может состоять только из пробелов.
5. Перед сохранением пробелы в начале и конце должны быть удалены.

### Валидация `description`

1. Поле опциональное.
2. Если поле указано, минимальная длина: 1 символ.
3. Если поле указано, максимальная длина: 1000 символов.
4. Строка не может состоять только из пробелов.
5. Перед сохранением пробелы в начале и конце должны быть удалены.

### Бизнес-правила `completed` и `completed_at`

1. Если `completed = true`, то `completed_at` должен быть заполнен.
2. Если `completed = true`, то `completed_at >= created_at`.
3. Если `completed = false`, то `completed_at` должен быть `null`.
4. Клиент не должен напрямую управлять `completed_at`.
5. При переводе задачи из `completed = false` в `completed = true` приложение автоматически выставляет `completed_at = now()`.
6. При переводе задачи из `completed = true` в `completed = false` приложение автоматически сбрасывает `completed_at = null`.
7. Если задача уже выполнена и клиент повторно отправляет `completed = true`, приложение не должно менять `completed_at`, если остальные поля не требуют изменения. Это защищает исходное время завершения от случайной перезаписи.

## 6. Схема базы данных

## 6.1 PostgreSQL schema

Рекомендуется создать отдельную схему:

```sql
CREATE SCHEMA IF NOT EXISTS todoapp;
```

## 6.2 Таблица `users`

```sql
CREATE TABLE IF NOT EXISTS todoapp.users (
    id          UUID PRIMARY KEY,
    name        VARCHAR(100) NOT NULL CHECK (char_length(trim(name)) >= 3),
    email       VARCHAR(255) NOT NULL UNIQUE,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT now(),

    CONSTRAINT users_email_format_check CHECK (
        char_length(trim(email)) >= 6
        AND email ~* '^[A-Z0-9._%+\-]+@[A-Z0-9.\-]+\.[A-Z]{2,}$'
    )
);
```

Примечание: регулярное выражение для email в БД является упрощённой проверкой учебного проекта. Основная пользовательская валидация должна выполняться в приложении.

## 6.3 Таблица `tasks`

```sql
CREATE TABLE IF NOT EXISTS todoapp.tasks (
    id            UUID PRIMARY KEY,
    user_id       UUID NOT NULL REFERENCES todoapp.users(id) ON DELETE CASCADE,
    title         VARCHAR(100) NOT NULL CHECK (char_length(trim(title)) >= 1),
    description   VARCHAR(1000),
    completed     BOOLEAN NOT NULL DEFAULT false,
    created_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    updated_at    TIMESTAMPTZ NOT NULL DEFAULT now(),
    completed_at  TIMESTAMPTZ,

    CONSTRAINT tasks_description_check CHECK (
        description IS NULL OR char_length(trim(description)) >= 1
    ),

    CONSTRAINT tasks_completed_at_check CHECK (
        (completed = false AND completed_at IS NULL)
        OR
        (completed = true AND completed_at IS NOT NULL AND completed_at >= created_at)
    )
);
```

## 6.4 Индексы

```sql
CREATE INDEX IF NOT EXISTS idx_tasks_user_id
    ON todoapp.tasks(user_id);

CREATE INDEX IF NOT EXISTS idx_tasks_created_at
    ON todoapp.tasks(created_at);

CREATE INDEX IF NOT EXISTS idx_tasks_completed
    ON todoapp.tasks(completed);

CREATE INDEX IF NOT EXISTS idx_tasks_completed_at
    ON todoapp.tasks(completed_at);

CREATE INDEX IF NOT EXISTS idx_tasks_user_created_at
    ON todoapp.tasks(user_id, created_at);
```

## 6.5 `updated_at`

В учебной версии можно обновлять `updated_at` из приложения при каждом `PATCH`.

В расширенной версии можно добавить PostgreSQL trigger:

```sql
CREATE OR REPLACE FUNCTION todoapp.set_updated_at()
RETURNS TRIGGER AS $$
BEGIN
    NEW.updated_at = now();
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER users_set_updated_at
BEFORE UPDATE ON todoapp.users
FOR EACH ROW
EXECUTE FUNCTION todoapp.set_updated_at();

CREATE TRIGGER tasks_set_updated_at
BEFORE UPDATE ON todoapp.tasks
FOR EACH ROW
EXECUTE FUNCTION todoapp.set_updated_at();
```

## 7. REST API

## 7.1 Общие правила API

### Base URL

```text
/api/v1
```

### Content-Type

Все запросы с телом должны использовать:

```http
Content-Type: application/json
```

Все ответы возвращаются в формате JSON.

### Формат успешного ответа для одного объекта

```json
{
  "data": {}
}
```

### Формат успешного ответа для списка

```json
{
  "data": [],
  "pagination": {
    "limit": 20,
    "offset": 0,
    "total": 100
  }
}
```

### Формат ошибки

```json
{
  "error": {
    "code": "validation_error",
    "message": "invalid request body",
    "details": [
      {
        "field": "name",
        "message": "must be at least 3 characters"
      }
    ]
  }
}
```

### HTTP status codes

|                         Код | Когда используется                                       |
| --------------------------: | -------------------------------------------------------- |
|                    `200 OK` | Успешное получение или изменение ресурса                 |
|               `201 Created` | Успешное создание ресурса                                |
|            `204 No Content` | Успешное удаление ресурса                                |
|           `400 Bad Request` | Невалидный JSON, неправильные query-параметры            |
|             `404 Not Found` | Ресурс не найден                                         |
|              `409 Conflict` | Конфликт бизнес-правил или ограничений БД                |
|  `422 Unprocessable Entity` | Данные синтаксически корректны, но не проходят валидацию |
| `500 Internal Server Error` | Неожиданная ошибка сервера                               |

## 7.2 Pagination

Для endpoints, возвращающих список, поддерживаются query-параметры:

| Параметр | Тип | Default | Min | Max | Описание           |
| -------- | --- | ------: | --: | --: | ------------------ |
| `limit`  | int |      20 |   1 | 100 | Количество записей |
| `offset` | int |       0 |   0 |   — | Смещение           |

Пример:

```http
GET /api/v1/tasks?limit=20&offset=40
```

Если `limit` или `offset` имеют нечисловое значение, API возвращает `400 Bad Request`.

## 8. Users API

## 8.1 Создать пользователя

```http
POST /api/v1/users
```

### Request body

```json
{
  "name": "Ivan Petrov",
  "email": "ivan.petrov@example.com"
}
```

Поле `email` обязательно.

### Response `201 Created`

```json
{
  "data": {
    "id": "2cb8ff0e-7b2a-4dbb-9283-e98b4a2f6d1a",
    "name": "Ivan Petrov",
    "email": "ivan.petrov@example.com",
    "created_at": "2026-05-07T10:00:00Z",
    "updated_at": "2026-05-07T10:00:00Z"
  }
}
```

### Ошибки

| Ситуация                                      |   Код |
| --------------------------------------------- | ----: |
| Невалидный JSON                               | `400` |
| `name` отсутствует                            | `422` |
| `name` короче 3 символов                      | `422` |
| `email` отсутствует                           | `422` |
| `email` не соответствует формату              | `422` |
| `email` уже используется другим пользователем | `409` |

## 8.2 Получить пользователя по ID

```http
GET /api/v1/users/{user_id}
```

### Response `200 OK`

```json
{
  "data": {
    "id": "2cb8ff0e-7b2a-4dbb-9283-e98b4a2f6d1a",
    "name": "Ivan Petrov",
    "email": "ivan.petrov@example.com",
    "created_at": "2026-05-07T10:00:00Z",
    "updated_at": "2026-05-07T10:00:00Z"
  }
}
```

### Ошибки

| Ситуация               |   Код |
| ---------------------- | ----: |
| `user_id` не UUID      | `400` |
| Пользователь не найден | `404` |

## 8.3 Получить список пользователей

```http
GET /api/v1/users?limit=20&offset=0
```

### Response `200 OK`

```json
{
  "data": [
    {
      "id": "2cb8ff0e-7b2a-4dbb-9283-e98b4a2f6d1a",
      "name": "Ivan Petrov",
      "email": "ivan.petrov@example.com",
      "created_at": "2026-05-07T10:00:00Z",
      "updated_at": "2026-05-07T10:00:00Z"
    }
  ],
  "pagination": {
    "limit": 20,
    "offset": 0,
    "total": 1
  }
}
```

## 8.4 Обновить пользователя

```http
PATCH /api/v1/users/{user_id}
```

PATCH должен поддерживать частичное обновление. Можно передать только одно поле.

### Request body

```json
{
  "name": "Ivan Ivanovich Petrov",
  "email": "ivan.petrov.updated@example.com"
}
```

Email удалить нельзя, потому что поле является обязательным. Можно только заменить его на другой валидный и свободный email.

### Response `200 OK`

```json
{
  "data": {
    "id": "2cb8ff0e-7b2a-4dbb-9283-e98b4a2f6d1a",
    "name": "Ivan Ivanovich Petrov",
    "email": "ivan.petrov.updated@example.com",
    "created_at": "2026-05-07T10:00:00Z",
    "updated_at": "2026-05-07T10:10:00Z"
  }
}
```

### Ошибки

| Ситуация               |             Код |
| ---------------------- | --------------: |
| Пустое тело `{}`       | `400` или `422` |
| `user_id` не UUID      |           `400` |
| Пользователь не найден |           `404` |
| Невалидные поля        |           `422` |

## 8.5 Удалить пользователя

```http
DELETE /api/v1/users/{user_id}
```

### Response `204 No Content`

Тело ответа отсутствует.

### Бизнес-правило

При удалении пользователя все его задачи удаляются автоматически через `ON DELETE CASCADE`.

### Ошибки

| Ситуация               |   Код |
| ---------------------- | ----: |
| `user_id` не UUID      | `400` |
| Пользователь не найден | `404` |

## 9. Tasks API

## 9.1 Создать задачу

```http
POST /api/v1/tasks
```

### Request body

```json
{
  "user_id": "2cb8ff0e-7b2a-4dbb-9283-e98b4a2f6d1a",
  "title": "Prepare Go project",
  "description": "Implement REST API for ToDo List",
  "completed": false
}
```

Поле `description` можно не передавать или передать как `null`.

Поле `completed` в строгой версии API обязательно. Если в учебной реализации удобнее упростить API, допускается default `false`, но это должно быть явно отражено в README и тестах.

### Response `201 Created`

```json
{
  "data": {
    "id": "6d66eab7-77bb-4910-95a3-2051abf1fd1a",
    "user_id": "2cb8ff0e-7b2a-4dbb-9283-e98b4a2f6d1a",
    "title": "Prepare Go project",
    "description": "Implement REST API for ToDo List",
    "completed": false,
    "created_at": "2026-05-07T10:00:00Z",
    "updated_at": "2026-05-07T10:00:00Z",
    "completed_at": null
  }
}
```

Если задача создаётся сразу как выполненная:

```json
{
  "user_id": "2cb8ff0e-7b2a-4dbb-9283-e98b4a2f6d1a",
  "title": "Already done task",
  "completed": true
}
```

то приложение автоматически выставляет `completed_at = now()`.

### Ошибки

| Ситуация                                     |   Код |
| -------------------------------------------- | ----: |
| Невалидный JSON                              | `400` |
| `user_id` не UUID                            | `400` |
| Пользователь не найден                       | `404` |
| `title` отсутствует                          | `422` |
| `title` длиннее 100 символов                 | `422` |
| `description` длиннее 1000 символов          | `422` |
| `completed` отсутствует в строгой версии API | `422` |

## 9.2 Получить задачу по ID

```http
GET /api/v1/tasks/{task_id}
```

### Response `200 OK`

```json
{
  "data": {
    "id": "6d66eab7-77bb-4910-95a3-2051abf1fd1a",
    "user_id": "2cb8ff0e-7b2a-4dbb-9283-e98b4a2f6d1a",
    "title": "Prepare Go project",
    "description": "Implement REST API for ToDo List",
    "completed": false,
    "created_at": "2026-05-07T10:00:00Z",
    "updated_at": "2026-05-07T10:00:00Z",
    "completed_at": null
  }
}
```

### Ошибки

| Ситуация          |   Код |
| ----------------- | ----: |
| `task_id` не UUID | `400` |
| Задача не найдена | `404` |

## 9.3 Получить список задач

```http
GET /api/v1/tasks?limit=20&offset=0
```

### Дополнительные фильтры

| Параметр       | Тип              | Описание                     |
| -------------- | ---------------- | ---------------------------- |
| `completed`    | boolean          | Фильтр по статусу выполнения |
| `user_id`      | UUID             | Фильтр по пользователю       |
| `created_from` | RFC3339 datetime | Дата создания от             |
| `created_to`   | RFC3339 datetime | Дата создания до             |

Пример:

```http
GET /api/v1/tasks?user_id=2cb8ff0e-7b2a-4dbb-9283-e98b4a2f6d1a&completed=true&limit=10&offset=0
```

### Response `200 OK`

```json
{
  "data": [
    {
      "id": "6d66eab7-77bb-4910-95a3-2051abf1fd1a",
      "user_id": "2cb8ff0e-7b2a-4dbb-9283-e98b4a2f6d1a",
      "title": "Prepare Go project",
      "description": "Implement REST API for ToDo List",
      "completed": false,
      "created_at": "2026-05-07T10:00:00Z",
      "updated_at": "2026-05-07T10:00:00Z",
      "completed_at": null
    }
  ],
  "pagination": {
    "limit": 20,
    "offset": 0,
    "total": 1
  }
}
```

## 9.4 Получить задачи конкретного пользователя

Допускаются два варианта маршрута. Рекомендуемый основной вариант:

```http
GET /api/v1/users/{user_id}/tasks?limit=20&offset=0
```

Альтернативный вариант через query-фильтр:

```http
GET /api/v1/tasks?user_id={user_id}
```

В рамках проекта достаточно реализовать один из вариантов, но он должен быть описан в README и покрыт тестами.

## 9.5 Обновить задачу

```http
PATCH /api/v1/tasks/{task_id}
```

PATCH должен поддерживать частичное обновление.

### Request body

```json
{
  "title": "Prepare final Go project",
  "description": "Finish handlers, service layer and tests",
  "completed": true
}
```

Можно передать только одно поле:

```json
{
  "completed": true
}
```

Чтобы удалить описание:

```json
{
  "description": null
}
```

### Response `200 OK`

```json
{
  "data": {
    "id": "6d66eab7-77bb-4910-95a3-2051abf1fd1a",
    "user_id": "2cb8ff0e-7b2a-4dbb-9283-e98b4a2f6d1a",
    "title": "Prepare final Go project",
    "description": "Finish handlers, service layer and tests",
    "completed": true,
    "created_at": "2026-05-07T10:00:00Z",
    "updated_at": "2026-05-07T11:00:00Z",
    "completed_at": "2026-05-07T11:00:00Z"
  }
}
```

### Бизнес-правила обновления

1. Если `completed` меняется с `false` на `true`, выставить `completed_at = now()`.
2. Если `completed` меняется с `true` на `false`, выставить `completed_at = null`.
3. Если `completed` не передан, не менять `completed` и `completed_at`.
4. Если переданы `title` или `description`, обновить только эти поля.
5. `user_id` задачи менять нельзя.
6. `created_at`, `updated_at`, `completed_at` напрямую через API менять нельзя.

### Ошибки

| Ситуация          |             Код |
| ----------------- | --------------: |
| Пустое тело `{}`  | `400` или `422` |
| `task_id` не UUID |           `400` |
| Задача не найдена |           `404` |
| Невалидные поля   |           `422` |

## 9.6 Удалить задачу

```http
DELETE /api/v1/tasks/{task_id}
```

### Response `204 No Content`

Тело ответа отсутствует.

### Ошибки

| Ситуация          |   Код |
| ----------------- | ----: |
| `task_id` не UUID | `400` |
| Задача не найдена | `404` |

## 10. Statistics API

## 10.1 Получить статистику по задачам

```http
GET /api/v1/statistics/tasks
```

### Query-параметры

| Параметр       | Тип              | Обязательный | Описание                                              |
| -------------- | ---------------- | -----------: | ----------------------------------------------------- |
| `user_id`      | UUID             |          Нет | Статистика только по задачам конкретного пользователя |
| `created_from` | RFC3339 datetime |          Нет | Учитывать задачи, созданные начиная с этой даты       |
| `created_to`   | RFC3339 datetime |          Нет | Учитывать задачи, созданные до этой даты              |

Примеры:

```http
GET /api/v1/statistics/tasks
```

```http
GET /api/v1/statistics/tasks?user_id=2cb8ff0e-7b2a-4dbb-9283-e98b4a2f6d1a
```

```http
GET /api/v1/statistics/tasks?created_from=2026-05-01T00:00:00Z&created_to=2026-05-07T23:59:59Z
```

### Response `200 OK`

```json
{
  "data": {
    "total_created": 100,
    "total_completed": 75,
    "completion_percentage": 75.0,
    "average_completion_time_seconds": 86400
  }
}
```

### Расчёт метрик

| Метрика                           | Расчёт                                                                    |
| --------------------------------- | ------------------------------------------------------------------------- |
| `total_created`                   | Количество задач, подходящих под фильтр                                   |
| `total_completed`                 | Количество задач с `completed = true`, подходящих под фильтр              |
| `completion_percentage`           | `total_completed / total_created * 100`                                   |
| `average_completion_time_seconds` | Среднее значение `completed_at - created_at` только для завершённых задач |

### Поведение при пустой выборке

Если задач нет:

```json
{
  "data": {
    "total_created": 0,
    "total_completed": 0,
    "completion_percentage": 0,
    "average_completion_time_seconds": null
  }
}
```

### Ошибки

| Ситуация                                    |   Код |
| ------------------------------------------- | ----: |
| `user_id` не UUID                           | `400` |
| `created_from` или `created_to` не RFC3339  | `400` |
| `created_from > created_to`                 | `422` |
| `user_id` указан, но пользователь не найден | `404` |

## 11. Конфигурация

Приложение должно конфигурироваться через переменные окружения.

### `.env.example`

```env
APP_NAME=todo-list-api
APP_ENV=local
HTTP_HOST=0.0.0.0
HTTP_PORT=8080

POSTGRES_HOST=localhost
POSTGRES_PORT=5432
POSTGRES_DB=todoapp
POSTGRES_USER=todoapp
POSTGRES_PASSWORD=todoapp
POSTGRES_SSLMODE=disable

LOG_LEVEL=debug
SHUTDOWN_TIMEOUT=10s
```

### Правила конфигурации

1. Приложение не должно хранить реальные секреты в репозитории.
2. `.env` должен быть добавлен в `.gitignore`.
3. В репозитории должен быть `.env.example`.
4. При отсутствии обязательной переменной приложение должно завершаться с понятной ошибкой.
5. Connection string для PostgreSQL может собираться из отдельных переменных или задаваться одной переменной `DATABASE_URL`.

## 12. Логирование

Приложение должно логировать:

1. Старт приложения.
2. Адрес, на котором запущен HTTP-сервер.
3. Успешное подключение к PostgreSQL.
4. Ошибки подключения к PostgreSQL.
5. Каждый HTTP request: method, path, status code, duration.
6. Ошибки handlers/service/repository.
7. Graceful shutdown.

### Требования к логам

1. Логи должны быть структурированными или легко читаемыми.
2. В production-like режиме лучше использовать JSON-формат.
3. В логах не должно быть секретов: паролей, токенов, connection string с паролем.

Пример JSON-лога:

```json
{
  "level": "info",
  "msg": "http request completed",
  "method": "POST",
  "path": "/api/v1/tasks",
  "status": 201,
  "duration_ms": 12
}
```

## 13. Работа с ошибками

В проекте должны быть доменные ошибки:

```go
var (
    ErrUserNotFound = errors.New("user not found")
    ErrTaskNotFound = errors.New("task not found")
    ErrInvalidInput = errors.New("invalid input")
)
```

Repository должен преобразовывать ошибки PostgreSQL в ошибки приложения.

Handler должен преобразовывать ошибки приложения в HTTP status code.

### Маппинг ошибок

| Ошибка                                    |     HTTP status |
| ----------------------------------------- | --------------: |
| Невалидный JSON                           |           `400` |
| Невалидный UUID                           |           `400` |
| Ошибка валидации                          |           `422` |
| Пользователь не найден                    |           `404` |
| Задача не найдена                         |           `404` |
| Foreign key violation при создании задачи | `404` или `422` |
| Ошибка базы данных                        |           `500` |
| Неожиданная ошибка                        |           `500` |

## 14. Транзакции и конкурентный доступ

### Обязательные требования

1. Создание, чтение, обновление и удаление одного ресурса может выполняться без явной транзакции, если используется один SQL-запрос.
2. Операции, состоящие из нескольких SQL-запросов, должны выполняться в транзакции.
3. Обновление `completed` и `completed_at` должно быть атомарным.
4. Нельзя сначала прочитать задачу, потом вне транзакции принять решение и затем обновить, если между чтением и обновлением возможна гонка.

### Рекомендуемый подход для PATCH задачи

Для защиты от lost update можно использовать один из подходов:

1. Один атомарный SQL `UPDATE ... RETURNING`.
2. Транзакция + `SELECT ... FOR UPDATE`.
3. Optimistic locking через поле `updated_at` или `version` в расширенной версии.

Для учебной версии достаточно атомарного `UPDATE ... RETURNING` или транзакции с `SELECT FOR UPDATE`.

## 15. Миграции

В проекте должны быть миграции:

```text
migrations/
├── 000001_init.up.sql
└── 000001_init.down.sql
```

### `up` миграция должна создавать:

1. Схему `todoapp`.
2. Таблицу `users`.
3. Таблицу `tasks`.
4. Индексы.
5. При необходимости функцию и trigger для `updated_at`.

### `down` миграция должна удалять:

1. Triggers.
2. Function `set_updated_at`.
3. Таблицу `tasks`.
4. Таблицу `users`.
5. Схему `todoapp`, если она пустая.

Удалять нужно в обратном порядке, чтобы не нарушить foreign key constraints.

## 16. Docker и Docker Compose

## 16.1 `docker-compose.yml`

Должны быть сервисы:

1. `postgres` — база данных.
2. `migrate` — запуск миграций, опционально.
3. `app` — приложение, опционально для полной контейнеризации.

Минимальный учебный вариант может запускать в Docker только PostgreSQL, а приложение запускать локально через `go run`.

## 16.2 Требования к PostgreSQL контейнеру

1. Использовать официальный образ PostgreSQL.
2. Настроить database, user, password через env.
3. Пробросить порт `5432`.
4. Использовать volume для хранения данных.
5. Добавить healthcheck.

## 16.3 Dockerfile приложения

В расширенной версии должен быть multi-stage Dockerfile:

1. Stage build: сборка Go binary.
2. Stage runtime: минимальный образ для запуска binary.
3. Не запускать приложение от root, если используется production-like подход.

## 17. Makefile

Makefile должен содержать команды:

```makefile
.PHONY: run build test lint fmt env-up env-down env-clean migrate-up migrate-down migrate-create

run:
	go run ./cmd/todoapp

build:
	go build -o ./bin/todoapp ./cmd/todoapp

test:
	go test ./...

fmt:
	gofmt -w .

env-up:
	docker compose up -d postgres

env-down:
	docker compose down

env-clean:
	docker compose down -v

migrate-up:
	migrate -path ./migrations -database "$${DATABASE_URL}" up

migrate-down:
	migrate -path ./migrations -database "$${DATABASE_URL}" down

migrate-create:
	migrate create -ext sql -dir ./migrations -seq $(name)
```

Конкретный синтаксис может отличаться, но команды должны быть документированы в README.

## 18. README

README должен содержать:

1. Название проекта.
2. Краткое описание.
3. Стек технологий.
4. Требования для запуска: Go, Docker, Docker Compose, migrate.
5. Инструкцию запуска окружения.
6. Инструкцию применения миграций.
7. Инструкцию запуска приложения.
8. Примеры curl-запросов.
9. Описание переменных окружения.
10. Описание API endpoints.
11. Инструкцию запуска тестов.

## 19. Тестирование

## 19.1 Unit tests

Покрыть тестами:

1. Валидацию пользователя.
2. Валидацию задачи.
3. Бизнес-правила `completed` / `completed_at`.
4. Расчёт статистики.
5. Маппинг ошибок в HTTP responses, если handlers тестируются отдельно.

## 19.2 Integration tests

Желательно покрыть:

1. Создание пользователя в PostgreSQL.
2. Получение пользователя.
3. Создание задачи для существующего пользователя.
4. Ошибку создания задачи для несуществующего пользователя.
5. Обновление `completed` с автоматическим выставлением `completed_at`.
6. Удаление пользователя и каскадное удаление задач.
7. Расчёт статистики по тестовым данным.

## 19.3 Тестовые сценарии приёмки

### Сценарий 1: создание пользователя

1. Отправить `POST /api/v1/users` с валидным `name` и `email`.
2. Получить `201 Created`.
3. Проверить, что в ответе есть `id`, `created_at`, `updated_at`.

### Сценарий 2: ошибка создания пользователя без email

1. Отправить `POST /api/v1/users` только с `name`.
2. Получить `422 Unprocessable Entity`.
3. Проверить, что в ответе есть ошибка по полю `email`.

### Сценарий 3: ошибка валидации пользователя

1. Отправить `POST /api/v1/users` с `name = "ab"`.
2. Получить `422 Unprocessable Entity`.

### Сценарий 4: создание задачи

1. Создать пользователя.
2. Создать задачу с `completed = false`.
3. Получить `201 Created`.
4. Проверить, что `completed_at = null`.

### Сценарий 5: завершение задачи

1. Создать задачу с `completed = false`.
2. Отправить `PATCH /api/v1/tasks/{task_id}` с `completed = true`.
3. Получить `200 OK`.
4. Проверить, что `completed = true` и `completed_at != null`.

### Сценарий 6: возврат задачи в незавершённое состояние

1. Создать завершённую задачу.
2. Отправить `PATCH /api/v1/tasks/{task_id}` с `completed = false`.
3. Проверить, что `completed = false` и `completed_at = null`.

### Сценарий 7: статистика

1. Создать 10 задач.
2. Отметить 4 задачи как выполненные.
3. Запросить `GET /api/v1/statistics/tasks`.
4. Проверить:

   * `total_created = 10`
   * `total_completed = 4`
   * `completion_percentage = 40`

## 20. Нефункциональные требования

## 20.1 Надёжность

1. Приложение должно корректно обрабатывать ошибки БД.
2. Приложение не должно падать при невалидном JSON.
3. Приложение должно использовать graceful shutdown.
4. При завершении приложения HTTP-сервер должен перестать принимать новые запросы и дождаться завершения активных запросов в пределах `SHUTDOWN_TIMEOUT`.

## 20.2 Производительность

Для учебной версии достаточно:

1. Время ответа CRUD endpoints на локальной машине: до 100 ms при нормальной работе БД.
2. Время ответа statistics endpoint: до 500 ms на тестовой выборке до 100 000 задач.
3. Использование индексов для фильтрации по `user_id`, `created_at`, `completed`.

## 20.3 Безопасность

Несмотря на отсутствие авторизации, приложение должно:

1. Использовать parameterized SQL queries.
2. Не конкатенировать пользовательский ввод в SQL.
3. Не логировать пароли и секреты.
4. Не раскрывать внутренние ошибки БД клиенту.
5. Ограничивать размер request body, например 1 MB.
6. Валидировать все path и query parameters.

## 20.4 Совместимость

1. API должен быть версионирован через `/api/v1`.
2. Время в API должно возвращаться в формате RFC3339.
3. Все timestamps должны храниться как `TIMESTAMPTZ`.
4. Рекомендуется использовать UTC в ответах API.

## 21. Критерии готовности проекта

Проект считается готовым, если:

1. Приложение запускается локально.
2. PostgreSQL запускается через Docker Compose.
3. Миграции применяются командой из Makefile.
4. Реализован CRUD для пользователей.
5. Реализован CRUD для задач.
6. Задачи связаны с пользователями через foreign key.
7. Удаление пользователя удаляет его задачи.
8. Реализована пагинация для списков.
9. Реализован endpoint статистики.
10. Реализована валидация входных данных.
11. Ошибки возвращаются в едином JSON-формате.
12. Реализовано логирование HTTP-запросов.
13. Есть `.env.example`.
14. Есть README с инструкцией запуска.
15. Команда `go test ./...` выполняется успешно.
16. Команда `go fmt` или `make fmt` не оставляет неформатированный код.
17. Приложение корректно завершает работу по `SIGINT` / `SIGTERM`.

## 22. Рекомендуемый порядок реализации

### Этап 1: Инициализация проекта

1. Создать Git-репозиторий.
2. Создать Go module.
3. Создать базовую структуру директорий.
4. Добавить `.gitignore`.
5. Добавить `.env.example`.
6. Добавить Makefile.

### Этап 2: Окружение

1. Настроить `docker-compose.yml` с PostgreSQL.
2. Добавить команды `env-up`, `env-down`, `env-clean`.
3. Проверить подключение к БД.

### Этап 3: Миграции

1. Создать первую миграцию.
2. Описать таблицы `users` и `tasks`.
3. Добавить constraints и индексы.
4. Добавить команды `migrate-up`, `migrate-down`, `migrate-create`.

### Этап 4: Конфигурация и запуск приложения

1. Реализовать загрузку env.
2. Реализовать подключение к PostgreSQL.
3. Реализовать HTTP server.
4. Реализовать graceful shutdown.

### Этап 5: Users API

1. Domain model.
2. Repository.
3. Service.
4. Handler.
5. Валидация.
6. Тесты.

### Этап 6: Tasks API

1. Domain model.
2. Repository.
3. Service.
4. Handler.
5. Бизнес-логика `completed_at`.
6. Тесты.

### Этап 7: Statistics API

1. Repository query для агрегации.
2. Service для расчёта метрик.
3. Handler.
4. Фильтры.
5. Тесты.

### Этап 8: Документация и финальная проверка

1. Заполнить README.
2. Добавить curl-примеры.
3. Проверить запуск с нуля.
4. Проверить миграции up/down.
5. Запустить тесты.
6. Проверить форматирование.

## 23. Дополнительные задания после MVP

Эти задачи не обязательны для первой версии, но полезны для развития проекта:

1. Добавить OpenAPI/Swagger спецификацию.
2. Добавить Dockerfile и запуск всего приложения через Docker Compose.
3. Добавить healthcheck endpoint:

```http
GET /healthz
```

4. Добавить readiness endpoint:

```http
GET /readyz
```

5. Добавить optimistic locking через `version`.
6. Добавить soft delete для пользователей и задач.
7. Добавить поиск задач по title.
8. Добавить сортировку задач по `created_at`, `updated_at`, `completed_at`.
9. Добавить GitHub Actions CI.
10. Добавить нагрузочный тест через `wrk` или `k6`.
11. Добавить метрики Prometheus.
12. Добавить авторизацию JWT как отдельный этап.

## 24. Минимальный MVP

Если нужно сделать проект максимально быстро, MVP должен включать:

1. PostgreSQL через Docker Compose.
2. Миграции `users` и `tasks`.
3. `POST /users`, `GET /users/{id}`, `GET /users`.
4. `POST /tasks`, `GET /tasks/{id}`, `GET /tasks`, `PATCH /tasks/{id}`, `DELETE /tasks/{id}`.
5. `GET /statistics/tasks`.
6. Валидацию обязательных полей.
7. Автоматическое управление `completed_at`.
8. Единый формат ошибок.
9. README с запуском.

## 25. Definition of Done

Задача разработки считается завершённой, если:

1. Код реализован и отформатирован.
2. Ошибки обрабатываются явно.
3. SQL-запросы параметризованы.
4. Валидация покрывает граничные случаи.
5. Добавлены или обновлены тесты.
6. README обновлён, если изменилось поведение API.
7. Приложение запускается после `make env-up`, `make migrate-up`, `make run`.
8. Все тесты проходят.
