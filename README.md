# Запуск

```shell
$ docker-compose up
```

Порт по умолчанию: 8080

# Содержимое контейнера:

1. Mysql - база данных, порт 3306
2. Phpmyadmin, порт 3003
3. Сервер на Go, порт 8080

# Volumes:

1. go_media_stream_mysql - хранилище для бд
2. go_media_stream_uploads - хранилище для загруженых фильмов

# Страницы

- / - Главная страница с фильмами
- /auth - Страница регистрации / авторизации
- /upload - Страница загрузки фильмов
- /video/{id} - Страница для просмотра фильма
- /stream/video/{id} - Видео поток
- /stream/audio/{id} - Аудио поток
