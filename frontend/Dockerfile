FROM node:22-alpine AS builder

WORKDIR /app
ARG VITE_PUBLIC_SITE_URL=https://rostmebel.shop
ENV VITE_PUBLIC_SITE_URL=$VITE_PUBLIC_SITE_URL
COPY package*.json ./
RUN npm ci

COPY . .
RUN npm run build

FROM nginx:alpine
COPY --from=builder /app/dist /usr/share/nginx/html
COPY nginx.conf /etc/nginx/conf.d/default.conf

EXPOSE 80
CMD ["nginx", "-g", "daemon off;"]
