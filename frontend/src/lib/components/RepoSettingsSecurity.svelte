<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from "svelte";
  import { toast } from "$lib/toast.svelte";
  import { dialog } from "$lib/dialog.svelte";
  import Checkbox from "./ui/Checkbox.svelte";
  import Button from "./ui/Button.svelte";
  import Input from "./ui/Input.svelte";
  import Card from "./ui/Card.svelte";
  import { i18n } from "$lib/i18n.svelte";

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
        fetchData();
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

<Card padding="none" class="settings-premium-card">
  <header class="section-header-premium">
    <div class="header-flex">
      <div class="title-with-icon">
        <span class="material-symbols-rounded hero-icon" class:active={data.enabled}>shield</span>
        <div>
          <h2>Secret Guardian</h2>
          <p class="subtitle">Protección proactiva contra la subida de secretos, API keys y archivos sensibles en tus commits.</p>
        </div>
      </div>
      <div class="header-action">
        <Checkbox 
          label="Habilitar Guardián" 
          bind:checked={data.enabled} 
          onchange={() => saveSettings()}
        />
      </div>
    </div>
  </header>

  <div class="settings-content">
    {#if loading}
      <div class="skeleton-form"></div>
    {:else if data.enabled}
      <div class="security-stack">
        <div class="mode-island island">
          <Checkbox 
            label="Modo Estricto (Bloquear Push)" 
            description="Si se detecta un secreto en los cambios entrantes, el push será rechazado inmediatamente." 
            bind:checked={data.strict}
            onchange={saveSettings}
          />
        </div>

        <div class="patterns-container">
          <div class="patterns-header">
            <span class="count-tag">{data.patterns.length} Patrones Activos</span>
          </div>

          <div class="patterns-list">
            {#each data.patterns as pattern}
              <div class="pattern-tile tonal" class:global={pattern.IsGlobal}>
                <div class="pattern-core">
                  <div class="pattern-identity">
                    <span class="pattern-name">{pattern.Name}</span>
                    {#if pattern.IsGlobal}
                      <span class="badge global">Sistema</span>
                    {:else}
                      <span class="badge custom">Personalizado</span>
                    {/if}
                  </div>
                  <code class="pattern-regex">{pattern.Regex}</code>
                  {#if pattern.Description}
                    <p class="pattern-desc">{pattern.Description}</p>
                  {/if}
                </div>
                {#if !pattern.IsGlobal}
                  <Button variant="ghost" size="sm" icon="delete" class="danger-btn" onclick={() => deletePattern(pattern.ID)} />
                {/if}
              </div>
            {/each}
          </div>

          <div class="add-pattern-box tonal">
            <header class="add-header">
              <span class="material-symbols-rounded">add_moderator</span>
              <strong>Añadir Patrón Personalizado</strong>
            </header>
            <div class="form-grid">
              <Input label="Nombre del Patrón" bind:value={newPattern.name} placeholder="p.ej. Mi Clave API Personalizada" />
              <Input label="Expresión Regular (Regex)" bind:value={newPattern.regex} placeholder="p.ej. key-[0-9a-f]{32}" />
              <Input label="Descripción" bind:value={newPattern.description} placeholder="¿Qué detecta este patrón?" class="full-width" />
            </div>
            <div class="form-actions">
              <Button variant="primary" icon="add" onclick={addPattern} disabled={!newPattern.name || !newPattern.regex}>
                Registrar Patrón
              </Button>
            </div>
          </div>
        </div>
      </div>
    {:else}
      <div class="empty-state-guardian island">
        <span class="material-symbols-rounded">verified_user</span>
        <h3>Guardián Desactivado</h3>
        <p>Activa el Guardián para proteger este repositorio de fugas accidentales de información sensible.</p>
        <Button variant="primary" onclick={() => { data.enabled = true; saveSettings(); }}>Activar Protección</Button>
      </div>
    {/if}
  </div>
</Card>

<style>
  .section-header-premium { padding: 2.5rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); }
  .header-flex { display: flex; justify-content: space-between; align-items: center; gap: 2rem; }
  .title-with-icon { display: flex; align-items: center; gap: 2rem; }
  .hero-icon { font-size: 3rem; color: var(--primary-color); background: rgba(var(--primary-rgb), 0.1); padding: 1.25rem; border-radius: 2rem; flex-shrink: 0; transition: all 0.3s; }
  .hero-icon.active { background: var(--primary-color); color: var(--on-primary); }
  
  h2 { margin: 0; font-size: 2rem; font-weight: 950; letter-spacing: -0.04em; }
  .subtitle { margin: 0.5rem 0 0 0; color: var(--text-muted); font-size: 1.1rem; }

  @media (max-width: 600px) {
    .section-header-premium { padding: 1.5rem; }
    .header-flex { flex-direction: column; align-items: stretch; gap: 1.5rem; }
    h2 { font-size: 1.25rem; }
    .hero-icon { font-size: 1.25rem; padding: 0.5rem; border-radius: 0.75rem; }
    .title-with-icon { gap: 0.75rem; }
    .subtitle { display: none; }
  }

  .settings-content { padding: 2.5rem; }
  @media (max-width: 600px) { .settings-content { padding: 1.5rem; } }

  .security-stack { display: flex; flex-direction: column; gap: 2.5rem; }
  
  .mode-island { padding: 2rem; border-radius: 2rem; border: 2px solid var(--primary-color); background: rgba(var(--primary-rgb), 0.05); }

  .patterns-container { display: flex; flex-direction: column; gap: 1.5rem; }
  .patterns-header { border-bottom: 2px solid var(--md-sys-color-outline-variant); padding-bottom: 0.75rem; }
  .count-tag { font-size: 0.75rem; font-weight: 950; text-transform: uppercase; color: var(--primary-color); letter-spacing: 0.1em; }

  .patterns-list { display: grid; grid-template-columns: repeat(auto-fill, minmax(400px, 1fr)); gap: 1rem; }
  @media (max-width: 850px) { .patterns-list { grid-template-columns: 1fr; } }

  .pattern-tile { display: flex; justify-content: space-between; align-items: flex-start; padding: 1.5rem; border-radius: 1.5rem; border: 1px solid var(--md-sys-color-outline-variant); transition: all 0.2s; }
  .pattern-tile:hover { transform: translateY(-2px); border-color: var(--primary-color); }
  .pattern-tile.global { border-left: 6px solid var(--primary-color); }
  
  .pattern-core { display: flex; flex-direction: column; gap: 0.75rem; flex: 1; min-width: 0; }
  .pattern-identity { display: flex; align-items: center; gap: 0.75rem; }
  .pattern-name { font-weight: 900; font-size: 1.1rem; color: var(--text-color); }
  
  .pattern-regex { font-family: 'JetBrains Mono', monospace; font-size: 0.85rem; color: var(--primary-color); background: rgba(var(--primary-rgb), 0.08); padding: 0.4rem 0.75rem; border-radius: 0.75rem; word-break: break-all; width: fit-content; }
  .pattern-desc { margin: 0; font-size: 0.9rem; color: var(--text-muted); opacity: 0.8; line-height: 1.4; }

  .badge { font-size: 0.65rem; font-weight: 950; text-transform: uppercase; padding: 0.2rem 0.6rem; border-radius: 100px; }
  .badge.global { background: var(--primary-color); color: var(--on-primary); }
  .badge.custom { background: var(--md-sys-color-surface-container-highest); color: var(--text-muted); }

  :global(.danger-btn:hover) { color: var(--error) !important; background: rgba(var(--error-rgb), 0.1) !important; }

  .add-pattern-box { padding: 2.5rem; border-radius: 2rem; border: 1px dashed var(--md-sys-color-outline); display: flex; flex-direction: column; gap: 2rem; }
  .add-header { display: flex; align-items: center; gap: 0.75rem; color: var(--primary-color); }
  .add-header span { font-size: 1.75rem; }
  .add-header strong { font-weight: 950; font-size: 0.85rem; text-transform: uppercase; letter-spacing: 0.1em; }

  .form-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 1.5rem; }
  :global(.full-width) { grid-column: span 2; }
  .form-actions { display: flex; justify-content: flex-end; }

  .empty-state-guardian { padding: 5rem 2rem; text-align: center; display: flex; flex-direction: column; align-items: center; gap: 1.5rem; }
  .empty-state-guardian span { font-size: 5rem; color: var(--primary-color); opacity: 0.2; }
  .empty-state-guardian h3 { margin: 0; font-size: 1.75rem; font-weight: 950; }
  .empty-state-guardian p { max-width: 500px; margin: 0; color: var(--text-muted); font-size: 1.1rem; line-height: 1.5; }

  .skeleton-form { height: 400px; background: var(--md-sys-color-surface-variant); border-radius: 2rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 50% { opacity: 0.5; } }

  @media (max-width: 600px) {
    .form-grid { grid-template-columns: 1fr; }
    :global(.full-width) { grid-column: span 1; }
    .pattern-tile { flex-direction: column; gap: 1rem; }
    .pattern-tile :global(button) { align-self: flex-end; }
  }
</style>
