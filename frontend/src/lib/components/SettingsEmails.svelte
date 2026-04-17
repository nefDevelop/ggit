<script lang="ts">
  import { onMount } from 'svelte';
  import { apiFetch } from '$lib/api';
  
  let emails = $state<any[]>([]);
  let newEmail = $state('');
  let loading = $state(true);
  let saving = $state(false);
  let message = $state('');
  let isError = $state(false);

  async function fetchEmails() {
    loading = true;
    try {
      const res = await apiFetch('/api/v2/user/settings/emails');
      if (res.ok) {
        emails = await res.json();
      }
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function addEmail() {
    if (!newEmail) return;
    saving = true;
    message = '';
    try {
      const res = await apiFetch('/api/v2/user/settings/emails', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email: newEmail })
      });
      if (res.ok) {
        message = 'Email added successfully!';
        newEmail = '';
        await fetchEmails();
      } else {
        const body = await res.json();
        message = body.message || 'Failed to add email.';
        isError = true;
      }
    } catch (e) {
      message = 'Network error.';
      isError = true;
    } finally {
      saving = false;
    }
  }

  async function deleteEmail(email: string) {
    if (!confirm(`Are you sure you want to delete ${email}?`)) return;
    try {
      const res = await apiFetch('/api/v2/user/settings/emails', {
        method: 'DELETE',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email })
      });
      if (res.ok) {
        await fetchEmails();
      }
    } catch (e) {
      console.error(e);
    }
  }

  async function makePrimary(email: string) {
    try {
      const res = await apiFetch('/api/v2/user/settings/emails/primary', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ email })
      });
      if (res.ok) {
        await fetchEmails();
      }
    } catch (e) {
      console.error(e);
    }
  }

  onMount(fetchEmails);
</script>

<div class="settings-card island">
  <header class="card-header">
    <h3 class="card-title">Manage Email Addresses</h3>
    <p class="card-desc">Your primary email address will be used for notifications and security alerts.</p>
  </header>

  <div class="email-list">
    {#if loading}
      <div class="skeleton-list">
        {#each Array(2) as _}<div class="skeleton-item"></div>{/each}
      </div>
    {:else}
      {#each emails as email}
        <div class="email-item {email.is_primary ? 'primary' : ''}">
          <div class="email-info">
            <span class="email-address">{email.email}</span>
            <div class="email-badges">
              {#if email.is_primary}
                <span class="badge primary-badge">Primary</span>
              {/if}
              {#if email.is_activated}
                <span class="badge active-badge">Activated</span>
              {:else}
                <span class="badge pending-badge">Pending</span>
              {/if}
            </div>
          </div>
          <div class="email-actions">
            {#if !email.is_primary && email.is_activated}
              <button class="btn btn-icon-text" onclick={() => makePrimary(email.email)}>
                <span class="material-symbols-rounded">star</span> Make Primary
              </button>
            {/if}
            {#if !email.is_primary}
              <button class="btn btn-icon-only btn-ghost error" onclick={() => deleteEmail(email.email)}>
                <span class="material-symbols-rounded">delete</span>
              </button>
            {/if}
          </div>
        </div>
      {/each}
    {/if}
  </div>

  <form class="add-email-form" onsubmit={(e) => { e.preventDefault(); addEmail(); }}>
    <div class="input-group">
      <input 
        type="email" 
        placeholder="Add new email address..." 
        bind:value={newEmail} 
        required 
      />
      <button type="submit" class="btn btn-primary" disabled={saving || !newEmail}>
        {saving ? 'Adding...' : 'Add'}
      </button>
    </div>
    {#if message}
      <p class="status-msg" class:error={isError}>{message}</p>
    {/if}
  </form>
</div>

<style>
  .settings-card { display: flex; flex-direction: column; }
  .card-header { padding: 1.5rem 2rem; border-bottom: 1px solid var(--border-color); }
  .card-title { margin: 0; font-size: 1.25rem; font-weight: 700; }
  .card-desc { margin: 0.5rem 0 0; font-size: 0.9rem; color: var(--text-muted); }

  .email-list { display: flex; flex-direction: column; }
  .email-item { 
    display: flex; 
    justify-content: space-between; 
    align-items: center; 
    padding: 1.5rem 2rem; 
    border-bottom: 1px solid var(--border-color);
    transition: background 0.2s;
  }
  .email-item:hover { background: rgba(255,255,255,0.02); }
  .email-item.primary { background: rgba(var(--primary-rgb), 0.03); }

  .email-info { display: flex; flex-direction: column; gap: 0.5rem; }
  .email-address { font-size: 1.05rem; font-weight: 600; color: #fff; }
  
  .email-badges { display: flex; gap: 0.5rem; }
  .badge { font-size: 0.7rem; font-weight: 800; text-transform: uppercase; padding: 0.15rem 0.5rem; border-radius: 4px; }
  .primary-badge { background: var(--primary-color); color: #000; }
  .active-badge { background: rgba(16, 185, 129, 0.2); color: #10b981; }
  .pending-badge { background: rgba(245, 158, 11, 0.2); color: #f59e0b; }

  .email-actions { display: flex; align-items: center; gap: 1rem; }

  .add-email-form { padding: 2rem; }
  .input-group { display: flex; gap: 1rem; }
  .input-group input { 
    flex: 1; 
    background: rgba(0,0,0,0.2); 
    border: 1px solid var(--border-color); 
    padding: 0.75rem 1rem; 
    border-radius: 0.75rem; 
    color: var(--text-color);
    font-size: 1rem;
    font-family: inherit;
  }
  .input-group input:focus { outline: none; border-color: var(--primary-color); background: rgba(0,0,0,0.4); }

  .status-msg { margin-top: 1rem; font-size: 0.9rem; color: #10b981; }
  .status-msg.error { color: #ef4444; }

  .skeleton-list { display: flex; flex-direction: column; }
  .skeleton-item { height: 80px; width: 100%; animation: pulse 1.5s infinite; border-bottom: 1px solid var(--border-color); }
  @keyframes pulse { 0%, 100% { opacity: 0.3; } 50% { opacity: 0.6; } }

  .btn-ghost.error { color: #ef4444; }
  .btn-ghost.error:hover { background: rgba(239, 68, 68, 0.1); }
</style>
