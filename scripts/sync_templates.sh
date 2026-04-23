#!/bin/bash

# Configuración
REPO_URL="https://github.com/github/gitignore.git"
TEMP_DIR=".temp_gitignore"
TARGET_DIR="conf/gitignore"

echo "🚀 Iniciando sincronización de plantillas Gitignore..."

# 1. Limpieza de instalaciones previas
rm -rf "$TEMP_DIR"
rm -rf "$TARGET_DIR"

# 2. Clonado ligero del repositorio oficial
echo "📥 Descargando desde $REPO_URL..."
git clone --depth 1 "$REPO_URL" "$TEMP_DIR" || { echo "❌ Error al clonar"; exit 1; }

# 3. Creación de estructura de directorios
mkdir -p "$TARGET_DIR/languages"
mkdir -p "$TARGET_DIR/global"

echo "📂 Reestructurando archivos..."

# 4. Mover lenguajes (los que están en la raíz)
find "$TEMP_DIR" -maxdepth 1 -name "*.gitignore" -exec cp {} "$TARGET_DIR/languages/" \;

# 5. Mover globales (los que están en la carpeta Global)
find "$TEMP_DIR/Global" -name "*.gitignore" -exec cp {} "$TARGET_DIR/global/" \;

# 6. Crear el archivo Ggit Security Base (_ggit_security.gitignore)
echo "🛡️ Generando archivo de seguridad base (_ggit_security.gitignore)..."

cat <<EOT > "$TARGET_DIR/global/_ggit_security.gitignore"
# ##########################################
# Ggit Security Base - Recomendado
# ##########################################

# Identificadores de entorno y secretos
.env
.env.local
.env.development.local
.env.test.local
.env.production.local
*.pem
*.key
auth.json

# Archivos de Sistema
.DS_Store
Thumbs.db
Desktop.ini

# Logs y Temporales
logs
*.log
npm-debug.log*
yarn-debug.log*
yarn-error.log*
.eslintcache
.stylelintcache

# Carpetas de dependencias comunes (fallback)
node_modules/
bower_components/
jspm_packages/

# End of Ggit Security Base
EOT

# 7. Limpieza
rm -rf "$TEMP_DIR"

echo "✅ Sincronización completada con éxito."
echo "📍 Las plantillas están disponibles en: $TARGET_DIR"
