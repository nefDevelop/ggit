#!/bin/bash

INPUT_DIR="conf/localeBORRAR"
OUTPUT_DIR="frontend/src/locales"

# --- Patching Step ---
echo "Parcheando archivos .ini si es necesario..."

# Lista de archivos y traducciones a añadir
declare -A patches
patches["conf/locale/locale_cs-CZ.ini"]="[admin]\naudit_log=Protokol auditu"
patches["conf/locale/locale_ro-RO.ini"]="[admin]\naudit_log=Jurnal de audit"
patches["conf/locale/locale_es-ES.ini"]="[admin]\naudit_log=Registro de auditoría"
patches["conf/locale/locale_gl-ES.ini"]="[admin]\naudit_log=Rexistro de auditoría"
patches["conf/locale/locale_en-US.ini"]="[admin]\naudit_log=Audit Log"

for file in "${!patches[@]}"; do
    if [ -f "$file" ]; then
        # Solo añade la sección si no existe para evitar duplicados
        if ! grep -q "\[admin\]" "$file"; then
            echo -e "\n${patches[$file]}" >> "$file"
            echo " -> Archivo '$file' parcheado."
        fi
    fi
done

echo "Parcheo finalizado."
# --- End Patching Step ---
# Crear el directorio de salida si no existe
mkdir -p "$OUTPUT_DIR"

# Crear un script temporal de Node.js para hacer el parseo de forma segura
cat << 'EOF' > /tmp/convert_ini_to_json.js
const fs = require('fs');

const inputFile = process.argv[2];
const outputFile = process.argv[3];

const content = fs.readFileSync(inputFile, 'utf-8');
const lines = content.split('\n');
const result = {};
let currentSection = '';

lines.forEach(line => {
    line = line.trim();
    // Ignorar líneas vacías y comentarios
    if (!line || line.startsWith(';') || line.startsWith('#')) return;
    
    // Detectar secciones (ej: [auth])
    if (line.startsWith('[') && line.endsWith(']')) {
        currentSection = line.slice(1, -1).trim();
        return;
    }
    
    const eqIndex = line.indexOf('=');
    if (eqIndex === -1) return;
    
    let key = line.slice(0, eqIndex).trim();
    let value = line.slice(eqIndex + 1).trim();
    
    // Limpiar las comillas invertidas que usa Gogs/Ggit (ej: require_error=` no puede estar vacío.`)
    if (value.startsWith('`') && value.endsWith('`')) {
        value = value.slice(1, -1);
    } else if (value.startsWith('`')) {
        value = value.slice(1);
    }
    
    // Anidar dentro de la sección actual (ej: form.UserName -> result.form.UserName)
    const fullKey = currentSection ? `${currentSection}.${key}` : key;
    const keys = fullKey.split('.');
    let current = result;
    
    for (let i = 0; i < keys.length - 1; i++) {
        current[keys[i]] = current[keys[i]] || {};
        current = current[keys[i]];
    }
    current[keys[keys.length - 1]] = value;
});

fs.writeFileSync(outputFile, JSON.stringify(result, null, 2));
EOF

echo "Iniciando conversión de archivos .ini a .json..."

for ini_file in "$INPUT_DIR"/locale_*.ini; do
    if [ ! -f "$ini_file" ]; then continue; fi
    
    filename=$(basename "$ini_file")
    lang_code=${filename#locale_} # Elimina 'locale_'
    lang_code=${lang_code%.ini}   # Elimina '.ini'
    
    json_file="$OUTPUT_DIR/${lang_code}.json"
    
    echo " -> Convirtiendo $filename a ${lang_code}.json"
    node /tmp/convert_ini_to_json.js "$ini_file" "$json_file"
done

# Limpiar script temporal
rm /tmp/convert_ini_to_json.js

echo "¡Conversión completada exitosamente! Los archivos están en $OUTPUT_DIR"