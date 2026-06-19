<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  import { i18n } from '$lib/i18n.svelte';
  import { dialog } from '$lib/dialog.svelte';
  import Card from './ui/Card.svelte';
  import Button from './ui/Button.svelte';
  
import { ShieldCheck, Copy, List, LockOpen, Cpu, Trash2, PlusCircle, Key } from "lucide-svelte";


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
    if (!await dialog.confirm('¿Seguro que quieres eliminar este token?')) return;
    try {
      const res = await apiFetch(`/api/v2/user/settings/tokens/${id}`, { method: 'DELETE' });
      if (res.ok) tokens = tokens.filter(t => t.id !== id);
    } catch (e) { console.error(e); }
  }

  function formatDate(dateStr: string) {
    if (!dateStr) return 'N/A';
    const d = new Date(dateStr);
    if (isNaN(d.getTime())) return 'N/A';
    return d.toLocaleDateString(undefined, {
      year: 'numeric',
      month: 'long',
      day: 'numeric'
    });
  }

  onMount(fetchTokens);
</script>

<div class="settings-tokens-container">

  <div class="flex-col gap-xl">
    {#if generatedToken}
      <div class="generated-alert break-words">
        <div class="alert-top">
           <ShieldCheck size={48} />
           <div class="alert-text">
              <strong>¡Token Generado con Éxito!</strong>
              <p>Copia este token ahora. Por seguridad, no podrás volver a verlo.</p>
           </div>
        </div>
        <div class="token-display">
           <code>{generatedToken}</code>
           <Button variant="ghost" icon={Copy} onclick={() => {
             navigator.clipboard.writeText(generatedToken);
             alert('¡Token copiado al portapapeles!');
           }} />
        </div>
        <Button variant="primary" onclick={() => generatedToken = ''}>He guardado el token</Button>
      </div>
    {/if}

    <!-- TOKENS LIST -->
    <div class="flex-col gap-md">
      <div class="module-tag">
        <List size={18} />
        Mis Tokens Activos
      </div>
      <div class="tonal-list">
        {#if loading}
          <div class="skeleton-tokens"></div>
        {:else if tokens.length === 0}
          <div class="empty-state-box">
             <LockOpen size={48} />
             <p>No has generado ningún token todavía.</p>
          </div>
        {:else}
          {#each tokens as token}
            <div class="token-row">
              <div class="token-main">
                <div class="token-icon-wrap tonal">
                   <Cpu size={28} />
                </div>
                <div class="token-info">
                   <span class="token-name">{token.name}</span>
                   <span class="token-date">Creado el {formatDate(token.created_at)}</span>
                </div>
              </div>
              <div class="token-actions">
                <Button variant="ghost" size="sm" icon={Trash2} class="danger-btn" onclick={() => deleteToken(token.id)} />
              </div>
            </div>
          {/each}
        {/if}
      </div>
    </div>

    <!-- GENERATE FORM -->
    <div class="flex-col gap-md">
      <div class="module-tag">
        <PlusCircle size={18} />
        Generar Nuevo Token
      </div>
      <div class="form-box">
        <form onsubmit={(e) => { e.preventDefault(); createToken(); }} class="premium-token-form">
          <div class="inline-input-layout">
             <div class="field">
                <label for="newTokenName">Nombre del Propósito</label>
                <input id="newTokenName" type="text" bind:value={newTokenName} placeholder="Ej: CI/CD Pipeline o App Móvil" required />
             </div>
             <div class="button-wrap">
               <Button type="submit" variant="primary" icon={Key} disabled={submitting || !newTokenName} class="w-full">
                  {submitting ? 'Generando...' : 'Generar Token'}
               </Button>
             </div>
          </div>
          {#if error}
             <p class="error-msg">{error}</p>
          {/if}
        </form>
      </div>
    </div>
  </div>
</div>

<style>
  


  .module-tag { font-size: 0.8rem; font-weight: 950; text-transform: uppercase; color: var(--primary-color); letter-spacing: 0.15em; display: flex; align-items: center; gap: 0.75rem; border-bottom: 2px solid var(--md-sys-color-outline-variant); padding-bottom: 0.75rem; width: fit-content; opacity: 0.8; }

  .generated-alert { padding: 3rem; border: 2px solid var(--primary-color); background: rgba(var(--primary-rgb), 0.05); border-radius: 1.25rem; display: flex; flex-direction: column; gap: 2rem; animation: slideDown 0.4s; }
  @media (max-width: 600px) { .generated-alert { padding: 1.5rem; gap: 1.5rem; } }
  @keyframes slideDown { from { opacity: 0; transform: translateY(-20px); } to { opacity: 1; transform: translateY(0); } }
  .alert-top { display: flex; align-items: center; gap: 1.5rem; }
  .alert-text strong { display: block; font-size: 1.25rem; font-weight: 950; }
  .alert-text p { margin: 0.25rem 0 0; color: var(--text-muted); }

  .token-display { display: flex; align-items: center; gap: 1.5rem; background: #08080c; border: 1px solid #1a1a24; padding: 1rem 1.5rem; border-radius: 0.75rem; }
  @media (max-width: 600px) { .token-display { padding: 0.75rem 1rem; border-radius: 1rem; gap: 0.75rem; } }
  .token-display code { flex: 1; font-family: 'JetBrains Mono', monospace; font-size: 1.4rem; color: #6ee7b7; word-break: break-all; }
  @media (max-width: 600px) { .token-display code { font-size: 1rem; } }

  .tonal-list { display: flex; flex-direction: column; background: rgba(0,0,0,0.15); border-radius: 1.25rem; overflow: hidden; border: 1px solid var(--md-sys-color-outline-variant); }
  @media (max-width: 600px) { .tonal-list { border-radius: 1.25rem; } }
  .token-row { display: flex; align-items: center; justify-content: space-between; padding: 2rem 2.5rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); gap: 3rem; transition: background 0.3s; }
  
  @media (max-width: 800px) {
    .token-row { flex-direction: column; align-items: stretch; padding: 1.5rem; gap: 1.5rem; }
    .token-actions { display: flex; justify-content: flex-end; }
  }

  @media (max-width: 600px) {
    .token-row { padding: 1rem 1.25rem; gap: 1rem; }
  }

  .token-row:last-child { border-bottom: none; }
  .token-row:hover { background: rgba(255,255,255,0.02); }

  .token-main { display: flex; align-items: center; gap: 1.5rem; min-width: 0; }
  .token-icon-wrap { width: 56px; height: 56px; border-radius: 0.75rem; display: flex; align-items: center; justify-content: center; background: rgba(var(--primary-rgb), 0.1); color: var(--primary-color); flex-shrink: 0; }

  .token-info { display: flex; flex-direction: column; gap: 0.4rem; min-width: 0; flex: 1; }
  .token-name { font-size: 1.25rem; font-weight: 900; color: var(--text-color); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  .token-date { font-size: 0.8rem; color: var(--text-muted); opacity: 0.6; font-weight: 700; }

  :global(.danger-btn:hover) { color: var(--error) !important; background: rgba(var(--error-rgb), 0.1) !important; }

  /* FORM */
  .form-box { padding: 3rem; background: rgba(255,255,255,0.01); border-radius: 1.25rem; }
  .inline-input-layout { display: flex; align-items: flex-end; gap: 1.5rem; width: 100%; }
  
  .field { flex: 1; display: flex; flex-direction: column; gap: 0.75rem; min-width: 0; }
  .button-wrap { flex-shrink: 0; }

  @media (max-width: 768px) {
    .form-box { padding: 1.5rem; }
    .inline-input-layout { flex-direction: column; align-items: stretch; }
    .button-wrap { width: 100%; }
  }
  
  .field label { font-size: 0.85rem; font-weight: 800; color: var(--text-color); margin-left: 0.5rem; }
  .field input { background: var(--surface-variant); border: 2px solid transparent; padding: 0.85rem 1.5rem; border-radius: 0.75rem; color: var(--text-color); font-size: 1.1rem; transition: all 0.3s; width: 100%; height: 3.5rem; box-sizing: border-box; }
  .field input:focus { border-color: var(--primary-color); background: rgba(var(--primary-rgb), 0.05); }

  .skeleton-tokens { height: 160px; background: var(--surface-variant); border-radius: 2rem; animation: pulse 1.5s infinite; }

  .error-msg { color: var(--error); background: rgba(var(--error-rgb), 0.1); padding: 1rem 2rem; border-radius: 1.25rem; font-size: 0.9rem; font-weight: 700; margin: 1.5rem 0 0; border: 1px solid rgba(var(--error-rgb), 0.2); }

  .empty-state-box { display: flex; flex-direction: column; align-items: center; justify-content: center; padding: 4rem 2rem; gap: 1.5rem; color: var(--text-muted); text-align: center; }
  .empty-state-box p { font-size: 1.1rem; font-weight: 600; opacity: 0.7; }
</style>
