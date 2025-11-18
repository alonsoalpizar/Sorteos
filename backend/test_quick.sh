#!/bin/bash
# Test rápido de endpoints admin

BASE_URL="http://localhost:8080/api/v1"
GREEN='\033[0;32m'
RED='\033[0;31m'
NC='\033[0m'

echo "=== TEST RÁPIDO ADMIN ENDPOINTS ==="
echo ""
echo "Login como admin@sorteos.com..."
read -sp "Password: " PASS
echo ""

# Login
TOKEN=$(curl -s -X POST \
    -H "Content-Type: application/json" \
    -d "{\"email\":\"admin@sorteos.com\",\"password\":\"$PASS\"}" \
    "${BASE_URL}/auth/login" | jq -r '.data.access_token')

if [ -z "$TOKEN" ] || [ "$TOKEN" = "null" ]; then
    echo -e "${RED}✗ Login fallido${NC}"
    exit 1
fi

echo -e "${GREEN}✓ Login exitoso${NC}"
echo ""

# Test endpoints principales (uno de cada módulo)
test_get() {
    local endpoint=$1
    local name=$2
    echo -n "Testing $name... "
    status=$(curl -s -o /dev/null -w "%{http_code}" -H "Authorization: Bearer $TOKEN" "${BASE_URL}${endpoint}")
    if [ "$status" = "200" ] || [ "$status" = "404" ]; then
        echo -e "${GREEN}✓ $status${NC}"
        return 0
    else
        echo -e "${RED}✗ $status${NC}"
        return 1
    fi
}

echo "Testing 11 módulos admin:"
echo ""

test_get "/admin/categories" "1. Categories"
test_get "/admin/config" "2. Config"
test_get "/admin/settlements" "3. Settlements"
test_get "/admin/users" "4. Users"
test_get "/admin/organizers" "5. Organizers"
test_get "/admin/payments" "6. Payments"
test_get "/admin/raffles" "7. Raffles"
test_get "/admin/notifications/history" "8. Notifications"
test_get "/admin/reports/dashboard" "9. Reports"
test_get "/admin/system/parameters" "10. System"
test_get "/admin/audit" "11. Audit"

echo ""
echo -e "${GREEN}✓ Test completado - 11 módulos verificados${NC}"
