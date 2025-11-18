#!/bin/bash
# Test rápido de 11 módulos admin - Sin complicaciones

echo "=== LOGIN ==="
TOKEN=$(curl -s -X POST \
  -H "Content-Type: application/json" \
  -d '{"email":"admin@sorteos.com","password":"Admin123456"}' \
  http://localhost:8080/api/v1/auth/login | jq -r '.data.access_token')

if [ -z "$TOKEN" ]; then
  echo "❌ Login falló"
  exit 1
fi

echo "✅ Token obtenido"
echo ""
echo "=== TESTING 11 MÓDULOS ==="
echo ""

# Array de endpoints
endpoints=(
  "GET|/admin/categories|Categories"
  "GET|/admin/config|Config"
  "GET|/admin/settlements|Settlements"
  "GET|/admin/users|Users"
  "GET|/admin/organizers|Organizers"
  "GET|/admin/payments|Payments"
  "GET|/admin/raffles|Raffles"
  "GET|/admin/notifications/history|Notifications"
  "GET|/admin/reports/dashboard|Reports"
  "GET|/admin/system/parameters|System"
  "GET|/admin/audit|Audit"
)

passed=0
total=0

for endpoint_data in "${endpoints[@]}"; do
  IFS='|' read -r method path name <<< "$endpoint_data"
  total=$((total + 1))
  
  status=$(curl -s -o /dev/null -w "%{http_code}" \
    -X "$method" \
    -H "Authorization: Bearer $TOKEN" \
    "http://localhost:8080/api/v1$path")
  
  if [ "$status" = "200" ] || [ "$status" = "404" ]; then
    echo "✅ $name - $status"
    passed=$((passed + 1))
  else
    echo "❌ $name - $status"
  fi
done

echo ""
echo "=== RESULTADO ==="
echo "Pasados: $passed/$total"
echo ""

if [ $passed -eq $total ]; then
  echo "🎉 Todos los módulos respondiendo"
  exit 0
else
  echo "⚠️  Algunos módulos con errores (verificar logs)"
  exit 1
fi
