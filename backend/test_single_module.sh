#!/bin/bash
# Test individual module with detailed error

MODULE=$1
TOKEN=$(curl -s -X POST -H "Content-Type: application/json" \
  -d '{"email":"admin@sorteos.com","password":"Admin123456"}' \
  http://localhost:8080/api/v1/auth/login | jq -r '.data.access_token')

echo "Testing: $MODULE"
echo "===================="
curl -s -H "Authorization: Bearer $TOKEN" \
  "http://localhost:8080/api/v1/admin/$MODULE" | jq .
echo ""
