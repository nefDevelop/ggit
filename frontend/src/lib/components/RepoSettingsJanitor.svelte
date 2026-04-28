<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from "svelte";
  import { toast } from "$lib/toast.svelte";
  import Button from "./ui/Button.svelte";
  import Badge from "./ui/Badge.svelte";
  import { tool } from "$lib/tool";

  const repo: any = getContext("repo");
  let loading = $state(true);
  let scanning = $state(false);
  let cleaning = $state(false);
  let report = $state<any>(null);
  let selectedPaths = $state<string[]>([]);

  async function fetchReport() {
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/janitor/report`);
      if (res.ok) {
        report = await res.json();
        selectedPaths = report.issues?.map((i: any) => i.path) || [];
      }
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function runScan() {
    scanning = true;
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/janitor/scan`, {
        method: "POST"
      });
      if (res.ok) {
        report = await res.json();
        selectedPaths = report.issues?.map((i: any) => i.path) || [];
        toast.success("Scan complete!");
      }
    } catch (e) {
      console.error(e);
    } finally {
      scanning = false;
    }
  }

  async function runCleanup() {
    if (selectedPaths.length === 0) return;
    cleaning = true;
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/janitor/cleanup`, {
        method: "POST",
        body: JSON.stringify({ paths: selectedPaths })
      });
      if (res.ok) {
        toast.success("Cleanup successful! A new commit has been created.");
        report = null;
        selectedPaths = [];
        fetchReport();
      }
    } catch (e) {
      console.error(e);
      toast.error("Failed to perform cleanup.");
    } finally {
      cleaning = false;
    }
  }

  function toggleSelect(path: string) {
    if (selectedPaths.includes(path)) {
      selectedPaths = selectedPaths.filter(p => p !== path);
    } else {
      selectedPaths = [...selectedPaths, path];
    }
  }

  onMount(fetchReport);
</script>

<div class="janitor-panel">
  <div class="settings-card">
    <header class="card-header">
      <div class="header-content">
        <span class="material-symbols-rounded icon">mop</span>
        <div>
          <h3>Ggit Janitor 🧹</h3>
          <p>Análisis de salud del repositorio para detectar basura técnica y archivos innecesarios.</p>
        </div>
      </div>
      <div class="header-actions">
        <Button variant="outline" icon="search" onclick={runScan} loading={scanning} disabled={cleaning}>
          {report ? 'Volver a Escanear' : 'Escanear Ahora'}
        </Button>
        {#if report && report.issues?.length > 0}
          <Button variant="primary" icon="cleaning_services" onclick={runCleanup} loading={cleaning} disabled={scanning || selectedPaths.length === 0}>
            Limpiar Seleccionados ({selectedPaths.length})
          </Button>
        {/if}
      </div>
    </header>

    <div class="card-body">
      {#if loading}
        <div class="loading-state">Cargando reporte...</div>
      {:else if !report || !report.issues || report.issues.length === 0}
        <div class="empty-state tonal">
          <span class="material-symbols-rounded">verified</span>
          <p>¡Tu repositorio está impecable! No se han detectado archivos innecesarios.</p>
        </div>
      {:else}
        <div class="summary-bar tonal">
          <div class="summary-item">
            <span class="value">{report.issues.length}</span>
            <span class="label">Asuntos encontrados</span>
          </div>
          {#if report.total_waste > 0}
            <div class="summary-item">
              <span class="value">{tool.fileSize(report.total_waste)}</span>
              <span class="label">Espacio desperdiciado</span>
            </div>
          {/if}
          <div class="summary-item">
            <span class="value">{new Date(report.generated_at * 1000).toLocaleString()}</span>
            <span class="label">Último análisis</span>
          </div>
        </div>

        <div class="issues-list">
          {#each report.issues as issue}
            <div 
              class="issue-item" 
              class:critical={issue.severity === 'critical'} 
              class:selected={selectedPaths.includes(issue.path)} 
              onclick={() => toggleSelect(issue.path)}
              onkeydown={(e) => { if (e.key === 'Enter' || e.key === ' ') { e.preventDefault(); toggleSelect(issue.path); } }}
              role="button"
              tabindex="0"
            >
              <div class="issue-selection">
                <span class="material-symbols-rounded">{selectedPaths.includes(issue.path) ? 'check_box' : 'check_box_outline_blank'}</span>
              </div>
              <div class="issue-type-icon">
                {#if issue.type === 'dependency'}
                  <span class="material-symbols-rounded">package_2</span>
                {:else if issue.type === 'build'}
                  <span class="material-symbols-rounded">construction</span>
                {:else if issue.type === 'system'}
                  <span class="material-symbols-rounded">desktop_windows</span>
                {:else}
                  <span class="material-symbols-rounded">file_present</span>
                {/if}
              </div>
              <div class="issue-info">
                <div class="issue-title">
                  <code class="path">{issue.path}</code>
                  <Badge variant={issue.severity === 'critical' ? 'danger' : 'warning'}>{issue.type}</Badge>
                </div>
                <p class="issue-desc">{issue.description}</p>
              </div>
              {#if issue.size > 0}
                <div class="issue-size">{tool.fileSize(issue.size)}</div>
              {/if}
            </div>
          {/each}
        </div>

        <div class="cleanup-hint">
          <span class="material-symbols-rounded">magic_button</span>
          <p>Al pulsar <b>Limpiar</b>, Ggit eliminará los archivos seleccionados y añadirá sugerencias a tu <code>.gitignore</code> automáticamente mediante un nuevo commit.</p>
        </div>
      {/if}
    </div>
  </div>
</div>

<style>
  .janitor-panel { display: flex; flex-direction: column; gap: 2rem; }
  .settings-card { background: var(--md-sys-color-surface-container); border-radius: 2rem; overflow: hidden; border: 1px solid var(--md-sys-color-outline-variant); }
  
  .card-header { padding: 2rem; display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid var(--md-sys-color-outline-variant); }
  .header-content { display: flex; gap: 1.5rem; align-items: center; }
  .header-actions { display: flex; gap: 1rem; }
  .icon { font-size: 3rem; color: var(--primary-color); }
  .card-header h3 { margin: 0; font-size: 1.5rem; font-weight: 900; }
  .card-header p { margin: 0.25rem 0 0; color: var(--text-muted); }

  .card-body { padding: 2.5rem; }

  .summary-bar { display: flex; gap: 3rem; padding: 1.5rem 2rem; border-radius: 1.5rem; margin-bottom: 2rem; }
  .summary-item { display: flex; flex-direction: column; }
  .summary-item .value { font-size: 1.25rem; font-weight: 900; color: var(--primary-color); }
  .summary-item .label { font-size: 0.8rem; color: var(--text-muted); text-transform: uppercase; letter-spacing: 0.05em; }

  .issues-list { display: flex; flex-direction: column; gap: 1rem; }
  .issue-item { display: flex; align-items: center; gap: 1.5rem; padding: 1.25rem; background: var(--md-sys-color-surface-container-high); border-radius: 1.25rem; border: 1px solid var(--md-sys-color-outline-variant); cursor: pointer; transition: all 0.2s; }
  .issue-item:hover { background: var(--md-sys-color-surface-container-highest); }
  .issue-item.selected { border-color: var(--primary-color); background: rgba(var(--primary-rgb), 0.05); }
  .issue-item.critical { border-left: 4px solid var(--md-sys-color-error); }
  
  .issue-selection { color: var(--primary-color); }
  .issue-type-icon { width: 40px; height: 40px; border-radius: 50%; background: rgba(var(--primary-rgb), 0.1); color: var(--primary-color); display: flex; align-items: center; justify-content: center; flex-shrink: 0; }
  .issue-info { flex: 1; }
  .issue-title { display: flex; align-items: center; gap: 1rem; margin-bottom: 0.25rem; }
  .path { font-family: var(--font-mono); font-weight: 600; font-size: 0.95rem; }
  .issue-desc { margin: 0; font-size: 0.9rem; color: var(--text-muted); }
  .issue-size { font-weight: 800; color: var(--text-muted); font-size: 0.9rem; }
  .empty-state { display: flex; flex-direction: column; align-items: center; gap: 1rem; padding: 4rem; border-radius: 2rem; text-align: center; }
  .empty-state span { font-size: 4rem; color: var(--md-sys-color-primary); }
  .empty-state p { margin: 0; font-size: 1.1rem; color: var(--text-muted); font-weight: 600; }

  .cleanup-hint { margin-top: 2rem; display: flex; gap: 1rem; align-items: center; padding: 1.25rem; background: rgba(var(--primary-rgb), 0.05); border-radius: 1rem; color: var(--text-muted); font-size: 0.9rem; }
  .cleanup-hint span { color: var(--primary-color); }
  .cleanup-hint code { background: rgba(0,0,0,0.1); padding: 0.1rem 0.3rem; border-radius: 0.3rem; }

  @media (max-width: 600px) {
    .summary-bar { flex-direction: column; gap: 1.5rem; }
    .issue-item { flex-direction: column; align-items: flex-start; gap: 1rem; }
    .issue-size { align-self: flex-end; }
  }
</style>
