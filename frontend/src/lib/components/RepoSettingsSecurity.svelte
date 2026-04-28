<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from "svelte";
  import { toast } from "$lib/toast.svelte";
  import { dialog } from "$lib/dialog.svelte";
  import Checkbox from "./ui/Checkbox.svelte";
  import Button from "./ui/Button.svelte";
  import Input from "./ui/Input.svelte";

  const repo: any = getContext("repo");
  let loading = $state(true);
  let saving = $state(false);
  let data = $state<any>({
    enabled: false,
    strict: true,
    patterns: []
  });

  let newPattern = $state({
    name: "",
    regex: "",
    description: ""
  });

  async function fetchData() {
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/security/guardian`);
      if (res.ok) {
        const json = await res.json();
        data = { ...json, patterns: json.patterns || [] };
      }
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function saveSettings() {
    saving = true;
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/security/guardian`, {
        method: "PATCH",
        body: JSON.stringify({
          enabled: data.enabled,
          strict: data.strict
        })
      });
      if (res.ok) {
        toast.success("Security settings updated.");
      }
    } catch (e) {
      console.error(e);
    } finally {
      saving = false;
    }
  }

  async function addPattern() {
    if (!newPattern.name || !newPattern.regex) return;
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/security/guardian/patterns`, {
        method: "POST",
        body: JSON.stringify(newPattern)
      });
      if (res.ok) {
        toast.success("Pattern added.");
        newPattern = { name: "", regex: "", description: "" };
        fetchData();
      }
    } catch (e) {
      console.error(e);
    }
  }

  async function deletePattern(id: number) {
    if (!await dialog.confirm("¿Estás seguro de eliminar este patrón personalizado?")) return;
    try {
      const res = await apiFetch(`/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/security/guardian/patterns/${id}`, {
        method: "DELETE"
      });
      if (res.ok) {
        toast.success("Pattern deleted.");
        fetchData();
      }
    } catch (e) {
      console.error(e);
    }
  }

  onMount(fetchData);
</script>

<div class="security-settings">
  <div class="settings-card">
    <header class="card-header">
      <div class="header-content">
        <span class="material-symbols-rounded guardian-icon">shield</span>
        <div>
          <h3>Secret Guardian 🛡️</h3>
          <p>Protección proactiva contra la subida de secretos, API keys y archivos sensibles.</p>
        </div>
      </div>
      <div class="header-action">
        <Checkbox 
          label="Habilitar Guardián" 
          bind:checked={data.enabled} 
          onchange={saveSettings}
        />
      </div>
    </header>

    {#if data.enabled}
      <div class="card-body">
        <div class="mode-selector tonal">
          <Checkbox 
            label="Modo Estricto (Bloquear Push)" 
            description="Si se detecta un secreto, el push será rechazado inmediatamente." 
            bind:checked={data.strict}
            onchange={saveSettings}
          />
        </div>

        <section class="patterns-section">
          <header class="section-header">
            <h4>Patrones de Detección</h4>
            <p>Se listan los patrones globales y personalizados aplicados a este repositorio.</p>
          </header>

          <div class="patterns-list">
            {#each data.patterns as pattern}
              <div class="pattern-item" class:global={pattern.IsGlobal}>
                <div class="pattern-info">
                  <div class="pattern-name">
                    {pattern.Name}
                    {#if pattern.IsGlobal}
                      <span class="badge global">Global</span>
                    {/if}
                  </div>
                  <code class="pattern-regex">{pattern.Regex}</code>
                  <p class="pattern-desc">{pattern.Description || ''}</p>
                </div>
                {#if !pattern.IsGlobal}
                  <Button variant="ghost" icon="delete" onclick={() => deletePattern(pattern.ID)} />
                {/if}
              </div>
            {/each}
          </div>

          <div class="add-pattern-form tonal">
            <h5>Añadir Patrón Personalizado</h5>
            <div class="form-grid">
              <Input label="Nombre" bind:value={newPattern.name} placeholder="p.ej. My Custom API Key" />
              <Input label="Regex" bind:value={newPattern.regex} placeholder="p.ej. key-[0-9a-f]{32}" />
              <Input label="Descripción" bind:value={newPattern.description} placeholder="Breve explicación..." class="full-width" />
            </div>
            <div class="form-actions">
              <Button variant="primary" icon="add" onclick={addPattern} disabled={!newPattern.name || !newPattern.regex}>
                Añadir Patrón
              </Button>
            </div>
          </div>
        </section>
      </div>
    {/if}
  </div>
</div>

<style>
  .security-settings { display: flex; flex-direction: column; gap: 2rem; }
  .settings-card { background: var(--md-sys-color-surface-container); border-radius: 2rem; overflow: hidden; border: 1px solid var(--md-sys-color-outline-variant); }
  
  .card-header { padding: 2rem; display: flex; justify-content: space-between; align-items: center; border-bottom: 1px solid var(--md-sys-color-outline-variant); background: rgba(var(--primary-rgb), 0.03); }
  .header-content { display: flex; gap: 1.5rem; align-items: center; }
  .guardian-icon { font-size: 3rem; color: var(--primary-color); }
  .card-header h3 { margin: 0; font-size: 1.5rem; font-weight: 900; }
  .card-header p { margin: 0.25rem 0 0; color: var(--text-muted); }

  .card-body { padding: 2.5rem; display: flex; flex-direction: column; gap: 2.5rem; }

  .mode-selector { padding: 1.5rem; border-radius: 1.5rem; border: 1px solid var(--md-sys-color-outline-variant); }

  .section-header { margin-bottom: 1.5rem; }
  .section-header h4 { margin: 0; font-size: 1.25rem; font-weight: 800; }
  .section-header p { margin: 0.5rem 0 0; color: var(--text-muted); }

  .patterns-list { display: flex; flex-direction: column; gap: 1rem; margin-bottom: 2rem; }
  .pattern-item { display: flex; justify-content: space-between; align-items: flex-start; padding: 1.25rem; background: var(--md-sys-color-surface-container-high); border-radius: 1.25rem; border: 1px solid var(--md-sys-color-outline-variant); }
  .pattern-item.global { border-left: 4px solid var(--primary-color); }
  
  .pattern-info { display: flex; flex-direction: column; gap: 0.5rem; }
  .pattern-name { font-weight: 800; display: flex; align-items: center; gap: 0.75rem; }
  .pattern-regex { font-family: var(--font-mono); font-size: 0.85rem; color: var(--primary-color); background: rgba(var(--primary-rgb), 0.1); padding: 0.25rem 0.5rem; border-radius: 0.5rem; word-break: break-all; }
  .pattern-desc { margin: 0; font-size: 0.9rem; color: var(--text-muted); }

  .badge { font-size: 0.7rem; padding: 0.2rem 0.6rem; border-radius: 1rem; font-weight: 900; text-transform: uppercase; }
  .badge.global { background: var(--md-sys-color-primary-container); color: var(--md-sys-color-on-primary-container); }

  .add-pattern-form { padding: 2rem; border-radius: 1.5rem; border: 1px dotted var(--md-sys-color-outline); }
  .add-pattern-form h5 { margin: 0 0 1.5rem; font-size: 1.1rem; font-weight: 800; }
  .form-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }
  :global(.full-width) { grid-column: span 2; }
  .form-actions { display: flex; justify-content: flex-end; margin-top: 1.5rem; }

  @media (max-width: 600px) {
    .form-grid { grid-template-columns: 1fr; }
    :global(.full-width) { grid-column: span 1; }
  }
</style>
