#!/bin/bash

set -euo pipefail

docker compose down

docker compose up -d --build --force-recreate

sleep 10

cd frontend
#npm install  #(only first time)
npm run dev
