# PROMPT: Production-Ready Website — РОСТ Мебель

---

## КОНТЕКСТ И ЦЕЛЬ

Ты — Senior Full-Stack Engineer с глубокой экспертизой в Go, Vue.js, PostgreSQL, Domain-Driven Design и облачном деплое. Твоя задача — написать полноценный, production-ready веб-сайт для мебельной компании **РОСТ Мебель** с нуля. Код должен быть чистым, расширяемым, безопасным и готовым к реальной нагрузке. Никаких заглушек, никаких TODO без реализации — только рабочий, компилируемый код.

---

## ТЕХНОЛОГИЧЕСКИЙ СТЕК

### Backend
- **Язык:** Go 1.22+
- **Архитектура:** Domain-Driven Design (DDD) — строгое разделение на слои: `domain`, `application`, `infrastructure`, `interfaces`
- **HTTP-роутер:** `chi` или `fiber` (на твоё усмотрение, обоснуй выбор)
- **База данных:** PostgreSQL 16 — работа ТОЛЬКО через `pgxpool/v5` (никаких ORM, никакого GORM, никакого sqlx)
- **Миграции:** `golang-migrate/migrate`
- **AI-интеграция:** Google AI Studio — модель **Gemini 2.5 Flash** через REST API (`generativelanguage.googleapis.com`) — для умного подбора мебели по текстовым предпочтениям клиента
- **Аутентификация:** JWT (access + refresh токены) для админ-панели
- **Кэш / Rate Limiting:** Redis (через `go-redis/v9`) — для rate limiting заявок и кэширования AI-ответов
- **Валидация:** `go-playground/validator`
- **Логирование:** `slog` (structured logging, стандартная библиотека Go 1.21+)
- **Конфигурация:** `.env` + `os.Getenv` через конфиг-struct с валидацией при старте
- **Swagger:** `swaggo/swag` — автогенерация OpenAPI документации

### Frontend
- **Фреймворк:** Vue 3 (Composition API, `<script setup>`)
- **Сборщик:** Vite 5
- **Роутинг:** Vue Router 4
- **Стейт:** Pinia
- **UI:** Tailwind CSS 3 + HeadlessUI — кастомный дизайн, без UI-кит библиотек
- **HTTP-клиент:** Axios с интерцепторами (автообновление токенов)
- **Анимации:** GSAP (плавные переходы, scroll-анимации) + Vue Transition
- **Иконки:** Lucide Vue
- **SEO:** Vue Meta / unhead для мета-тегов

### Инфраструктура
- **Контейнеризация:** Docker + Docker Compose v2
- **Reverse Proxy:** Nginx (SSL termination, gzip, статика, проксирование API)
- **SSL:** Let's Encrypt через Certbot (в инструкции по деплою)
- **БД-бэкапы:** pg_dump по крону внутри отдельного контейнера

---

## АРХИТЕКТУРА BACKEND (DDD)

