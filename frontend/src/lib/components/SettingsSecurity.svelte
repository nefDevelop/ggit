<script lang="ts">
  import { onMount } from 'svelte';
  import { apiFetch } from '$lib/api';
  
  let oldPassword = $state('');
  let newPassword = $state('');
  let repeatPassword = $state('');
  let saving = $state(false);
  let message = $state('');
  let isError = $state(false);

  // 2FA State
  let twoFactorEnabled = $state(false);
  let enrollment = $state(null as any);
  let passcode = $state('');
  let recoveryCodes = $state([] as string[]);
  let showRecovery = $state(false);
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
      const res = await apiFetch('/api/v2/user/settings/password', {
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

  async function fetch2FAStatus() {
    try {
      const res = await apiFetch('/api/v2/user/settings/security/two_factor');
      if (res.ok) {
        const data = await res.json();
        twoFactorEnabled = data.enabled;
      }
    } catch (e) {
      console.error(e);
    }
  }

  async function initiate2FA() {
    try {
      const res = await apiFetch('/api/v2/user/settings/security/two_factor/enable', { method: 'POST' });
      if (res.ok) {
        enrollment = await res.json();
      }
    } catch (e) {
      console.error(e);
    }
  }

  async function confirm2FA() {
    try {
      const res = await apiFetch('/api/v2/user/settings/security/two_factor/confirm', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({ passcode })
      });
      if (res.ok) {
        const data = await res.json();
        recoveryCodes = data.recovery_codes;
        twoFactorEnabled = true;
        enrollment = null;
        showRecovery = true;
      } else {
        alert('Invalid passcode');
      }
    } catch (e) {
      console.error(e);
    }
  }

  async function disable2FA() {
    if (!confirm('Are you sure you want to disable 2FA? This will make your account less secure.')) return;
    try {
      const res = await apiFetch('/api/v2/user/settings/security/two_factor', { method: 'DELETE' });
      if (res.ok) {
        twoFactorEnabled = false;
        showRecovery = false;
      }
    } catch (e) {
      console.error(e);
    }
  }

  onMount(() => {
    fetch2FAStatus();
  });
</script>

<div class="settings-stack">
  <div class="settings-card island island-interactive">
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

  <div class="settings-card island island-interactive {twoFactorEnabled ? 'tonal-success' : 'tonal-warning'}">
    <header class="card-header">
      <div class="header-with-badge">
        <h3 class="card-title">Two-Factor Authentication</h3>
        <span class="status-badge {twoFactorEnabled ? 'enabled' : 'disabled'}">
          {twoFactorEnabled ? 'Protected' : 'Unprotected'}
        </span>
      </div>
      <p class="card-desc">Add an extra layer of security to your account by requiring more than just a password to log in.</p>
    </header>

    <div class="card-body">
      {#if showRecovery}
        <div class="recovery-box island island-interactive">
          <h4>Success! 2FA is enabled.</h4>
          <p>Please save these recovery codes in a safe place. You can use them to log in if you lose access to your device.</p>
          <div class="recovery-codes">
            {#each recoveryCodes as code}
              <code>{code}</code>
            {/each}
          </div>
          <button class="btn btn-tonal" onclick={() => showRecovery = false}>I've saved them</button>
        </div>
      {:else if enrollment}
        <div class="enrollment-flow island island-interactive">
          <div class="qr-section">
            <img src={enrollment.qr_code} alt="QR Code" class="qr-image" />
            <div class="secret-text">
              <p>Scan this QR code with your authenticator app (e.g., Google Authenticator, Authy).</p>
              <p>Or enter this secret manually:</p>
              <code>{enrollment.secret}</code>
            </div>
          </div>
          <div class="confirm-section">
            <label for="passcode">Enter 6-digit passcode to verify:</label>
            <div class="input-group">
              <input type="text" id="passcode" bind:value={passcode} placeholder="000000" maxlength="6" />
              <button class="btn btn-primary" onclick={confirm2FA} disabled={passcode.length < 6}>Verify & Enable</button>
            </div>
            <button class="btn btn-ghost" onclick={() => enrollment = null}>Cancel</button>
          </div>
        </div>
      {:else if twoFactorEnabled}
        <div class="enabled-state">
          <p>You have two-factor authentication enabled via an authenticator app.</p>
          <div class="actions">
             <button class="btn btn-outline" onclick={async () => {
                const res = await apiFetch('/api/v2/user/settings/security/two_factor/recovery_codes');
                if (res.ok) {
                   const data = await res.json();
                   recoveryCodes = data.recovery_codes;
                   showRecovery = true;
                }
             }}>View Recovery Codes</button>
             <button class="btn btn-ghost error" onclick={disable2FA}>Disable 2FA</button>
          </div>
        </div>
      {:else}
        <div class="disabled-state">
          <p>Two-factor authentication is not currently enabled for your account.</p>
          <button class="btn btn-primary" onclick={initiate2FA}>Enable 2FA</button>
        </div>
      {/if}
    </div>
  </div>
</div>

<style>
  .settings-stack { display: flex; flex-direction: column; gap: 2rem; }
  .settings-card { display: flex; flex-direction: column; }
  .card-header { padding: 1.5rem 2rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); }
  .card-title { margin: 0; font-size: 1.25rem; font-weight: 700; }
  .card-desc { font-size: 0.9rem; color: var(--text-muted); margin: 0.5rem 0 0; }
  
  .header-with-badge { display: flex; align-items: center; justify-content: space-between; }
  .status-badge { font-size: 0.7rem; font-weight: 800; text-transform: uppercase; padding: 0.2rem 0.6rem; border-radius: 4px; }
  .status-badge.enabled { background: #10b981; color: #000; }
  .status-badge.disabled { background: #f59e0b; color: #000; }

  .card-body { padding: 2rem; }

  .settings-form { display: flex; flex-direction: column; gap: 1.5rem; }
  .form-group { display: flex; flex-direction: column; gap: 0.5rem; }
  .form-group label { font-size: 0.85rem; font-weight: 800; text-transform: uppercase; color: var(--text-muted); padding-left: 0.25rem; }
  
  input {
    background: rgba(0,0,0,0.2);
    border: 1px solid var(--md-sys-color-outline-variant);
    padding: 0.75rem 1rem;
    border-radius: 0.75rem;
    color: var(--text-color);
    font-family: inherit;
    font-size: 1rem;
    transition: all 0.2s;
  }
  input:focus { outline: none; border-color: var(--primary-color); background: rgba(0,0,0,0.4); }

  .form-actions { display: flex; justify-content: flex-end; align-items: center; gap: 1.5rem; margin-top: 1rem; padding-top: 1.5rem; border-top: 1px solid var(--md-sys-color-outline-variant); }
  .status-msg { font-size: 0.9rem; color: #10b981; }
  .status-msg.error { color: #ef4444; }

  /* 2FA Styles */
  .recovery-box { padding: 2rem; background: rgba(16, 185, 129, 0.05); border: 1px solid rgba(16, 185, 129, 0.2); }
  .recovery-codes { display: grid; grid-template-columns: repeat(2, 1fr); gap: 0.5rem; margin: 1.5rem 0; }
  .recovery-codes code { background: #000; padding: 0.5rem; border-radius: 0.4rem; text-align: center; color: var(--primary-color); font-weight: 700; letter-spacing: 0.1em; }

  .enrollment-flow { display: flex; flex-direction: column; gap: 2rem; }
  .qr-section { display: flex; align-items: center; gap: 2rem; }
  .qr-image { width: 180px; height: 180px; background: #fff; padding: 10px; border-radius: 1rem; }
  .secret-text code { display: block; background: #000; padding: 0.75rem; border-radius: 0.5rem; margin-top: 0.5rem; font-size: 1.1rem; color: var(--primary-color); }

  .confirm-section { display: flex; flex-direction: column; gap: 1rem; border-top: 1px solid var(--md-sys-color-outline-variant); pt: 1.5rem; }
  .input-group { display: flex; gap: 1rem; }
  .input-group input { width: 150px; text-align: center; font-size: 1.5rem; letter-spacing: 0.2em; font-weight: 800; }

  .actions { display: flex; gap: 1rem; margin-top: 1rem; }
  .btn-ghost.error { color: #ef4444; }
  .btn-ghost.error:hover { background: rgba(239, 68, 68, 0.1); }
</style>
