# Chat Service REST API

RESTful веб-сервис для простого чата, построенный на Go с использованием чистой архитектуры (Clean Architecture). Этот проект демонстрирует современные практики разработки backend-приложений на Go, включая контейнеризацию, тестирование и документирование API.

## 🌟 Обзор проекта

Этот сервис предоставляет REST API для регистрации пользователей, аутентификации и обмена текстовыми сообщениями. Он следует принципам Чистой архитектуры Роберта Мартина (Uncle Bob), обеспечивая разделение ответственностей, тестируемость и независимость от фреймворков и баз данных.

### Основные возможности

- **Полноценный REST API** для управления пользователями и сообщениями.
- **Аутентификация и авторизация** с использованием JWT токенов.
- **Полная контейнеризация** с помощью Docker и Docker Compose для лёгкого развертывания.
- **Автоматические миграции базы данных** для управления схемой.
- **Интерактивная API документация** через Swagger UI/OpenAPI.
- **Структурированное логирование** для мониторинга и отладки.
- **Конфигурация через файлы и переменные окружения**.
- **Graceful shutdown** для корректного завершения работы.
- **Unit-тесты** для ключевых слоёв бизнес-логики (`usecase`, `service`).
- **Проверка кода с помощью `go vet` и `staticcheck`**.

## 🏗️ Архитектура

Проект строго следует принципам **Чистой архитектуры**:

1.  **Entity (Entities):** Доменные модели (`User`, `Message`, `Session`). Это ядро системы, независимое от фреймворков.
2.  **Use Case (Usecase):** Бизнес-логика приложения. Определяет, что система *может* делать. Зависит от Entity и интерфейсов Repository/Service.
3.  **Interface Adapters (Handler, Adapter):**
    *   `Handler`: Реализует HTTP API (Gin). Преобразует HTTP-запросы в вызовы Use Case.
    *   `Adapter`: Реализует интерфейсы Repository и внешние зависимости (PostgreSQL, JWT). Это "порт" для взаимодействия с внешним миром.
4.  **Frameworks & Drivers (Frameworks):** Внешние библиотеки и фреймворки (Gin, pgx, logrus, Viper и т.д.).

**Ключевые принципы:**
- **Зависимости направлены внутрь:** Внутренние слои не зависят от внешних.
- **Интерфейсы для связи:** Взаимодействие между слоями происходит через интерфейсы, определённые в слое Use Case.
- **Независимость от фреймворков:** Основная логика не завязана на конкретные технологии.

### Структура проекта

```
chat-service/
├── cmd/
│   └── server/
│       └── main.go                 # Точка входа в приложение
├── configs/
│   └── config.yaml                 # Основной конфигурационный файл
├── internal/                       # Основной код приложения (не для внешнего импорта)
│   ├── app/
│   │   └── app.go                  # Логика запуска и graceful shutdown приложения
│   ├── entity/
│   │   ├── user.go                 # Доменная модель User
│   │   ├── message.go              # Доменная модель Message
│   │   └── session.go              # Доменная модель Session
│   ├── usecase/
│   │   ├── interfaces.go           # Интерфейсы Repository и Service
│   │   ├── user.go                 # Бизнес-логика для пользователей
│   │   ├── message.go              # Бизнес-логика для сообщений
│   │   └── session.go              # Бизнес-логика для сессий
│   ├── service/
│   │   ├── interfaces.go           # Интерфейсы внутренних сервисов
│   │   ├── hash.go                 # Реализация хэширования паролей (bcrypt)
│   │   └── jwt.go                  # Реализация работы с JWT
│   ├── handler/
│   │   ├── response.go             # Структуры HTTP ответов и обработка ошибок
│   │   ├── middleware.go           # Middleware (Auth, CORS, Logging)
│   │   ├── user.go                 # HTTP обработчики для пользователей
│   │   ├── message.go              # HTTP обработчики для сообщений
│   │   └── handler.go              # Настройка маршрутов Gin
│   ├── adapter/
│   │   └── postgres/
│   │       ├── postgres.go         # Базовый адаптер для PostgreSQL (пул соединений, логирование)
│   │       ├── query_builder.go    # Настройка Squirrel
│   │       ├── user.go             # Реализация UserRepository
│   │       ├── message.go          # Реализация MessageRepository
│   │       └── session.go          # Реализация SessionRepository
├── pkg/                            # Общие пакеты (могут быть переиспользованы)
│   ├── logger/
│   │   └── logger.go               # Инициализация Logrus
│   └── config/
│       └── config.go               # Загрузка и парсинг конфигурации (Viper)
├── migrations/
│   ├── 000001_*.up.sql             # SQL миграции "вверх"
│   └── 000001_*.down.sql           # SQL миграции "вниз"
├── docs/
│   └── swagger.go                  # Аннотации для Swag
├── docker/
│   ├── Dockerfile                  # Dockerfile для сборки приложения
│   └── docker-compose.yml          # Docker Compose конфигурация
├── Makefile                        # Скрипты для сборки, тестирования, запуска
├── go.mod                          # Зависимости Go
└── go.sum                          # Контрольные суммы зависимостей
```

