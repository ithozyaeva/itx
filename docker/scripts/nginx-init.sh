#!/bin/sh
# Генерация nginx конфигурации из шаблона

apk add --no-cache gettext

DOMAIN="${DOMAIN}"
HTTPS_PORT="${HTTPS_PORT:-}"
BACKEND_CONTAINER="${BACKEND_CONTAINER:-backend}"
ENVIRONMENT="${ENVIRONMENT:-prod}"

# Формируем строку для редиректа с портом (если порт указан)
if [ -n "${HTTPS_PORT}" ]; then
  REDIRECT_PORT=":${HTTPS_PORT}"
else
  REDIRECT_PORT=""
fi

if [ "${ENVIRONMENT}" = "dev" ]; then
  echo 'Waiting for certificate to be created...'
  MAX_WAIT=30
  WAIT_COUNT=0
  
  while [ ! -f "/etc/letsencrypt/live/${DOMAIN}/fullchain.pem" ] && [ ${WAIT_COUNT} -lt ${MAX_WAIT} ]; do
    sleep 1
    WAIT_COUNT=$((WAIT_COUNT + 1))
  done
  
  if [ ! -f "/etc/letsencrypt/live/${DOMAIN}/fullchain.pem" ]; then
    echo 'ERROR: Certificate not found after waiting. Please check certbot-init-selfsigned logs.'
    exit 1
  fi
else
  echo 'Production environment: checking for existing certificate...'
  if [ ! -f "/etc/letsencrypt/live/${DOMAIN}/fullchain.pem" ]; then
    echo 'Warning: Certificate not found. It will be created by certbot-init.'
  fi
fi

echo 'Generating nginx configuration...'
export DOMAIN REDIRECT_PORT BACKEND_CONTAINER
envsubst '${DOMAIN} ${REDIRECT_PORT} ${BACKEND_CONTAINER}' < /tmp/nginx.conf.template > /etc/nginx/conf.d/default.conf
echo "Nginx configuration updated with domain: ${DOMAIN}"

