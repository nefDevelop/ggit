<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from "svelte";
  import { i18n } from '$lib/i18n.svelte';
  import Card from './ui/Card.svelte';
  import Button from './ui/Button.svelte';

  let keys = $state<any[]>([]);
  let loading = $state(true);
  let newKeyTitle = $state("");
  let newKeyContent = $state("");
  let submitting = $state(false);
  let error = $state("");

  async function fetchKeys() {
    loading = true;
    try {
      const res = await apiFetch("/api/v2/user/settings/ssh_keys");
      if (res.ok) keys = await res.json();
    } catch (e) { console.error(e); } finally { loading = false; }
  }

  async function addKey() {
    if (!newKeyTitle || !newKeyContent) return;
    submitting = true;
    try {
      const res = await apiFetch("/api/v2/user/settings/ssh_keys", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ title: newKeyTitle, content: newKeyContent }),
      });
      if (res.ok) {
        newKeyTitle = ""; newKeyContent = "";
        error = "";
        await fetchKeys();
      } else {
        const data = await res.json();
        error = data.message || 'Error al añadir la clave SSH';
      }
    } catch (e) { 
      error = "Error de conexión al servidor";
      console.error(e); 
    } finally { submitting = false; }
  }

  async function deleteKey(id: number) {
    if (!confirm("¿Seguro que quieres eliminar esta clave SSH?")) return;
    try {
      const res = await apiFetch(`/api/v2/user/settings/ssh_keys/${id}`, { method: "DELETE" });
      if (res.ok) fetchKeys();
    } catch (e) { console.error(e); }
  }

  onMount(fetchKeys);
</script>

