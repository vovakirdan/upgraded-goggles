# User-Post-System

## Структура проекта
```
user-post-system/           # Корневая папка проекта
│── api/                    # gRPC API
│   ├── proto/              # Протобуферы
│   │   ├── user.proto      # gRPC API для пользователей
│   │   ├── post.proto      # gRPC API для постов
│   │   ├── gateway.proto   # gRPC API Gateway
│   │   └── common.proto    # Общие структуры и сообщения
│   ├── gateway/            # Реализация API Gateway
│   │   ├── main.go
│   │   ├── config.go
│   │   ├── middleware.go
│   │   ├── routes.go
│   │   └── swagger/        # Файлы для Swagger
│   │       ├── docs.go
│   │       ├── swagger.json
│   │       ├── swagger.yaml
│   │       └── index.html
│   └── grpc/               # Реализация gRPC-серверов
│       ├── user_service.go
│       ├── post_service.go
│       └── server.go
│── cmd/                    # Точки входа для сервисов
│   ├── user-service/       # Запуск сервиса пользователей
│   │   ├── main.go
│   ├── post-service/       # Запуск сервиса постов
│   │   ├── main.go
│── config/                 # Файлы конфигурации
│   ├── config.yaml
│   ├── config.go
│── internal/               # Основная бизнес-логика
│   ├── user/               # Логика работы с пользователями
│   │   ├── repository.go   # Работа с БД
│   │   ├── service.go      # Бизнес-логика
│   │   ├── handler.go      # gRPC обработчики
│   │   ├── models.go       # Модели пользователей
│   │   ├── validator.go    # Валидация данных
│   ├── post/               # Логика работы с постами
│   │   ├── repository.go
│   │   ├── service.go
│   │   ├── handler.go
│   │   ├── models.go
│   │   ├── validator.go
│   ├── auth/               # Авторизация и аутентификация
│   │   ├── jwt.go          # Генерация и верификация токенов
│   │   ├── middleware.go   # Middleware для аутентификации
│   ├── logger/             # Логирование
│   │   ├── logger.go       # Логирование в БД и файлы
│── migrations/             # Миграции базы данных
│   ├── 0001_init.up.sql
│   ├── 0001_init.down.sql
│── pkg/                    # Утилиты и вспомогательные пакеты
│   ├── database/           # Подключение к БД
│   │   ├── db.go
│   ├── errors/             # Общие ошибки
│   │   ├── errors.go
│── test/                   # Тесты проекта
│   ├── user_test.go
│   ├── post_test.go
│── Dockerfile              # Docker-образ
│── docker-compose.yaml     # Docker Compose конфигурация
│── README.md               # Документация
│── go.mod
│── go.sum
```