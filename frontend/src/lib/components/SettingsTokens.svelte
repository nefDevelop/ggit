<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  import { i18n } from '$lib/i18n.svelte';
  import Card from './ui/Card.svelte';
  import Button from './ui/Button.svelte';

  let tokens = $state<any[]>([]);
  let loading = $state(true);
  let newTokenName = $state('');
  let generatedToken = $state('');
  let submitting = $state(false);
  let error = $state("");

  async function fetchTokens() {
    loading = true;
    try {
      const res = await apiFetch('/api/v2/user/settings/tokens');
      if (res.ok) tokens = await res.json();
    } catch (e) { console.error(e); } finally { loading = false; }
  }

  async function createToken() {
    if (!newTokenName) return;
    submitting = true;
    try {
      const res = await apiFetch('/api/v2/user/settings/tokens', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ name: newTokenName })
      });
      if (res.ok) {
        const token = await res.json();
        generatedToken = token.sha1; 
        tokens = [token, ...tokens];
        newTokenName = '';
        error = '';
      } else {
        const data = await res.json();
        error = data.message || 'Error al generar el token';
      }
    } catch (e) {
      error = "Error de conexión al servidor";
      console.error(e);
    } finally {
      submitting = false;
    }
  }

  async function deleteToken(id: number) {
    if (!confirm('¿Seguro que quieres eliminar este token?')) return;
    try {
      const res = await apiFetch(`/api/v2/user/settings/tokens/${id}`, { method: 'DELETE' });
      if (res.ok) tokens = tokens.filter(t => t.id !== id);
    } catch (e) { console.error(e); }
  }

  onMount(fetchTokens);
</script>

