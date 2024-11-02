#!/bin/bash

if [ -f .env ]; then
    source .env
fi

cd sql/schema
goose turso $DATABASE_URL up
# goose sqlite3 ./../../$DATABASE_URL up
