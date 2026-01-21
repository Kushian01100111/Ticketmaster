#!/usr/bin/env bash
set -euo pipefail

COMPOSE_FILE="${COMPOSE_FILE:-docker-compose.yml}"
MONGO_SERVICE="${MONGO_SERVICE:-mongodb}"

# Asegura carpeta de datos (bind mount)
mkdir -p ./data

echo "==> Levantando MongoDB..."
docker compose -f "$COMPOSE_FILE" up -d "$MONGO_SERVICE"

echo "==> Esperando a que MongoDB esté listo..."
# En tu compose definís root user/pass, eso autentica contra 'admin'
# Probamos un ping hasta que responda.
for i in {1..60}; do
  if docker compose -f "$COMPOSE_FILE" exec -T "$MONGO_SERVICE" \
      mongosh "mongodb://user:pass@localhost:27017/admin" \
      --quiet --eval "db.adminCommand({ ping: 1 })" >/dev/null 2>&1; then
    echo "==> MongoDB listo."
    break
  fi
  sleep 1
  if [[ $i -eq 60 ]]; then
    echo "ERROR: MongoDB no respondió a tiempo." >&2
    docker compose -f "$COMPOSE_FILE" logs "$MONGO_SERVICE" >&2 || true
    exit 1
  fi
done

echo "==> Verificación rápida (show dbs):"
docker compose -f "$COMPOSE_FILE" exec -T "$MONGO_SERVICE" \
  mongosh "mongodb://user:pass@localhost:27017/admin" \
  --eval "show dbs"