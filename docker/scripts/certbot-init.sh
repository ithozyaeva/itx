#!/bin/sh
# Инициализация Let's Encrypt сертификата

DOMAIN="${DOMAIN}"
CERTBOT_EMAIL="${CERTBOT_EMAIL:-admin@${DOMAIN}}"
ENVIRONMENT="${ENVIRONMENT:-prod}"

echo 'Waiting for nginx to be ready...'
sleep 15

if [ "${ENVIRONMENT}" = "dev" ]; then
  if [ -f "/etc/letsencrypt/live/${DOMAIN}/fullchain.pem" ]; then
    CERT_ISSUER=$(openssl x509 -in "/etc/letsencrypt/live/${DOMAIN}/fullchain.pem" -noout -issuer 2>/dev/null)
    CERT_SUBJECT=$(openssl x509 -in "/etc/letsencrypt/live/${DOMAIN}/fullchain.pem" -noout -subject 2>/dev/null)
    
    if [ "${CERT_ISSUER}" = "${CERT_SUBJECT}" ]; then
      echo 'Self-signed certificate detected. Requesting Let'\''s Encrypt certificate...'
      certbot certonly --webroot --webroot-path=/var/www/certbot \
        --email "${CERTBOT_EMAIL}" \
        --agree-tos --no-eff-email \
        -d "${DOMAIN}" \
        -d "www.${DOMAIN}" \
        --non-interactive --force-renewal && \
        echo 'Let'\''s Encrypt certificate obtained successfully' || \
        echo 'Let'\''s Encrypt certificate request failed, using self-signed'
    else
      echo 'Valid Let'\''s Encrypt certificate already exists'
    fi
  else
    echo 'No certificate found, requesting Let'\''s Encrypt certificate...'
    certbot certonly --webroot --webroot-path=/var/www/certbot \
      --email "${CERTBOT_EMAIL}" \
      --agree-tos --no-eff-email \
      -d "${DOMAIN}" \
      -d "www.${DOMAIN}" \
      --non-interactive || echo 'Let'\''s Encrypt certificate request failed'
  fi
else
  if [ ! -f "/etc/letsencrypt/live/${DOMAIN}/fullchain.pem" ]; then
    echo 'No certificate found for production, requesting Let'\''s Encrypt certificate...'
    certbot certonly --webroot --webroot-path=/var/www/certbot \
      --email "${CERTBOT_EMAIL}" \
      --agree-tos --no-eff-email \
      -d "${DOMAIN}" \
      -d "www.${DOMAIN}" \
      --non-interactive && \
      echo 'Let'\''s Encrypt certificate obtained successfully' || \
      echo 'Let'\''s Encrypt certificate request failed'
  else
    # Проверяем, является ли существующий сертификат самоподписанным
    CERT_ISSUER=$(openssl x509 -in "/etc/letsencrypt/live/${DOMAIN}/fullchain.pem" -noout -issuer 2>/dev/null)
    CERT_SUBJECT=$(openssl x509 -in "/etc/letsencrypt/live/${DOMAIN}/fullchain.pem" -noout -subject 2>/dev/null)
    
    if [ "${CERT_ISSUER}" = "${CERT_SUBJECT}" ]; then
      echo 'Self-signed certificate detected in production. Replacing with Let'\''s Encrypt certificate...'
      certbot certonly --webroot --webroot-path=/var/www/certbot \
        --email "${CERTBOT_EMAIL}" \
        --agree-tos --no-eff-email \
        -d "${DOMAIN}" \
        -d "www.${DOMAIN}" \
        --non-interactive --force-renewal && \
        echo 'Let'\''s Encrypt certificate obtained successfully' || \
        echo 'Let'\''s Encrypt certificate request failed, keeping temporary certificate'
    else
      echo "Valid Let's Encrypt certificate already exists for domain: ${DOMAIN}"
    fi
  fi
fi

echo 'Certificate initialization complete'