```
backend/
├── cmd/
│   └── server/
│       └── main.go                  # точка входа, wire всех зависимостей
├── internal/
│   ├── domain/
│   │   ├── product/
│   │   │   ├── entity.go            # Product entity, Value Objects
│   │   │   ├── repository.go        # интерфейс репозитория
│   │   │   └── service.go           # доменный сервис
│   │   ├── order/
│   │   │   ├── entity.go            # Order (заявка) entity
│   │   │   ├── repository.go
│   │   │   └── service.go
│   │   └── admin/
│   │       ├── entity.go            # Admin entity
│   │       └── repository.go
│   ├── application/
│   │   ├── product/
│   │   │   └── usecase.go           # ProductUseCase (CRUD + AI-поиск)
│   │   ├── order/
│   │   │   └── usecase.go           # OrderUseCase (создать, список, обновить статус)
│   │   └── admin/
│   │       └── usecase.go           # AdminUseCase (логин, рефреш)
│   ├── infrastructure/
│   │   ├── postgres/
│   │   │   ├── pool.go              # pgxpool setup
│   │   │   ├── product_repo.go      # реализация ProductRepository
│   │   │   ├── order_repo.go
│   │   │   └── admin_repo.go
│   │   ├── redis/
│   │   │   └── client.go            # rate limiter + cache
│   │   ├── gemini/
│   │   │   └── client.go            # Google Gemini REST client
│   │   └── migrations/
│   │       ├── 001_init.up.sql
│   │       ├── 001_init.down.sql
│   │       └── ...
│   └── interfaces/
│       ├── http/
│       │   ├── server.go            # HTTP-сервер, middleware setup
│       │   ├── middleware/
│       │   │   ├── auth.go          # JWT middleware
│       │   │   ├── ratelimit.go     # Redis rate limiter middleware
│       │   │   ├── cors.go
│       │   │   └── logger.go
│       │   └── handler/
│       │       ├── product.go
│       │       ├── order.go
│       │       └── admin.go
│       └── dto/
│           ├── product.go           # Request/Response DTO
│           ├── order.go
│           └── admin.go
├── migrations/                       # SQL-файлы миграций
├── Dockerfile
├── .env.example
└── go.mod
```

---

## АРХИТЕКТУРА FRONTEND

```
frontend/
├── src/
│   ├── assets/                      # шрифты, SVG, изображения
│   ├── components/
│   │   ├── common/                  # Button, Input, Modal, Badge, Spinner
│   │   ├── catalog/                 # ProductCard, ProductGrid, FilterBar
│   │   ├── ai/                      # AISearchPanel, AIResultCard
│   │   ├── order/                   # OrderForm, OrderSuccess
│   │   └── admin/                   # AdminSidebar, AdminTable, AdminModal
│   ├── composables/
│   │   ├── useAuth.ts
│   │   ├── useProducts.ts
│   │   ├── useOrders.ts
│   │   └── useAISearch.ts
│   ├── pages/
│   │   ├── HomePage.vue             # главная страница
│   │   ├── CatalogPage.vue          # каталог с фильтрами
│   │   ├── ProductPage.vue          # карточка товара
│   │   ├── ContactPage.vue          # контакты + форма заявки
│   │   └── admin/
│   │       ├── LoginPage.vue
│   │       ├── DashboardPage.vue
│   │       ├── ProductsPage.vue
│   │       └── OrdersPage.vue
│   ├── router/
│   │   └── index.ts
│   ├── stores/
│   │   ├── auth.ts
│   │   ├── products.ts
│   │   └── orders.ts
│   ├── api/
│   │   ├── client.ts                # Axios instance + интерцепторы
│   │   ├── products.ts
│   │   ├── orders.ts
│   │   └── admin.ts
│   ├── types/
│   │   └── index.ts
│   └── main.ts
├── index.html
├── vite.config.ts
├── tailwind.config.js
├── tsconfig.json
├── Dockerfile
└── nginx.conf
```

---

## ФУНКЦИОНАЛЬНЫЕ ТРЕБОВАНИЯ

### 1. ПУБЛИЧНАЯ ЧАСТЬ (клиентский сайт)

#### Главная страница
- Hero-секция: полноэкранный баннер с анимированным слоганом, CTA-кнопка "Смотреть каталог"
- **AI-поиск мебели** (ключевая фича): текстовое поле с placeholder "Опишите вашу комнату или предпочтения..." → запрос к Gemini API → возврат подобранных товаров из БД
  - Пример: пользователь пишет "хочу светлую скандинавскую спальню до 50000 рублей" → AI анализирует запрос, извлекает параметры, находит товары из каталога
  - Skeleton-loader пока идёт запрос
  - Красивая анимация появления карточек результатов
- Секция "Почему РОСТ Мебель": иконки + описание преимуществ
- Секция "Хиты продаж": горизонтальный скролл карточек товаров
- Секция "Новинки сезона": grid карточек
- Секция "Как мы работаем": numbered steps с иконками
- Секция "Отзывы": карусель (без внешней библиотеки, на CSS scroll-snap)
- Форма заявки (inline): имя, телефон, комментарий
- Footer: контакты, ссылки, соцсети

