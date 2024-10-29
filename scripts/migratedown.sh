#!/bin/bash

if [ -f .env ]; then
    source .env
fi

cd sql/schema
goose $GOOSE_DRIVER $DATABASE_URL down
