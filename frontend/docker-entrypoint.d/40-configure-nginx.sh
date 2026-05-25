#!/bin/sh
set -eu

PUBLIC_SITE_URL="${PUBLIC_SITE_URL:-https://rostmebel.shop}"
DOMAIN="${DOMAIN:-}"
ENABLE_TLS="$(printf '%s' "${ENABLE_TLS:-auto}" | tr '[:upper:]' '[:lower:]')"
TLS_CERT_PATH="${TLS_CERT_PATH:-/host-ssl/certs/fullchain.crt}"
TLS_KEY_PATH="${TLS_KEY_PATH:-/host-ssl/private/private.key}"
CONFIG_ROOT="/opt/rostmebel-nginx"

infer_domain_from_url() {
    host="${1#*://}"
    host="${host%%/*}"
    host="${host%%\?*}"
    host="${host%%\#*}"
    host="${host%%:*}"
    host="${host#www.}"
    printf '%s' "$host"
}

# Keep PUBLIC_SITE_URL as the fallback source of truth when DOMAIN is not set.
DOMAIN="${DOMAIN#www.}"
if [ -z "$DOMAIN" ]; then
    DOMAIN="$(infer_domain_from_url "$PUBLIC_SITE_URL")"
fi
if [ -z "$DOMAIN" ]; then
    DOMAIN="rostmebel.shop"
fi

mkdir -p /etc/nginx/conf.d /etc/nginx/includes
cp "$CONFIG_ROOT/includes/"*.inc /etc/nginx/includes/

TEMPLATE="$CONFIG_ROOT/nginx.http.conf"
MODE="http"

if [ "$ENABLE_TLS" != "false" ] && [ "$ENABLE_TLS" != "0" ] && [ "$ENABLE_TLS" != "off" ]; then
    if [ -f "$TLS_CERT_PATH" ] && [ -f "$TLS_KEY_PATH" ]; then
        TEMPLATE="$CONFIG_ROOT/nginx.conf"
        MODE="https"
    else
        echo "TLS is enabled, but certificate files are missing at $TLS_CERT_PATH and $TLS_KEY_PATH."
        echo "Starting in HTTP mode."
    fi
fi

export DOMAIN TLS_CERT_PATH TLS_KEY_PATH
envsubst '${DOMAIN} ${TLS_CERT_PATH} ${TLS_KEY_PATH}' < "$TEMPLATE" > /etc/nginx/conf.d/default.conf

echo "Rendered Nginx config in $MODE mode for domain $DOMAIN"
nginx -t
