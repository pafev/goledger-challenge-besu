#!/bin/bash

if [ -f "../../.env" ]; then
  CWD_PATH="../.."
elif [ -f "../.env" ]; then
  CWD_PATH=".."
elif [ -f "./.env" ]; then
  CWD_PATH="."
else
  echo ".env file not found in expected locations."
  exit 1
fi

docker compose -f "$CWD_PATH/scripts/db/docker-compose.yml" down || docker-compose -f "$CWD_PATH/scripts/db/docker-compose.yml" down
docker compose -f "$CWD_PATH/scripts/db/docker-compose.yml" down -v || docker-compose -f "$CWD_PATH/scripts/db/docker-compose.yml" down -v

docker compose -f "$CWD_PATH/scripts/db/docker-compose.yml" up -d || docker-compose -f "$CWD_PATH/scripts/db/docker-compose.yml" up -d

while ! docker exec goledger_challenge_db pg_isready -U postgres -d goledger_challenge >/dev/null 2>&1; do
  sleep 2
done

sed -i "s|^DATABASE_URL=.*|DATABASE_URL=postgresql://user:pass@localhost:5432/goledger_challenge?sslmode=disable|" "$CWD_PATH/.env"

echo ""
echo "URL = postgresql://user:pass@localhost:5432/goledger_challenge"
echo ""
echo "to interact to psql shell:"
echo "  docker exec -it goledger_challenge_db psql -U user -d goledger_challenge"
echo ""
echo "to view log:"
echo "  docker compose logs -f postgres"
echo ""
echo "to stop:"
echo "  docker compose down"
