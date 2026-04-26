<script lang="ts">
  import { onMount } from "svelte";
  import { apiFetch } from "$lib/api";
  import { i18n } from '$lib/i18n.svelte';
  import Card from './ui/Card.svelte';
  import Button from './ui/Button.svelte';

  let emails = $state<any[]>([]);
  let newEmail = $state("");
  let loading = $state(true);
  let saving = $state(false);
  let message = $state("");
  let isError = $state(false);

  async function fetchEmails() {
    loading = true;
    try {
      const res = await apiFetch("/api/v2/user/settings/emails");
      if (res.ok) emails = await res.json();
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function addEmail() {
    if (!newEmail) return;
    saving = true;
    message = "";
    isError = false;
    try {
      const res = await apiFetch("/api/v2/user/settings/emails", {
        method: "POST",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email: newEmail }),
      });
      if (res.ok) {
        message = "¡Correo añadido correctamente!";
        newEmail = "";
        await fetchEmails();
      } else {
        const body = await res.json();
        message = body.message || "Error al añadir el correo.";
        isError = true;
      }
    } catch (e) {
      message = "Error de red.";
      isError = true;
    } finally {
      saving = false;
    }
  }

  async function deleteEmail(email: string) {
    if (!confirm(`¿Seguro que quieres eliminar ${email}?`)) return;
    try {
      const res = await apiFetch("/api/v2/user/settings/emails", {
        method: "DELETE",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email }),
      });
      if (res.ok) fetchEmails();
    } catch (e) { console.error(e); }
  }

  async function makePrimary(email: string) {
    try {
      const res = await apiFetch("/api/v2/user/settings/emails/primary", {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({ email }),
      });
      if (res.ok) fetchEmails();
    } catch (e) { console.error(e); }
  }

  onMount(fetchEmails);
</script>