## 🚀 Быстрый старт

### Требования

- Go 1.21 или выше
- Docker & Docker Compose
- Git
- (Опционально) Make

### Запуск с Docker (Рекомендуется)

Это самый простой способ запустить проект, так как все зависимости (PostgreSQL, миграции) будут запущены в контейнерах.

1.  **Клонируйте репозиторий:**
    ```bash
    git clone <URL_вашего_репозитория>
    cd chat-service
    ```
2.  **Соберите и запустите сервисы:**
    ```bash
    # Используя Makefile
    make docker-dev
    # Или напрямую с Docker Compose
    docker-compose -f docker/docker-compose.yml up -d
    ```
3.  **Проверьте статус контейнеров:**
    ```bash
    docker-compose -f docker/docker-compose.yml ps
    # Все контейнеры должны быть в статусе "Up" или "Healthy"
    ```
4.  **Проверьте работу:**
    - Health check: `curl http://localhost:8080/health`
    - Swagger UI: Откройте `http://localhost:8080/swagger/index.html` в браузере.

### Локальная разработка (без Docker)

Если вы хотите запускать сервис локально, вам нужно будет иметь установленный PostgreSQL.

1.  **Установите зависимости:**
    ```bash
    go mod tidy
    # Установите CLI инструменты
    make install-deps
    make install-migrate
    make install-swag # Для генерации документации
    ```
2.  **Настройте PostgreSQL:**
    - Убедитесь, что PostgreSQL запущен.
    - Создайте базу данных и пользователя, соответствующие настройкам в `configs/config.yaml` или `.env`.
3.  **Примените миграции:**
    ```bash
    make migrate-up
    # Или вручную:
    # migrate -path ./migrations -database "postgresql://user:pass@localhost:5432/dbname?sslmode=disable" up
    ```
4.  **Запустите приложение:**
    ```bash
    make run
    # Или напрямую:
    # go run cmd/server/main.go
    ```

## 🧪 Тестирование

Проект включает unit-тесты для слоёв `usecase` и `service`. Тесты для `adapter` (PostgreSQL) были начаты, но могут требовать доработки (например, с использованием `dockertest`).

```bash
# Запустить все тесты
make test
# go test -v ./...

# Запустить тесты с покрытием
make test-coverage
# go test -coverprofile=coverage.out ./... && go tool cover -html=coverage.out

# Запустить тесты для конкретных слоёв
make test-usecase
# go test -v ./internal/usecase/...

make test-service
# go test -v ./internal/service/...
```

## 📚 API Документация

### Swagger/OpenAPI

После запуска приложения интерактивная документация доступна по адресу:
- **Swagger UI**: `http://localhost:8080/swagger/index.html`
- **Swagger JSON**: `http://localhost:8080/swagger/doc.json`

### Основные endpoints

#### Аутентификация
- `POST /api/v1/register`
  - **Описание:** Регистрация нового пользователя.
  - **Тело запроса:** `{"username": "string", "email": "string", "password": "string"}`
  - **Ответ:** Объект пользователя и сессии (с JWT токеном).
- `POST /api/v1/login`
  - **Описание:** Вход в систему.
  - **Тело запроса:** `{"email": "string", "password": "string"}`
  - **Ответ:** Объект пользователя и сессии (с JWT токеном).

#### Профиль пользователя
*(Требуется `Authorization: Bearer <token>` заголовок)*
- `GET /api/v1/profile`
  - **Описание:** Получить информацию о текущем пользователе.
- `PUT /api/v1/profile`
  - **Описание:** Обновить профиль текущего пользователя.
  - **Тело запроса:** `{"username": "string", "email": "string"}`
