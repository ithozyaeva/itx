#!/bin/sh
# Инициализация Let's Encrypt сертификата

apk add --no-cache docker-cli curl

DOMAIN="${DOMAIN}"
CERTBOT_EMAIL="${CERTBOT_EMAIL:-admin@${DOMAIN}}"
ENVIRONMENT="${ENVIRONMENT:-prod}"
NGINX_CONTAINER="${NGINX_CONTAINER:-nginx}"

echo 'Waiting for nginx to be ready...'
sleep 15

reload_nginx() {
  echo 'Reloading nginx to apply new certificate...'
  if command -v docker >/dev/null 2>&1 && docker ps | grep -q "${NGINX_CONTAINER}"; then
    if docker exec "${NGINX_CONTAINER}" nginx -s reload 2>/dev/null; then
      echo 'Nginx reloaded successfully'
    else
      echo 'Warning: Failed to reload nginx'
    fi
  else
    echo 'Warning: Could not reload nginx - please reload manually'
  fi
}

if [ "${ENVIRONMENT}" = "dev" ]; then
  if [ -f "/etc/letsencrypt/live/${DOMAIN}/fullchain.pem" ]; then
    CERT_ISSUER=$(openssl x509 -in "/etc/letsencrypt/live/${DOMAIN}/fullchain.pem" -noout -issuer 2>/dev/null)
    CERT_SUBJECT=$(openssl x509 -in "/etc/letsencrypt/live/${DOMAIN}/fullchain.pem" -noout -subject 2>/dev/null)
    
    if [ "${CERT_ISSUER}" = "${CERT_SUBJECT}" ]; then
      echo 'Self-signed certificate detected. Requesting Let'\''s Encrypt certificate...'
      if certbot certonly --webroot --webroot-path=/var/www/certbot \
        --email "${CERTBOT_EMAIL}" \
        --agree-tos --no-eff-email \
        -d "${DOMAIN}" \
        -d "www.${DOMAIN}" \
        --non-interactive --force-renewal; then
        echo 'Let'\''s Encrypt certificate obtained successfully'
        reload_nginx
      else
        echo 'Let'\''s Encrypt certificate request failed, using self-signed'
      fi
    else
      echo 'Valid Let'\''s Encrypt certificate already exists'
    fi
  else
    echo 'No certificate found, requesting Let'\''s Encrypt certificate...'
    if certbot certonly --webroot --webroot-path=/var/www/certbot \
      --email "${CERTBOT_EMAIL}" \
      --agree-tos --no-eff-email \
      -d "${DOMAIN}" \
      -d "www.${DOMAIN}" \
      --non-interactive; then
      echo 'Let'\''s Encrypt certificate obtained successfully'
      reload_nginx
    else
      echo 'Let'\''s Encrypt certificate request failed'
    fi
  fi
else
  if [ ! -f "/etc/letsencrypt/live/${DOMAIN}/fullchain.pem" ]; then
    echo 'No certificate found for production, requesting Let'\''s Encrypt certificate...'
    if certbot certonly --webroot --webroot-path=/var/www/certbot \
      --email "${CERTBOT_EMAIL}" \
      --agree-tos --no-eff-email \
      -d "${DOMAIN}" \
      -d "www.${DOMAIN}" \
      --non-interactive; then
      echo 'Let'\''s Encrypt certificate obtained successfully'
      reload_nginx
    else
      echo 'Let'\''s Encrypt certificate request failed'
    fi
  else
    # Проверяем, является ли существующий сертификат самоподписанным
    CERT_ISSUER=$(openssl x509 -in "/etc/letsencrypt/live/${DOMAIN}/fullchain.pem" -noout -issuer 2>/dev/null)
    CERT_SUBJECT=$(openssl x509 -in "/etc/letsencrypt/live/${DOMAIN}/fullchain.pem" -noout -subject 2>/dev/null)
    
    if [ "${CERT_ISSUER}" = "${CERT_SUBJECT}" ]; then
      echo 'Self-signed certificate detected in production. Replacing with Let'\''s Encrypt certificate...'
      if certbot certonly --webroot --webroot-path=/var/www/certbot \
        --email "${CERTBOT_EMAIL}" \
        --agree-tos --no-eff-email \
        -d "${DOMAIN}" \
        -d "www.${DOMAIN}" \
        --non-interactive --force-renewal; then
        echo 'Let'\''s Encrypt certificate obtained successfully'
        reload_nginx
      else
        echo 'Let'\''s Encrypt certificate request failed, keeping temporary certificate'
      fi
    else
      echo "Valid Let's Encrypt certificate already exists for domain: ${DOMAIN}"
    fi
  fi
fi

echo 'Certificate initialization complete'
