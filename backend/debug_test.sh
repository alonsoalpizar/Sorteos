#!/bin/bash
TOKEN=$(curl -s -X POST -H "Content-Type: application/json" -d '{"email":"admin@sorteos.com","password":"Admin123456"}' http://localhost:8080/api/v1/auth/login | jq -r '.data.access_token')

echo "Testing categories..."
curl -H "Authorization: Bearer $TOKEN" "http://localhost:8080/api/v1/admin/categories" 2>&1

echo ""
echo "Checking logs..."
sudo journalctl -u sorteos-api --since "10 seconds ago" --no-pager | tail -30
