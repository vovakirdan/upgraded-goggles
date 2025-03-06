# user-post-service: Сервис Управления Пользователями и Постами

## Содержание
- [Введение](#введение)
- [Структура Проекта](#структура-проекта)
- [Установка и Настройка](#установка-и-настройка)
- [Генерация .proto файлов](#генерация-proto-файлов)
- [Конфигурация](#конфигурация)
- [Обзор Сервисов](#обзор-сервисов)
- [API Шлюз](#api-шлюз)
- [gRPC Сервисы](#grpc-сервисы)
- [Схема Базы Данных](#схема-базы-данных)
- [Логирование](#логирование)
- [Аутентификация и Безопасность](#аутентификация-и-безопасность)
- [Документация Swagger](#документация-swagger)
- [Тестирование](#тестирование)
- [Развертывание](#развертывание)

---

## Введение
Проект **user-post-service** — это многофункциональный сервис, предназначенный для управления пользователями и постами. Проект интегрирует аутентификацию, управление постами, логирование и коммуникацию сервисов на основе gRPC. API Шлюз предоставляет HTTP интерфейс для внешнего доступа, в то время как Swagger используется для документации API.

### Ключевые Особенности:
1. **Управление Пользователями и Постами**
   - Регистрация и аутентификация пользователей (на основе JWT)
   - Операции CRUD для постов
   - Валидация данных
2. **Система Логирования**
   - Логирует все входящие и исходящие запросы
   - Захватывает ошибки и исключения
   - Хранит логи в файлах
3. **gRPC и API Шлюз**
   - Реализует gRPC для межсервисной коммуникации
   - API Шлюз для маршрутизации HTTP запросов к gRPC сервисам
   - Промежуточное ПО для авторизации в API Шлюзе
4. **Документация API Swagger**
   - Автоматическая генерация документации на основе обработчиков API
   - Удобный интерфейс для тестирования API
5. **Комплексное Тестирование и Документация**
   - Интеграционные и модульные тесты для всех компонентов
   - Четкая документация по настройке и использованию

---

## Структура Проекта
```
upgraded-goggles/
├── README.md
├── api/
│   ├── gateway/
│   │   ├── config.go
│   │   ├── main.go
│   │   ├── middleware.go
│   │   ├── routes.go
│   │   └── swagger/
│   │       ├── index.html
│   │       └── swagger.yaml
│   ├── grpc/
│   │   ├── post_service.go
│   │   ├── server.go
│   │   └── user_service.go
│   └── proto/
│       ├── common.proto
│       ├── gateway.proto
│       ├── post.proto
│       └── user.proto
├── cmd/
│   ├── post-service/
│   │   └── main.go
│   └── user-service/
│       └── main.go
├── internal/
│   ├── logger/
│   │   └── logger.go
│   ├── post/
│   │   ├── models.go
│   │   ├── repository.go
│   │   ├── service.go
│   │   └── validator.go
│   └── user/
│       ├── models.go
│       ├── repository.go
│       ├── service.go
│       └── validator.go
├── pkg/
│   └── database/
│       └── db.go
├── test/
│   ├── integration_test.go
│   ├── middleware_test.go
│   ├── post_test.go
│   └── user_test.go
├── go.mod
├── go.sum
├── Dockerfile
└── docker-compose.yml
```

---

## Установка и Настройка
### Предварительные Условия
- Go 1.23+
- Docker и Docker Compose
- PostgreSQL

### Шаги
1. Клонируйте репозиторий:
   ```sh
   git clone https://github.com/your-repo/upgraded-goggles.git
   cd upgraded-goggles
   ```
2. Настройте переменные окружения:
   ```sh
   export HTTP_PORT=":8080"
   export USER_SERVICE_ADDRESS="localhost:50051"
   export POST_SERVICE_ADDRESS="localhost:50052"
   ```
   Или создайте файл `.env` в корневой директории с переменными:
   ```
   # Database configuration
   DATABASE_URL=postgres://user:password@localhost:5432/dbname?sslmode=disable

   # API Gateway configuration
   HTTP_PORT=:8080
   USER_SERVICE_ADDRESS=localhost:50051
   POST_SERVICE_ADDRESS=localhost:50052
   ```
3. Соберите проект:
   ```sh
   go mod download
   ```
4. Запустите сервисы с помощью Docker:
   ```sh
   docker-compose up --build
   ```
5. Проверьте сервисы:
   ```sh
   curl http://localhost:8080/health
   ```

--- 

## Генерация .proto файлов
Для генерации .proto файлов, выполните следующие команды:
1. Установите плагины `protoc` и `protoc-gen-go`:
   ```sh
   go install google.golang.org/protobuf/cmd/protoc-gen-go@latest
   go install google.golang.org/grpc/cmd/protoc-gen-go-grpc@latest
   ```
2. Сгенерируйте .proto файлы:
   * Common
   ```bash
   protoc -I api/proto \
  --go_out=paths=source_relative:api/proto/common \
   --go-grpc_out=paths=source_relative:api/proto/common \
   api/proto/common/common.proto
   ```
   * Gateway
   ```bash
   protoc -I api/proto \
   --go_out=paths=source_relative:api/proto/gateway \
   --go-grpc_out=paths=source_relative:api/proto/gateway \
   api/proto/gateway/gateway.proto
   ```
   * Post
   ```bash
   protoc -I api/proto \
   --go_out=paths=source_relative:api/proto/post \
   --go-grpc_out=paths=source_relative:api/proto/post \
   api/proto/post/post.proto
   ```
   * User
   ```bash
   protoc -I api/proto \
   --go_out=paths=source_relative:api/proto/user \
   --go-grpc_out=paths=source_relative:api/proto/user \
   api/proto/user/user.proto
   ```

---

---

## Конфигурация
Конфигурации загружаются из переменных окружения:
| Переменная               | Описание                             | Значение по умолчанию |
|--------------------------|-------------------------------------|-----------------------|
| `HTTP_PORT`              | Порт для API Шлюза                 | `:8080`               |
| `USER_SERVICE_ADDRESS`   | Адрес gRPC сервиса пользователей    | `localhost:50051`     |
| `POST_SERVICE_ADDRESS`   | Адрес gRPC сервиса постов           | `localhost:50052`     |
| `DATABASE_URL`           | Строка подключения к PostgreSQL      | -                     |

---

## Обзор Сервисов
### API Шлюз
Обрабатывает HTTP запросы и перенаправляет их к gRPC сервисам.
- Использует промежуточное ПО для аутентификации JWT.
- Обслуживает статическую документацию Swagger.
- Предоставляет конечные точки для пользователей и постов.

### Сервис Пользователей
Управляет регистрацией, аутентификацией и получением пользователей.
- Использует bcrypt для хеширования паролей.
- Возвращает JWT токен после успешного входа.

### Сервис Постов
Управляет операциями CRUD для постов.
- Валидирует данные поста перед сохранением.
- Поддерживает обновление и удаление постов.

---

## API Шлюз
Маршрутизирует API запросы к gRPC сервисам.
### Конечные Точки:
| Метод | Путь                 | Описание               |
|-------|----------------------|-----------------------|
| POST  | `/v1/users/register` | Регистрация нового пользователя |
| POST  | `/v1/users/login`    | Вход пользователя      |
| GET   | `/v1/users/{id}`     | Получить данные пользователя |
| POST  | `/v1/posts`          | Создать пост          |
| GET   | `/v1/posts/{id}`     | Получить данные поста |
| PUT   | `/v1/posts/{id}`     | Обновить пост         |
| DELETE| `/v1/posts/{id}`     | Удалить пост          |

---

## gRPC Сервисы
### UserService
- `RegisterUser(RegisterRequest) -> RegisterResponse`
- `LoginUser(LoginRequest) -> LoginResponse`
- `GetUser(UserRequest) -> User`

### PostService
- `CreatePost(CreatePostRequest) -> CreatePostResponse`
- `GetPost(GetPostRequest) -> GetPostResponse`
- `UpdatePost(UpdatePostRequest) -> UpdatePostResponse`
- `DeletePost(DeletePostRequest) -> DeletePostResponse`

---

## Логирование
- Использует `log.Logger` для логирования запросов и ошибок.
- Логи хранятся в `logs/app.log`.

---

## Аутентификация и Безопасность
- Использует промежуточное ПО для аутентификации JWT.
- Пароли хешируются с использованием bcrypt.
- Проверки авторизации на защищенных маршрутах.

---

## Документация Swagger
Swagger UI доступен по адресу:
```
http://localhost:8080/swagger/index.html
```

---

## Тестирование
Модульные тесты доступны в директории `test/`.
Запустите тесты с помощью:
```sh
go test ./test/...
```

---

## Развертывание
Для развертывания с помощью Docker:
```sh
docker-compose up --build -d
```
