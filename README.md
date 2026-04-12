# РОСТ Мебель

Production-ready сайт и админ-панель для мебельной компании: публичный каталог проектов, заявки, отзывы, загрузка изображений, AI-поиск по портфолио, Telegram-уведомления и административные инструменты для управления контентом.

Проект рассчитан на запуск через Docker Compose: фронтенд обслуживается Nginx, API работает на Go, данные хранятся в PostgreSQL, Redis используется для быстрых служебных операций и ограничений.

## Возможности

- Публичные страницы: главная, каталог, карточка проекта, контакты, избранное.
- Каталог проектов с категориями, ценами, изображениями, характеристиками и отзывами.
- Форма заявки с защитой от спама и лимитами.
- Модерация отзывов.
- Админ-панель: проекты, заявки, отзывы, статистика, экспорт данных.
- AI-поиск по естественному запросу пользователя через Google Gemini.
- Детерминированный fallback для AI-поиска, если Gemini недоступен, выключен или вернул ошибку.
- Единый формат API-ошибок с кодом, сообщением и метаданными.
- Единый outbound HTTP-клиент с поддержкой прокси для Gemini и Telegram.
- Healthcheck/readiness endpoints для Docker, CI/CD и мониторинга.

## Стек

| Слой | Технологии |
| --- | --- |
| Frontend | Vue 3, TypeScript, Vite 8, Pinia, Vue Router, Tailwind CSS, GSAP |
| Backend | Go, chi, pgx, PostgreSQL, Redis, JWT, slog, golang-migrate |
| AI | Google Gemini API, Postgres full-text search, серверный fallback |
| Infra | Docker Compose, Nginx, GitHub Actions |
| Тесты и проверки | Go tests, TypeScript build, Vite build, Docker build, compose validation |

## Архитектура

```text
Browser
  |
  v
Nginx frontend container
  |-- static Vue SPA
  |-- /api/v1/*  -----------> Go backend
  |-- /uploads/* -----------> shared uploads volume
  |
  v
Go backend
  |-- PostgreSQL: проекты, заявки, отзывы, пользователи, миграции
  |-- Redis: быстрые проверки и лимиты
  |-- Gemini API: AI-ранжирование найденных проектов
  |-- Telegram API: уведомления о заявках
       |
       v
optional outbound proxy
```

Серверная часть разделена на слои:

- `backend/internal/domain` — доменные модели и доменные ошибки.
- `backend/internal/application` — use case-логика.
- `backend/internal/infrastructure` — PostgreSQL, Redis, Gemini, Telegram, outbound HTTP-клиент.
- `backend/internal/interfaces` — HTTP handlers, DTO, middleware.
- `backend/cmd/server` — сборка приложения, миграции, wiring зависимостей.

Фронтенд разделен на:

- `frontend/src/pages` — страницы публичной части и админки.
- `frontend/src/components` — переиспользуемые UI-блоки.
- `frontend/src/stores` — Pinia-хранилища.
- `frontend/src/api` — API-клиенты.
- `frontend/src/router` — маршруты и SEO meta.
- `frontend/src/composables` — переиспользуемая клиентская логика.

## Быстрый запуск через Docker

Требования:

- Docker
- Docker Compose

Создайте локальный `.env`:

```powershell
Copy-Item .env.example .env
```

Для Bash:

```bash
cp .env.example .env
```

Заполните обязательные значения:

```env
JWT_SECRET=replace-with-64-random-characters
ADMIN_USERNAME=replace-with-admin-login
ADMIN_PASSWORD=replace-with-strong-password
```

Запустите весь стек:

```bash
docker compose up --build -d
```

Проверьте состояние контейнеров:

```bash
docker compose ps
```

Адреса по умолчанию:

| Назначение | URL |
| --- | --- |
| Сайт | http://localhost |
| Админ-панель | http://localhost/admin/login |
| Backend напрямую | http://localhost:8081 |
| Liveness | http://localhost/healthz |
| Readiness | http://localhost/readyz |

Логи backend:

```bash
docker compose logs --tail=100 backend
```

Остановить стек:

```bash
docker compose down
```

Остановить стек и удалить локальные volume с данными:

```bash
docker compose down -v
```

Команду с `-v` используйте только если точно не нужны локальная база, Redis и загруженные файлы.

## Переменные окружения

Корневой `.env` используется Docker Compose. Файл не должен попадать в Git.

