<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from 'svelte';
  import { authState } from '$lib/auth.svelte';
  
  let profile = $state<any>({
    username: '',
    full_name: '',
    website: '',
    location: '',
    avatar_url: '',
    use_custom_avatar: false,
    enable_dicebear_avatar: false
  });

  let loading = $state(true);
  let saving = $state(false);
  let message = $state('');
  let error = $state('');
  let avatarCacheBuster = $state(Date.now());

  // DiceBear Logic
  const dicebearStyles = ['avataaars', 'bottts', 'pixel-art', 'lorelei-neutral', 'micah', 'open-peeps', 'notionists'];
  let currentDicebearStyle = $state('avataaars');
  let dicebearSeed = $state('');
  let showDiceBear = $state(false);
  const dicebearUrl = $derived(`https://api.dicebear.com/7.x/${currentDicebearStyle}/png?seed=${dicebearSeed || profile.username}`);

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
    error = '';
    try {
      const res = await apiFetch('/api/v2/user/settings/profile', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          name: profile.username,
          full_name: profile.full_name,
          website: profile.website,
          location: profile.location
        })
      });
      if (res.ok) {
        message = 'Profile updated successfully!';
        if (authState.user) {
          authState.user.username = profile.username;
          authState.user.full_name = profile.full_name;
        }
        setTimeout(() => message = '', 3000);
      }
    } catch (e) {
      error = 'Failed to update profile.';
    } finally {
      saving = false;
    }
  }

  async function handleAvatarUpload(e: Event) {
    const file = (e.target as HTMLInputElement).files?.[0];
    if (!file) return;
    saving = true;
    try {
      const formData = new FormData();
      formData.append('avatar', file);
      const res = await apiFetch('/api/v2/user/settings/avatar', {
        method: 'POST',
        body: formData
      });
      if (res.ok) {
        message = 'Avatar updated!';
        avatarCacheBuster = Date.now();
        await fetchProfile();
      }
    } catch (e) {
      error = 'Upload failed.';
    } finally {
      saving = false;
    }
  }

  async function deleteAvatar() {
    saving = true;
    try {
      const res = await apiFetch('/api/v2/user/settings/avatar', { method: 'DELETE' });
      if (res.ok) {
        message = 'Avatar removed.';
        avatarCacheBuster = Date.now();
        fetchProfile();
      }
    } catch (e) {
      error = 'Failed to delete avatar.';
    } finally {
      saving = false;
    }
  }

  async function saveDicebearAvatar() {
    saving = true;
    try {
      const imgRes = await fetch(dicebearUrl);
      const blob = await imgRes.blob();
      const formData = new FormData();
      formData.append('avatar', new File([blob], 'avatar.png', { type: 'image/png' }));
      const res = await apiFetch('/api/v2/user/settings/avatar', { method: 'POST', body: formData });
      if (res.ok) {
        message = 'Magic avatar saved!';
        avatarCacheBuster = Date.now();
        await fetchProfile();
      }
    } catch (e) { error = 'Failed to save magic avatar.'; }
    finally { saving = false; }
  }

  onMount(fetchProfile);
</script>

