#!/bin/bash
TOKEN=$(curl -s -X POST -H "Content-Type: application/json" -d '{"email":"admin@sorteos.com","password":"Admin123456"}' http://localhost:8080/api/v1/auth/login | jq -r '.data.access_token')

curl -s -H "Authorization: Bearer $TOKEN" "http://localhost:8080/api/v1/admin/categories" > /dev/null
sleep 1
sudo journalctl -u sorteos-api -n 30 --no-pager | grep -A 5 "ERROR"