| Переменная | Обязательна | Значение по умолчанию | Назначение |
| --- | --- | --- | --- |
| `JWT_SECRET` | Да | нет | Секрет для подписи JWT. В production используйте длинное случайное значение. |
| `ADMIN_USERNAME` | Да | нет | Логин первого администратора. |
| `ADMIN_PASSWORD` | Да | нет | Пароль первого администратора. |
| `GEMINI_API_KEY` | Нет | пусто | Ключ Google Gemini. Если пустой, AI-поиск работает через fallback. |
| `GEMINI_MODEL` | Нет | `gemini-2.5-flash` | Основная модель Gemini. |
| `GEMINI_FALLBACK_MODELS` | Нет | `gemini-2.5-flash-lite` | Список моделей через запятую для fallback при ошибках 400/404 от Gemini. |
| `TELEGRAM_TOKEN` | Нет | пусто | Токен Telegram-бота для уведомлений. |
| `TELEGRAM_CHAT_ID` | Нет | пусто | Chat ID, куда отправлять уведомления. |
| `ORDER_LIMIT_ENABLED` | Нет | `true` | Включает серверные лимиты на отправку заявок. |
| `POSTGRES_HOST_PORT` | Нет | `55432` | Host-порт PostgreSQL для локального доступа с машины разработчика. Внутри Docker backend использует `postgres:5432`. |
| `OUTBOUND_PROXY_SCHEME` | Нет | `http` | Схема outbound-прокси: `http`, `https`, `socks5`, `socks5h`. |
| `OUTBOUND_PROXY_HOST` | Нет | пусто | Host прокси. Если host или port пустые, прокси выключен. |
| `OUTBOUND_PROXY_PORT` | Нет | пусто | Port прокси. |
| `OUTBOUND_PROXY_USERNAME` | Нет | пусто | Логин прокси, если нужен. |
| `OUTBOUND_PROXY_PASSWORD` | Нет | пусто | Пароль прокси, если нужен. |

Внутри backend также поддерживаются:

| Переменная | Значение по умолчанию | Назначение |
| --- | --- | --- |
| `ENV` | `development` | Режим приложения. В Docker Compose установлен `production`. |
| `PORT` | `8080` | Порт backend внутри контейнера или локального процесса. |
| `DATABASE_URL` | `postgres://user:password@localhost:55432/rostmebel?sslmode=disable` | Подключение к PostgreSQL при локальном запуске backend вне Docker. |
| `REDIS_URL` | `localhost:6379` | Адрес Redis. |
| `REDIS_PASSWORD` | пусто | Пароль Redis, если используется. |
| `JWT_ACCESS_TTL` | `15m` | Время жизни access token. |
| `JWT_REFRESH_TTL` | `720h` | Время жизни refresh token. |
| `ALLOWED_ORIGINS` | `http://localhost:5173,http://localhost:80` | CORS origins через запятую. |

## AI-поиск и Gemini

AI-поиск находится за endpoint:

```http
POST /api/v1/ai/search
```

Пайплайн:

1. Backend собирает кандидатов из PostgreSQL по тексту запроса, категориям и популярности.
2. Кандидаты отправляются в Gemini компактным JSON.
3. Gemini возвращает ранжированный список проектов и человекочитаемый ответ.
4. Если Gemini недоступен, выключен, отвечает невалидно или модель не поддерживается, backend возвращает серверный fallback по кандидатам.

Для production используйте актуальные модели:

```env
GEMINI_MODEL=gemini-2.5-flash
GEMINI_FALLBACK_MODELS=gemini-2.5-flash-lite
```

Если в логах есть ошибка вида `models/... is not found`, проверьте `GEMINI_MODEL` и `GEMINI_FALLBACK_MODELS` в `.env` на сервере. Старые или несуществующие модели вроде `gemma-4-31b` нужно заменить.

Если Google возвращает ошибку о скомпрометированном ключе, кодом это не чинится: ключ нужно отозвать в Google Cloud/AI Studio, создать новый и обновить `.env` на сервере.

## Outbound proxy

Gemini и Telegram используют общий outbound HTTP-клиент. Прокси включается только если заданы оба значения: `OUTBOUND_PROXY_HOST` и `OUTBOUND_PROXY_PORT`.

Пример SOCKS5-прокси:

```env
OUTBOUND_PROXY_SCHEME=socks5
OUTBOUND_PROXY_HOST=127.0.0.1
OUTBOUND_PROXY_PORT=1080
OUTBOUND_PROXY_USERNAME=
OUTBOUND_PROXY_PASSWORD=
```

Для DNS-резолва на стороне прокси используйте:

```env
OUTBOUND_PROXY_SCHEME=socks5h
```

Если прокси не нужен, оставьте host и port пустыми:

```env
OUTBOUND_PROXY_HOST=
OUTBOUND_PROXY_PORT=
```

## API

Базовый префикс API:

```text
/api/v1
```

