#!/usr/bin/env bash
# Siembra datos de ejemplo en la instancia de demo (puerto 3007).
# Se ejecuta con la base recien creada por "make demo", asi que no necesita
# ser idempotente: falla rapido si algo no va bien.
set -euo pipefail

BASE="${1:-http://localhost:3007}"
AUTH="demo:demo1234"

api() {
    # $1 = metodo, $2 = ruta, $3 = cuerpo JSON (opcional)
    local method="$1" path="$2" data="${3:-}"
    if [ -n "$data" ]; then
        curl -fsS -u "$AUTH" -H 'Content-Type: application/json' \
            -X "$method" -d "$data" "$BASE$path" 
    else
        curl -fsS -u "$AUTH" -H 'Content-Type: application/json' \
            -X "$method" "$BASE$path"
    fi
}

jget() {
    # Lee JSON por stdin y extrae una expresion python (ej: d["id"]).
    python3 -c "import sys,json; d=json.load(sys.stdin); print($1)"
}

echo "== demo: usuarios extra =="
api POST /api/v2/admin/users '{"username":"alice","email":"alice@example.com","password":"alice1234","full_name":"Alice Ejemplo","activated":true}' > /dev/null
api POST /api/v2/admin/users '{"username":"bob","email":"bob@example.com","password":"bob1234","full_name":"Bob Ejemplo","activated":true}' > /dev/null

echo "== demo: organizacion acme =="
api POST /api/v2/user/settings/orgs '{"username":"acme","full_name":"ACME Corp"}' > /dev/null
ORG_ID=$(api GET /api/v2/orgs/acme | jget 'd["id"]')
echo "   org acme id=$ORG_ID"

echo "== demo: equipos de acme =="
TEAM_ID=$(api POST /api/v2/orgs/acme/teams '{"name":"developers","description":"Equipo de desarrollo","permission":"write"}' | jget 'd["ID"] or d["id"]')
echo "   team developers id=$TEAM_ID"
api POST /api/v2/orgs/teams/$TEAM_ID/members/alice > /dev/null

echo "== demo: repositorios =="
api POST /api/v2/user/repos '{"name":"ggit-demo","description":"Repositorio de ejemplo para capturas","private":false,"auto_init":true,"default_branch":"master"}' > /dev/null
api POST /api/v2/user/repos '{"name":"api-ejemplo","description":"API de ejemplo en Go","private":false,"auto_init":true,"default_branch":"master"}' > /dev/null
api POST /api/v2/user/repos '{"name":"notas-diarias","description":"Notas personales (privado)","private":true,"auto_init":true,"default_branch":"master"}' > /dev/null
api POST /api/v2/user/repos "{\"name\":\"infra\",\"description\":\"Infraestructura de ACME (Terraform)\",\"private\":false,\"auto_init\":true,\"default_branch\":\"master\",\"uid\":$ORG_ID}" > /dev/null
api POST /api/v2/orgs/teams/$TEAM_ID/repos/infra > /dev/null

echo "== demo: colaboradora en ggit-demo =="
api POST /api/v2/demo/ggit-demo/settings/collaborators/alice > /dev/null

echo "== demo: archivos de ggit-demo =="
api POST /api/v2/demo/ggit-demo/editor/edit '{"tree_path":"docs/guia.md","content":"# Guia de uso\n\nRepositorio de ejemplo de Ggit.\n\n## Flujo rapido\n\n1. Crea un repositorio.\n2. Clonalo: `git clone <url>`.\n3. Haz push de tu rama y abre un PR.\n","commit_summary":"docs: add usage guide","is_new_file":true,"branch":"master"}' > /dev/null
api POST /api/v2/demo/ggit-demo/editor/edit '{"tree_path":"docs/adr/0001-stack.md","content":"# ADR 0001: Stack\n\n- Backend: Go + Gin + SQLite.\n- Frontend: Svelte.\n- Decision: priorizar binario pequeno y bajo consumo.\n","commit_summary":"docs: add first ADR","is_new_file":true,"branch":"master"}' > /dev/null
api POST /api/v2/demo/ggit-demo/editor/edit '{"tree_path":"src/main.go","content":"package main\n\nimport \"fmt\"\n\nfunc main() {\n\tfmt.Println(\"hola desde ggit-demo\")\n}\n","commit_summary":"feat: add hello example","is_new_file":true,"branch":"master"}' > /dev/null
api POST /api/v2/demo/ggit-demo/editor/edit '{"tree_path":"CHANGELOG.md","content":"# Changelog\n\n## v0.1.0\n\n- Guia de uso inicial.\n- Ejemplo en Go.\n- Tablero e issues de muestra.\n","commit_summary":"docs: add changelog","is_new_file":true,"branch":"master"}' > /dev/null

