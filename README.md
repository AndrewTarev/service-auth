# Auth Service

## Основной функционал

1. Регистрация пользователей:
   - Создание нового пользователя с сохранением данных в базу.
   - Хэширование паролей.

2. Логин:
   - Проверка логина и пароля.
   - Генерация JWT токенов (access и refresh) через (private, public) certs.
   - Запись токенов в куки.
   - Запись refresh в redis.

3. Обновление токенов (access, refresh):
   - Проверка валидности refresh токена.
   - Проверка наличия refresh токена в redis.
   - Выдача нового access токена.
   - Обновление refresh в redis

4. Отзыв refresh токена:
   - Удаление refresh токена из redis

## Структура проекта

```
service-auth
├── cmd
├── docs                            # Swagger документация
├── internal
│   ├── app
│   │   ├── delivery                # Слой хэндлеров
│   │   │   ├── http
│   │   │   └── middleware          # Обработка ошибок тут
│   │   ├── errs
│   │   ├── models                  # Сущности
│   │   ├── repository              # Слой работы с БД
│   │   │   └── migrations
│   │   ├── service                 # Слой сервисов
│   │   │   └── mocks
│   │   └── utils             
│   ├── configs
│   └── server
├── pkg                             # Настройки БД, логгера, редиса
│   ├── db
│   ├── logger
│   └── redis_client
└── test                            # Тесты
```

## Установка приложения:

1. Склонируйте репозиторий себе на компьютер
   - git clone https://github.com/AndrewTarev/service-auth.git

2. Установите свои переменные в .env файл
3. Сгенерируйте новые сертификаты по пути ./internal/app/utils/README.md
4. Замените старые сертификаты ./internal/certs
5. Запустите сборку контейнеров
   - docker-compose up --build

API документация (Swagger/OpenAPI) доступна по пути http://localhost:8080/swagger/index.html