<Card class="settings-premium-card" padding="lg">
  <header class="section-header">
    <div class="title-with-icon">
      <span class="material-symbols-rounded hero-icon">key</span>
      <div>
        <h2>Claves SSH</h2>
        <p class="subtitle">Gestiona las llaves públicas vinculadas a tu cuenta para acceso seguro vía Git</p>
      </div>
    </div>
  </header>

  <div class="modular-stack">
    <!-- KEYS LIST -->
    <div class="config-module">
      <div class="module-tag">
        <span class="material-symbols-rounded">vpn_key</span>
        Llaves Activas
      </div>
      <div class="island tonal-list">
        {#if loading}
           <div class="skeleton-keys"></div>
        {:else if keys.length === 0}
           <div class="empty-keys-box">
              <span class="material-symbols-rounded">no_encryption</span>
              <p>No hay claves SSH configuradas.</p>
           </div>
        {:else}
          {#each keys as key}
            <div class="key-row">
              <div class="key-main">
                <div class="key-type-icon tonal">
                   <span class="material-symbols-rounded">fingerprint</span>
                </div>
                <div class="key-info">
                   <span class="key-title">{key.title}</span>
                   <code class="key-fingerprint">{key.fingerprint || 'Huella no disponible'}</code>
                   <span class="key-date">Añadida el {new Date(key.created_at).toLocaleDateString()}</span>
                </div>
              </div>
              <div class="key-actions">
                <Button variant="ghost" size="sm" icon="delete" class="danger-btn" onclick={() => deleteKey(key.id)} />
              </div>
            </div>
          {/each}
        {/if}
      </div>
    </div>

    <!-- ADD KEY FORM -->
    <div class="config-module">
      <div class="module-tag">
        <span class="material-symbols-rounded">add_link</span>
        Añadir Nueva Clave
      </div>
      <div class="island form-box">
        <form onsubmit={(e) => { e.preventDefault(); addKey(); }} class="premium-key-form">
           <div class="field">
              <label>Título de la Clave</label>
              <input type="text" bind:value={newKeyTitle} placeholder="Ej: Portátil Trabajo" required />
           </div>
           <div class="field">
              <label>Contenido (Key Content)</label>
              <textarea 
                bind:value={newKeyContent} 
                placeholder="Comienza con 'ssh-rsa', 'ssh-ed25519'..." 
                required
              ></textarea>
           </div>
           {#if error}
              <p class="error-msg">{error}</p>
           {/if}
           <div class="form-footer">
              <Button type="submit" variant="primary" icon="key" disabled={submitting || !newKeyContent}>
                 {submitting ? 'Añadiendo...' : 'Vincular Clave SSH'}
              </Button>
           </div>
        </form>
      </div>
    </div>
  </div>
</Card>

<style>
  
  .section-header { margin-bottom: 4rem; }
  .title-with-icon { display: flex; align-items: center; gap: 2rem; }
  .hero-icon { font-size: 3rem; color: var(--primary-color); background: rgba(var(--primary-rgb), 0.1); padding: 1.25rem; border-radius: 2rem; }
  h2 { margin: 0; font-size: 2.25rem; font-weight: 950; letter-spacing: -0.04em; }
  .subtitle { margin: 0.5rem 0 0 0; color: var(--text-muted); font-size: 1.1rem; }

  @media (max-width: 600px) {
    .section-header { margin-bottom: 2rem; }
    h2 { font-size: 1.5rem; }
    .hero-icon { font-size: 1.5rem; padding: 0.6rem; border-radius: 1rem; }
    .title-with-icon { gap: 0.75rem; }
    .subtitle { font-size: 0.9rem; }
  }

  .modular-stack { display: flex; flex-direction: column; gap: 4rem; }
  .config-module { display: flex; flex-direction: column; gap: 1.5rem; }

  .module-tag { font-size: 0.8rem; font-weight: 950; text-transform: uppercase; color: var(--primary-color); letter-spacing: 0.15em; display: flex; align-items: center; gap: 0.75rem; border-bottom: 2px solid var(--md-sys-color-outline-variant); padding-bottom: 0.75rem; width: fit-content; opacity: 0.8; }
  .module-tag span { font-size: 1.25rem; }

  .tonal-list { display: flex; flex-direction: column; background: rgba(0,0,0,0.15); border-radius: 2.25rem; overflow: hidden; border: 1px solid var(--md-sys-color-outline-variant); }
  .key-row { display: flex; align-items: center; justify-content: space-between; padding: 2rem 2.5rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); gap: 3rem; transition: background 0.3s; }
  
  @media (max-width: 800px) {
    .key-row { flex-direction: column; align-items: stretch; padding: 1.5rem; gap: 1.5rem; }
    .key-actions { display: flex; justify-content: flex-end; }
  }

  .key-row:last-child { border-bottom: none; }
  .key-row:hover { background: rgba(255,255,255,0.02); }

  .key-main { display: flex; align-items: center; gap: 1.5rem; min-width: 0; }
  .key-type-icon { width: 56px; height: 56px; border-radius: 1.25rem; display: flex; align-items: center; justify-content: center; background: rgba(var(--primary-rgb), 0.1); color: var(--primary-color); flex-shrink: 0; }
  .key-type-icon span { font-size: 2rem; }

  .key-info { display: flex; flex-direction: column; gap: 0.4rem; overflow: hidden; flex: 1; }
  .key-title { font-size: 1.25rem; font-weight: 900; color: var(--text-color); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .key-fingerprint { font-family: 'JetBrains Mono', monospace; font-size: 0.8rem; color: #888; background: #111; padding: 0.3rem 0.75rem; border-radius: 8px; border: 1px solid #222; overflow: hidden; text-overflow: ellipsis; white-space: nowrap; max-width: 100%; }
  .key-date { font-size: 0.8rem; color: var(--text-muted); opacity: 0.6; font-weight: 700; }

  /* ADD FORM */
  .form-box { padding: 3rem; background: rgba(255,255,255,0.01); }
  @media (max-width: 600px) {
    .form-box { padding: 1.5rem; }
    .form-footer button { width: 100%; }
  }
  .premium-key-form { display: flex; flex-direction: column; gap: 2rem; }
  .field { display: flex; flex-direction: column; gap: 0.75rem; }
  .field label { font-size: 0.85rem; font-weight: 800; color: var(--text-color); margin-left: 0.5rem; }
  .field input, .field textarea { background: var(--surface-variant); border: 2px solid transparent; padding: 1rem 1.5rem; border-radius: 1.25rem; color: var(--text-color); font-size: 1.1rem; transition: all 0.3s; }
  .field textarea { font-family: 'JetBrains Mono', monospace; min-height: 160px; font-size: 0.95rem; line-height: 1.6; resize: vertical; }
  .field input:focus, .field textarea:focus { border-color: var(--primary-color); background: rgba(var(--primary-rgb), 0.05); }

  .form-footer { display: flex; justify-content: flex-end; padding-top: 1.5rem; border-top: 1px solid var(--md-sys-color-outline-variant); }

  .skeleton-keys { height: 200px; background: var(--surface-variant); border-radius: 2rem; animation: pulse 1.5s infinite; }

  .error-msg { color: var(--error); background: rgba(var(--error-rgb), 0.1); padding: 1rem 1.5rem; border-radius: 1rem; font-size: 0.9rem; font-weight: 700; margin-bottom: 0.5rem; border: 1px solid rgba(var(--error-rgb), 0.2); }
</style>
