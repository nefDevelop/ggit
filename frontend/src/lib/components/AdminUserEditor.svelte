<script lang="ts">
  import { onMount } from 'svelte';
  
  let { userId = null, onSave, onCancel } = $props<{
    userId?: number | null;
    onSave: () => void;
    onCancel: () => void;
  }>();

  let user = $state<any>({
    username: '',
    email: '',
    password: '',
    full_name: '',
    activated: true,
    admin: false
  });

  let loading = $state(false);
  let submitting = $state(false);
  let error = $state('');

  async function fetchUser() {
    if (!userId) return;
    loading = true;
    try {
      const res = await fetch(`/api/v2/admin/users/${userId}`);
      if (res.ok) {
        const data = await res.json();
        user = {
          username: data.username,
          email: data.email,
          full_name: data.full_name,
          activated: data.is_activated,
          admin: data.is_admin,
          password: '' // Don't fetch password
        };
      }
    } catch (e) {
      error = 'Failed to load user data.';
    } finally {
      loading = false;
    }
  }

  async function save() {
    submitting = true;
    error = '';
    const method = userId ? 'PATCH' : 'POST';
    const url = userId ? `/api/v2/admin/users/${userId}` : '/api/v2/admin/users';
    
    // For update, only send changed fields. For create, password is required.
    const payload = { ...user };
    if (userId && !payload.password) delete payload.password;
    if (userId) delete payload.username; // Cannot change username in this view typically

    try {
      const res = await fetch(url, {
        method,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      });
      if (res.ok) {
        onSave();
      } else {
        const data = await res.json();
        error = data.error || 'Failed to save user.';
      }
    } catch (e) {
      error = 'A network error occurred.';
    } finally {
      submitting = false;
    }
  }

  onMount(fetchUser);
</script>

<div class="user-editor island tonal">
  <header>
    <h3>{userId ? 'Edit User' : 'Create New User'}</h3>
    <p>{userId ? `Modifying account: ${user.username}` : 'Add a new member to the platform.'}</p>
  </header>

  {#if loading}
    <div class="spinner"></div>
  {:else}
    <form onsubmit={(e) => { e.preventDefault(); save(); }}>
      <div class="form-grid">
        {#if !userId}
          <div class="field">
            <label for="username">Username</label>
            <input id="username" type="text" bind:value={user.username} required />
          </div>
        {/if}
        <div class="field">
          <label for="email">Email Address</label>
          <input id="email" type="email" bind:value={user.email} required />
        </div>
        <div class="field">
          <label for="password">{userId ? 'New Password (leave blank to keep current)' : 'Password'}</label>
          <input id="password" type="password" bind:value={user.password} required={!userId} />
        </div>
        <div class="field">
          <label for="full_name">Full Name</label>
          <input id="full_name" type="text" bind:value={user.full_name} />
        </div>
      </div>

      <div class="checkbox-group">
        <label class="checkbox">
          <input type="checkbox" bind:checked={user.activated} />
          <span>Account Activated</span>
        </label>
        <label class="checkbox">
          <input type="checkbox" bind:checked={user.admin} />
          <span>Administrator Privileges</span>
        </label>
      </div>

      {#if error}
        <div class="error-msg">{error}</div>
      {/if}

      <div class="actions">
        <button type="button" class="btn btn-secondary" onclick={onCancel}>Cancel</button>
        <button type="submit" class="btn btn-primary" disabled={submitting}>
          {submitting ? 'Saving...' : 'Save User'}
        </button>
      </div>
    </form>
  {/if}
</div>

<style>
  .user-editor { display: flex; flex-direction: column; gap: 1.5rem; max-width: 600px; }
  h3 { margin: 0; font-size: 1.25rem; }
  p { margin: 0; color: var(--text-muted); font-size: 0.9rem; }

  form { display: flex; flex-direction: column; gap: 1.5rem; }
  .form-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 1rem; }
  
  .field { display: flex; flex-direction: column; gap: 0.4rem; }
  .field label { font-size: 0.75rem; font-weight: 800; text-transform: uppercase; color: var(--text-muted); }
  .field input { background: rgba(0,0,0,0.2); border: 1px solid var(--border-color); color: #fff; padding: 0.6rem 0.8rem; border-radius: 0.5rem; }
  
  .checkbox-group { display: flex; gap: 2rem; }
  .checkbox { display: flex; align-items: center; gap: 0.5rem; cursor: pointer; font-size: 0.9rem; }
  
  .error-msg { color: #ef4444; font-size: 0.85rem; font-weight: 700; }
  
  .actions { display: flex; justify-content: flex-end; gap: 1rem; margin-top: 1rem; }
</style>
