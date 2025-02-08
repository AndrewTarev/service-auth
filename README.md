# Auth Service

## Основной функционал

1. Регистрация пользователей:
   - Создание нового пользователя с сохранением данных в базу.
   - Хэширование паролей.

2. Логин:
   - Проверка логина и пароля.
   - Генерация JWT токенов (access и refresh).

3. Обновление токена (refresh):
   - Проверка валидности refresh токена.
   - Выдача нового access токена.

4. Валидация токена:
   - Проверка переданных токенов на валидность.

## Структура проекта
```
service-auth
├── Dockerfile
├── Makefile
├── README.md
├── authREADME.md
├── cmd
│   └── main.go
├── docker-compose.yml
├── go.mod
├── go.sum
└── internal
    ├── app
    │   ├── delivery                      # Контроллеры HTTP (handler)
    │   │   └── http
    │   │       ├── auth.go
    │   │       ├── handler.go
    │   │       └── response.go
    │   ├── domain                        # Сущности
    │   │   ├── token_payload.go
    │   │   └── user.go
    │   ├── repository                    # Слой работы с БД
    │   │   ├── auth_postgres.go
    │   │   ├── auth_redis.go
    │   │   └── repository.go
    │   └── service                       # Бизнес логика
    │       ├── auth_service.go
    │       └── service.go
    ├── configs
    │   ├── config.go
    │   └── config.yaml
    ├── ifrastructure                        
    │   ├── auth_helper
    │   │   └── jwt_manager.go
    │   ├── db                            # Настройки подключения к Postgres
    │   │   ├── db.go
    │   │   └── migrations
    │   │       ├── 000001_init.down.sql
    │   │       └── 000001_init.up.sql
    │   ├── logging                       # Конфиг логгера
    │   │   └── logger.go
    │   └── redis_client                  # Настройки подключения к Redis
    │       └── redis_connection.go
    └── server                                  # Настройки сервера
        └── server.go
```