echo "== demo: archivos de api-ejemplo =="
api POST /api/v2/demo/api-ejemplo/editor/edit '{"tree_path":"internal/server.go","content":"package internal\n\n// Server expone la API HTTP de ejemplo.\nfunc NewServer(addr string) *Server { return &Server{Addr: addr} }\n\ntype Server struct{ Addr string }\n","commit_summary":"feat: add server skeleton","is_new_file":true,"branch":"master"}' > /dev/null
api POST /api/v2/demo/api-ejemplo/editor/edit '{"tree_path":"docs/api.md","content":"# API de ejemplo\n\n- `GET /health` -> 200 ok\n- `GET /api/v2/version` -> version\n","commit_summary":"docs: add api doc","is_new_file":true,"branch":"master"}' > /dev/null

echo "== demo: archivos de acme/infra =="
api POST /api/v2/acme/infra/editor/edit '{"tree_path":"terraform/main.tf","content":"# Ejemplo minimo\nresource \"null_resource\" \"demo\" {}\n","commit_summary":"feat: add terraform skeleton","is_new_file":true,"branch":"master"}' > /dev/null
api POST /api/v2/acme/infra/editor/edit '{"tree_path":"docs/runbook.md","content":"# Runbook\n\n1. `terraform init`\n2. `terraform plan`\n3. Revisar el PR antes de aplicar.\n","commit_summary":"docs: add runbook","is_new_file":true,"branch":"master"}' > /dev/null

echo "== demo: etiquetas y milestones de ggit-demo =="
BUG_ID=$(api POST /api/v2/demo/ggit-demo/labels '{"name":"bug","color":"#d73a4a"}' | jget 'd["id"]')
ENH_ID=$(api POST /api/v2/demo/ggit-demo/labels '{"name":"enhancement","color":"#a2eeef"}' | jget 'd["id"]')
DOCS_ID=$(api POST /api/v2/demo/ggit-demo/labels '{"name":"docs","color":"#0075ca"}' | jget 'd["id"]')
PRI_ID=$(api POST /api/v2/demo/ggit-demo/labels '{"name":"priority::high","color":"#b60205"}' | jget 'd["id"]')
MS1_ID=$(api POST /api/v2/demo/ggit-demo/milestones '{"title":"v1.0-demo","description":"Hito para la demo con capturas"}' | jget 'd["id"]')
MS2_ID=$(api POST /api/v2/demo/ggit-demo/milestones '{"title":"v1.1-siguiente","description":"Ideas para despues de la demo"}' | jget 'd["id"]')
echo "   labels bug=$BUG_ID enh=$ENH_ID docs=$DOCS_ID pri=$PRI_ID ms1=$MS1_ID ms2=$MS2_ID"

echo "== demo: issues de ggit-demo =="
api POST /api/v2/demo/ggit-demo/issues "{\"title\":\"Anadir capturas al README\",\"content\":\"El README queda vacio en las capturas. Anadir quickstart + captura del dashboard.\",\"label_ids\":[$DOCS_ID,$PRI_ID]}" > /dev/null
api POST /api/v2/demo/ggit-demo/issues "{\"title\":\"El buscador no filtra por lenguaje\",\"content\":\"Pasos: abrir un repo, buscar sin filtro por lenguaje. Severidad: media.\",\"label_ids\":[$BUG_ID]}" > /dev/null
api POST /api/v2/demo/ggit-demo/issues "{\"title\":\"Publicar release v0.1.0\",\"content\":\"checklist: guia, tag, release.\",\"label_ids\":[$ENH_ID]}" > /dev/null
# El POST con milestone_id falla en el backend (UPDATE sin WHERE en el issue aun sin ID),
# asi que el hito se asigna en un segundo paso.
api PUT /api/v2/demo/ggit-demo/issues/1 "{\"milestone_id\":$MS1_ID}" > /dev/null
api PUT /api/v2/demo/ggit-demo/issues/2 "{\"milestone_id\":$MS1_ID}" > /dev/null
api PUT /api/v2/demo/ggit-demo/issues/3 "{\"milestone_id\":$MS2_ID}" > /dev/null
api POST /api/v2/demo/ggit-demo/issues/1/comments '{"body":"De acuerdo, preparo un README con quickstart y dos capturas (dashboard + repo)."}' > /dev/null
api POST /api/v2/demo/ggit-demo/issues/2/comments '{"body":"Lo reproduzco en local. Propongo anadir `language:go` como primer filtro."}' > /dev/null
api PUT /api/v2/demo/ggit-demo/issues/3 '{"is_closed":true}' > /dev/null