<div class="profile-modular-card island">
  <div class="profile-layout">
    <!-- AVATAR COLUMN -->
    <div class="avatar-sidebar">
      <div class="avatar-shield">
        <img src={`${profile.avatar_url}${profile.avatar_url?.includes('?') ? '&' : '?'}t=${avatarCacheBuster}`} alt="Avatar" class="main-avatar" />
        <div class="avatar-actions-overlay">
          <label class="action-btn" title="Upload New">
            <span class="material-symbols-rounded">cloud_upload</span>
            <input type="file" accept="image/*" onchange={handleAvatarUpload} hidden />
          </label>
          {#if profile.use_custom_avatar}
            <button class="action-btn delete" onclick={deleteAvatar} title="Remove">
              <span class="material-symbols-rounded">delete</span>
            </button>
          {/if}
        </div>
      </div>

      {#if profile.enable_dicebear_avatar}
        <button class="btn btn-tonal magic-trigger" class:active={showDiceBear} onclick={() => showDiceBear = !showDiceBear}>
          <span class="material-symbols-rounded">auto_fix_high</span>
          Magic Avatar
        </button>
      {/if}
    </div>

    <!-- INFO COLUMN -->
    <div class="info-content">
      <header class="content-header">
        <h3 class="title">Public Profile</h3>
        <p class="subtitle">Update your personal information and identity.</p>
      </header>

      {#if loading}
        <div class="skeleton-form"></div>
      {:else}
        <form class="precision-form" onsubmit={(e) => { e.preventDefault(); saveProfile(); }}>
          <div class="form-grid">
            <div class="form-group">
              <label for="username">Username</label>
              <input type="text" id="username" bind:value={profile.username} placeholder="Username" />
              <p class="hint">Your unique identifier on Ggit.</p>
            </div>

            <div class="form-group">
              <label for="full_name">Full Name</label>
              <input type="text" id="full_name" bind:value={profile.full_name} placeholder="Your display name" />
            </div>

            <div class="form-group">
              <label for="website">Website</label>
              <input type="url" id="website" bind:value={profile.website} placeholder="https://..." />
            </div>

            <div class="form-group">
              <label for="location">Location</label>
              <input type="text" id="location" bind:value={profile.location} placeholder="City, Country" />
            </div>
          </div>

          <div class="form-footer">
            {#if message}<span class="status success">{message}</span>{/if}
            {#if error}<span class="status error">{error}</span>{/if}
            <button type="submit" class="btn btn-action-primary" disabled={saving}>
              {saving ? 'Saving...' : 'Save Changes'}
            </button>
          </div>
        </form>
      {/if}
    </div>
  </div>

  {#if showDiceBear}
    <div class="dicebear-panel">
       <div class="dicebear-preview-box">
          <img src={dicebearUrl} alt="Preview" />
          <button class="shuffle-icon" onclick={() => dicebearSeed = Math.random().toString(36).substring(7)}>
            <span class="material-symbols-rounded">casino</span>
          </button>
       </div>
       <div class="dicebear-controls">
         <div class="styles-row">
           {#each dicebearStyles as style}
             <button class="style-item" class:active={currentDicebearStyle === style} onclick={() => currentDicebearStyle = style}>
               {style}
             </button>
           {/each}
         </div>
         <button class="btn btn-primary btn-sm" onclick={saveDicebearAvatar} disabled={saving}>Confirm Magic Avatar</button>
       </div>
    </div>
  {/if}
</div>

<style>
  .profile-modular-card { padding: 0; overflow: hidden; }
  .profile-layout { display: flex; gap: 4rem; padding: 3rem; }

  /* Avatar Sidebar */
  .avatar-sidebar { display: flex; flex-direction: column; align-items: center; gap: 2rem; width: 180px; }
  .avatar-shield { position: relative; width: 160px; height: 160px; border-radius: 4rem; overflow: hidden; background: var(--surface-variant); border: 4px solid var(--border-color); box-shadow: var(--shadow-lg); }
  .main-avatar { width: 100%; height: 100%; object-fit: cover; }
  .avatar-actions-overlay { position: absolute; inset: 0; background: rgba(0,0,0,0.6); display: flex; align-items: center; justify-content: center; gap: 0.75rem; opacity: 0; transition: opacity 0.2s; backdrop-filter: blur(4px); }
  .avatar-shield:hover .avatar-actions-overlay { opacity: 1; }
  .action-btn { width: 40px; height: 40px; border-radius: 50%; background: white; color: var(--primary-color); display: flex; align-items: center; justify-content: center; border: none; cursor: pointer; transition: transform 0.2s; }
  .action-btn:hover { transform: scale(1.1); }
  .action-btn.delete { color: var(--error); }
  .magic-trigger { border-radius: 100px; font-weight: 700; width: 100%; }

  /* Info Content */
  .info-content { flex: 1; }
  .content-header { margin-bottom: 2.5rem; }
  .title { margin: 0; font-size: 1.5rem; font-weight: 900; }
  .subtitle { margin: 0.5rem 0 0 0; color: var(--text-muted); font-size: 0.95rem; }

  .precision-form { display: flex; flex-direction: column; gap: 2rem; }
  .form-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 2rem; }
  .form-group { display: flex; flex-direction: column; gap: 0.5rem; }
  .form-group label { font-size: 0.75rem; font-weight: 800; text-transform: uppercase; color: var(--text-muted); letter-spacing: 0.05em; }
  .form-group input { background: var(--surface-variant); border: 1px solid var(--border-color); padding: 0.85rem 1.25rem; border-radius: 1rem; color: var(--text-color); font-size: 1rem; transition: all 0.2s; }
  .form-group input:focus { border-color: var(--primary-color); outline: none; background: rgba(var(--primary-rgb), 0.05); }
  .hint { margin: 0.25rem 0 0 0; font-size: 0.75rem; color: var(--text-muted); opacity: 0.7; }

  .form-footer { display: flex; justify-content: flex-end; align-items: center; gap: 2rem; margin-top: 1rem; padding-top: 2rem; border-top: 1px solid var(--border-color); }
  .status { font-size: 0.9rem; font-weight: 600; }
  .status.success { color: var(--success); }
  .status.error { color: var(--error); }

  /* DiceBear Panel */
  .dicebear-panel { background: var(--surface-variant); padding: 2rem; border-top: 1px solid var(--border-color); display: flex; gap: 3rem; align-items: center; animation: slideIn 0.4s ease-out; }
  @keyframes slideIn { from { opacity: 0; transform: translateY(20px); } to { opacity: 1; transform: translateY(0); } }
  .dicebear-preview-box { position: relative; width: 120px; height: 120px; border-radius: 2rem; background: white; padding: 10px; box-shadow: var(--shadow-md); flex-shrink: 0; }
  .dicebear-preview-box img { width: 100%; height: 100%; }
  .shuffle-icon { position: absolute; bottom: -8px; right: -8px; width: 32px; height: 32px; border-radius: 50%; background: var(--primary-color); color: white; border: 2px solid white; display: flex; align-items: center; justify-content: center; cursor: pointer; }
  .dicebear-controls { flex: 1; display: flex; flex-direction: column; gap: 1.5rem; }
  .styles-row { display: flex; flex-wrap: wrap; gap: 0.5rem; }
  .style-item { padding: 0.4rem 1rem; border-radius: 100px; background: rgba(255,255,255,0.05); border: 1px solid var(--border-color); font-size: 0.8rem; font-weight: 600; cursor: pointer; color: var(--text-muted); }
  .style-item.active { background: var(--primary-color); color: white; border-color: var(--primary-color); }

  @media (max-width: 900px) {
    .profile-layout { flex-direction: column; align-items: center; gap: 3rem; }
    .avatar-sidebar { width: 100%; }
    .form-grid { grid-template-columns: 1fr; }
    .dicebear-panel { flex-direction: column; }
  }
</style>
