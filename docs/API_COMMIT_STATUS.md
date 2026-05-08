# Commit Status API

La Commit Status API permite a servicios externos (como CI/CD, herramientas de análisis estático o linters) reportar el estado de calidad y construcción de un commit específico.

## Estados Soportados
Los estados permitidos son:
- `pending`: La tarea está en ejecución.
- `success`: La tarea se completó con éxito.
- `error`: Hubo un error técnico al intentar ejecutar la tarea.
- `failure`: La tarea se completó pero el resultado fue negativo (ej. tests fallidos).
- `warning`: La tarea se completó con advertencias menores.

## Endpoints de la API (v2)

### 1. Reportar un Estado
`POST /api/v2/repos/:owner/:repo/statuses/:sha`

**Cuerpo (JSON):**
```json
{
  "state": "success",
  "target_url": "https://ci.example.com/build/123",
  "description": "All tests passed successfully",
  "context": "continuous-integration/jenkins"
}
```
- `context`: Identificador único para el servicio (ej. `ci/build`, `lint/style`). Si se envía un nuevo estado con el mismo contexto para el mismo SHA, solo el último será considerado para el estado combinado.

### 2. Obtener Estados de un Commit
`GET /api/v2/repos/:owner/:repo/statuses/:sha`

Retorna la lista de los últimos estados reportados para cada contexto único.

## Integración en la Interfaz (Frontend)

El sistema calcula automáticamente un **Estado Combinado** basado en la precedencia:
`failure/error` > `pending` > `warning` > `success`.

### Lugares donde se muestra:
1. **Lista de Commits**: Aparece un icono (✅/❌/⚠️) junto a cada mensaje de commit.
2. **Explorador de Archivos**: Se muestra en la barra de "Último commit" de la carpeta o archivo actual.
3. **Detalle del Commit**:
   - Icono de estado general en la cabecera.
   - Panel detallado con el desglose de cada servicio (Context), su descripción y enlace directo a los logs externos.
4. **Pull Requests**:
   - El "Merge Box" muestra si el CI ha pasado o si la fusión está bloqueada por fallos en las pruebas.

## Uso Técnico (Backend)

Para obtener el estado en el código Go:
```go
// Obtener estado combinado (string)
status := database.Handle.CommitStatuses().GetCombinedStatusBySHA(ctx, repoID, sha)

// Obtener lista detallada
statuses, err := database.Handle.CommitStatuses().GetLatestBySHA(ctx, repoID, sha)
```
