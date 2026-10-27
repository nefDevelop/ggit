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
        })
      });
      if (res.ok) {
        message = 'Profile updated successfully.';
        if (authState.user) {
          authState.user.username = profile.username;
          authState.user.full_name = profile.full_name;
          authState.user.avatar_url = profile.avatar_url;
        }
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
        message = 'Avatar deleted. Using default avatar.';
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
  import CloudUpload from "lucide-svelte/icons/cloud-upload";
  import Trash2 from "lucide-svelte/icons/trash-2";
  import WandSparkles from "lucide-svelte/icons/wand-sparkles";
  import Dices from "lucide-svelte/icons/dices";
  import Check from "lucide-svelte/icons/check";
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
        <div class="section-header-with-avatar">
          <div class="avatar-shield-header">
            <img src={`${profile.avatar_url}${profile.avatar_url.includes('?') ? '&' : '?'}t=${avatarCacheBuster}`} alt="Current Avatar" class="profile-avatar" />
            <div class="avatar-overlay">
              <label class="overlay-action-mini" title="Upload Photo">
                <CloudUpload size={18} />
                <input type="file" accept="image/*" onchange={handleAvatarUpload} hidden />
              </label>
              {#if profile.use_custom_avatar}
                <button class="overlay-action-mini delete" onclick={deleteAvatar} title="Remove Custom Avatar" aria-label="Quitar avatar personalizado">
                  <Trash2 size={18} />
                </button>
              {/if}
            </div>
          </div>
          <div class="section-info">
            <h2 class="section-title">Public Profile</h2>
            {#if profile.enable_dicebear_avatar}
              <button class="magic-btn-mini" class:active={showDiceBear} onclick={() => showDiceBear = !showDiceBear}>
                <WandSparkles size={18} />
                {showDiceBear ? 'Hide Generator' : 'Magic Avatar'}
              </button>
            {/if}
          </div>
        </div>
        
        <div class="info-column">
            
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

        {#if showDiceBear && profile.enable_dicebear_avatar}
          <div class="dicebear-expansion">
            <div class="dicebear-tool">
              <div class="dicebear-preview">
                <img src={dicebearUrl} alt="DiceBear Preview" class="preview-img" />
                <button class="shuffle-btn" onclick={shuffleDicebear} title="Shuffle Seed" aria-label="Cambiar semilla">
                  <Dices size={18} />
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
                  <Check size={18} /> Use this Avatar
                </button>
              </div>
            </div>
          </div>
        {/if}
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


  .success-message { color: var(--status-success); font-size: 0.9rem; margin-top: 1rem; font-weight: 600; }
  .error-message { color: var(--status-danger-soft); font-size: 0.9rem; margin-top: 1rem; font-weight: 600; }

  /* New Profile Card Layout */
  /* New Profile Card Layout */
  .profile-main-card { padding: 2.5rem; }
  .section-header-with-avatar { display: flex; align-items: center; gap: 1.5rem; margin-bottom: 2.5rem; padding-bottom: 2rem; border-bottom: 1px solid var(--border-color); }
  
  .avatar-shield-header { position: relative; width: 80px; height: 80px; border-radius: 2rem; overflow: hidden; border: 3px solid var(--border-color); background: var(--md3-surface-1); flex-shrink: 0; }
  .profile-avatar { width: 100%; height: 100%; object-fit: cover; }
  
  .avatar-overlay { position: absolute; inset: 0; background: rgba(0,0,0,0.7); opacity: 0; transition: all 0.3s ease; display: flex; align-items: center; justify-content: center; gap: 0.5rem; backdrop-filter: blur(2px); }
  .avatar-shield-header:hover .avatar-overlay { opacity: 1; }

  .overlay-action-mini { width: 32px; height: 32px; border-radius: 50%; background: white; color: var(--primary-color); display: flex; align-items: center; justify-content: center; cursor: pointer; border: none; transition: transform 0.2s; }
  .overlay-action-mini:hover { transform: scale(1.1); }
  .overlay-action-mini.delete { color: var(--status-danger); }

  .section-info { flex: 1; display: flex; align-items: center; justify-content: space-between; gap: 1rem; }
  .section-title { margin-bottom: 0 !important; }

  .magic-btn-mini { display: flex; align-items: center; gap: 0.5rem; background: var(--md3-surface-3); border: 1px solid var(--border-color); color: var(--text-color); padding: 0.5rem 1rem; border-radius: 100px; font-weight: 700; font-size: 0.8rem; cursor: pointer; transition: all 0.2s; }
  .magic-btn-mini:hover { background: var(--primary-color); color: white; border-color: var(--primary-color); }
  .magic-btn-mini.active { background: var(--primary-color); color: white; border-color: var(--primary-color); }

  .info-column { width: 100%; }
  .form-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 1.5rem; }

  @media (max-width: 850px) {
    .section-header-with-avatar { flex-direction: column; text-align: center; gap: 1.5rem; }
    .section-info { flex-direction: column; }
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
