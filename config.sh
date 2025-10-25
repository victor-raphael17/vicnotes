#!/bin/bash

set -euo pipefail

cd backend

npm init -y

npm install express pg cors dotenv uuid

npm install --save-dev nodemon
