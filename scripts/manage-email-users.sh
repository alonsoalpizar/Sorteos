#!/bin/bash

#═══════════════════════════════════════════════════════════════
# Script de Gestión de Usuarios de Email - sorteos.club
#═══════════════════════════════════════════════════════════════

set -e

# Colores para output
RED='\033[0;31m'
GREEN='\033[0;32m'
YELLOW='\033[1;33m'
BLUE='\033[0;34m'
NC='\033[0m' # No Color

# Función para mostrar banner
show_banner() {
    echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
    echo -e "${BLUE}  📧 Gestor de Usuarios de Email - sorteos.club${NC}"
    echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
    echo ""
}

# Función para crear usuario
create_user() {
    echo -e "${YELLOW}📝 Crear Nuevo Usuario de Email${NC}"
    echo ""

    # Solicitar nombre de usuario
    read -p "Nombre de usuario (sin @sorteos.club): " username

    if [ -z "$username" ]; then
        echo -e "${RED}❌ Error: El nombre de usuario no puede estar vacío${NC}"
        exit 1
    fi

    # Verificar si el usuario ya existe
    if id "$username" &>/dev/null; then
        echo -e "${RED}❌ Error: El usuario '$username' ya existe${NC}"
        exit 1
    fi

    # Solicitar contraseña
    echo ""
    read -s -p "Contraseña para $username@sorteos.club: " password
    echo ""
    read -s -p "Confirmar contraseña: " password2
    echo ""

    if [ "$password" != "$password2" ]; then
        echo -e "${RED}❌ Error: Las contraseñas no coinciden${NC}"
        exit 1
    fi

    if [ -z "$password" ]; then
        echo -e "${RED}❌ Error: La contraseña no puede estar vacía${NC}"
        exit 1
    fi

    echo ""
    echo -e "${BLUE}🔧 Creando usuario $username@sorteos.club...${NC}"

    # Crear usuario del sistema
    useradd -m -s /usr/sbin/nologin "$username"

    # Establecer contraseña
    echo "$username:$password" | chpasswd

    # Crear estructura Maildir
    mkdir -p "/home/$username/Maildir"/{cur,new,tmp}
    mkdir -p "/home/$username/Maildir/.Drafts"/{cur,new,tmp}
    mkdir -p "/home/$username/Maildir/.Sent"/{cur,new,tmp}
    mkdir -p "/home/$username/Maildir/.Trash"/{cur,new,tmp}
    mkdir -p "/home/$username/Maildir/.Junk"/{cur,new,tmp}

    # Establecer permisos
    chown -R "$username:$username" "/home/$username/Maildir"
    chmod -R 700 "/home/$username/Maildir"

    # Guardar credenciales
    CREDS_FILE="/opt/Sorteos/mail-server-docs/mail-server-credentials.txt"
    echo "$username@sorteos.club:$password" >> "$CREDS_FILE"
    chmod 600 "$CREDS_FILE"

    echo ""
    echo -e "${GREEN}✅ Usuario creado exitosamente!${NC}"
    echo ""
    echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
    echo -e "📧 Email: ${GREEN}$username@sorteos.club${NC}"
    echo -e "🔑 Contraseña: ${GREEN}$password${NC}"
    echo -e "🌐 Webmail: ${GREEN}https://webmail.sorteos.club${NC}"
    echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
    echo ""
}

# Función para listar usuarios
list_users() {
    echo -e "${YELLOW}📋 Usuarios de Email Configurados${NC}"
    echo ""
    echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"

    # Leer archivo de credenciales
    CREDS_FILE="/opt/Sorteos/mail-server-docs/mail-server-credentials.txt"

    if [ -f "$CREDS_FILE" ]; then
        cat "$CREDS_FILE" | while IFS=: read -r email password; do
            username=$(echo "$email" | cut -d'@' -f1)
            echo -e "📧 ${GREEN}$email${NC}"
            echo -e "   🔑 Contraseña: $password"
            echo -e "   📁 Maildir: /home/$username/Maildir"
            echo ""
        done
    else
        echo -e "${RED}❌ No se encontró el archivo de credenciales${NC}"
    fi

    echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
}

