<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  import { authState } from '$lib/auth.svelte';
  import ThemeSelector from '$lib/components/ThemeSelector.svelte';

  let profile = $state({
    username: '',
    full_name: '',
    email: '',
    website: '',
    location: '',
    disable_news_feed: false,
    enable_dicebear_avatar: false,
    avatar_url: '',
    use_custom_avatar: false
  });

  const dicebearStyles = ['avataaars', 'bottts', 'pixel-art', 'lorelei-neutral', 'micah', 'open-peeps', 'notionists'];
  let currentDicebearStyle = $state('avataaars');
  let dicebearSeed = $state('');

  const dicebearUrl = $derived(`https://api.dicebear.com/7.x/${currentDicebearStyle}/png?seed=${dicebearSeed || profile.username}`);
  let avatarCacheBuster = $state(Date.now());

  async function saveDicebearAvatar() {
    saving = true;
    try {
      const imgRes = await fetch(dicebearUrl);
      const blob = await imgRes.blob();
      
      const formData = new FormData();
      const file = new File([blob], 'avatar.png', { type: 'image/png' });
      formData.append('avatar', file);

      const res = await apiFetch('/api/v2/user/settings/avatar', {
        method: 'POST',
        body: formData
      });

      if (res.ok) {
        message = 'DiceBear avatar saved successfully.';
        avatarCacheBuster = Date.now();
        fetchProfile();
      }
    } catch (e) {
      error = 'Failed to save DiceBear avatar.';
    } finally {
      saving = false;
    }
  }

  function shuffleDicebear() {
    dicebearSeed = Math.random().toString(36).substring(7);
  }

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
      const res = await apiFetch('/api/v2/user/settings/profile');
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
      const res = await apiFetch('/api/v2/user/settings/profile', {
        method: 'PUT',
        headers: { 
          'Content-Type': 'application/json'
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

  async function handleAvatarUpload(e: Event) {
    const file = (e.target as HTMLInputElement).files?.[0];
    if (!file) return;

    saving = true;
    error = '';
    try {
      const formData = new FormData();
      formData.append('avatar', file);

      const res = await apiFetch('/api/v2/user/settings/avatar', {
        method: 'POST',
        body: formData
      });

      if (res.ok) {
        message = 'Avatar updated successfully.';
        avatarCacheBuster = Date.now();
        fetchProfile();
      } else {
        error = 'Failed to upload avatar.';
      }
    } catch (e) {
      error = 'Error uploading file.';
    } finally {
      saving = false;
    }
  }

  async function deleteAvatar() {
    saving = true;
    try {
      const res = await apiFetch('/api/v2/user/settings/avatar', {
        method: 'DELETE'
      });
      if (res.ok) {
        message = 'Avatar deleted. Using Gravatar/Identicon.';
        avatarCacheBuster = Date.now();
        fetchProfile();
      }
    } catch (e) {
      error = 'Failed to delete avatar.';
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

      <section class="settings-section island">
        <header class="section-top">
          <h2 class="section-title">Avatar Management</h2>
          {#if profile.enable_dicebear_avatar}
            <div class="section-badge mini">DiceBear Active</div>
          {/if}
        </header>

        <div class="avatar-management-grid">
          <!-- Current Avatar Section -->
          <div class="current-avatar-box">
             <label>Current Photo</label>
             <div class="avatar-display">
                <img src={`${profile.avatar_url}${profile.avatar_url.includes('?') ? '&' : '?'}t=${avatarCacheBuster}`} alt="Current Avatar" class="large-avatar" />
                <div class="avatar-overlay-controls">
                  <label class="icon-btn" title="Upload Photo">
                    <span class="material-symbols-rounded">cloud_upload</span>
                    <input type="file" accept="image/*" onchange={handleAvatarUpload} hidden />
                  </label>
                  {#if profile.use_custom_avatar}
                    <button class="icon-btn delete" onclick={deleteAvatar} title="Remove Custom Avatar">
                      <span class="material-symbols-rounded">delete</span>
                    </button>
                  {/if}
                </div>
             </div>
             <p class="tiny-hint">Manual upload or drag & drop</p>
          </div>

          <!-- DiceBear Section (Optional) -->
          {#if profile.enable_dicebear_avatar}
            <div class="dicebear-box">
              <label>Generate with DiceBear</label>
              <div class="dicebear-workspace">
                 <div class="dicebear-preview">
                    <img src={dicebearUrl} alt="DiceBear Preview" class="dicebear-img" />
                    <button class="shuffle-btn" onclick={shuffleDicebear} title="Shuffle Seed">
                      <span class="material-symbols-rounded">casino</span>
                    </button>
                 </div>
                 <div class="dicebear-styles">
                   <div class="style-grid">
                    {#each dicebearStyles as style}
                      <button 
                        class="style-chip" 
                        class:active={currentDicebearStyle === style}
                        onclick={() => currentDicebearStyle = style}
                      >
                        {style}
                      </button>
                    {/each}
                   </div>
                   <button class="btn btn-primary btn-sm save-dice-btn" onclick={saveDicebearAvatar} disabled={saving}>
                      <span class="material-symbols-rounded">save</span> Use this Avatar
                   </button>
                 </div>
              </div>
            </div>
          {/if}
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

  .avatar-management-grid { display: grid; grid-template-columns: auto 1fr; gap: 3rem; align-items: start; }
  .section-top { display: flex; justify-content: space-between; align-items: center; margin-bottom: 2rem; }
  .section-badge.mini { position: static; font-size: 0.6rem; }

  .current-avatar-box { display: flex; flex-direction: column; gap: 1rem; align-items: center; }
  .avatar-display { position: relative; width: 120px; height: 120px; }
  .large-avatar { width: 100%; height: 100%; border-radius: 2rem; object-fit: cover; border: 4px solid var(--md3-surface-2); box-shadow: var(--shadow-lg); transition: transform 0.3s cubic-bezier(0.4, 0, 0.2, 1); }
  
  .avatar-overlay-controls { position: absolute; top: 0; left: 0; width: 100%; height: 100%; display: flex; align-items: center; justify-content: center; gap: 0.5rem; background: rgba(0,0,0,0.4); border-radius: 2rem; opacity: 0; transition: opacity 0.2s; }
  .avatar-display:hover .avatar-overlay-controls { opacity: 1; }
  .avatar-display:hover .large-avatar { transform: scale(0.95); }

  .icon-btn { width: 36px; height: 36px; border-radius: 50%; background: white; color: var(--primary-color); display: flex; align-items: center; justify-content: center; cursor: pointer; border: none; transition: transform 0.2s; }
  .icon-btn:hover { transform: scale(1.1); }
  .icon-btn.delete { color: var(--error); }

  .dicebear-box { background: var(--md3-surface-2); padding: 1.5rem; border-radius: 1.5rem; display: flex; flex-direction: column; gap: 1rem; }
  .dicebear-workspace { display: flex; gap: 1.5rem; align-items: center; }
  
  .dicebear-preview { position: relative; width: 100px; height: 100px; flex-shrink: 0; }
  .dicebear-img { width: 100%; height: 100%; border-radius: 1rem; background: white; padding: 5px; }
  .shuffle-btn { position: absolute; bottom: -8px; right: -8px; width: 32px; height: 32px; border-radius: 50%; background: var(--primary-color); color: white; border: 2px solid var(--md3-surface-2); display: flex; align-items: center; justify-content: center; cursor: pointer; font-size: 14px; }

  .dicebear-styles { flex: 1; display: flex; flex-direction: column; gap: 1rem; }
  .style-grid { display: flex; flex-wrap: wrap; gap: 0.5rem; }
  .style-chip { padding: 0.3rem 0.8rem; border-radius: 2rem; background: var(--md3-surface-3); border: 1px solid var(--border-color); font-size: 0.75rem; font-weight: 700; cursor: pointer; color: var(--text-muted); transition: all 0.2s; }
  .style-chip:hover { background: var(--border-color); }
  .style-chip.active { background: var(--primary-color); color: white; border-color: var(--primary-color); }
  
  .save-dice-btn { align-self: flex-start; }
  .btn-sm { padding: 0.4rem 1rem; font-size: 0.8rem; }

  .skeleton-island { height: 300px; display: flex; flex-direction: column; gap: 1rem; padding: 2rem; justify-content: center; }
  .skeleton-line { height: 40px; background: var(--md3-surface-3); border-radius: 0.5rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 0.5; } 50% { opacity: 0.8; } }
</style>