echo "== demo: rama + PR + review en ggit-demo =="
api POST /api/v2/demo/ggit-demo/editor/edit '{"tree_path":"docs/guia.md","content":"# Guia de uso\n\nRepositorio de ejemplo de Ggit.\n\n## Flujo rapido\n\n1. Crea un repositorio.\n2. Clonalo: `git clone <url>`.\n3. Haz push de tu rama y abre un PR.\n\n## Capturas\n\n- Dashboard con actividad reciente.\n- Vista de repo con README.\n","commit_summary":"docs: expand guide on feature branch","branch":"feature/guia-ampliada"}' > /dev/null
api POST /api/v2/demo/ggit-demo/pulls '{"title":"Ampliar la guia de uso","content":"Anade seccion de capturas para la demo.","base_branch":"master","head_branch":"feature/guia-ampliada"}' > /dev/null
# NOTA: no se crea review de aprobacion porque la tabla pull_review no existe en
# instalaciones frescas (bug del backend: falta en Tables + las migraciones no
# corren en BD nueva). El PR queda abierto y visible igualmente.

echo "== demo: CI statuses sobre master =="
MASTER_SHA=$(api GET "/api/v2/demo/ggit-demo/branches" | python3 -c "import sys,json; d=json.load(sys.stdin); b=d.get('default') or {}; c=b.get('commit') or {}; print(c.get('id',''))")
if [ -n "$MASTER_SHA" ]; then
    api POST /api/v2/demo/ggit-demo/statuses/$MASTER_SHA '{"state":"success","context":"ci/build","description":"build ok","target_url":"http://localhost:3007/demo/ggit-demo"}' > /dev/null
    api POST /api/v2/demo/ggit-demo/statuses/$MASTER_SHA '{"state":"pending","context":"ci/e2e","description":"e2e en curso","target_url":"http://localhost:3007/demo/ggit-demo"}' > /dev/null
fi

echo "== demo: tag + release de ggit-demo =="
api POST /api/v2/demo/ggit-demo/tags '{"name":"v0.1.0","target":"master","message":"primera demo"}' > /dev/null
api POST /api/v2/demo/ggit-demo/releases '{"tag_name":"v0.1.0","target":"master","name":"v0.1.0 - demo","body":"Primera release de la demo: guia, ejemplo en Go e issues de muestra."}' > /dev/null

echo "== demo: proteccion de master =="
api PUT /api/v2/demo/ggit-demo/branches/master/protection '{"protected":true,"require_pull_request":true,"required_approvals":1}' > /dev/null

echo "== demo: datos ligeros en otros repos =="
api POST /api/v2/demo/api-ejemplo/labels '{"name":"bug","color":"#d73a4a"}' > /dev/null
api POST /api/v2/demo/api-ejemplo/issues '{"title":"Exponer /version","content":"Devolver la version del binario en JSON."}' > /dev/null
api POST /api/v2/demo/api-ejemplo/issues/1/comments '{"body":"Lo marco para v1.1."}' > /dev/null
api POST /api/v2/acme/infra/issues '{"title":"Revisar plan de terraform","content":"Adjuntar la salida de `terraform plan` al PR."}' > /dev/null
api POST /api/v2/acme/infra/issues/1/comments '{"body":"He subido el plan al PR #4 de ejemplo."}' > /dev/null

echo "== demo: gist publico =="
api POST /api/v2/gists '{"description":"Snippet de ejemplo para la demo","is_public":true,"files":[{"filename":"hola.go","content":"package main\n\nimport \"fmt\"\n\nfunc main() { fmt.Println(\"hola\") }\n","language":"Go"}]}' > /dev/null

echo ""
echo "Demo sembrada:"
echo "  http://localhost:3007/demo/ggit-demo"
echo "  http://localhost:3007/demo/api-ejemplo"
echo "  http://localhost:3007/acme/infra"
echo "  http://localhost:3007/org/acme"
