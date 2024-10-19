#!/bin/bash

curl -X POST \
    -H "Content-Type: application/json" \
    -H "X-REQUEST-ID: requestId" \
    -d '{}' \
    http://localhost:10001/api/v1/users/login

echo ""

curl -X POST \
    -H "Content-Type: application/json" \
    -H "X-REQUEST-ID: requestId" \
    -d '{"email": "email@email.com", "password": "password@A1"}' \
    http://localhost:10001/api/v1/users/login