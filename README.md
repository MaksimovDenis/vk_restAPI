# vk_restAPI 🚀

Веб-сервер для приложения "Фильмотека", который представляет REST API для управления базой фильмов.

## Интсрукция по установке и использованию
- Для запуска приложения, необходимо из корневой директории проекта выполнить команду:  
`git clone https://github.com/MaksimovDenis/vk_restAPI.git`  
`docker-compose build`  
`docker-compose up`  
- Команда запустит все контейнеры описанные в файле `docker-compose.yaml`, включая приложение и БД PostgreSQL.

- После запуска приложения swagger документация доступна по ссылке:  
[http://localhost:8000/swagger/index.html](URL)
![Swagger](https://github.com/MaksimovDenis/vk_restAPI/assets/44647373/c1c63b72-ce61-4d3a-bef5-350eef336253)  
![Swagger](https://github.com/MaksimovDenis/vk_restAPI/assets/44647373/30023cdc-9d6b-4b19-8c66-7c8ba2a31a8d)    

Для авторизации необходимо зарегестрироваться, залогиниться и получить JWT токен  

## Технологии и зависимости
- Язык программирование: Golang 1.21.6  
- PostgreSQL latest  
- Для реализации http сервера использовались стандартные библиотеки go.  

Инструменты:
- Docker: Для контейнеризации и управления средой разработки.  
- Docker Compose: Для управления многоконтейнерными приложениями.
  
Основные зависимости:
- [github.com/DATA-DOG/go-sqlmock v1.5.2](https://github.com/DATA-DOG/go-sqlmock): Библиотека для мокирования SQL запросов в тестах.
- [github.com/aws/smithy-go v1.13.3](https://github.com/aws/smithy-go): Библиотека для работы с AWS SDK.
- [github.com/golang-jwt/jwt/v4 v4.5.0](https://github.com/golang-jwt/jwt): Библиотека для работы с JSON Web Tokens (JWT).
- [github.com/golang/mock v1.6.0](https://github.com/golang/mock): Библиотека для генерации моков.
- [github.com/jmoiron/sqlx v1.3.5](https://github.com/jmoiron/sqlx): Расширение стандартного пакета базы данных SQL в Go.
- [github.com/joho/godotenv v1.5.1](https://github.com/joho/godotenv): Библиотека для загрузки переменных окружения из файла .env.
- [github.com/sirupsen/logrus v1.9.3](https://github.com/sirupsen/logrus): Библиотека для логирования в Go.
- [github.com/stretchr/testify v1.9.0](https://github.com/stretchr/testify): Библиотека с удобными утверждениями (assertions) в тестах.
- [github.com/swaggo/http-swagger/v2 v2.0.2](https://github.com/swaggo/http-swagger): Пакет для генерации Swagger документации из аннотаций в коде.
- [github.com/swaggo/swag v1.16.3](https://github.com/swaggo/swag): Инструмент для автоматической генерации Swagger документации в формате JSON из аннотаций в коде.
- [github.com/golang-migrate/migrate/v4 v4.17.0](https://github.com/golang-migrate/migrate): Инструмент для миграции базы данных.

## Описание функционала:
Приложение поддерживает следующие функции:  
1.Управление актёрами: добавление новоого актёра, просмотр, редактирование информации и удаление информации о актёрах.  
2.Управление фильмами: добавление нового фильма, просмотр информации о фильмах с возможностью сортировки по рейтингу, дате и названию, поиск по фрагменту имени актёра или по названию фильма, редактирование информации и удаление информации о фильме.  
3.Возможность связи между актёрами и фильмами, в которых они учавствовали.   

API приложения закрыт авторизацией. 
Приложение поддерживает 2 роли - админинстратор и пользователь. В зависимости от роли - меняется доступный функционал для клиента. 