#### Каталог
- Категории мебели в виде горизонтальных chips/tabs
- Фильтры: цена (range slider), материал, цвет, комната
- Сортировка: по цене, по новизне, по популярности
- Адаптивный grid карточек (2 col mobile, 3 col tablet, 4 col desktop)
- Пагинация (cursor-based на бэкенде)
- URL отражает состояние фильтров (query params)

#### Карточка товара
- Галерея фотографий с миниатюрами
- Характеристики в таблице
- Кнопка "Оставить заявку" → модальное окно с формой

#### Форма заявки (с защитой от спама)
- Поля: имя (обязательно), телефон (обязательно, маска), email (необязательно), комментарий, выбранный товар
- Honeypot поле (скрытое, для ботов)
- Rate limiting: не более 3 заявок с одного IP за 24 часа (Redis)
- Fingerprinting: хэш User-Agent + IP для дополнительной защиты
- Cooldown на фронтенде: кнопка блокируется на 60 секунд после отправки
- Валидация телефона по маске +7 (XXX) XXX-XX-XX

### 2. ADMIN ПАНЕЛЬ

#### Аутентификация
- JWT access token (15 минут) + refresh token (30 дней) в httpOnly cookie
- Защищённые роуты (Vue Router guard)
- Автоматическое обновление access token через интерцептор Axios

#### Дашборд
- Статистика: количество товаров, новые заявки сегодня, общее количество заявок, топ-3 запрашиваемых товара
- График заявок по дням (последние 30 дней) — нарисованный через SVG или Canvas (без Chart.js)
- Быстрые действия: добавить товар, посмотреть новые заявки

#### Управление товарами
- Таблица с сортировкой и поиском
- Добавление/редактирование товара:
  - Название, описание (rich-text через contenteditable или простой textarea)
  - Категория (выпадающий список с возможностью создать новую)
  - Цена, цена со скидкой
  - Характеристики: динамические key-value пары (добавить/удалить)
  - **AI-теги**: поле для ключевых слов, по которым Gemini будет находить товар (например: "скандинавский стиль, дуб, светлый, спальня, минимализм")
  - Загрузка изображений: drag-and-drop, превью, множественная загрузка (до 10 фото)
  - Статус: опубликован / черновик / архив
- Массовые действия: удалить выбранные, изменить статус
- Мягкое удаление (soft delete — поле `deleted_at`)

#### Управление заявками
- Таблица заявок с фильтром по статусу: Новая / В обработке / Завершена / Отклонена
- Изменение статуса одним кликом (dropdown)
- Просмотр деталей заявки: товар, контакты, комментарий, дата, IP
- Пометить как спам (блокировка IP в Redis на 7 дней)
- Экспорт в CSV

#### Управление категориями
- CRUD категорий с иконкой (emoji или SVG) и slug

### 3. AI-ПОИСК (ключевая функция)

**Логика работы:**

1. Пользователь вводит текст предпочтений на фронтенде
2. Frontend отправляет POST `/api/v1/ai/search` с телом `{ "query": "..." }`
3. Backend формирует prompt для Gemini:

```
Ты — ИИ-консультант мебельного магазина РОСТ Мебель.
Пользователь написал: "{user_query}"

Доступные товары (JSON):
{products_json} // все опубликованные товары с их AI-тегами, категорией, ценой

Верни ТОЛЬКО валидный JSON массив product_id тех товаров, которые наиболее подходят под запрос пользователя.
Отсортируй по релевантности. Верни не более 8 товаров. Если ни один не подходит — верни пустой массив [].
Формат ответа: {"ids": [1, 5, 12]}
```

4. Backend парсит ответ Gemini, достаёт товары из БД по ID
5. Возвращает отсортированный список товаров клиенту
6. Ответ кэшируется в Redis на 5 минут по хэшу запроса

