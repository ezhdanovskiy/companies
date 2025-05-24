# Companies Microservice

[![codecov](https://codecov.io/gh/ezhdanovskiy/companies/branch/master/graph/badge.svg)](https://codecov.io/gh/ezhdanovskiy/companies)

## Обзор проекта

Companies - это микросервис для управления информацией о компаниях, построенный на Go с использованием слоистой архитектуры. Сервис предоставляет REST API для выполнения CRUD операций над компаниями, публикует события изменений в Apache Kafka и поддерживает JWT-аутентификацию для защищенных эндпоинтов.

### Основные возможности
- 🏢 Полный CRUD для управления компаниями
- 🔐 JWT-аутентификация для защищенных операций
- 📨 Асинхронная публикация событий в Kafka
- 🗄️ PostgreSQL для хранения данных
- 🐳 Docker-контейнеризация
- ✅ Покрытие unit и интеграционными тестами

## Быстрый старт

### Предварительные требования
- Go 1.19+
- Docker и Docker Compose
- Make
- curl (для тестирования API)

### Локальный запуск

1. Склонируйте репозиторий:
```bash
git clone https://github.com/ezhdanovskiy/companies.git
cd companies
```

2. Запустите приложение с инфраструктурой:
```bash
make run/local
```
Эта команда:
- Запустит PostgreSQL и Kafka в Docker
- Создаст Kafka топик `companies-mutations`
- Применит миграции БД
- Соберет и запустит приложение

3. В отдельном терминале протестируйте API:
```bash
make company/lifecycle
```
Это выполнит полный CRUD цикл операций над компанией.

4. Для просмотра событий в Kafka:
```bash
make kafka/topic/consume
```

### Запуск тестов

```bash
# Unit тесты
make test

# Интеграционные тесты
make test/int

# Полный цикл тестирования с docker-compose
make test/int/docker-compose
```

## CI/CD и покрытие кода

### GitHub Actions
Проект использует GitHub Actions для непрерывной интеграции. Workflow запускается при каждом push и pull request в ветки master/main и включает:
- Запуск unit и интеграционных тестов
- Генерацию отчетов о покрытии кода
- Загрузку покрытия в Codecov

### Настройка Codecov
Для включения отчетов о покрытии кода:

1. **Получите токен Codecov**:
   - Перейдите на [Codecov](https://app.codecov.io/gh/ezhdanovskiy/companies/settings)
   - Скопируйте токен репозитория

2. **Добавьте токен в GitHub**:
   - Перейдите в Settings → Secrets and variables → Actions вашего репозитория
   - Нажмите "New repository secret"
   - Имя: `CODECOV_TOKEN`
   - Значение: вставьте токен из Codecov

3. **Покрытие будет отображаться автоматически** после слияния workflow в ветку master

## Структура проекта

```
.
├── cmd/
│   └── companies/          # Точка входа приложения
│       └── main.go
├── internal/              # Внутренние пакеты приложения
│   ├── application/       # Инициализация и оркестрация
│   │   ├── application.go
│   │   └── logger.go
│   ├── auth/             # JWT аутентификация
│   │   └── jwt.go
│   ├── config/           # Конфигурация приложения
│   │   └── config.go
│   ├── http/             # HTTP слой (Gin)
│   │   ├── handlers.go
│   │   ├── server.go
│   │   ├── dependencies.go
│   │   ├── requests/     # DTO для запросов
│   │   └── mocks/        # Моки для тестов
│   ├── kafka/            # Kafka producer
│   │   ├── producer.go
│   │   └── message.go
│   ├── middlewares/      # HTTP middlewares
│   │   └── auth.go
│   ├── models/           # Доменные модели
│   │   ├── company.go
│   │   └── errors.go
│   ├── repository/       # Слой работы с БД
│   │   ├── repository.go
│   │   ├── repository_test.go
│   │   └── entities.go
│   ├── service/          # Бизнес-логика
│   │   ├── service.go
│   │   ├── service_test.go
│   │   ├── dependencies.go
│   │   └── mocks/
│   └── tests/            # Интеграционные тесты
│       └── integration_test.go
├── migrations/           # SQL миграции
├── docker-compose.yml    # Конфигурация Docker
├── Dockerfile           # Образ приложения
├── Makefile            # Команды разработки
├── go.mod              # Go модуль
├── codecov.yml         # Конфигурация Codecov
└── CLAUDE.md           # Инструкции для Claude AI
```

## Архитектура

Приложение построено с использованием слоистой архитектуры (Layered Architecture):

![Диаграмма зависимостей пакетов](docs/diagrams/package-dependencies.png)

### Слои приложения

1. **HTTP Layer** (`internal/http/`)
   - Обработка HTTP запросов с использованием Gin framework
   - Валидация входных данных
   - Маршрутизация и middleware

2. **Service Layer** (`internal/service/`)
   - Реализация бизнес-логики
   - Публикация событий в Kafka
   - Координация между репозиторием и внешними сервисами

3. **Repository Layer** (`internal/repository/`)
   - Работа с PostgreSQL через Bun ORM
   - Инкапсуляция логики работы с БД

4. **Application Layer** (`internal/application/`)
   - Инициализация компонентов
   - Управление жизненным циклом приложения
   - Настройка логирования

### Внешние зависимости

- **PostgreSQL** - основное хранилище данных
- **Apache Kafka** - брокер сообщений для асинхронных событий
- **Zookeeper** - координатор для Kafka

### API Endpoints

#### Публичные эндпоинты
- `GET /api/v1/companies/:uuid` - получение информации о компании

#### Защищенные эндпоинты (требуют JWT токен)
- `POST /api/v1/secured/companies` - создание новой компании
- `PATCH /api/v1/secured/companies/:uuid` - обновление компании
- `DELETE /api/v1/secured/companies/:uuid` - удаление компании

### Модель данных

```go
type Company struct {
    ID              uuid.UUID
    Name            string    // уникальное, max 15 символов
    Description     string    // max 3000 символов
    EmployeesAmount int
    Registered      bool
    Type            CompanyType
    CreatedAt       time.Time
    UpdatedAt       time.Time
}

// CompanyType - типы компаний
type CompanyType string

const (
    Corporations       CompanyType = "Corporations"
    NonProfit         CompanyType = "NonProfit"
    Cooperative       CompanyType = "Cooperative"
    SoleProprietorship CompanyType = "Sole Proprietorship"
)
```

## Конфигурация

Приложение конфигурируется через переменные окружения. Все настройки загружаются через Viper.

### Переменные окружения

#### База данных
- `DB_HOST` - хост PostgreSQL (по умолчанию: `localhost`)
- `DB_PORT` - порт PostgreSQL (по умолчанию: `5432`)
- `DB_USER` - пользователь БД (по умолчанию: `db`)
- `DB_PASSWORD` - пароль БД (по умолчанию: `db`)
- `DB_NAME` - имя БД (по умолчанию: `db`)

#### Kafka
- `KAFKA_ADDR` - адрес Kafka брокера (по умолчанию: `localhost:9092`)
- `KAFKA_TOPIC` - топик для событий (по умолчанию: `companies-mutations`)

#### HTTP сервер
- `HTTP_PORT` - порт HTTP сервера (по умолчанию: `8080`)

#### Аутентификация
- `JWT_KEY` - секретный ключ для JWT токенов

#### Логирование
- `LOG_LEVEL` - уровень логирования (debug, info, warn, error)
- `LOG_ENCODING` - формат логов (json, console)

## Доступные команды

### Сборка и запуск
```bash
make build                # Сборка бинарного файла
make run                  # Запуск собранного бинарника
make run/local            # Полный локальный запуск с инфраструктурой
```

### Тестирование
```bash
make test                 # Запуск unit тестов
make test/int             # Запуск unit и интеграционных тестов
make test/int/docker-compose  # Полный цикл интеграционного тестирования
```

### Качество кода
```bash
make lint                 # Запуск golangci-lint
make fmt                  # Форматирование кода
make generate             # Генерация моков для тестов
```

### Docker и инфраструктура
```bash
make up                   # Запуск PostgreSQL и Kafka
make down                 # Остановка контейнеров
make kafka/topic/create   # Создание Kafka топика
make kafka/topic/consume  # Просмотр сообщений в топике
```

### Миграции БД
```bash
make migrate/up           # Применение миграций
make migrate/down         # Откат последней миграции
```

### Тестирование API
```bash
make company/lifecycle    # Полный CRUD цикл через curl
make company/create       # Создание компании
make company/get          # Получение компании
make company/patch        # Обновление компании
make company/delete       # Удаление компании
```

### Диаграммы
```bash
make diagrams             # Генерация диаграмм из DOT файлов
```

## Особенности разработки

### Соглашения по коду
- Используется стандартное форматирование Go (gofmt)
- Линтер: golangci-lint с настройками по умолчанию
- Моки генерируются с помощью gomock

### Тестирование
- Unit тесты находятся рядом с кодом (`*_test.go`)
- Интеграционные тесты в `internal/tests/`
- Для запуска интеграционных тестов используется тег `integration`
- Тесты требуют запущенные PostgreSQL и Kafka

### События Kafka
Все мутирующие операции (CREATE, UPDATE, DELETE) публикуют события в топик `companies-mutations`:

```json
{
  "type": "CREATE|UPDATE|DELETE",
  "companyId": "uuid",
  "data": {
    // данные компании
  },
  "timestamp": "2024-01-01T00:00:00Z"
}
```

### JWT аутентификация
- Алгоритм: HS256
- Токен передается в заголовке: `Authorization: Bearer <token>`
- Защищенные эндпоинты требуют валидный токен

## Разработка

### Добавление новой функциональности

1. Определите модель в `internal/models/`
2. Добавьте методы репозитория в `internal/repository/`
3. Реализуйте бизнес-логику в `internal/service/`
4. Создайте HTTP handlers в `internal/http/`
5. Напишите тесты для каждого слоя
6. Обновите документацию

### Генерация моков для тестов

```bash
make generate
```

Это создаст моки для интерфейсов, помеченных комментарием:
```go
//go:generate mockgen -source=file.go -destination=mocks/file_mock.go
```

## Развертывание

### Docker

Для сборки Docker образа:
```bash
docker build -t companies:latest .
```

### Docker Compose

Полное развертывание с инфраструктурой:
```bash
docker-compose up -d
```

## Мониторинг и логи

Приложение использует структурированное логирование через Zap. Логи выводятся в stdout в формате JSON (production) или console (development).

Просмотр логов:
```bash
# Логи приложения
docker-compose logs -f companies

# Логи всех сервисов
docker-compose logs -f
```

## Troubleshooting

### Проблемы с подключением к БД
1. Проверьте, что PostgreSQL запущен: `docker-compose ps`
2. Проверьте переменные окружения
3. Убедитесь, что миграции применены: `make migrate/up`

### Проблемы с Kafka
1. Проверьте, что Kafka и Zookeeper запущены
2. Убедитесь, что топик создан: `make kafka/topic/create`
3. Проверьте логи Kafka: `docker-compose logs kafka`

### Проблемы с тестами
1. Для интеграционных тестов требуется запущенная инфраструктура
2. Используйте `make test/int/docker-compose` для полного цикла
3. Проверьте, что порты 5432 (PostgreSQL) и 9092 (Kafka) свободны

## Лицензия

Этот проект распространяется под лицензией MIT. См. файл [LICENSE](LICENSE) для подробностей.