<Card class="settings-premium-card" padding="lg">
  <header class="section-header">
    <div class="title-with-icon">
      <span class="material-symbols-rounded hero-icon">token</span>
      <div>
        <h2>Tokens de Acceso Personal</h2>
        <p class="subtitle">Claves seguras para autenticarte vía API o Git sin usar tu contraseña</p>
      </div>
    </div>
  </header>

  <div class="modular-stack">
    {#if generatedToken}
      <div class="generated-alert island">
        <div class="alert-top">
           <span class="material-symbols-rounded">security_update_good</span>
           <div class="alert-text">
              <strong>¡Token Generado con Éxito!</strong>
              <p>Copia este token ahora. Por seguridad, no podrás volver a verlo.</p>
           </div>
        </div>
        <div class="token-display">
           <code>{generatedToken}</code>
           <Button variant="ghost" icon="content_copy" onclick={() => {
             navigator.clipboard.writeText(generatedToken);
             alert('¡Token copiado al portapapeles!');
           }} />
        </div>
        <Button variant="primary" onclick={() => generatedToken = ''}>He guardado el token</Button>
      </div>
    {/if}

    <!-- TOKENS LIST -->
    <div class="config-module">
      <div class="module-tag">
        <span class="material-symbols-rounded">format_list_bulleted</span>
        Mis Tokens Activos
      </div>
      <div class="island tonal-list">
        {#if loading}
          <div class="skeleton-tokens"></div>
        {:else if tokens.length === 0}
          <div class="empty-tokens-box">
             <span class="material-symbols-rounded">lock_open</span>
             <p>No has generado ningún token todavía.</p>
          </div>
        {:else}
          {#each tokens as token}
            <div class="token-row">
              <div class="token-main">
                <div class="token-icon-wrap tonal">
                   <span class="material-symbols-rounded">api</span>
                </div>
                <div class="token-info">
                   <span class="token-name">{token.name}</span>
                   <span class="token-date">Creado el {new Date(token.created_at).toLocaleDateString()}</span>
                </div>
              </div>
              <div class="token-actions">
                <Button variant="ghost" size="sm" icon="delete" class="danger-btn" onclick={() => deleteToken(token.id)} />
              </div>
            </div>
          {/each}
        {/if}
      </div>
    </div>

    <!-- GENERATE FORM -->
    <div class="config-module">
      <div class="module-tag">
        <span class="material-symbols-rounded">add_circle</span>
        Generar Nuevo Token
      </div>
      <div class="island form-box">
        <form onsubmit={(e) => { e.preventDefault(); createToken(); }} class="premium-token-form">
          <div class="inline-input-layout">
             <div class="field">
                <label>Nombre del Propósito</label>
                <input type="text" bind:value={newTokenName} placeholder="Ej: CI/CD Pipeline o App Móvil" required />
             </div>
             <Button type="submit" variant="primary" icon="key" disabled={submitting || !newTokenName}>
                {submitting ? 'Generando...' : 'Generar Token'}
             </Button>
          </div>
          {#if error}
             <p class="error-msg">{error}</p>
          {/if}
        </form>
      </div>
    </div>
  </div>
</Card>

<style>
  
  .section-header { margin-bottom: 4rem; }
  .header-flex { display: flex; justify-content: space-between; align-items: center; gap: 2rem; }
  @media (max-width: 600px) {
    .section-header { margin-bottom: 2rem; }
    h2 { font-size: 1.5rem; }
    .hero-icon { font-size: 1.5rem; padding: 0.6rem; border-radius: 1rem; }
    .title-with-icon { gap: 0.75rem; }
    .subtitle { font-size: 0.9rem; }
  }

  .title-with-icon { display: flex; align-items: center; gap: 2rem; }
  .hero-icon { font-size: 3rem; color: var(--primary-color); background: rgba(var(--primary-rgb), 0.1); padding: 1.25rem; border-radius: 2rem; flex-shrink: 0; }
  h2 { margin: 0; font-size: 2.25rem; font-weight: 950; letter-spacing: -0.04em; }
  .subtitle { margin: 0.5rem 0 0 0; color: var(--text-muted); font-size: 1.1rem; }
  .config-module { display: flex; flex-direction: column; gap: 1.5rem; }

  .module-tag { font-size: 0.8rem; font-weight: 950; text-transform: uppercase; color: var(--primary-color); letter-spacing: 0.15em; display: flex; align-items: center; gap: 0.75rem; border-bottom: 2px solid var(--md-sys-color-outline-variant); padding-bottom: 0.75rem; width: fit-content; opacity: 0.8; }
  .module-tag span { font-size: 1.25rem; }

  .generated-alert { padding: 3rem; border: 2px solid var(--primary-color); background: rgba(var(--primary-rgb), 0.05); display: flex; flex-direction: column; gap: 2rem; animation: slideDown 0.4s; }
  @keyframes slideDown { from { opacity: 0; transform: translateY(-20px); } to { opacity: 1; transform: translateY(0); } }
  .alert-top { display: flex; align-items: center; gap: 1.5rem; }
  .alert-top span { font-size: 3rem; color: var(--primary-color); }
  .alert-text strong { display: block; font-size: 1.25rem; font-weight: 950; }
  .alert-text p { margin: 0.25rem 0 0; color: var(--text-muted); }

  .token-display { display: flex; align-items: center; gap: 1.5rem; background: #08080c; border: 1px solid #1a1a24; padding: 1rem 1.5rem; border-radius: 1.5rem; }
  .token-display code { flex: 1; font-family: 'JetBrains Mono', monospace; font-size: 1.4rem; color: #6ee7b7; word-break: break-all; }

  .tonal-list { display: flex; flex-direction: column; background: rgba(0,0,0,0.15); border-radius: 2.25rem; overflow: hidden; border: 1px solid var(--md-sys-color-outline-variant); }
  .token-row { display: flex; align-items: center; justify-content: space-between; padding: 2rem 2.5rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); gap: 3rem; transition: background 0.3s; }
  
  @media (max-width: 800px) {
    .token-row { flex-direction: column; align-items: stretch; padding: 1.5rem; gap: 1.5rem; }
    .token-actions { display: flex; justify-content: flex-end; }
  }

  .token-row:last-child { border-bottom: none; }
  .token-row:hover { background: rgba(255,255,255,0.02); }

  .token-main { display: flex; align-items: center; gap: 1.5rem; min-width: 0; }
  .token-icon-wrap { width: 56px; height: 56px; border-radius: 1.25rem; display: flex; align-items: center; justify-content: center; background: rgba(var(--primary-rgb), 0.1); color: var(--primary-color); flex-shrink: 0; }
  .token-icon-wrap span { font-size: 2rem; }

  .token-info { display: flex; flex-direction: column; gap: 0.4rem; min-width: 0; flex: 1; }
  .token-name { font-size: 1.25rem; font-weight: 900; color: var(--text-color); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .token-date { font-size: 0.8rem; color: var(--text-muted); opacity: 0.6; font-weight: 700; }

  :global(.danger-btn:hover) { color: var(--error) !important; background: rgba(var(--error-rgb), 0.1) !important; }
  .empty-tokens-box { padding: 4rem; text-align: center; color: var(--text-muted); display: flex; flex-direction: column; align-items: center; gap: 1rem; }
  .empty-tokens-box span { font-size: 3rem; opacity: 0.2; }

  /* FORM */
  .form-box { padding: 3rem; background: rgba(255,255,255,0.01); border-radius: 2.5rem; }
  .inline-input-layout { display: flex; align-items: flex-end; gap: 1.5rem; }
  
  @media (max-width: 600px) {
    .form-box { padding: 1.5rem; }
    .inline-input-layout { flex-direction: column; align-items: stretch; }
  }
  .field { flex: 1; display: flex; flex-direction: column; gap: 0.75rem; }
  .field label { font-size: 0.85rem; font-weight: 800; color: var(--text-color); margin-left: 0.5rem; }
  .field input { background: var(--surface-variant); border: 2px solid transparent; padding: 0.85rem 1.5rem; border-radius: 1.25rem; color: var(--text-color); font-size: 1.1rem; transition: all 0.3s; width: 100%; height: 3.5rem; }
  .field input:focus { border-color: var(--primary-color); background: rgba(var(--primary-rgb), 0.05); }

  .skeleton-tokens { height: 160px; background: var(--surface-variant); border-radius: 2rem; animation: pulse 1.5s infinite; }

  .error-msg { color: var(--error); background: rgba(var(--error-rgb), 0.1); padding: 1rem 2rem; border-radius: 1.25rem; font-size: 0.9rem; font-weight: 700; margin: 1.5rem 0 0; border: 1px solid rgba(var(--error-rgb), 0.2); }
</style>