<Card class="settings-premium-card" padding="lg">
  <header class="section-header">
    <div class="title-with-icon">
      <span class="material-symbols-rounded hero-icon">mail</span>
      <div>
        <h2>Configuración de Email</h2>
        <p class="subtitle">Gestiona tus direcciones y elige tu correo principal de contacto</p>
      </div>
    </div>
  </header>

  <div class="modular-stack">
    <!-- EMAIL LIST -->
    <div class="config-module">
      <div class="module-tag">
        <span class="material-symbols-rounded">list_alt</span>
        Direcciones Vinculadas
      </div>
      <div class="island tonal-list">
        {#if loading}
          <div class="skeleton-emails"></div>
        {:else}
          {#each emails as email}
            <div class="email-row" class:is-primary={email.is_primary}>
              <div class="email-core">
                <span class="material-symbols-rounded email-icon">{email.is_primary ? 'verified' : 'alternate_email'}</span>
                <div class="email-meta">
                  <span class="address">{email.email}</span>
                  <div class="badge-row">
                    {#if email.is_primary}
                      <span class="badge primary">Principal</span>
                    {/if}
                    {#if email.is_activated}
                      <span class="badge active">Activado</span>
                    {:else}
                      <span class="badge pending">Pendiente de verificación</span>
                    {/if}
                  </div>
                </div>
              </div>
              <div class="email-actions">
                {#if !email.is_primary && email.is_activated}
                   <Button variant="ghost" size="sm" icon="star" onclick={() => makePrimary(email.email)}>Principal</Button>
                {/if}
                {#if !email.is_primary}
                   <Button variant="ghost" size="sm" icon="delete" class="danger-btn" onclick={() => deleteEmail(email.email)} />
                {/if}
              </div>
            </div>
          {/each}
        {/if}
      </div>
    </div>

    <!-- ADD EMAIL -->
    <div class="config-module">
       <div class="module-tag">
         <span class="material-symbols-rounded">add_circle</span>
         Añadir Nueva Dirección
       </div>
       <div class="island add-form-box">
          <form onsubmit={(e) => { e.preventDefault(); addEmail(); }} class="add-email-form">
            <div class="premium-input-group">
               <span class="material-symbols-rounded">mark_email_read</span>
               <input 
                 type="email" 
                 placeholder="ejemplo@correo.com" 
                 bind:value={newEmail} 
                 required 
               />
               <Button type="submit" variant="primary" disabled={saving || !newEmail}>
                  {saving ? 'Añadiendo...' : 'Vincular Correo'}
               </Button>
            </div>
          </form>
          {#if message}
            <div class="status-alert" class:is-error={isError}>
               <span class="material-symbols-rounded">{isError ? 'error' : 'check_circle'}</span>
               {message}
            </div>
          {/if}
       </div>
    </div>
  </div>
</Card>

<style>
  
  .section-header { margin-bottom: 4rem; }
  .title-with-icon { display: flex; align-items: center; gap: 2rem; }
  .hero-icon { font-size: 3rem; color: var(--primary-color); background: rgba(var(--primary-rgb), 0.1); padding: 1.25rem; border-radius: 2rem; flex-shrink: 0; }
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
  @media (max-width: 600px) {
    .modular-stack { gap: 2rem; }
  }
  .config-module { display: flex; flex-direction: column; gap: 1.5rem; }

  .module-tag { font-size: 0.8rem; font-weight: 950; text-transform: uppercase; color: var(--primary-color); letter-spacing: 0.15em; display: flex; align-items: center; gap: 0.75rem; border-bottom: 2px solid var(--md-sys-color-outline-variant); padding-bottom: 0.75rem; width: fit-content; opacity: 0.8; }
  .module-tag span { font-size: 1.25rem; }

  .tonal-list { display: flex; flex-direction: column; background: rgba(0,0,0,0.15); border-radius: 2.25rem; overflow: hidden; border: 1px solid var(--md-sys-color-outline-variant); }
  @media (max-width: 600px) {
    .tonal-list { border-radius: 1.25rem; }
  }
  
  .email-row { display: flex; align-items: center; justify-content: space-between; padding: 2rem 2.5rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); gap: 3rem; transition: background 0.3s; }
  
  @media (max-width: 800px) {
    .email-row { flex-direction: column; align-items: stretch; padding: 1.5rem; gap: 1.5rem; }
    .email-actions { justify-content: flex-end; }
  }

  @media (max-width: 600px) {
    .email-row { padding: 1rem 1.25rem; gap: 1rem; }
  }

  .email-row:last-child { border-bottom: none; }
  .email-row.is-primary { background: rgba(var(--primary-rgb), 0.04); }
  .email-row:hover { background: rgba(255,255,255,0.02); }

  .email-core { display: flex; align-items: center; gap: 1.5rem; min-width: 0; }
  .email-icon { font-size: 2rem; color: var(--primary-color); opacity: 0.5; flex-shrink: 0; }
  .is-primary .email-icon { opacity: 1; }
  
  .email-meta { display: flex; flex-direction: column; gap: 0.4rem; min-width: 0; }
  .address { font-size: 1.25rem; font-weight: 900; color: var(--text-color); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }
  
  .badge-row { display: flex; gap: 0.75rem; flex-wrap: wrap; }
  .badge { font-size: 0.7rem; font-weight: 900; text-transform: uppercase; padding: 0.25rem 0.75rem; border-radius: 8px; background: rgba(255,255,255,0.05); color: var(--text-muted); }
  .badge.primary { background: var(--primary-color); color: white; }
  .badge.active { background: rgba(16, 185, 129, 0.15); color: #10b981; }
  .badge.pending { background: rgba(245, 158, 11, 0.15); color: #f59e0b; }

  .email-actions { display: flex; gap: 0.75rem; }
  :global(.danger-btn:hover) { color: var(--error) !important; background: rgba(var(--error-rgb), 0.1) !important; }

  .add-form-box { padding: 3rem; border-radius: 2.5rem; border: 1px solid var(--md-sys-color-outline-variant); background: rgba(255,255,255,0.01); }
  @media (max-width: 600px) { .add-form-box { padding: 1.5rem; border-radius: 1.5rem; } }

  .premium-input-group { display: flex; align-items: center; gap: 1rem; background: var(--surface-variant); padding: 0.75rem 1.25rem; border-radius: 1.5rem; border: 2px solid transparent; transition: all 0.3s; }
  @media (max-width: 600px) { .premium-input-group { flex-direction: column; align-items: stretch; border-radius: 1rem; } }
  .premium-input-group:focus-within { border-color: var(--primary-color); background: rgba(var(--primary-rgb), 0.05); }
  .premium-input-group span { color: var(--primary-color); opacity: 0.6; }
  .premium-input-group input { flex: 1; background: transparent; border: none; outline: none; color: var(--text-color); font-size: 1.1rem; padding: 0.5rem 0; }

  .status-alert { margin-top: 1.5rem; display: flex; align-items: center; gap: 0.75rem; padding: 1rem 1.5rem; border-radius: 1.25rem; background: rgba(16, 185, 129, 0.1); color: #10b981; font-weight: 700; border: 1px solid rgba(16, 185, 129, 0.2); animation: fadeInUp 0.4s; }
  .status-alert.is-error { background: rgba(var(--error-rgb), 0.1); color: var(--error); border-color: rgba(var(--error-rgb), 0.2); }
  @keyframes fadeInUp { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }

  .skeleton-emails { height: 200px; background: var(--surface-variant); border-radius: 2rem; animation: pulse 1.5s infinite; }
</style>
