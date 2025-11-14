#!/bin/bash

# Script de prueba para verificar configuración de emails
# Uso: ./test_email.sh [sendgrid|smtp]

set -e

# Colores
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

echo -e "${BLUE}========================================${NC}"
echo -e "${BLUE}  Test de Configuración de Emails${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""

# Verificar que existe .env
if [ ! -f ".env" ]; then
    echo -e "${RED}❌ Error: No se encontró archivo .env${NC}"
    echo "Copia .env.example o .env.smtp.example a .env"
    exit 1
fi

# Cargar variables de entorno
source .env

# Determinar proveedor
PROVIDER=${1:-$CONFIG_EMAIL_PROVIDER}
if [ -z "$PROVIDER" ]; then
    PROVIDER="sendgrid"
fi

echo -e "${BLUE}Proveedor configurado:${NC} $PROVIDER"
echo ""

# ====================
# Test SendGrid
# ====================
if [ "$PROVIDER" = "sendgrid" ]; then
    echo -e "${YELLOW}Testing SendGrid...${NC}"

    # Verificar API Key
    if [ -z "$CONFIG_SENDGRID_API_KEY" ]; then
        echo -e "${RED}❌ CONFIG_SENDGRID_API_KEY no configurado${NC}"
        exit 1
    fi

    if [[ "$CONFIG_SENDGRID_API_KEY" == *"your_sendgrid"* ]]; then
        echo -e "${RED}❌ CONFIG_SENDGRID_API_KEY es un placeholder${NC}"
        echo "Obtén tu API key en https://app.sendgrid.com/"
        exit 1
    fi

    echo -e "${GREEN}✓${NC} API Key configurado"

    # Verificar FROM Email
    if [ -z "$CONFIG_SENDGRID_FROM_EMAIL" ]; then
        echo -e "${RED}❌ CONFIG_SENDGRID_FROM_EMAIL no configurado${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓${NC} From Email: $CONFIG_SENDGRID_FROM_EMAIL"

    # Verificar FROM Name
    if [ -z "$CONFIG_SENDGRID_FROM_NAME" ]; then
        echo -e "${YELLOW}⚠${NC} CONFIG_SENDGRID_FROM_NAME no configurado (usando default)"
    else
        echo -e "${GREEN}✓${NC} From Name: $CONFIG_SENDGRID_FROM_NAME"
    fi

    # Test de conectividad con SendGrid API
    echo ""
    echo -e "${YELLOW}Probando conectividad con SendGrid API...${NC}"

    HTTP_CODE=$(curl -s -o /dev/null -w "%{http_code}" \
        -X GET "https://api.sendgrid.com/v3/user/profile" \
        -H "Authorization: Bearer $CONFIG_SENDGRID_API_KEY")

    if [ "$HTTP_CODE" = "200" ]; then
        echo -e "${GREEN}✓${NC} Conexión exitosa con SendGrid API"
    elif [ "$HTTP_CODE" = "401" ]; then
        echo -e "${RED}❌ API Key inválido (HTTP 401)${NC}"
        exit 1
    else
        echo -e "${YELLOW}⚠${NC} HTTP Code: $HTTP_CODE"
    fi

