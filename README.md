# РОСТ Мебель — Production-Ready Website

A modern, full-stack website for a furniture company featuring AI-powered furniture search using Google Gemini 1.5 Flash.

## 🚀 Features
- **AI-Powered Search**: Describe your preferences (e.g., "warm light bedroom under 50k") and get relevant results.
- **DDD Architecture**: Clean and maintainable Go backend using Domain-Driven Design.
- **Modern Frontend**: Fast and beautiful Vue 3 + Tailwind CSS interface.
- **Admin Panel**: Manage products, orders, and view stats.
- **Performance & Security**: Redis rate limiting, JWT auth, Postgres full-text search fallback.

## 🛠 Tech Stack
- **Backend**: Go (Chi, pgx, slog), PostgreSQL, Redis, Google Gemini API.
- **Frontend**: Vue 3 (Composition API, Pinia), Vite, Tailwind CSS, GSAP.
- **Infrastructure**: Docker, Nginx.

## 🏁 Quick Start
1. Ensure you have Docker and Docker Compose installed.
2. Fill in the `.env` file (see `.env.example`).
3. Run the project:
   ```bash
   docker compose up -d
   ```
4. Access the site at `http://localhost`.
5. Access the admin panel at `http://localhost/admin/login`.

## 📜 Development
### Backend
```bash
cd backend
go run cmd/server/main.go
```

### Frontend
```bash
cd frontend
npm install
npm run dev
```

## ⚖️ License
MIT
