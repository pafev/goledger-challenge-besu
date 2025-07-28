#! /bin/bash

POSTGRES_USER="user"
POSTGRES_PASSWORD="pass"
POSTGRES_DB="database"

docker rm -f postgres-db-goledger-challenge

docker pull postgres
docker run --name postgres-db-goledger-challenge \
  -e POSTGRES_PASSWORD=$POSTGRES_PASSWORD \
  -e POSTGRES_USER=$POSTGRES_USER \
  -e POSTGRES_DB=$POSTGRES_DB \
  -p 5432:5432 \
  -v db_data:/var/lib/postgresql/data \
  -d postgres

if [ -f "../../.env" ]; then
  ENV_PATH="../../.env"
elif [ -f "../.env" ]; then
  ENV_PATH="../.env"
elif [ -f "./.env" ]; then
  ENV_PATH="./.env"
else
  echo ".env file not found in expected locations."
  exit 1
fi

sed -i "s|^DATABASE_URL=.*|DATABASE_URL=postgresql://$POSTGRES_USER:$POSTGRES_PASSWORD@localhost:5432/$POSTGRES_DB?sslmode=disable|" $ENV_PATH
