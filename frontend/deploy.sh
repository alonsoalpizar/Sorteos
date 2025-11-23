#!/bin/bash
# Script de deploy seguro para el frontend
# Preserva la carpeta /uploads para no perder imágenes subidas

set -e

DEST="/var/www/sorteos.club"
SRC="/opt/Sorteos/frontend/dist"

echo "Compilando frontend..."
cd /opt/Sorteos/frontend
npm run build

echo "Desplegando a $DEST..."

# Crear backup de uploads si existe
if [ -d "$DEST/uploads" ]; then
    echo "Preservando carpeta uploads..."
    sudo mv "$DEST/uploads" /tmp/sorteos_uploads_backup
fi

# Limpiar destino (excepto uploads que ya movimos)
sudo rm -rf "$DEST"/*

# Copiar nuevos archivos
sudo cp -r "$SRC"/* "$DEST/"

# Restaurar uploads
if [ -d "/tmp/sorteos_uploads_backup" ]; then
    sudo mv /tmp/sorteos_uploads_backup "$DEST/uploads"
    echo "Carpeta uploads restaurada"
else
    sudo mkdir -p "$DEST/uploads/raffles"
fi

# Permisos
sudo chown -R www-data:www-data "$DEST"

echo "Deploy completado exitosamente"
