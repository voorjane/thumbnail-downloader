# Сервис для сохранения превью видео с YouTube

Сервис содержит клиент и сервер для сохранения превью (thumbnail) любого YouTube ролика. При повторной загрузке превью изображения скачивается с кэша сервера.

## Запуск

1. Склонируйте репозиторий и установите зависимости:

    `git clone <url>` \
    `go mod tidy`

2. Код сервиса gRPC уже в репозитории, при желании его можно перегенерировать
   
    `sh generate.sh`

3. Запустите PostgreSQL и создайте файл `.env` в корне репозитория
    - Поля для файла окружения \
      `DATABASE_HOST` \
      `DATABASE_PORT` \
      `POSTGRES_USER` \
      `POSTGRES_PASSWORD` \
      `POSTGRES_DB`

4. Запустите сервер используя `make` или `make run_server`, а затем запустите клиент:

   `go run cmd/client/main.go <async> <ссылки>` или `make run_client <ссылки>` (флаг async не работает с make)

5. Превью сохраняются в папке downloads в корне репозитория

- Замечание

  Не обрабатываются видео без превью (пример: https://www.youtube.com/watch?v=jNQXAC9IVRw)

## Тесты

`make test`
