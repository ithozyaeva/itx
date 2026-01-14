#!/bin/sh
# Создание self-signed сертификата для dev окружения или временного для production

apk add --no-cache openssl
DOMAIN="${DOMAIN}"
ENVIRONMENT="${ENVIRONMENT:-prod}"

# Проверяем, нужен ли сертификат
if [ ! -f "/etc/letsencrypt/live/${DOMAIN}/fullchain.pem" ] || ! openssl x509 -checkend 86400 -noout -in "/etc/letsencrypt/live/${DOMAIN}/fullchain.pem" 2>/dev/null; then
  if [ "${ENVIRONMENT}" = "dev" ]; then
    echo 'Creating self-signed certificate for dev environment...'
  else
    echo 'Creating temporary self-signed certificate for production (will be replaced by Let'\''s Encrypt)...'
  fi
  
  mkdir -p "/etc/letsencrypt/live/${DOMAIN}"
  openssl req -x509 -nodes -newkey rsa:2048 -days 365 \
    -keyout "/etc/letsencrypt/live/${DOMAIN}/privkey.pem" \
    -out "/etc/letsencrypt/live/${DOMAIN}/fullchain.pem" \
    -subj "/CN=${DOMAIN}" 2>/dev/null
  
  if [ -f "/etc/letsencrypt/live/${DOMAIN}/fullchain.pem" ]; then
    if [ "${ENVIRONMENT}" = "dev" ]; then
      echo "Self-signed certificate created for domain: ${DOMAIN}"
    else
      echo "Temporary self-signed certificate created for domain: ${DOMAIN} (will be replaced by certbot-init)"
    fi
  else
    echo 'ERROR: Failed to create certificate'
    exit 1
  fi
else
  echo "Certificate already exists for domain: ${DOMAIN}"
fi