**Fallback:** если Gemini недоступен → full-text поиск через PostgreSQL `to_tsvector`

---

## API ENDPOINTS

```
# Public API
GET    /api/v1/products                    # список товаров (с фильтрами, пагинацией)
GET    /api/v1/products/:id                # товар по ID
GET    /api/v1/categories                  # список категорий
POST   /api/v1/orders                      # создать заявку
POST   /api/v1/ai/search                   # AI-поиск товаров

# Admin API (требует JWT)
POST   /api/v1/admin/auth/login            # логин
POST   /api/v1/admin/auth/refresh          # обновить токен
POST   /api/v1/admin/auth/logout           # выход

GET    /api/v1/admin/products              # список всех товаров (включая черновики)
POST   /api/v1/admin/products              # создать товар
PUT    /api/v1/admin/products/:id          # обновить товар
DELETE /api/v1/admin/products/:id          # удалить (soft delete)
POST   /api/v1/admin/products/:id/images   # загрузить изображения

GET    /api/v1/admin/categories            # список категорий
POST   /api/v1/admin/categories            # создать категорию
PUT    /api/v1/admin/categories/:id        # обновить категорию
DELETE /api/v1/admin/categories/:id        # удалить категорию

GET    /api/v1/admin/orders                # список заявок (с фильтрами)
GET    /api/v1/admin/orders/:id            # детали заявки
PATCH  /api/v1/admin/orders/:id/status     # изменить статус
POST   /api/v1/admin/orders/:id/spam       # пометить как спам
GET    /api/v1/admin/orders/export         # экспорт в CSV

GET    /api/v1/admin/stats                 # статистика дашборда

# Swagger UI
GET    /swagger/*                          # документация API
```

---

## БАЗА ДАННЫХ (PostgreSQL 16)

### Таблицы

```sql
-- Категории
CREATE TABLE categories (
    id          BIGSERIAL PRIMARY KEY,
    name        VARCHAR(100) NOT NULL,
    slug        VARCHAR(100) NOT NULL UNIQUE,
    icon        VARCHAR(50),
    sort_order  INT NOT NULL DEFAULT 0,
    created_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at  TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Товары
CREATE TABLE products (
    id              BIGSERIAL PRIMARY KEY,
    category_id     BIGINT REFERENCES categories(id) ON DELETE SET NULL,
    name            VARCHAR(255) NOT NULL,
    slug            VARCHAR(255) NOT NULL UNIQUE,
    description     TEXT,
    price           NUMERIC(12, 2) NOT NULL,
    price_old       NUMERIC(12, 2),
    images          JSONB NOT NULL DEFAULT '[]', -- [{url, is_main}]
    specs           JSONB NOT NULL DEFAULT '{}', -- {"Материал": "Дуб", "Размер": "120x60"}
    ai_tags         TEXT,                        -- "скандинавский, дуб, светлый, спальня"
    status          VARCHAR(20) NOT NULL DEFAULT 'draft', -- published|draft|archived
    views_count     INT NOT NULL DEFAULT 0,
    orders_count    INT NOT NULL DEFAULT 0,
    search_vector   TSVECTOR,                    -- для full-text поиска
    deleted_at      TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_products_status ON products(status) WHERE deleted_at IS NULL;
CREATE INDEX idx_products_category ON products(category_id) WHERE deleted_at IS NULL;
CREATE INDEX idx_products_search ON products USING GIN(search_vector);
CREATE INDEX idx_products_price ON products(price) WHERE deleted_at IS NULL AND status = 'published';

-- Триггер для автообновления search_vector
CREATE OR REPLACE FUNCTION update_product_search_vector()
RETURNS TRIGGER AS $$
BEGIN
    NEW.search_vector := to_tsvector('russian',
        COALESCE(NEW.name, '') || ' ' ||
        COALESCE(NEW.description, '') || ' ' ||
        COALESCE(NEW.ai_tags, '')
    );
    RETURN NEW;
END;
$$ LANGUAGE plpgsql;

CREATE TRIGGER trg_product_search_vector
BEFORE INSERT OR UPDATE ON products
FOR EACH ROW EXECUTE FUNCTION update_product_search_vector();

-- Заявки
CREATE TABLE orders (
    id              BIGSERIAL PRIMARY KEY,
    product_id      BIGINT REFERENCES products(id) ON DELETE SET NULL,
    client_name     VARCHAR(100) NOT NULL,
    client_phone    VARCHAR(20) NOT NULL,
    client_email    VARCHAR(255),
    comment         TEXT,
    status          VARCHAR(20) NOT NULL DEFAULT 'new', -- new|processing|done|rejected|spam
    ip_address      INET NOT NULL,
    user_agent      TEXT,
    fingerprint     VARCHAR(64),
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    updated_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

CREATE INDEX idx_orders_status ON orders(status);
CREATE INDEX idx_orders_created_at ON orders(created_at DESC);
CREATE INDEX idx_orders_ip ON orders(ip_address);

-- Администраторы
CREATE TABLE admins (
    id              BIGSERIAL PRIMARY KEY,
    username        VARCHAR(50) NOT NULL UNIQUE,
    password_hash   VARCHAR(255) NOT NULL,
    refresh_token   VARCHAR(512),
    last_login_at   TIMESTAMPTZ,
    created_at      TIMESTAMPTZ NOT NULL DEFAULT NOW()
);

-- Блокировки IP (дублируется в Redis, но сохраняется в БД для аудита)
CREATE TABLE ip_blocks (
    id          BIGSERIAL PRIMARY KEY,
    ip_address  INET NOT NULL UNIQUE,
    reason      VARCHAR(50) NOT NULL, -- spam|rate_limit
    blocked_at  TIMESTAMPTZ NOT NULL DEFAULT NOW(),
    expires_at  TIMESTAMPTZ NOT NULL
);
```