# ====================
# Test SMTP
# ====================
elif [ "$PROVIDER" = "smtp" ]; then
    echo -e "${YELLOW}Testing SMTP...${NC}"

    # Verificar configuración
    if [ -z "$CONFIG_SMTP_HOST" ]; then
        echo -e "${RED}❌ CONFIG_SMTP_HOST no configurado${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓${NC} Host: $CONFIG_SMTP_HOST"

    if [ -z "$CONFIG_SMTP_PORT" ]; then
        echo -e "${RED}❌ CONFIG_SMTP_PORT no configurado${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓${NC} Port: $CONFIG_SMTP_PORT"

    if [ -z "$CONFIG_SMTP_FROM_EMAIL" ]; then
        echo -e "${RED}❌ CONFIG_SMTP_FROM_EMAIL no configurado${NC}"
        exit 1
    fi
    echo -e "${GREEN}✓${NC} From Email: $CONFIG_SMTP_FROM_EMAIL"

    # Test de conectividad
    echo ""
    echo -e "${YELLOW}Probando conectividad con servidor SMTP...${NC}"

    # Verificar si telnet está instalado
    if ! command -v telnet &> /dev/null; then
        echo -e "${YELLOW}⚠${NC} telnet no está instalado. Instalando..."
        sudo apt-get install -y telnet > /dev/null 2>&1
    fi

    # Test de conectividad con timeout
    timeout 5 bash -c "cat < /dev/null > /dev/tcp/$CONFIG_SMTP_HOST/$CONFIG_SMTP_PORT" 2>/dev/null

    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓${NC} Servidor SMTP accesible en $CONFIG_SMTP_HOST:$CONFIG_SMTP_PORT"
    else
        echo -e "${RED}❌ No se puede conectar a $CONFIG_SMTP_HOST:$CONFIG_SMTP_PORT${NC}"
        echo "Verifica firewall, DNS y que el servidor SMTP esté corriendo"
        exit 1
    fi

    # Verificar DNS si es un dominio propio
    if [[ "$CONFIG_SMTP_HOST" == *"sorteos"* ]] || [[ "$CONFIG_SMTP_HOST" != "smtp.gmail.com" ]] && [[ "$CONFIG_SMTP_HOST" != "smtp.office365.com" ]]; then
        echo ""
        echo -e "${YELLOW}Verificando registros DNS...${NC}"

        # Verificar MX record
        MX=$(dig +short MX ${CONFIG_SMTP_HOST#mail.} | head -1)
        if [ -z "$MX" ]; then
            echo -e "${YELLOW}⚠${NC} No se encontró registro MX"
        else
            echo -e "${GREEN}✓${NC} MX Record: $MX"
        fi

        # Verificar SPF
        SPF=$(dig +short TXT ${CONFIG_SMTP_HOST#mail.} | grep "v=spf1")
        if [ -z "$SPF" ]; then
            echo -e "${YELLOW}⚠${NC} No se encontró registro SPF (recomendado para deliverability)"
        else
            echo -e "${GREEN}✓${NC} SPF Record encontrado"
        fi

        # Verificar DKIM
        DKIM=$(dig +short TXT default._domainkey.${CONFIG_SMTP_HOST#mail.} | grep "v=DKIM1")
        if [ -z "$DKIM" ]; then
            echo -e "${YELLOW}⚠${NC} No se encontró registro DKIM (recomendado para deliverability)"
        else
            echo -e "${GREEN}✓${NC} DKIM Record encontrado"
        fi
    fi

else
    echo -e "${RED}❌ Proveedor desconocido: $PROVIDER${NC}"
    echo "Usa 'sendgrid' o 'smtp'"
    exit 1
fi

# ====================
# Verificaciones comunes
# ====================
echo ""
echo -e "${YELLOW}Verificando configuración común...${NC}"

# Frontend URL
if [ -z "$CONFIG_FRONTEND_URL" ]; then
    echo -e "${YELLOW}⚠${NC} CONFIG_FRONTEND_URL no configurado (usando default)"
else
    echo -e "${GREEN}✓${NC} Frontend URL: $CONFIG_FRONTEND_URL"
fi

# Verificar backend compilado
if [ ! -f "sorteos-api" ]; then
    echo -e "${YELLOW}⚠${NC} Backend no compilado. Compilando..."
    go build -o sorteos-api cmd/api/main.go
    if [ $? -eq 0 ]; then
        echo -e "${GREEN}✓${NC} Backend compilado exitosamente"
    else
        echo -e "${RED}❌ Error compilando backend${NC}"
        exit 1
    fi
else
    echo -e "${GREEN}✓${NC} Backend compilado"
fi

# ====================
# Resumen
# ====================
echo ""
echo -e "${BLUE}========================================${NC}"
echo -e "${GREEN}✓ Configuración de emails verificada${NC}"
echo -e "${BLUE}========================================${NC}"
echo ""
echo -e "Proveedor: ${GREEN}$PROVIDER${NC}"

if [ "$PROVIDER" = "sendgrid" ]; then
    echo -e "API Key: ${GREEN}Configurado${NC}"
    echo -e "From: ${GREEN}$CONFIG_SENDGRID_FROM_NAME <$CONFIG_SENDGRID_FROM_EMAIL>${NC}"
else
    echo -e "SMTP Host: ${GREEN}$CONFIG_SMTP_HOST:$CONFIG_SMTP_PORT${NC}"
    echo -e "From: ${GREEN}$CONFIG_SMTP_FROM_NAME <$CONFIG_SMTP_FROM_EMAIL>${NC}"
fi

echo ""
echo -e "${BLUE}Próximos pasos:${NC}"
echo "1. Reiniciar backend: sudo systemctl restart sorteos-api"
echo "2. Probar registro de usuario:"
echo ""
echo -e "${YELLOW}curl -X POST http://localhost:8080/api/v1/auth/register \\"
echo "  -H 'Content-Type: application/json' \\"
echo "  -d '{"
echo "    \"email\": \"test@example.com\","
echo "    \"password\": \"Password123!@#\","
echo "    \"accepted_terms\": true,"
echo "    \"accepted_privacy\": true"
echo "  }'${NC}"
echo ""
echo "3. Verificar logs: sudo journalctl -u sorteos-api -f"
echo ""