Публичные endpoints:

| Method | Endpoint | Назначение |
| --- | --- | --- |
| `GET` | `/projects` | Список проектов. |
| `GET` | `/projects/{id}` | Карточка проекта. |
| `GET` | `/projects/{id}/reviews` | Отзывы проекта. |
| `GET` | `/categories` | Список категорий. |
| `POST` | `/orders` | Создание заявки. |
| `POST` | `/reviews` | Создание отзыва. |
| `POST` | `/ai/search` | AI-поиск по портфолио. |
| `POST` | `/uploads/images` | Загрузка изображения. |

Админские endpoints:

| Method | Endpoint | Назначение |
| --- | --- | --- |
| `POST` | `/admin/auth/login` | Вход в админку. |
| `POST` | `/admin/auth/refresh` | Обновление токенов. |
| `POST` | `/admin/auth/logout` | Выход. |
| `GET` | `/admin/stats` | Статистика. |
| `GET` | `/admin/projects` | Список проектов. |
| `POST` | `/admin/projects` | Создание проекта. |
| `PUT` | `/admin/projects/{id}` | Обновление проекта. |
| `DELETE` | `/admin/projects/{id}` | Удаление проекта. |
| `GET` | `/admin/projects/export` | Экспорт проектов. |
| `GET` | `/admin/orders` | Список заявок. |
| `PATCH` | `/admin/orders/{id}/status` | Обновление статуса заявки. |
| `POST` | `/admin/orders/{id}/spam` | Пометить заявку как спам. |
| `GET` | `/admin/orders/export` | Экспорт заявок. |
| `GET` | `/admin/reviews` | Список отзывов для модерации. |
| `PATCH` | `/admin/reviews/{id}/status` | Обновление статуса отзыва. |
| `DELETE` | `/admin/reviews/{id}` | Удаление отзыва. |
| `POST` | `/admin/upload` | Загрузка изображения из админки. |

Защищенные админские endpoints требуют JWT access token в заголовке:

```http
Authorization: Bearer <access_token>
```

## Формат ошибок API

Backend возвращает ошибки в структурированном формате:

```json
{
  "error": {
    "code": "VALIDATION_FAILED",
    "message": "Request validation failed",
    "meta": {
      "fields": [
        {
          "field": "phone",
          "rule": "required",
          "param": ""
        }
      ]
    }
  }
}
```

Frontend должен ориентироваться на `error.code` и `error.meta`, а не парсить текст сообщения. `message` нужен как fallback для логов и редких непредвиденных состояний.

## Локальная разработка

### Backend

Требования:

- Go версии из `backend/go.mod`
- PostgreSQL
- Redis

Можно поднять только инфраструктурные сервисы:

```bash
docker compose up -d postgres redis
```

Если backend запускается локально вне Docker, используйте `localhost:55432` для PostgreSQL, если не меняли `POSTGRES_HOST_PORT`.

Запуск backend локально:

```bash
cd backend
go run ./cmd/server
```

Тесты:

```bash
cd backend
go test ./...
```

Миграции применяются автоматически при старте backend из папки `backend/migrations`.

### Frontend

Требования:

- Node.js 22.12 или новее
- npm

Vite 8 чувствителен к версии Node. Если локально стоит Node 20.18 или ниже, сборка может падать. Используйте Node 22, как в Dockerfile и GitHub Actions.

Установка и dev-сервер:

```bash
cd frontend
npm ci
npm run dev
```

Production build:

```bash
cd frontend
npm run build
```

Smoke-тесты, если установлен Playwright browser:

```bash
cd frontend
npm run test:smoke
```

Установка Chromium для Playwright:

```bash
cd frontend
npm run test:smoke:install
```

## Проверка качества перед merge

Минимальный набор локальных проверок:

```bash
cd backend
go test ./...
```

```bash
cd frontend
npm ci
npm run build
```

```bash
docker compose config --quiet
docker compose build
docker compose up --build -d
docker compose ps
```

Проверка health endpoints после запуска:

```bash
curl http://localhost/healthz
curl http://localhost/readyz
```

Если нужно проверить AI-поиск вручную:

```bash
curl -X POST http://localhost/api/v1/ai/search \
  -H "Content-Type: application/json" \
  -d "{\"query\":\"кухня до 400000 светлая\"}"
```

## CI/CD

В репозитории настроены GitHub Actions:

- `.github/workflows/quality.yml` запускается на pull request и push в `main`.
- `.github/workflows/deploy.yml` запускается на push в `main`.

Quality workflow проверяет:

- `go test ./...` для backend.
- `npm ci` и `npm run build` для frontend.
- `docker compose build` для всего стека.