- `DELETE /api/v1/profile`
  - **Описание:** Удалить аккаунт текущего пользователя.
- `POST /api/v1/logout`
  - **Описание:** Завершить текущую сессию (удалить токен на клиенте).

#### Сообщения
*(Требуется `Authorization: Bearer <token>` заголовок для всех, кроме GET /api/v1/messages)*
- `POST /api/v1/messages`
  - **Описание:** Создать новое сообщение.
  - **Тело запроса:** `{"content": "string"}`
- `GET /api/v1/messages`
  - **Описание:** Получить все сообщения (публичный endpoint).
- `GET /api/v1/messages/my`
  - **Описание:** Получить все сообщения текущего пользователя.
- `GET /api/v1/messages/{id}`
  - **Описание:** Получить конкретное сообщение по его UUID.
- `DELETE /api/v1/messages/{id}`
  - **Описание:** Удалить конкретное сообщение по его UUID (только если оно принадлежит пользователю).

#### Health Check
- `GET /health`
  - **Описание:** Проверка состояния сервиса.
  - **Ответ:** `{"success": true, "message": "Service is running", "data": {"status": "ok"}}`

## 🐳 Docker и Docker Compose

### Структура `docker-compose.yml`

Файл `docker/docker-compose.yml` определяет три сервиса:
1.  **`postgres`**: Контейнер с PostgreSQL 15.
2.  **`chat-service`**: Контейнер с вашим Go приложением.
3.  **`migrate`**: Временный контейнер, который применяет миграции к БД при старте.

### Команды Docker

```bash
# Сборка образов
make docker-build
# docker-compose -f docker/docker-compose.yml build

# Запуск всех сервисов в фоне (-d)
make docker-up
# docker-compose -f docker/docker-compose.yml up -d

# Остановка и удаление контейнеров
make docker-down
# docker-compose -f docker/docker-compose.yml down

# Просмотр логов всех сервисов
make docker-logs
# docker-compose -f docker/docker-compose.yml logs -f

# Просмотр логов конкретного сервиса
docker-compose -f docker/docker-compose.yml logs -f chat-service

# Просмотр запущенных контейнеров
docker-compose -f docker/docker-compose.yml ps

# Выполнение команды внутри контейнера
docker-compose -f docker/docker-compose.yml exec chat-service sh
docker-compose -f docker/docker-compose.yml exec postgres psql -U chatuser -d chatdb
```

## 🗃️ Миграции базы данных

Миграции управляются с помощью инструмента `golang-migrate`. SQL файлы находятся в директории `migrations/`.

### Структура миграций

Каждая миграция состоит из двух файлов:
- `000001_description.up.sql`: SQL команды для применения изменений.
- `000001_description.down.sql`: SQL команды для отката изменений.

### Работа с миграциями

Если вы используете Docker Compose, миграции применяются автоматически сервисом `migrate`.

Для локальной работы:
```bash
# Применить все неприменённые миграции "вверх"
make migrate-up
# migrate -path ./migrations -database "$DATABASE_URL" up

# Откатить последнюю миграцию "вниз"
make migrate-down
# migrate -path ./migrations -database "$DATABASE_URL" down 1

# Сброс базы данных (откат всех миграций)
make migrate-reset
# migrate -path ./migrations -database "$DATABASE_URL" down -all

# Показать текущую версию миграции
make migrate-version
# migrate -path ./migrations -database "$DATABASE_URL" version
```

## ⚙️ Конфигурация

Конфигурация загружается из файла `configs/config.yaml` и может быть переопределена переменными окружения с префиксом `CHAT_`.

### Файл `configs/config.yaml`

