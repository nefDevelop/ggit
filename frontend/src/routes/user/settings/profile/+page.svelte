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
  let showDiceBear = $state(false);

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
        await fetchProfile();
        // Update global authState
        if (authState.user) {
          authState.user.avatar_url = profile.avatar_url;
        }
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
        authState.user.full_name = profile.full_name;
        authState.user.avatar_url = profile.avatar_url;
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
        await fetchProfile();
        // Update global authState
        if (authState.user) {
          authState.user.avatar_url = profile.avatar_url;
        }
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
    <div class="skeleton-island tonal">
      <div class="skeleton-line"></div>
      <div class="skeleton-line"></div>
      <div class="skeleton-line"></div>
    </div>
  {:else}
    <div class="settings-grid">
      <section class="settings-section tonal profile-main-card">
        <div class="profile-layout">
          <div class="avatar-column">
            <div class="avatar-container">
              <img src={`${profile.avatar_url}${profile.avatar_url.includes('?') ? '&' : '?'}t=${avatarCacheBuster}`} alt="Current Avatar" class="profile-avatar" />
              <div class="avatar-overlay">
                <label class="overlay-action" title="Upload Photo">
                  <span class="material-symbols-rounded">cloud_upload</span>
                  <input type="file" accept="image/*" onchange={handleAvatarUpload} hidden />
                </label>
                {#if profile.use_custom_avatar}
                  <button class="overlay-action delete" onclick={deleteAvatar} title="Remove Custom Avatar">
                    <span class="material-symbols-rounded">delete</span>
                  </button>
                {/if}
              </div>
            </div>
            {#if profile.enable_dicebear_avatar}
              <button class="magic-btn" class:active={showDiceBear} onclick={() => showDiceBear = !showDiceBear}>
                <span class="material-symbols-rounded">auto_fix_high</span>
                {showDiceBear ? 'Hide Generator' : 'Magic Avatar'}
              </button>
            {/if}
          </div>

          <div class="info-column">
            <h2 class="section-title">Public Profile</h2>
            
            <div class="form-grid">
              <div class="form-group">
                <label for="username">Username</label>
                <input type="text" id="username" bind:value={profile.username} placeholder="Username" />
                <p class="field-hint">Your unique identifier on Ggit.</p>
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
            </div>

            <div class="form-actions">
              <button class="btn btn-primary" onclick={handleUpdate} disabled={saving}>
                {saving ? 'Saving...' : 'Save Changes'}
              </button>
              {#if message}
                <p class="success-message">{message}</p>
              {/if}
              {#if error}
                <p class="error-message">{error}</p>
              {/if}
            </div>
          </div>
        </div>

        {#if showDiceBear && profile.enable_dicebear_avatar}
          <div class="dicebear-expansion">
            <div class="dicebear-tool">
              <div class="dicebear-preview">
                <img src={dicebearUrl} alt="DiceBear Preview" class="preview-img" />
                <button class="shuffle-btn" onclick={shuffleDicebear} title="Shuffle Seed">
                  <span class="material-symbols-rounded">casino</span>
                </button>
              </div>
              <div class="dicebear-config">
                <div class="style-list">
                  {#each dicebearStyles as style}
                    <button 
                      class="style-pill" 
                      class:active={currentDicebearStyle === style}
                      onclick={() => currentDicebearStyle = style}
                    >
                      {style}
                    </button>
                  {/each}
                </div>
                <button class="btn btn-secondary btn-sm" onclick={saveDicebearAvatar} disabled={saving}>
                  <span class="material-symbols-rounded">check</span> Use this Avatar
                </button>
              </div>
            </div>
          </div>
        {/if}
      </section>

      <section class="settings-section tonal preferences-card">
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

  /* New Profile Card Layout */
  .profile-main-card { padding: 2.5rem; }
  .profile-layout { display: flex; gap: 3rem; }
  
  .avatar-column { display: flex; flex-direction: column; align-items: center; gap: 1.5rem; }
  .avatar-container { position: relative; width: 160px; height: 160px; border-radius: 3.5rem; overflow: hidden; border: 4px solid var(--md3-surface-2); box-shadow: var(--shadow-xl); background: var(--md3-surface-1); }
  .profile-avatar { width: 100%; height: 100%; object-fit: cover; transition: transform 0.6s cubic-bezier(0.34, 1.56, 0.64, 1); }
  
  .avatar-overlay { position: absolute; inset: 0; background: rgba(0,0,0,0.5); opacity: 0; transition: all 0.3s ease; display: flex; align-items: center; justify-content: center; gap: 0.75rem; backdrop-filter: blur(2px); }
  
  .avatar-container:hover .avatar-overlay { opacity: 1; }
  .avatar-container:hover .profile-avatar { transform: scale(1.15); }

  .overlay-action { width: 42px; height: 42px; border-radius: 50%; background: white; color: var(--primary-color); display: flex; align-items: center; justify-content: center; cursor: pointer; border: none; transition: transform 0.2s; box-shadow: 0 4px 12px rgba(0,0,0,0.2); }
  .overlay-action:hover { transform: scale(1.15); }
  .overlay-action.delete { color: #ef4444; }

  .magic-btn { display: flex; align-items: center; gap: 0.5rem; background: var(--md3-surface-3); border: 1px solid var(--border-color); color: var(--text-color); padding: 0.6rem 1.2rem; border-radius: 100px; font-weight: 700; font-size: 0.85rem; cursor: pointer; transition: all 0.2s; }
  .magic-btn:hover { background: var(--primary-color); color: white; border-color: var(--primary-color); transform: translateY(-2px); }
  .magic-btn.active { background: var(--primary-color); color: white; border-color: var(--primary-color); box-shadow: 0 0 15px var(--primary-color); }

  .info-column { flex: 1; }
  .form-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 1.5rem; margin-top: 1rem; }

  @media (max-width: 850px) {
    .profile-layout { flex-direction: column; align-items: center; text-align: center; gap: 2rem; }
    .avatar-column { width: 100%; }
    .info-column { width: 100%; }
    .form-actions { justify-content: center; }
  }

  @media (max-width: 600px) {
    .form-grid { grid-template-columns: 1fr; }
    .profile-main-card { padding: 1.5rem; }
  }

  .form-actions { margin-top: 2rem; padding-top: 2rem; border-top: 1px solid var(--border-color); display: flex; align-items: center; gap: 1.5rem; }

  /* DiceBear Expansion */
  .dicebear-expansion { margin-top: 2rem; padding: 2rem; background: var(--md3-surface-1); border-radius: 2rem; border: 1px dashed var(--border-color); animation: slideDown 0.4s cubic-bezier(0.2, 0.8, 0.2, 1); }
  @keyframes slideDown { from { opacity: 0; transform: translateY(-20px); } to { opacity: 1; transform: translateY(0); } }

  .dicebear-tool { display: flex; gap: 2rem; align-items: center; }
  @media (max-width: 600px) { .dicebear-tool { flex-direction: column; text-align: center; } }
  
  .dicebear-preview { position: relative; width: 120px; height: 120px; flex-shrink: 0; }
  .preview-img { width: 100%; height: 100%; border-radius: 1.5rem; background: white; padding: 8px; box-shadow: var(--shadow-md); }
  .shuffle-btn { position: absolute; bottom: -10px; right: -10px; width: 36px; height: 36px; border-radius: 50%; background: var(--primary-color); color: white; border: 3px solid var(--md3-surface-1); display: flex; align-items: center; justify-content: center; cursor: pointer; transition: all 0.3s; }
  .shuffle-btn:hover { transform: scale(1.1) rotate(180deg); }

  .dicebear-config { flex: 1; display: flex; flex-direction: column; gap: 1.5rem; }
  .style-list { display: flex; flex-wrap: wrap; gap: 0.6rem; }
  @media (max-width: 600px) { .style-list { justify-content: center; } .btn-secondary { align-self: center; } }

  .style-pill { padding: 0.4rem 1.2rem; border-radius: 100px; background: var(--md3-surface-3); border: 1px solid var(--border-color); font-size: 0.8rem; font-weight: 700; cursor: pointer; color: var(--text-muted); transition: all 0.2s; }
  .style-pill:hover { background: var(--md3-surface-4); transform: translateY(-1px); }
  .style-pill.active { background: var(--primary-color); color: white; border-color: var(--primary-color); }

  .btn-secondary { background: var(--md3-surface-4); color: var(--text-color); border: 1px solid var(--border-color); }
  .btn-secondary:hover { background: var(--md3-surface-5); }

  .skeleton-island { height: 300px; display: flex; flex-direction: column; gap: 1rem; padding: 2rem; justify-content: center; }
  .skeleton-line { height: 40px; background: var(--md3-surface-3); border-radius: 0.5rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0%, 100% { opacity: 0.5; } 50% { opacity: 0.8; } }
</style>