Deploy workflow сначала выполняет preflight:

- backend tests;
- frontend build;
- `docker compose config --quiet`.

После preflight выполняется SSH-деплой на VPS:

```bash
git pull --ff-only origin main
docker compose up --build -d
docker compose ps
docker compose exec -T backend wget -qO- http://localhost:8080/readyz
docker compose exec -T frontend wget -qO- http://127.0.0.1/healthz
docker image prune -f
```

Для деплоя нужны GitHub Secrets:

| Secret | Назначение |
| --- | --- |
| `SSH_HOST` | IP или host VPS. |
| `SSH_USER` | Пользователь для SSH. |
| `SSH_KEY` | Приватный SSH-ключ. |
| `PROJECT_PATH` | Путь к проекту на сервере. |

На сервере должен существовать production `.env` с актуальными секретами. Этот файл не хранится в Git.

## Данные и файлы

Docker Compose создает volume:

| Volume | Назначение |
| --- | --- |
| `postgres_data` | Данные PostgreSQL. |
| `redis_data` | Данные Redis. |
| `uploads_data` | Загруженные изображения. |

Frontend контейнер читает `uploads_data` как read-only и отдает файлы по `/uploads/`.

Перед production-обновлениями делайте backup PostgreSQL и важных загруженных файлов.

## Безопасность и эксплуатация

- Не коммитьте `.env`, реальные ключи, токены и дампы базы.
- Используйте сильные `JWT_SECRET`, `ADMIN_USERNAME`, `ADMIN_PASSWORD`.
- Если Gemini API key попал в лог, GitHub или был отмечен Google как leaked, его нужно отозвать и заменить.
- Оставляйте `ORDER_LIMIT_ENABLED=true` в production.
- Ограничивайте CORS через `ALLOWED_ORIGINS`, если backend доступен не только через Nginx.
- Следите за логами `backend` после деплоя: ошибки Gemini/Telegram не должны ломать основной пользовательский сценарий.
- Не храните локальные бинарники в репозитории. Для backend build artifacts вроде `*.exe` должны оставаться вне Git.

## Troubleshooting

### Gemini возвращает 404 по модели

Проверьте `.env` на сервере:

```env
GEMINI_MODEL=gemini-2.5-flash
GEMINI_FALLBACK_MODELS=gemini-2.5-flash-lite
```

После изменения:

```bash
docker compose up --build -d backend
docker compose logs --tail=100 backend
```

### Gemini или Telegram не работают из-за сетевых ограничений

Заполните `OUTBOUND_PROXY_*`. Для SOCKS-прокси часто нужен `socks5h`, чтобы DNS-запросы тоже шли через прокси.

### Frontend build падает из-за Node

Проверьте версию:

```bash
node -v
```

Используйте Node 22.12 или новее. Docker build уже использует Node 22.

### Backend не стартует из-за `.env`

В Docker Compose обязательные переменные помечены через `${VAR:?message}`. Проверьте, что заданы:

```env
JWT_SECRET=
ADMIN_USERNAME=
ADMIN_PASSWORD=
```

### Deploy падает на `.git/FETCH_HEAD: Permission denied`

Это означает, что пользователь из `SSH_USER` не владеет проектом на VPS или не может писать в `.git`. Обычно так происходит, если репозиторий когда-то клонировали или обновляли через `root`/`sudo`.

Исправление на VPS:

```bash
cd /path/to/rostmebel
sudo chown -R "$(id -un):$(id -gn)" .
```

Если деплой выполняется отдельным пользователем, подставьте его явно:

```bash
sudo chown -R deploy:deploy /path/to/rostmebel
```

### Контейнеры запущены, но сайт не отвечает

Проверьте:

```bash
docker compose ps
docker compose logs --tail=100 frontend
docker compose logs --tail=100 backend
curl http://localhost/healthz
curl http://localhost/readyz
```

Если порт `80` занят другим процессом, измените mapping frontend в `docker-compose.yml`.

## Рабочий процесс разработки

Рекомендуемый порядок для задач:

1. Сделать маленький, проверяемый набор изменений.
2. Запустить backend tests, если затронут backend.
3. Запустить frontend build, если затронут frontend.
4. Проверить `docker compose config --quiet`.
5. Поднять стек через `docker compose up --build -d`.
6. Проверить `docker compose ps`, `/healthz`, `/readyz`.
7. Посмотреть логи backend/frontend.
8. Для пользовательских сценариев проверить сайт вручную в браузере.

Такой порядок снижает риск, что локально все выглядит рабочим, но в контейнере или на VPS сломается из-за окружения, Node/Go версии, `.env` или Nginx proxy.
