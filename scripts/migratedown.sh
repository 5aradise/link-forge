#!/bin/bash

if [ -f .env ]; then
    source .env
fi

cd sql/schema
goose turso $DATABASE_URL down
# goose sqlite3 ./../../$DATABASE_URL down
