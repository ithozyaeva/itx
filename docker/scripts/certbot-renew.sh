#!/bin/sh
# Автоматическое обновление SSL сертификатов

apk add --no-cache docker-cli curl

NGINX_CONTAINER="${NGINX_CONTAINER:-nginx}"

trap exit TERM

while :; do
  if certbot renew --quiet; then
    echo 'Certificates renewed successfully'
    
    if command -v docker >/dev/null 2>&1 && docker ps | grep -q "${NGINX_CONTAINER}"; then
      if docker exec "${NGINX_CONTAINER}" nginx -s reload 2>/dev/null; then
        echo 'Nginx reloaded successfully'
      else
        echo 'Warning: Failed to reload nginx via docker exec, trying to restart container'
        docker restart "${NGINX_CONTAINER}" 2>/dev/null || \
          echo 'Warning: Could not restart nginx container'
      fi
    else
      echo 'Warning: Docker not available or nginx container not found, certificates renewed but nginx not reloaded'
      echo "Please reload nginx manually: docker exec ${NGINX_CONTAINER} nginx -s reload"
    fi
  else
    echo 'No certificates renewed'
  fi
  
  sleep 12h & wait ${!}
done

