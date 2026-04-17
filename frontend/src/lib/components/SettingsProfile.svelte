<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from 'svelte';
  
  let profile = $state<any>({});
  let loading = $state(true);
  let saving = $state(false);
  let message = $state('');

  async function fetchProfile() {
    loading = true;
    try {
      const res = await apiFetch('/api/v2/user/settings/profile');
      if (res.ok) profile = await res.json();
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function saveProfile() {
    saving = true;
    message = '';
    try {
      const res = await apiFetch('/api/v2/user/settings/profile', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(profile)
      });
      if (res.ok) {
        message = 'Profile updated successfully!';
        setTimeout(() => message = '', 3000);
      }
    } catch (e) {
      message = 'Failed to update profile.';
    } finally {
      saving = false;
    }
  }

  onMount(fetchProfile);
</script>

<div class="settings-card island">
  <header class="card-header">
    <h3 class="card-title">Public Profile</h3>
  </header>
  
  {#if loading}
    <div class="skeleton-form"></div>
  {:else}
    <form class="settings-form" onsubmit={(e) => { e.preventDefault(); saveProfile(); }}>
      <div class="form-group">
        <label for="full_name">Full Name</label>
        <input type="text" id="full_name" bind:value={profile.full_name} placeholder="Your name" />
      </div>

      <div class="form-group">
        <label for="website">Website</label>
        <input type="url" id="website" bind:value={profile.website} placeholder="https://example.com" />
      </div>

      <div class="form-group">
        <label for="location">Location</label>
        <input type="text" id="location" bind:value={profile.location} placeholder="City, Country" />
      </div>

      <div class="form-group checkbox-group">
        <label class="checkbox-container">
          <input type="checkbox" bind:checked={profile.disable_news_feed} />
          <span class="checkmark"></span>
          <div class="checkbox-text">
            <strong>Hide Activity Feed</strong>
            <p>Hide the technical audit log from your dashboard.</p>
          </div>
        </label>
      </div>

      <div class="form-actions">
        {#if message}
           <span class="status-msg" class:success={message.includes('success')}>{message}</span>
        {/if}
        <button type="submit" class="btn btn-primary" disabled={saving}>
          {saving ? 'Saving...' : 'Update Profile'}
        </button>
      </div>
    </form>
  {/if}
</div>

<style>
  .settings-card { display: flex; flex-direction: column; }
  .card-header { padding: 1.5rem 2rem; border-bottom: 1px solid var(--border-color); }
  .card-title { margin: 0; font-size: 1.25rem; font-weight: 700; }

  .settings-form { padding: 2rem; display: flex; flex-direction: column; gap: 1.5rem; }

  .form-group { display: flex; flex-direction: column; gap: 0.5rem; }
  .form-group label { font-size: 0.85rem; font-weight: 800; text-transform: uppercase; color: var(--text-muted); padding-left: 0.25rem; }
  
  input[type="text"], input[type="url"], input[type="email"] {
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

  .checkbox-group { margin-top: 1rem; padding: 1rem; background: rgba(255,255,255,0.02); border-radius: 1rem; }
  .checkbox-container { display: flex; align-items: flex-start; gap: 1rem; cursor: pointer; position: relative; }
  .checkbox-container input { visibility: hidden; width: 0; height: 0; position: absolute; }
  
  .checkmark { width: 1.5rem; height: 1.5rem; background: rgba(255,255,255,0.05); border: 2px solid var(--border-color); border-radius: 0.5rem; flex-shrink: 0; position: relative; transition: all 0.2s; }
  .checkbox-container input:checked + .checkmark { background: var(--primary-color); border-color: var(--primary-color); }
  .checkmark::after { content: 'check'; font-family: 'Material Symbols Rounded'; position: absolute; top: 50%; left: 50%; transform: translate(-50%, -50%); color: white; opacity: 0; font-size: 1rem; }
  .checkbox-container input:checked + .checkmark::after { opacity: 1; }

  .checkbox-text strong { display: block; font-size: 0.95rem; margin-bottom: 0.25rem; }
  .checkbox-text p { font-size: 0.85rem; color: var(--text-muted); margin: 0; }

  .form-actions { display: flex; justify-content: flex-end; align-items: center; gap: 1.5rem; margin-top: 2rem; padding-top: 1.5rem; border-top: 1px solid var(--border-color); }
  .status-msg { font-size: 0.9rem; color: #ef4444; }
  .status-msg.success { color: #10b981; }

  .skeleton-form { height: 300px; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 0.3; } 50% { opacity: 0.6; } }
</style>
