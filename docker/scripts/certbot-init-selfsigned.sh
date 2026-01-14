#!/bin/sh
# Создание self-signed сертификата для dev окружения

if [ "${ENVIRONMENT:-prod}" != "dev" ]; then
  echo 'Skipping self-signed certificate (not dev environment)'
  exit 0
fi

apk add --no-cache openssl
DOMAIN="${DOMAIN}"

if [ ! -f "/etc/letsencrypt/live/${DOMAIN}/fullchain.pem" ] || ! openssl x509 -checkend 86400 -noout -in "/etc/letsencrypt/live/${DOMAIN}/fullchain.pem" 2>/dev/null; then
  echo 'Creating self-signed certificate for initial nginx startup...'
  mkdir -p "/etc/letsencrypt/live/${DOMAIN}"
  openssl req -x509 -nodes -newkey rsa:2048 -days 365 \
    -keyout "/etc/letsencrypt/live/${DOMAIN}/privkey.pem" \
    -out "/etc/letsencrypt/live/${DOMAIN}/fullchain.pem" \
    -subj "/CN=${DOMAIN}" 2>/dev/null
  
  if [ -f "/etc/letsencrypt/live/${DOMAIN}/fullchain.pem" ]; then
    echo "Self-signed certificate created for domain: ${DOMAIN}"
  else
    echo 'ERROR: Failed to create certificate'
    exit 1
  fi
else
  echo "Certificate already exists for domain: ${DOMAIN}"
fi

