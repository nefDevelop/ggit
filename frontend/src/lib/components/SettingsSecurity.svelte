<script lang="ts">
  import { onMount } from 'svelte';
  
  let oldPassword = $state('');
  let newPassword = $state('');
  let repeatPassword = $state('');
  let saving = $state(false);
  let message = $state('');
  let isError = $state(false);

  async function updatePassword() {
    if (newPassword !== repeatPassword) {
      message = 'Passwords do not match.';
      isError = true;
      return;
    }

    saving = true;
    message = '';
    isError = false;
    try {
      const res = await fetch('/api/v2/user/settings/password', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ 
          old_password: oldPassword, 
          password: newPassword, 
          retype: repeatPassword 
        })
      });
      if (res.ok) {
        message = 'Password updated successfully!';
        oldPassword = '';
        newPassword = '';
        repeatPassword = '';
      } else {
        const body = await res.json();
        message = body.message || 'Failed to update password.';
        isError = true;
      }
    } catch (e) {
      message = 'Network error.';
      isError = true;
    } finally {
      saving = false;
    }
  }
</script>

<div class="settings-card island">
  <header class="card-header">
    <h3 class="card-title">Security & Password</h3>
  </header>
  
  <form class="settings-form" onsubmit={(e) => { e.preventDefault(); updatePassword(); }}>
    <div class="form-group">
      <label for="old_pass">Current Password</label>
      <input type="password" id="old_pass" bind:value={oldPassword} required />
    </div>

    <div class="form-group">
      <label for="new_pass">New Password</label>
      <input type="password" id="new_pass" bind:value={newPassword} required />
    </div>

    <div class="form-group">
      <label for="rep_pass">Repeat New Password</label>
      <input type="password" id="rep_pass" bind:value={repeatPassword} required />
    </div>

    <div class="form-actions">
      {#if message}
         <span class="status-msg" class:error={isError}>{message}</span>
      {/if}
      <button type="submit" class="btn btn-primary" disabled={saving || !oldPassword || !newPassword}>
        {saving ? 'Updating...' : 'Change Password'}
      </button>
    </div>
  </form>
</div>

<style>
  .settings-card { display: flex; flex-direction: column; }
  .card-header { padding: 1.5rem 2rem; border-bottom: 1px solid var(--border-color); }
  .card-title { margin: 0; font-size: 1.25rem; font-weight: 700; }

  .settings-form { padding: 2rem; display: flex; flex-direction: column; gap: 1.5rem; }

  .form-group { display: flex; flex-direction: column; gap: 0.5rem; }
  .form-group label { font-size: 0.85rem; font-weight: 800; text-transform: uppercase; color: var(--text-muted); padding-left: 0.25rem; }
  
  input[type="password"] {
    background: rgba(0,0,0,0.2);
    border: 1px solid var(--border-color);
    padding: 0.75rem 1rem;
    border-radius: 0.75rem;
    color: var(--text-color);
    font-family: inherit;
    font-size: 1rem;
    transition: all 0.2s;
  }
  input:focus { outline: none; border-color: var(--primary-color); background: rgba(0,0,0,0.4); }

  .form-actions { display: flex; justify-content: flex-end; align-items: center; gap: 1.5rem; margin-top: 1rem; padding-top: 1.5rem; border-top: 1px solid var(--border-color); }
  .status-msg { font-size: 0.9rem; color: #10b981; }
  .status-msg.error { color: #ef4444; }
</style>
