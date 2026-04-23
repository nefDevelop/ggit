<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from "svelte";
  import { i18n } from '$lib/i18n.svelte';
  import Card from './ui/Card.svelte';
  import Button from './ui/Button.svelte';

  let sources = $state<any[]>([]);
  let loading = $state(true);

  async function fetchAuths() {
    loading = true;
    try {
      const res = await apiFetch("/api/v2/admin/auths");
      if (res.ok) sources = await res.json();
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function deleteAuth(id: number) {
    if (!confirm("Are you sure you want to delete this authentication source?"))
      return;
    try {
      const res = await apiFetch(`/api/v2/admin/auths/${id}`, {
        method: "DELETE",
      });
      if (res.ok) fetchAuths();
    } catch (e) { console.error(e); }
  }

  onMount(fetchAuths);
</script>

<Card class="admin-view-card" padding="lg">
  <header class="section-header">
    <div class="header-flex">
      <div class="title-with-icon">
        <span class="material-symbols-rounded hero-icon">lock</span>
        <div>
          <h2>Orígenes de Autenticación</h2>
          <p class="subtitle">Gestiona métodos externos como LDAP, OAuth2 o PAM</p>
        </div>
      </div>
      <Button variant="primary" icon="add_moderator" class="hero-btn">Nueva Fuente</Button>
    </div>
  </header>

  {#if loading}
    <div class="skeleton-grid"></div>
  {:else if sources.length === 0}
    <div class="empty-state-card island">
      <span class="material-symbols-rounded">no_encryption</span>
      <h3>Sin fuentes de autenticación</h3>
      <p>Configura LDAP u otros servicios para permitir el acceso externo.</p>
      <Button variant="outline" icon="add">Añadir Primera Fuente</Button>
    </div>
  {:else}
    <div class="auth-grid">
      {#each sources as source}
        <div class="auth-source-card island tonal-hover">
          <div class="source-top">
            <div class="type-capsule" style="--accent: {source.type === 2 || source.type === 5 ? 'var(--primary-color)' : '#9c27b0'}">
              <span class="material-symbols-rounded">{source.type === 2 || source.type === 5 ? "hub" : "security"}</span>
            </div>
            <div class="source-meta">
              <div class="name-row">
                <span class="name">{source.name}</span>
                <span class="status-dot" class:active={source.is_actived}></span>
              </div>
              <span class="type-label">{source.type_name || "Protocolo Ggit"}</span>
            </div>
          </div>
          
          <div class="source-footer">
            <div class="active-badge" class:active={source.is_actived}>
              {source.is_actived ? 'Activo' : 'Deshabilitado'}
            </div>
            <div class="actions">
              <Button variant="ghost" size="sm" icon="settings" />
              <Button variant="ghost" size="sm" icon="delete" class="danger-btn" onclick={() => deleteAuth(source.id)} />
            </div>
          </div>
        </div>
      {/each}
    </div>
  {/if}
</Card>

<style>
  .admin-view-card { border-radius: 3rem; background: rgba(var(--surface-rgb), 0.5); backdrop-filter: blur(20px); border: 1px solid var(--md-sys-color-outline-variant); }
  .section-header { margin-bottom: 3.5rem; }
  .header-flex { display: flex; justify-content: space-between; align-items: center; gap: 2rem; }
  .title-with-icon { display: flex; align-items: center; gap: 2rem; }
  .hero-icon { font-size: 3rem; color: var(--primary-color); background: rgba(var(--primary-rgb), 0.1); padding: 1.25rem; border-radius: 2rem; }
  h2 { margin: 0; font-size: 2.25rem; font-weight: 950; letter-spacing: -0.04em; }
  .subtitle { margin: 0.5rem 0 0 0; color: var(--text-muted); font-size: 1.1rem; }

  .hero-btn { border-radius: 1.5rem !important; padding: 1rem 2rem !important; font-weight: 800 !important; }

  .auth-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(340px, 1fr)); gap: 1.25rem; }
  .auth-source-card { padding: 1.75rem; border-radius: 2.25rem; border: 1px solid var(--md-sys-color-outline-variant); display: flex; flex-direction: column; gap: 2rem; transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1); background: rgba(255,255,255,0.01); }
  .auth-source-card:hover { transform: translateY(-6px); border-color: var(--primary-color); background: rgba(var(--primary-rgb), 0.05); }

  .source-top { display: flex; align-items: center; gap: 1.25rem; }
  .type-capsule { width: 56px; height: 56px; border-radius: 1.25rem; background: rgba(var(--primary-rgb), 0.1); color: var(--accent); display: flex; align-items: center; justify-content: center; }
  .type-capsule span { font-size: 2.25rem; }

  .source-meta { display: flex; flex-direction: column; gap: 0.25rem; flex: 1; }
  .name-row { display: flex; align-items: center; gap: 0.75rem; justify-content: space-between; }
  .name { font-size: 1.25rem; font-weight: 900; color: var(--text-color); }
  .status-dot { width: 10px; height: 10px; border-radius: 50%; background: #444; }
  .status-dot.active { background: #10b981; box-shadow: 0 0 10px #10b98180; }
  .type-label { font-size: 0.75rem; text-transform: uppercase; font-weight: 800; color: var(--text-muted); letter-spacing: 0.1em; }

  .source-footer { display: flex; align-items: center; justify-content: space-between; padding-top: 1.5rem; border-top: 1px solid var(--md-sys-color-outline-variant); }
  .active-badge { font-size: 0.7rem; font-weight: 900; text-transform: uppercase; color: var(--text-muted); background: rgba(0,0,0,0.2); padding: 0.3rem 0.75rem; border-radius: 100px; }
  .active-badge.active { color: #10b981; background: rgba(16, 185, 129, 0.1); }

  .actions { display: flex; gap: 0.5rem; }
  :global(.danger-btn:hover) { color: var(--error) !important; background: rgba(var(--error-rgb), 0.1) !important; }

  .empty-state-card { padding: 5rem; text-align: center; border-radius: 3rem; display: flex; flex-direction: column; align-items: center; gap: 1.5rem; border: 2px dashed var(--md-sys-color-outline-variant); }
  .empty-state-card span { font-size: 4rem; opacity: 0.2; }
  .empty-state-card h3 { margin: 0; font-size: 1.5rem; font-weight: 900; }
  .empty-state-card p { margin: 0; color: var(--text-muted); font-size: 1.1rem; }

  .skeleton-grid { height: 400px; background: var(--surface-variant); border-radius: 2rem; animation: pulse 1.5s infinite; }
</style>
