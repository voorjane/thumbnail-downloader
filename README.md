# Тестовое задание для компании "Эшелон Технологии"

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

## CICD

Я работал только с GitLab CI/CD, вот как мог бы выглядеть пайплайн:

`.gitlab-ci.yml`
```
stages:
    - build
    - style
    - test
    - deploy

build_project:
  stage: build
  tags:
    - build
  script:
    - make ci_build
  artifacts:
    paths:
      - ./.bin/server
      - ./.bin/client
    expire_in: 30 days

tests:
  stage: test
  tags:
    - test
  script:
    - make test

deploy:
  stage: deploy
  tags:
    - deploy
  when: manual
  script:
    # скрипт для деплоя
```
