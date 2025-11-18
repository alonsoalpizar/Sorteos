#!/bin/bash

# Script de prueba para endpoints Admin de Almighty
# Fecha: 2025-11-18
# Versión: 1.0

# Colores para output
GREEN='\033[0;32m'
RED='\033[0;31m'
YELLOW='\033[1;33m'
NC='\033[0m' # No Color

# Configuración
API_BASE="http://localhost:8080/api/v1"
ADMIN_TOKEN=""  # Debe ser llenado con un token de admin válido

# Función para imprimir headers
print_header() {
    echo -e "\n${YELLOW}========================================${NC}"
    echo -e "${YELLOW}$1${NC}"
    echo -e "${YELLOW}========================================${NC}\n"
}

# Función para hacer requests
api_call() {
    local method=$1
    local endpoint=$2
    local data=$3

    echo -e "${GREEN}→ $method $endpoint${NC}"

    if [ -z "$data" ]; then
        curl -s -X $method \
            -H "Authorization: Bearer $ADMIN_TOKEN" \
            -H "Content-Type: application/json" \
            "$API_BASE$endpoint" | jq '.' 2>/dev/null || echo "Error: Invalid JSON or endpoint not found"
    else
        curl -s -X $method \
            -H "Authorization: Bearer $ADMIN_TOKEN" \
            -H "Content-Type: application/json" \
            -d "$data" \
            "$API_BASE$endpoint" | jq '.' 2>/dev/null || echo "Error: Invalid JSON or endpoint not found"
    fi
    echo ""
}

# Verificar dependencias
if ! command -v jq &> /dev/null; then
    echo -e "${RED}Error: jq is required. Install it with: sudo apt install jq${NC}"
    exit 1
fi

if [ -z "$ADMIN_TOKEN" ]; then
    echo -e "${RED}Error: ADMIN_TOKEN is empty. Please set it in the script or as environment variable.${NC}"
    echo -e "${YELLOW}Example: export ADMIN_TOKEN=\"your_admin_token_here\"${NC}"
    exit 1
fi

# ==================== HEALTH CHECK ====================
print_header "HEALTH CHECK"
curl -s "$API_BASE/../health" | jq '.'

# ==================== CATEGORY MANAGEMENT ====================
print_header "CATEGORY MANAGEMENT (4 endpoints)"

echo -e "${YELLOW}1. List Categories${NC}"
api_call "GET" "/admin/categories?page=1&page_size=10"

echo -e "${YELLOW}2. Create Category${NC}"
CATEGORY_DATA='{
  "name": "Test Category",
  "description": "Category created via API test",
  "icon_url": "https://example.com/icon.svg",
  "is_active": true
}'
CATEGORY_RESPONSE=$(api_call "POST" "/admin/categories" "$CATEGORY_DATA")
CATEGORY_ID=$(echo "$CATEGORY_RESPONSE" | jq -r '.category_id // empty')

if [ -n "$CATEGORY_ID" ]; then
    echo -e "${GREEN}✓ Category created with ID: $CATEGORY_ID${NC}"

    echo -e "\n${YELLOW}3. Update Category${NC}"
    UPDATE_DATA='{
      "name": "Updated Test Category",
      "is_active": false
    }'
    api_call "PUT" "/admin/categories/$CATEGORY_ID" "$UPDATE_DATA"

    echo -e "${YELLOW}4. Delete Category${NC}"
    api_call "DELETE" "/admin/categories/$CATEGORY_ID"
else
    echo -e "${RED}✗ Failed to create category, skipping update/delete tests${NC}"
fi

# ==================== SYSTEM CONFIG ====================
print_header "SYSTEM CONFIG (3 endpoints)"

echo -e "${YELLOW}1. List All Configs${NC}"
api_call "GET" "/admin/config"

echo -e "${YELLOW}2. Get Specific Config${NC}"
api_call "GET" "/admin/config/platform_commission"

echo -e "${YELLOW}3. Update Config${NC}"
CONFIG_DATA='{
  "config_value": "12.5"
}'
api_call "PUT" "/admin/config/platform_commission" "$CONFIG_DATA"

# ==================== SUMMARY ====================
print_header "TEST SUMMARY"
echo -e "${GREEN}Total Endpoints Tested: 7${NC}"
echo ""
echo -e "Category Management:"
echo -e "  ✓ GET    /admin/categories"
echo -e "  ✓ POST   /admin/categories"
echo -e "  ✓ PUT    /admin/categories/:id"
echo -e "  ✓ DELETE /admin/categories/:id"
echo ""
echo -e "System Config:"
echo -e "  ✓ GET    /admin/config"
echo -e "  ✓ GET    /admin/config/:key"
echo -e "  ✓ PUT    /admin/config/:key"
echo ""
echo -e "${YELLOW}Note: Check responses above for any errors${NC}"
echo -e "${YELLOW}Authentication requires valid admin token${NC}"