```yaml
server:
  host: "0.0.0.0"        # Хост, на котором слушает сервер
  port: 8080             # Порт сервера
  read_timeout: 10s      # Таймаут на чтение HTTP запроса
  write_timeout: 10s     # Таймаут на запись HTTP ответа
  idle_timeout: 60s      # Таймаут простоя соединения
  debug: false           # Режим отладки Gin (не используется)

database:
  host: "localhost"      # Хост БД (в Docker это "postgres")
  port: 5432             # Порт БД
  username: "chatuser"   # Имя пользователя БД
  password: "chatpass"   # Пароль пользователя БД
  name: "chatdb"         # Имя базы данных
  ssl_mode: "disable"    # Режим SSL
  max_connections: 25    # Максимальное количество соединений в пуле
  min_connections: 5     # Минимальное количество соединений в пуле
  max_conn_lifetime: 5m  # Максимальное время жизни соединения
  max_conn_idle_time: 2m # Максимальное время простоя соединения

jwt:
  secret_key: "..."      # Секретный ключ для подписи JWT
  expires_in: 24h        # Время жизни токена

logger:
  level: "info"          # Уровень логирования (debug, info, warn, error, fatal, panic)
  format: "json"         # Формат логов (json, text)
  output: "stdout"       # Куда писать логи (stdout, stderr, или путь к файлу)

app:
  name: "Chat Service"   # Название приложения
  version: "1.0.0"       # Версия
  environment: "development" # Окружение (development, staging, production)
```

### Переменные окружения

Переменные окружения имеют приоритет над значениями в `config.yaml`.
Создайте файл `.env` в корне проекта:

```env
# Database
POSTGRES_DB=chatdb
POSTGRES_USER=chatuser
POSTGRES_PASSWORD=chatpass
POSTGRES_PORT=5432

# Application
APP_PORT=8080
APP_ENVIRONMENT=development
JWT_SECRET_KEY=your-super-secret-jwt-key-change-in-production

# Для локального запуска миграций
DATABASE_URL=postgres://chatuser:chatpass@localhost:5432/chatdb?sslmode=disable
```

## 🛠️ Makefile

`Makefile` содержит удобные команды для разработки и развертывания.

```makefile
# --- Сборка и запуск ---
build:           # Собрать бинарный файл
run:             # Запустить собранное приложение
dev:             # install-deps + build + run

# --- Docker ---
docker-build:    # Собрать Docker образы
docker-up:       # Запустить контейнеры
docker-down:     # Остановить и удалить контейнеры
docker-dev:      # docker-build + docker-up
docker-logs:     # Показать логи контейнеров
docker-ps:       # Показать статус контейнеров

# --- Тестирование ---
test:            # Запустить все тесты
test-usecase:    # Запустить тесты usecase
test-service:    # Запустить тесты service
test-coverage:   # Запустить тесты с покрытием

# --- Миграции ---
migrate-up:      # Применить миграции
migrate-down:    # Откатить последнюю миграцию
migrate-reset:   # Сбросить миграции
migrate-version: # Показать версию миграций

# --- Документация ---
generate-docs:   # Сгенерировать Swagger документацию
swagger-init:    # Инициализировать Swag
swagger-fmt:     # Форматировать комментарии для Swag

# --- Зависимости ---
install-deps:    # Установить Go зависимости
install-migrate: # Установить CLI migrate
install-swag:    # Установить CLI swag
```

## 🔐 Безопасность

- **Пароли:** Хранятся в БД в виде хэшей, созданных с помощью `bcrypt`.
- **JWT:** Используется алгоритм подписи HS256. Токены имеют ограниченное время жизни.
- **Аутентификация:** Реализована через JWT Bearer токены в заголовке `Authorization`.
- **Логирование:** Все запросы и ошибки логируются, что помогает в аудите и отладке.
- **CORS:** Middleware для ограничения источников (в текущей реализации разрешены все `*`).
- **Валидация:** Входные данные валидируются на каждом уровне (Handler -> Use Case -> Entity).

## 📈 Мониторинг и наблюдаемость

- **Логирование:** Структурированное логирование через Logrus в формате JSON.
- **Health Check:** Endpoint `/health` для проверки состояния сервиса.
- **Graceful Shutdown:** При получении сигналов `SIGINT` или `SIGTERM` приложение корректно завершает обработку текущих запросов и закрывает ресурсы.

## 🤝 Вклад в проект

1.  Форкните репозиторий.
2.  Создайте ветку для новой функции (`git checkout -b feature/AmazingFeature`).
3.  Зафиксируйте изменения (`git commit -m 'Add some AmazingFeature'`).
4.  Запушьте ветку (`git push origin feature/AmazingFeature`).
5.  Откройте Pull Request.

## 📄 Лицензия

Этот проект лицензирован под MIT License 

## 👨‍💻 Автор

Автор - [loktev356@inbox.ru](mailto:loktev356@inbox.ru)

## 🆘 Поддержка

Если у вас есть вопросы или проблемы, пожалуйста, создайте issue в репозитории.

---