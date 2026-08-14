FROM node:22-alpine AS builder

WORKDIR /app
ARG VITE_PUBLIC_SITE_URL=https://rostmebel.shop
ARG VITE_YANDEX_METRIKA_ID=111598091
ENV VITE_PUBLIC_SITE_URL=$VITE_PUBLIC_SITE_URL
ENV VITE_YANDEX_METRIKA_ID=$VITE_YANDEX_METRIKA_ID
COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /opt/rostmebel-nginx/nginx.conf
COPY nginx.http.conf /opt/rostmebel-nginx/nginx.http.conf
COPY nginx/includes /opt/rostmebel-nginx/includes
COPY docker-entrypoint.d/40-configure-nginx.sh /docker-entrypoint.d/40-configure-nginx.sh
RUN chmod +x /docker-entrypoint.d/40-configure-nginx.sh

EXPOSE 80 443
CMD ["nginx", "-g", "daemon off;"]