# Función para eliminar usuario
delete_user() {
    echo -e "${YELLOW}🗑️  Eliminar Usuario de Email${NC}"
    echo ""

    # Listar usuarios actuales
    list_users

    # Solicitar nombre de usuario
    read -p "Nombre de usuario a eliminar (sin @sorteos.club): " username

    if [ -z "$username" ]; then
        echo -e "${RED}❌ Error: El nombre de usuario no puede estar vacío${NC}"
        exit 1
    fi

    # Verificar si el usuario existe
    if ! id "$username" &>/dev/null; then
        echo -e "${RED}❌ Error: El usuario '$username' no existe${NC}"
        exit 1
    fi

    # Confirmar eliminación
    echo ""
    echo -e "${RED}⚠️  ADVERTENCIA: Se eliminará el usuario $username@sorteos.club y TODOS sus emails${NC}"
    read -p "¿Estás seguro? (escribe 'SI' para confirmar): " confirm

    if [ "$confirm" != "SI" ]; then
        echo -e "${YELLOW}❌ Eliminación cancelada${NC}"
        exit 0
    fi

    echo ""
    echo -e "${BLUE}🗑️  Eliminando usuario $username@sorteos.club...${NC}"

    # Eliminar usuario y su directorio home
    userdel -r "$username" 2>/dev/null || true

    # Eliminar del archivo de credenciales
    CREDS_FILE="/opt/Sorteos/mail-server-docs/mail-server-credentials.txt"
    if [ -f "$CREDS_FILE" ]; then
        sed -i "/^$username@sorteos.club:/d" "$CREDS_FILE"
    fi

    echo ""
    echo -e "${GREEN}✅ Usuario eliminado exitosamente${NC}"
    echo ""
}

# Función para cambiar contraseña
change_password() {
    echo -e "${YELLOW}🔑 Cambiar Contraseña de Usuario${NC}"
    echo ""

    # Listar usuarios actuales
    list_users

    # Solicitar nombre de usuario
    read -p "Nombre de usuario (sin @sorteos.club): " username

    if [ -z "$username" ]; then
        echo -e "${RED}❌ Error: El nombre de usuario no puede estar vacío${NC}"
        exit 1
    fi

    # Verificar si el usuario existe
    if ! id "$username" &>/dev/null; then
        echo -e "${RED}❌ Error: El usuario '$username' no existe${NC}"
        exit 1
    fi

    # Solicitar nueva contraseña
    echo ""
    read -s -p "Nueva contraseña para $username@sorteos.club: " password
    echo ""
    read -s -p "Confirmar contraseña: " password2
    echo ""

    if [ "$password" != "$password2" ]; then
        echo -e "${RED}❌ Error: Las contraseñas no coinciden${NC}"
        exit 1
    fi

    if [ -z "$password" ]; then
        echo -e "${RED}❌ Error: La contraseña no puede estar vacía${NC}"
        exit 1
    fi

    echo ""
    echo -e "${BLUE}🔧 Cambiando contraseña...${NC}"

    # Cambiar contraseña
    echo "$username:$password" | chpasswd

    # Actualizar archivo de credenciales
    CREDS_FILE="/opt/Sorteos/mail-server-docs/mail-server-credentials.txt"
    if [ -f "$CREDS_FILE" ]; then
        sed -i "/^$username@sorteos.club:/d" "$CREDS_FILE"
        echo "$username@sorteos.club:$password" >> "$CREDS_FILE"
    fi

    echo ""
    echo -e "${GREEN}✅ Contraseña cambiada exitosamente!${NC}"
    echo ""
    echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
    echo -e "📧 Email: ${GREEN}$username@sorteos.club${NC}"
    echo -e "🔑 Nueva Contraseña: ${GREEN}$password${NC}"
    echo -e "${BLUE}════════════════════════════════════════════════════════════════${NC}"
    echo ""
}

# Menú principal
show_menu() {
    show_banner
    echo -e "${BLUE}Selecciona una opción:${NC}"
    echo ""
    echo "  1) 📝 Crear nuevo usuario"
    echo "  2) 📋 Listar usuarios existentes"
    echo "  3) 🔑 Cambiar contraseña de usuario"
    echo "  4) 🗑️  Eliminar usuario"
    echo "  5) 🚪 Salir"
    echo ""
    read -p "Opción: " option

    case $option in
        1) create_user ;;
        2) list_users ;;
        3) change_password ;;
        4) delete_user ;;
        5) echo -e "${GREEN}👋 Hasta luego!${NC}"; exit 0 ;;
        *) echo -e "${RED}❌ Opción inválida${NC}"; exit 1 ;;
    esac
}

# Verificar que se ejecuta como root
if [ "$EUID" -ne 0 ]; then
    echo -e "${RED}❌ Este script debe ejecutarse como root${NC}"
    echo -e "${YELLOW}Usa: sudo $0${NC}"
    exit 1
fi

# Ejecutar menú
show_menu
