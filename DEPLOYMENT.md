# РОСТ Мебель — Deployment Guide

## System Requirements
- Ubuntu 22.04 LTS or similar
- Docker 24.0+
- Docker Compose v2.0+
- Minimum 2GB RAM

## 1. Get Google Gemini API Key
1. Go to [Google AI Studio](https://aistudio.google.com/).
2. Create a new API Key.
3. Save it for the `.env` file.

## 2. Server Setup
1. Point your domain (e.g., `rostmebel.ru`) A-record to your server IP.
2. Clone this repository:
   ```bash
   git clone https://github.com/your-org/rostmebel.git
   cd rostmebel
   ```

## 3. Configuration
1. Create a `.env` file from the example:
   ```bash
   cp .env.example .env
   ```
2. Edit `.env` and fill in the following:
   - `GEMINI_API_KEY`: Your key from step 1.
   - `JWT_SECRET`: A long random string.
   - `ADMIN_PASSWORD`: A secure password for the first login.
   - `DATABASE_URL`: Ensure it points to the `postgres` service in docker-compose.

## 4. Run with Docker
Start all services:
```bash
docker compose up -d
```

Check logs to ensure everything started correctly:
```bash
docker compose logs -f backend
```

## 5. First Run & Migrations
The backend is configured to run migrations automatically on start.
To manually seed the admin user, ensure `ADMIN_USERNAME` and `ADMIN_PASSWORD` are set in `.env` before the first run.

## 6. SSL Configuration (Let's Encrypt)
On your host machine, install Certbot and get a certificate:
```bash
sudo apt install certbot
sudo certbot certonly --standalone -d rostmebel.ru
```
Then update your Nginx configuration (or use a separate Nginx container as a reverse proxy) to point to the certificates in `/etc/letsencrypt/live/rostmebel.ru/`.

## 7. Backups
To backup the database manually:
```bash
docker compose exec postgres pg_dump -U user rostmebel > backup_$(date +%F).sql
```

## 8. Troubleshooting
- **Database Connection**: Ensure the `postgres` healthcheck passes.
- **AI Search**: Check if `GEMINI_API_KEY` is valid. If not, the system will fallback to Postgres full-text search.
- **Image Uploads**: Images are stored in a Docker volume `uploads_data`. Ensure the backend has write access and Nginx has read access. If running on Linux, you may need to set permissions: `sudo chown -R 101:101 uploads` (or the appropriate UID/GID).
- **Rate Limit**: You can toggle the 1-order-per-day limit using `ORDER_LIMIT_ENABLED` in `.env`.