---

## ДИЗАЙН И UX ТРЕБОВАНИЯ

### Визуальный стиль
- **Палитра:** тёплые натуральные тона — кремовый (#FAF7F2), тёмно-коричневый (#2C1810), золотистый акцент (#C9A84C), белый (#FFFFFF), светло-серый (#F5F5F0)
- **Типографика:** Google Fonts — `Playfair Display` (заголовки, элегантный serif), `Inter` (текст, современный sans-serif)
- **Стиль:** премиальный, лаконичный, скандинавско-минималистичный с тёплыми акцентами
- **Borderradius:** мягкие скругления (8-16px), никаких острых углов

### Анимации
- При скролле: элементы плавно появляются (IntersectionObserver + CSS transitions)
- Hero: параллакс-эффект при скролле
- AI-поиск: анимация "думающего" состояния (пульсирующие точки)
- Карточки товаров: hover — лёгкое поднятие (transform: translateY(-4px)) + тень
- Переходы между страницами: fade через Vue Transition
- Skeleton loaders: везде, где идёт загрузка данных

### Адаптивность
- Мобильное меню: hamburger → боковое drawer-меню
- Все сетки адаптивны (mobile-first)
- Touch-friendly: кнопки минимум 44x44px

### Производительность
- Lazy loading изображений (IntersectionObserver)
- Изображения в WebP (конвертация через Go при загрузке)
- Virtual scroll для длинных списков в админке
- Vite code splitting по роутам

---

## БЕЗОПАСНОСТЬ

1. **SQL Injection:** параметризованные запросы pgx (никаких fmt.Sprintf в SQL)
2. **XSS:** Content Security Policy через nginx, экранирование на фронтенде
3. **CSRF:** SameSite=Strict cookie для refresh token
4. **Rate Limiting:** Redis sliding window — 3 заявки / IP / 24ч
5. **Honeypot:** скрытое поле `website` в форме заявки — если заполнено → игнор
6. **Password hashing:** bcrypt cost=12 для пароля администратора
7. **JWT:** RS256 (асимметричное шифрование), хранение refresh token в httpOnly cookie
8. **Headers:** Helmet-подобные заголовки через nginx (X-Frame-Options, X-Content-Type-Options, etc.)
9. **Upload validation:** проверка MIME-type файлов на бэкенде (не только расширения)
10. **IP блокировка:** автоблок при превышении rate limit + ручная блокировка через админку

---

## DOCKER COMPOSE СТРУКТУРА

```yaml
# docker-compose.yml (минимальная структура для понимания)
version: '3.9'
services:
  postgres:
    image: postgres:16-alpine
    volumes:
      - postgres_data:/var/lib/postgresql/data
    environment: ...
    healthcheck: ...

  redis:
    image: redis:7-alpine
    volumes:
      - redis_data:/data
    command: redis-server --appendonly yes

  backend:
    build: ./backend
    depends_on:
      postgres: { condition: service_healthy }
      redis: { condition: service_started }
    environment:
      - DATABASE_URL=...
      - REDIS_URL=...
      - GEMINI_API_KEY=...
      - JWT_PRIVATE_KEY=... (base64 encoded)
    volumes:
      - uploads:/app/uploads

  frontend:
    build: ./frontend
    depends_on: [backend]

  nginx:
    image: nginx:alpine
    ports:
      - "80:80"
      - "443:443"
    volumes:
      - ./nginx/nginx.conf:/etc/nginx/nginx.conf
      - ./nginx/ssl:/etc/nginx/ssl
      - uploads:/var/www/uploads
    depends_on: [backend, frontend]

  backup:
    image: postgres:16-alpine
    entrypoint: ["/bin/sh", "-c"]
    command: ["while true; do pg_dump ... | gzip > /backups/...; sleep 86400; done"]
    volumes:
      - backups:/backups
    depends_on: [postgres]

volumes:
  postgres_data:
  redis_data:
  uploads:
  backups:
```

Запуск всего одной командой:
```bash
docker compose up -d
```

---

## ДОПОЛНИТЕЛЬНЫЕ ФИЧИ (добавь сам)

1. **Telegram-уведомления при новой заявке** — бот отправляет сообщение с деталями заявки в Telegram-канал администратора (через Bot API)
2. **Избранное** — клиент может добавить товары в избранное (localStorage), счётчик в шапке
3. **Сравнение товаров** — выбрать до 3 товаров и сравнить характеристики в таблице side-by-side
4. **Автодополнение в поиске** — при вводе в AI-поле показываются популярные запросы (хранятся в Redis sorted set)
5. **История AI-запросов в сессии** — последние 5 запросов пользователя отображаются как chips под полем ввода
6. **Похожие товары** — на странице товара блок "Вам может понравиться" (запрос к Gemini с контекстом текущего товара)
7. **Счётчик просмотров** — атомарное обновление views_count через Redis и периодическая синхронизация с PostgreSQL (батч-апдейт раз в минуту)
8. **Sitemap.xml** — автогенерируемый эндпоинт `/sitemap.xml` для SEO
9. **Robots.txt** — корректный robots.txt
10. **Admin: смена пароля** — форма смены пароля в настройках профиля
11. **WebP конвертация** — при загрузке изображений в Go конвертировать в WebP + сохранять оригинал как fallback
12. **Lazy hydration** — компоненты ниже fold инициализируются только при скролле к ним

---

## ТРЕБОВАНИЯ К КАЧЕСТВУ КОДА

### Go (Backend)
- **Никаких глобальных переменных**, только dependency injection через конструкторы
- Все ошибки оборачивать с контекстом: `fmt.Errorf("productRepo.GetByID: %w", err)`
- **Интерфейсы в слое domain**, реализации в infrastructure — строгое соблюдение DIP
- Контекст (`context.Context`) первым аргументом везде
- Graceful shutdown: обработка SIGTERM/SIGINT с таймаутом 30 секунд
- Unit-тесты для use case слоя (с mock-репозиториями через интерфейсы)
- Никаких `panic` в production коде (кроме `main.go` при невозможности старта)
- Все SQL-запросы в отдельных константах / переменных, не inline в методах

### Vue (Frontend)
- TypeScript везде (строгий режим `strict: true` в tsconfig)
- Composables для всей бизнес-логики, компоненты — только UI
- Никаких прямых обращений к API из компонентов — только через composables/stores
- Все типы определены в `types/index.ts`
- ESLint + Prettier с pre-commit hook (husky)

### Общее
- `.env.example` с документацией каждой переменной
- `README.md` — полная инструкция по локальной разработке и production деплою
- Все секреты только через переменные окружения, никаких хардкодов
- Версионирование API (`/api/v1/...`)

---

## DEPLOYMENT.md — ИНСТРУКЦИЯ ПО ДЕПЛОЮ

Напиши подробный `DEPLOYMENT.md` со следующими разделами:

1. **Системные требования** (Ubuntu 22.04 LTS, Docker 24+, Docker Compose v2, минимум 2GB RAM)
2. **Получение API ключа Google AI Studio** (пошаговые скриншоты не нужны, но текстовые шаги)
3. **Настройка DNS** — А-записи домена на сервер
4. **Клонирование репозитория и настройка `.env`** — с объяснением каждой переменной
5. **Генерация RSA ключей для JWT** — команды `openssl`
6. **Первый запуск:**
   ```bash
   cp .env.example .env
   # nano .env  ← заполнить переменные
   docker compose up -d
   docker compose exec backend ./migrate up  # или автомиграции при старте
   docker compose exec backend ./seed        # создать первого администратора
   ```
7. **Настройка SSL через Let's Encrypt + Certbot**
8. **Обновление сайта** (zero-downtime deploy через docker compose pull + up)
9. **Мониторинг логов** (`docker compose logs -f backend`)
10. **Бэкап и восстановление БД**
11. **Troubleshooting** — частые проблемы и решения

---

## ВАЖНЫЕ УТОЧНЕНИЯ

- Начальный пароль администратора задаётся через переменную `ADMIN_PASSWORD` в `.env` и хэшируется при первом запуске / через seed-команду
- Файлы изображений хранятся в Docker volume `/uploads`, раздаются через nginx по пути `/uploads/...`
- Все временные метки в UTC, на фронтенде конвертируются в локальное время пользователя
- Пагинация cursor-based (по `id`), не offset — для производительности при больших объёмах
- Gemini API ключ — бесплатный тариф Google AI Studio (`aistudio.google.com`), модель `gemini-2.5-flash`
- Если `GEMINI_API_KEY` не задан — AI-поиск gracefully деградирует до PostgreSQL full-text поиска
- Никаких внешних UI-библиотек кроме Tailwind CSS, HeadlessUI и Lucide — всё пишем сами

---

## ФИНАЛЬНЫЙ ЧЕКЛИСТ

Перед тем как считать задачу выполненной, убедись:

- [ ] `docker compose up -d` поднимает всё без ошибок с нуля
- [ ] Автомиграции выполняются при старте бэкенда
- [ ] Публичный сайт доступен на `http://localhost`
- [ ] Админка доступна на `http://localhost/admin`
- [ ] Можно войти в админку, создать товар с фото, опубликовать его
- [ ] Товар отображается на сайте в каталоге
- [ ] AI-поиск находит товар по текстовому описанию
- [ ] Форма заявки работает, заявка появляется в админке
- [ ] После 3 заявок с одного IP — блокировка на 24 часа
- [ ] Swagger UI доступен на `http://localhost/swagger`
- [ ] Все Go-файлы компилируются без ошибок
- [ ] TypeScript компилируется без ошибок
- [ ] Нет console.error в браузере при нормальной работе

---

Приступай к реализации. Пиши весь код полностью, без сокращений. Каждый файл — от первой до последней строки. Не используй заглушки типа "// реализация аналогична" или "// TODO". Это production-ready проект.
