<script lang="ts">
  import { onMount } from 'svelte';
  import { authState } from '$lib/auth.svelte';
  import ThemeSelector from '$lib/components/ThemeSelector.svelte';

  let profile = $state({
    username: '',
    full_name: '',
    email: '',
    website: '',
    location: '',
    disable_news_feed: false
  });

  let loading = $state(true);
  let saving = $state(false);
  let message = $state('');
  let error = $state('');

  function getCsrfToken() {
    if (typeof document === 'undefined') return '';
    const match = document.cookie.match(new RegExp('(^| )XSRF-TOKEN=([^;]+)'));
    return match ? match[2] : '';
  }

  async function fetchProfile() {
    try {
      const res = await fetch('/api/v2/user/settings/profile');
      if (res.ok) {
        profile = await res.json();
      }
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function handleUpdate() {
    saving = true;
    message = '';
    error = '';
    try {
      const res = await fetch('/api/v2/user/settings/profile', {
        method: 'PUT',
        headers: { 
          'Content-Type': 'application/json',
          'X-Csrf-Token': getCsrfToken()
        },
        body: JSON.stringify({
          name: profile.username,
          full_name: profile.full_name,
          website: profile.website,
          location: profile.location,
          disable_news_feed: profile.disable_news_feed
        })
      });
      if (res.ok) {
        message = 'Profile updated successfully.';
        // Update global authState
        authState.user.username = profile.username;
        authState.user.disable_news_feed = profile.disable_news_feed;
      } else {
        const errData = await res.json();
        error = errData.message || 'Failed to update profile.';
      }
    } catch (e) {
      error = 'Internal server error.';
    } finally {
      saving = false;
    }
  }

  onMount(fetchProfile);
</script>

<svelte:head>
  <title>Settings - Ggit</title>
</svelte:head>

<div class="settings-container">
  <header class="settings-header">
    <h1 class="title">Account Settings</h1>
    <p class="subtitle">Manage your profile and dashboard preferences.</p>
  </header>

  {#if loading}
    <div class="island skeleton-island">
      <div class="skeleton-line"></div>
      <div class="skeleton-line"></div>
      <div class="skeleton-line"></div>
    </div>
  {:else}
    <div class="settings-grid">
      <section class="settings-section island">
        <h2 class="section-title">Public Profile</h2>
        
        <div class="form-group">
          <label for="username">Username</label>
          <input type="text" id="username" bind:value={profile.username} placeholder="Username" />
          <p class="field-hint">Your unique username on Ggit.</p>
        </div>

        <div class="form-group">
          <label for="full_name">Full Name</label>
          <input type="text" id="full_name" bind:value={profile.full_name} placeholder="Your legal or display name" />
        </div>

        <div class="form-group">
          <label for="website">Website</label>
          <input type="text" id="website" bind:value={profile.website} placeholder="https://example.com" />
        </div>

        <div class="form-group">
          <label for="location">Location</label>
          <input type="text" id="location" bind:value={profile.location} placeholder="City, Country" />
        </div>

        <button class="btn btn-primary" onclick={handleUpdate} disabled={saving}>
          {saving ? 'Saving...' : 'Save Changes'}
        </button>

        {#if message}
          <p class="success-message">{message}</p>
        {/if}
        {#if error}
          <p class="error-message">{error}</p>
        {/if}
      </section>

      <section class="settings-section island">
        <h2 class="section-title">Dashboard Preferences</h2>
        
        <div class="toggle-group">
          <div class="toggle-info">
            <span class="toggle-label">Enable Activity Feed</span>
            <p class="toggle-hint">Show pushes and merges on your dashboard.</p>
          </div>
          <button 
            class="toggle-switch" 
            class:active={!profile.disable_news_feed} 
            onclick={() => { profile.disable_news_feed = !profile.disable_news_feed; handleUpdate(); }}
          >
            <div class="toggle-handle"></div>
          </button>
        </div>
      </section>

      <ThemeSelector />
    </div>
  {/if}
</div>

<style>
  .settings-container {
    max-width: 900px;
    margin: 0 auto;
    padding: 2rem;
  }

  .settings-header {
    margin-bottom: 2.5rem;
  }

  .title { font-size: 2.25rem; font-weight: 900; letter-spacing: -0.02em; margin-bottom: 0.5rem; }
  .subtitle { color: var(--text-muted); font-size: 1.1rem; }

  .settings-grid {
    display: flex;
    flex-direction: column;
    gap: 2rem;
  }

  .settings-section {
    padding: 2rem;
  }

  .section-title {
    font-size: 1.25rem;
    font-weight: 800;
    margin-bottom: 1.5rem;
    color: var(--primary-color);
  }

  .form-group {
    margin-bottom: 1.5rem;
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  label { font-size: 0.85rem; font-weight: 700; color: var(--text-muted); text-transform: uppercase; letter-spacing: 0.05em; }
  
  input {
    background: rgba(255,255,255,0.03);
    border: 1px solid var(--border-color);
    padding: 0.75rem 1rem;
    border-radius: 0.75rem;
    color: var(--text-color);
    font-size: 1rem;
    transition: all 0.2s;
  }
  input:focus { outline: none; border-color: var(--primary-color); background: rgba(255,255,255,0.05); }

  .field-hint { font-size: 0.75rem; color: var(--text-muted); }

  .toggle-group {
    display: flex;
    justify-content: space-between;
    align-items: center;
    background: rgba(0,0,0,0.2);
    padding: 1.25rem;
    border-radius: 1rem;
  }

  .toggle-label { font-weight: 700; font-size: 1rem; }
  .toggle-hint { font-size: 0.8rem; color: var(--text-muted); }

  .toggle-switch {
    width: 50px;
    height: 28px;
    background: var(--md3-surface-5);
    border-radius: 14px;
    position: relative;
    cursor: pointer;
    border: none;
    transition: background 0.3s;
  }
  .toggle-switch.active { background: var(--primary-color); }
  
  .toggle-handle {
    width: 20px;
    height: 20px;
    background: white;
    border-radius: 50%;
    position: absolute;
    top: 4px;
    left: 4px;
    transition: transform 0.3s;
  }
  .toggle-switch.active .toggle-handle { transform: translateX(22px); }

  .success-message { color: #10b981; font-size: 0.9rem; margin-top: 1rem; font-weight: 600; }
  .error-message { color: #f87171; font-size: 0.9rem; margin-top: 1rem; font-weight: 600; }

  .skeleton-island { height: 300px; display: flex; flex-direction: column; gap: 1rem; padding: 2rem; justify-content: center; }
  .skeleton-line { height: 40px; background: var(--md3-surface-3); border-radius: 0.5rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 0.5; } 50% { opacity: 0.8; } }
</style>
