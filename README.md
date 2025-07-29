# Garantex Rate Service

gRPC сервис для получения курса USDT с биржи Grinex. Сервис предоставляет API для получения актуальных курсов валют и сохраняет историю курсов в PostgreSQL базе данных.

## Функциональность

- **GetRates** - получение текущего курса USDT/RUB с биржи Grinex
- **Healthcheck** - проверка работоспособности сервиса
- Автоматическое сохранение курсов в базу данных
- Graceful shutdown
- Логирование с помощью Zap
- Мониторинг с помощью Prometheus
- Трассировка с помощью OpenTelemetry
- Миграции базы данных

## Технологии

- **Go 1.24+**
- **gRPC** - для API
- **PostgreSQL** - для хранения данных
- **Zap** - для логирования
- **Prometheus** - для мониторинга
- **OpenTelemetry** - для трассировки
- **Docker** - для контейнеризации

## Быстрый старт

### Требования

- Go 1.24 или выше
- Docker и Docker Compose
- PostgreSQL (если запуск без Docker)

### Запуск с Docker Compose

1. Клонируйте репозиторий:
```bash
git clone <repository-url>
cd grinex-rate-service
```

2. Запустите сервис с помощью Docker Compose:
```bash
make run-docker
```

Сервис будет доступен на порту 8080, а PostgreSQL на порту 5460.

### Запуск локально

1. Запустите PostgreSQL:
```bash
make deps
```

2. Примените миграции:
```bash
make migrate-up
```

3. Запустите сервис:
```bash
make run
```

## Конфигурация

Сервис поддерживает конфигурацию через переменные окружения и флаги командной строки.

### Переменные окружения

| Переменная | Описание | Значение по умолчанию   |
|------------|----------|-------------------------|
| `SERVER_PORT` | Порт gRPC сервера | `8080`                  |
| `DB_HOST` | Хост PostgreSQL | `localhost`             |
| `DB_PORT` | Порт PostgreSQL | `5460`                  |
| `DB_USER` | Пользователь PostgreSQL | `db_admin`              |
| `DB_PASSWORD` | Пароль PostgreSQL | `3Qv@e8U0ImT`              |
| `DB_NAME` | Имя базы данных | `grinex_rates`          |
| `DB_SSLMODE` | SSL режим PostgreSQL | `disable`               |
| `GRINEX_BASE_URL` | Базовый URL API Grinex | `https://grinex.io`     |
| `GRINEX_TIMEOUT` | Таймаут запросов к API | `30s`                   |
| `GRINEX_USER_AGENT` | User-Agent для запросов | `GrinexRateService/1.0` |
| `LOG_LEVEL` | Уровень логирования | `info`                  |

### Флаги командной строки

```bash
./grinex-rate-service \
  --port=8080 \
  --db-host=localhost \
  --db-port=5432 \
  --db-user=postgres \
  --db-password=password \
  --db-name=grinex_rates \
  --db-sslmode=disable \
  --grinex-base-url=https://grinex.io \
  --grinex-timeout=30s \
  --log-level=info
```

## API

### GetRates

Получение текущего курса USDT/RUB.

**Request:**
```protobuf
message GetRatesReq {}
```

**Response:**
```protobuf
message GetRatesResp {
  string trading_pair = 1;
  double ask_price = 2;   
  double bid_price = 3;   
  google.protobuf.Timestamp timestamp = 4; 
}
```

### Healthcheck

Проверка работоспособности сервиса.

**Request:**
```protobuf
message HealthcheckReq {}
```

**Response:**
```protobuf
message HealthcheckResp {
  string status = 1;   // "healthy", "degraded", "unhealthy"
  string message = 2;  // status description
}
```

## Использование с grpcurl

```bash
# Получить текущий курс
grpcurl -plaintext localhost:8080 rateservice.v1.RateService/GetRates

# Проверить здоровье сервиса
grpcurl -plaintext localhost:8080 rateservice.v1.RateService/Healthcheck
```

## База данных

### Схема

```sql
CREATE TABLE rates (
    id BIGSERIAL PRIMARY KEY,
    trading_pair VARCHAR(20) NOT NULL,
    ask_price DECIMAL(20, 8) NOT NULL,
    bid_price DECIMAL(20, 8) NOT NULL,
    timestamp TIMESTAMP WITH TIME ZONE NOT NULL,
    created_at TIMESTAMP WITH TIME ZONE DEFAULT NOW()
);
```

### Миграции

```bash
# Применить миграции
make migrate-up

# Откатить миграции
make migrate-down
```

## Мониторинг

### Prometheus метрики

Сервис экспортирует метрики Prometheus на эндпоинте `/metrics` (если настроен HTTP сервер).

### Логирование

Логи выводятся в JSON формате с использованием Zap. Уровень логирования настраивается через переменную `LOG_LEVEL`.

## Разработка

### Сборка

```bash
# Сборка приложения
make build

# Сборка Docker образа
make docker-build
```

### Тестирование

```bash
# Запуск всех тестов
make test

# Запуск тестов с покрытием
make test-coverage
```

### Линтинг

```bash
# Запуск линтера
make lint
```

### Генерация protobuf

```bash
# Генерация Go файлов из proto
make proto
```

## Структура проекта

```
grinex-rate-service/
├── cmd/                    # Точка входа приложения
│   └── main.go
├── internal/               # Внутренние пакеты
│   ├── config/            # Конфигурация
│   ├── database/          # Работа с базой данных
│   └── service/           # Бизнес-логика
├── migrations/            # Миграции базы данных
├── pb/                    # Сгенерированные protobuf файлы
├── proto/                 # Proto файлы
├── server/                # gRPC сервер
├── Dockerfile             # Docker образ
├── docker-compose.yaml    # Docker Compose конфигурация
├── Makefile               # Команды для разработки
└── README.md              # Документация
```

## Docker

### Сборка образа

```bash
docker build -t grinex-rate-service .
```

### Запуск контейнера

```bash
docker run -d \
  --name grinex-rate-service \
  -p 8080:8080 \
  -e DB_HOST=postgres \
  -e DB_PORT=5460 \
  -e DB_USER=db_admin \
  -e DB_PASSWORD=3Qv@e8U0ImT \
  -e DB_NAME=grinex_rates \
  grinex-rate-service
```
