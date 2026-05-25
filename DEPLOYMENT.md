# РОСТ Мебель — Deployment Guide

## System Requirements
- Ubuntu 22.04 LTS or similar
- Docker 24.0+
- Docker Compose v2.0+
- Minimum 2GB RAM

## 1. DNS and repository
1. Point the domain A-record to your server IP.
2. Clone the repository:
   ```bash
   git clone https://github.com/your-org/rostmebel.git
   cd rostmebel
   ```

## 2. Configuration
Create `.env` from the example:

```bash
cp .env.example .env
```

Minimum required values for production:

```env
JWT_SECRET=replace-with-64-random-characters
ADMIN_USERNAME=replace-with-admin-login
ADMIN_PASSWORD=replace-with-strong-password
DOMAIN=rostmebel.shop
PUBLIC_SITE_URL=https://rostmebel.shop
ENABLE_TLS=auto
```

Optional but useful:

```env
GEMINI_API_KEY=
TELEGRAM_TOKEN=
TELEGRAM_CHAT_ID=
```

Use the existing certificate files from the host:

```env
HOST_SSL_CERTS_DIR=/etc/ssl/certs
HOST_SSL_PRIVATE_DIR=/etc/ssl/private
TLS_CERT_PATH=/host-ssl/certs/fullchain.crt
TLS_KEY_PATH=/host-ssl/private/private.key
```

Notes:

- `DOMAIN` is the canonical apex host. If empty, the frontend derives it from `PUBLIC_SITE_URL`.
- `PUBLIC_SITE_URL` must match the real public HTTPS URL.
- `TLS_CERT_PATH` and `TLS_KEY_PATH` point to the certificate and private key inside the frontend container.

## 3. First Docker start
Bring the stack up:

```bash
docker compose up --build -d
docker compose ps
```

Important behavior:

- if `ENABLE_TLS=false`, the frontend serves plain HTTP only;
- if `ENABLE_TLS=auto` or `ENABLE_TLS=true` and certificate files already exist, the frontend starts in HTTPS mode;
- if `ENABLE_TLS=auto` or `ENABLE_TLS=true` but certificate files are missing, the frontend starts in HTTP mode.

Check health:

```bash
docker compose exec -T backend wget -qO- http://localhost:8080/readyz
docker compose exec -T frontend wget -qO- http://127.0.0.1/healthz
docker compose exec -T frontend nginx -t
```

## 4. Verify mounted certificate files
Make sure the files exist on the host:

```bash
ls -l /etc/ssl/certs/fullchain.crt /etc/ssl/private/private.key
```

Then verify they are visible inside Docker:

```bash
docker compose exec -T frontend ls -l /host-ssl/certs/fullchain.crt /host-ssl/private/private.key
docker compose exec -T frontend nginx -t
```

## 5. Verify HTTPS
Quick checks:

```bash
curl -I http://rostmebel.shop
curl -I https://rostmebel.shop
curl -I https://www.rostmebel.shop
```

Expected behavior:

- `http://rostmebel.shop` -> `301` to `https://rostmebel.shop/...`
- `https://www.rostmebel.shop` -> `301` to `https://rostmebel.shop/...`
- `https://rostmebel.shop` -> `200`

Also verify inside Docker:

```bash
docker compose ps
docker compose exec -T backend wget -qO- http://localhost:8080/readyz >/dev/null
docker compose exec -T frontend wget -qO- http://127.0.0.1/healthz >/dev/null
```

## 6. First run and migrations
The backend runs migrations automatically on start.

Before the first production launch, make sure these values are set in `.env`:

```env
ADMIN_USERNAME=
ADMIN_PASSWORD=
JWT_SECRET=
```

## 7. Backups
Database backup:

```bash
docker compose exec postgres pg_dump -U user rostmebel > backup_$(date +%F).sql
```

You should also back up:

- `/etc/ssl/certs/fullchain.crt`
- `/etc/ssl/private/private.key`
- uploaded images from Docker volume `uploads_data`

## 8. Troubleshooting

### Frontend did not switch to HTTPS
Check whether the files exist:

```bash
ls -l /etc/ssl/certs/fullchain.crt /etc/ssl/private/private.key
```

Then restart the frontend:

```bash
docker compose restart frontend
docker compose exec -T frontend nginx -t
```

### Deploy succeeds but browser shows insecure site
Usually one of these:

- wrong certificate file is mounted
- private key does not match the certificate
- frontend was not restarted after certificate replacement
- `PUBLIC_SITE_URL` still points to the wrong domain or protocol
- an external proxy is terminating TLS with another certificate

### Docker is up, but HTTPS port is closed
Make sure `FRONTEND_HTTPS_PORT=443` is available on the host and not occupied by another process:

```bash
sudo ss -tulpn | grep :443
docker compose ps
```

### Backend or frontend healthchecks fail
Run:

```bash
docker compose logs --tail=100 backend
docker compose logs --tail=100 frontend
docker compose exec -T frontend nginx -t
curl http://localhost/readyz
curl http://localhost/healthz
```
