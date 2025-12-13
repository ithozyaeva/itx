#!/bin/sh
# Ожидание сертификата и конфига перед запуском nginx

DOMAIN="${DOMAIN:-dev.ithozyaeva.ru}"

echo 'Waiting for certificate and nginx config...'
MAX_WAIT=60
WAIT_COUNT=0

while ([ ! -f "/etc/letsencrypt/live/${DOMAIN}/fullchain.pem" ] || [ ! -f "/etc/nginx/conf.d/default.conf" ]) && [ ${WAIT_COUNT} -lt ${MAX_WAIT} ]; do
  sleep 1
  WAIT_COUNT=$((WAIT_COUNT + 1))
  if [ $((WAIT_COUNT % 5)) -eq 0 ]; then
    echo "Still waiting... (${WAIT_COUNT}/${MAX_WAIT} seconds)"
  fi
done

if [ ! -f "/etc/letsencrypt/live/${DOMAIN}/fullchain.pem" ]; then
  echo 'ERROR: Certificate not found after waiting.'
  exit 1
fi

if [ ! -f "/etc/nginx/conf.d/default.conf" ]; then
  echo 'ERROR: Nginx config not found after waiting.'
  exit 1
fi

echo 'Certificate and config found, starting nginx...'
exec /docker-entrypoint.sh nginx -g 'daemon off;'

