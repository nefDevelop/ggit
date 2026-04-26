<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from "svelte";
  import { authState } from "$lib/auth.svelte";
  import Avatar from "./Avatar.svelte";
  import { i18n } from "$lib/i18n.svelte";
  import Card from "./ui/Card.svelte";
  import Button from "./ui/Button.svelte";

  let profile = $state<any>({
    username: "",
    full_name: "",
    email: "",
    website: "",
    location: "",
    avatar_url: "",
    use_custom_avatar: false,
  });

  let loading = $state(true);
  let saving = $state(false);

  // DiceBear Logic
  const dicebearStyles = [
    "adventurer",
    "adventurer-neutral",
    "avataaars",
    "avataaars-neutral",
    "bottts",
    "bottts-neutral",
    "icons",
    "lorelei",
    "lorelei-neutral",
    "micah",
    "open-peeps",
    "notionists",
    "pixel-art",
    "shapes",
    "thumbs",
  ];
  let currentDicebearStyle = $state("avataaars");
  let dicebearSeed = $state("");
  let showDiceBear = $state(false);
  const dicebearUrl = $derived(
    `https://api.dicebear.com/7.x/${currentDicebearStyle}/png?seed=${dicebearSeed || profile.username || 'ggit'}`,
  );

  async function fetchProfile() {
    loading = true;
    try {
      const res = await apiFetch("/api/v2/user/settings/profile");
      if (res.ok) {
        const data = await res.json();
        profile = data;
      }
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function saveProfile() {
    saving = true;
    try {
      const res = await apiFetch("/api/v2/user/settings/profile", {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          full_name: profile.full_name,
          website: profile.website,
          location: profile.location,
        }),
      });
      if (res.ok) {
        // Sync with global auth state if needed
        if (authState.user) {
          authState.user.full_name = profile.full_name;
        }
      }
    } catch (e) {
      console.error(e);
    } finally {
      saving = false;
    }
  }

  async function handleAvatarUpload(e: Event) {
    const file = (e.target as HTMLInputElement).files?.[0];
    if (!file) return;

    const formData = new FormData();
    formData.append("avatar", file);

    try {
      const res = await apiFetch("/api/v2/user/settings/avatar", {
        method: "POST",
        body: formData,
      });
      if (res.ok) {
        const data = await res.json();
        profile.avatar_url = data.avatar_url + "?cb=" + Date.now();
        if (authState.user) authState.user.avatar_url = profile.avatar_url;
      }
    } catch (err) {
      console.error(err);
    }
  }

  async function deleteAvatar() {
    try {
      const res = await apiFetch("/api/v2/user/settings/avatar", {
        method: "DELETE",
      });
      if (res.ok) {
        const data = await res.json();
        profile.avatar_url = data.avatar_url + "?cb=" + Date.now();
        if (authState.user) authState.user.avatar_url = profile.avatar_url;
      }
    } catch (err) {
      console.error(err);
    }
  }

  async function saveDicebearAvatar() {
    saving = true;
    try {
      const imgRes = await fetch(dicebearUrl);
      const blob = await imgRes.blob();

      const formData = new FormData();
      const file = new File([blob], 'avatar.png', { type: 'image/png' });
      formData.append('avatar', file);

      const res = await apiFetch("/api/v2/user/settings/avatar/magic", {
        method: "POST",
        body: formData,
      });
      if (res.ok) {
        const data = await res.json();
        profile.avatar_url = data.avatar_url + "?cb=" + Date.now();
        if (authState.user) authState.user.avatar_url = profile.avatar_url;
        showDiceBear = false;
      }
    } catch (err) {
      console.error(err);
    } finally {
      saving = false;
    }
  }

  onMount(fetchProfile);
</script>

<Card class="settings-premium-card" padding="lg">
  <header class="section-header">
    <div class="header-flex">
      <div class="title-with-icon">
        <span class="material-symbols-rounded hero-icon">person</span>
        <div>
          <h2>{i18n.t("settings.public_profile")}</h2>
          <p class="subtitle">{i18n.t("settings.public_profile_desc")}</p>
        </div>
      </div>
      <Button
        variant="ghost"
        icon="magic_button"
        onclick={() => (showDiceBear = !showDiceBear)}
      >
        {showDiceBear ? "Cerrar selector" : "Avatar Mágico"}
      </Button>
    </div>
  </header>



  <div class="profile-layout-grid">
    <!-- AVATAR SIDE -->
    <div class="avatar-column">
      <div class="avatar-shield">
        <Avatar
          src={profile.avatar_url}
          username={profile.username}
          size={200}
        />
        <div class="avatar-overlay">
          <label class="overlay-btn" title="Subir foto">
            <span class="material-symbols-rounded">cloud_upload</span>
            <input
              type="file"
              accept="image/*"
              onchange={handleAvatarUpload}
              hidden
            />
          </label>
          {#if profile.avatar_url && !profile.avatar_url.includes("magic")}
            <button
              class="overlay-btn delete"
              onclick={deleteAvatar}
              title="Quitar foto"
            >
              <span class="material-symbols-rounded">delete</span>
            </button>
          {/if}
        </div>
      </div>
      <span class="username-tag">@{profile.username}</span>
    </div>

    <!-- FORM SIDE -->
    <div class="form-column">
      {#if loading}
        <div class="skeleton-form"></div>
      {:else}
        <form
          class="settings-form"
          onsubmit={(e) => {
            e.preventDefault();
            saveProfile();
          }}
        >
          <div class="form-grid">
            <div class="field">
              <label for="username">{i18n.t("username")}</label>
              <input
                id="username"
                type="text"
                bind:value={profile.username}
                disabled
              />
              <p class="hint">{i18n.t("settings.username_helper")}</p>
            </div>

            <div class="field">
              <label for="full_name">{i18n.t("settings.full_name")}</label>
              <input
                id="full_name"
                type="text"
                bind:value={profile.full_name}
                placeholder="Tu nombre real"
              />
            </div>

            <div class="field">
              <label for="email">{i18n.t("email")}</label>
              <input
                id="email"
                type="email"
                bind:value={profile.email}
                disabled
              />
            </div>

            <div class="field">
              <label for="website">{i18n.t("website")}</label>
              <input
                id="website"
                type="url"
                bind:value={profile.website}
                placeholder="https://tuweb.com"
              />
            </div>

            <div class="field full-width">
              <label for="location">{i18n.t("settings.location")}</label>
              <input
                id="location"
                type="text"
                bind:value={profile.location}
                placeholder="Ciudad, País"
              />
            </div>
          </div>

          <div class="form-footer">
            <Button type="submit" variant="primary" disabled={saving}>
              {saving ? i18n.t("settings.saving") : i18n.t("settings.update_profile")}
            </Button>
          </div>
        </form>
      {/if}
    </div>
  </div>

  {#if showDiceBear}
    <div class="magic-avatar-panel island tonal" style="margin-top: 4rem;">
      <div class="dicebear-preview">
        <div class="preview-shield">
          <img src={dicebearUrl} alt="DiceBear Preview" />
          <button
            class="shuffle-btn"
            onclick={() =>
              (dicebearSeed = Math.random().toString(36).substring(7))}
          >
            <span class="material-symbols-rounded">casino</span>
          </button>
        </div>
      </div>
      <div class="dicebear-settings">
        <div class="styles-grid">
          {#each dicebearStyles as style}
            <button
              class="style-chip"
              class:active={currentDicebearStyle === style}
              onclick={() => (currentDicebearStyle = style)}
            >
              {style.replace("-", " ")}
            </button>
          {/each}
        </div>
        <Button variant="primary" onclick={saveDicebearAvatar} disabled={saving}>
          Confirmar Avatar Mágico
        </Button>
      </div>
    </div>
  {/if}
</Card>

<style>
  .section-header {
    margin-bottom: 4rem;
  }
  .header-flex {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 1rem;
    flex-wrap: wrap;
    max-width: 100%;
    box-sizing: border-box;
  }
  .title-with-icon {
    display: flex;
    align-items: center;
    gap: 2rem;
    min-width: 0;
  }
  .hero-icon {
    font-size: 3rem;
    color: var(--primary-color);
    background: rgba(var(--primary-rgb), 0.1);
    padding: 1.25rem;
    border-radius: 2rem;
  }

  h2 {
    margin: 0;
    font-size: 2.25rem;
    font-weight: 950;
    letter-spacing: -0.04em;
  }
  .subtitle {
    margin: 0.5rem 0 0 0;
    color: var(--text-muted);
    font-size: 1.1rem;
  }

  @media (max-width: 600px) {
    .section-header { margin-bottom: 2rem; }
    h2 { font-size: 1.5rem; }
    .hero-icon { font-size: 1.5rem; padding: 0.6rem; border-radius: 1rem; }
    .title-with-icon { gap: 0.75rem; }
    .header-flex { flex-direction: column; align-items: flex-start; gap: 1rem; }
    .subtitle { font-size: 0.9rem; }
  }

  .profile-layout-grid {
    display: grid;
    grid-template-columns: 280px 1fr;
    gap: 4rem;
    max-width: 100%;
    width: 100%;
    box-sizing: border-box;
  }
  @media (max-width: 1200px) {
    .profile-layout-grid {
      grid-template-columns: 1fr;
      gap: 3rem;
    }
    .avatar-column { order: -1; }
  }
  @media (max-width: 600px) {
    .profile-layout-grid {
      gap: 1.5rem;
      padding: 0;
    }
  }

  .avatar-column {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1.5rem;
    max-width: 100%;
    box-sizing: border-box;
  }
  .avatar-shield {
    position: relative;
    width: 200px;
    height: 200px;
    max-width: 100%;
    border-radius: 5rem;
    overflow: hidden;
    border: 4px solid var(--md-sys-color-outline-variant);
    box-shadow: var(--shadow-xl);
    box-sizing: border-box;
    display: flex;
    align-items: center;
    justify-content: center;
  }
  @media (max-width: 600px) {
    .avatar-shield { width: 120px; height: 120px; border-radius: 3rem; }
  }

  .avatar-overlay {
    position: absolute;
    inset: 0;
    background: rgba(0, 0, 0, 0.6);
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 1rem;
    opacity: 0;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    backdrop-filter: blur(4px);
  }
  .avatar-shield:hover .avatar-overlay {
    opacity: 1;
  }

  @media (max-width: 1024px) {
    .avatar-overlay { opacity: 1; background: rgba(0,0,0,0.3); align-items: flex-end; padding-bottom: 1rem; }
    .overlay-btn { width: 40px; height: 40px; }
  }

  .overlay-btn {
    width: 48px;
    height: 48px;
    border-radius: 50%;
    background: white;
    color: #000;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    border: none;
    transition: transform 0.2s;
  }
  .overlay-btn:hover {
    transform: scale(1.1);
  }
  .overlay-btn.delete {
    color: var(--error);
  }

  .username-tag {
    font-size: 1.1rem;
    font-weight: 800;
    color: var(--primary-color);
    background: rgba(var(--primary-rgb), 0.1);
    padding: 0.4rem 1rem;
    border-radius: 100px;
    max-width: 100%;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    box-sizing: border-box;
  }

  .settings-form {
    display: flex;
    flex-direction: column;
    gap: 3rem;
    width: 100%;
    max-width: 100%;
    box-sizing: border-box;
  }
  .form-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(min(100%, 300px), 1fr));
    gap: 2rem;
    width: 100%;
    box-sizing: border-box;
  }
  @media (max-width: 800px) {
    .form-grid { grid-template-columns: 1fr; gap: 1.5rem; }
    .field.full-width { grid-column: auto; }
    .settings-form { gap: 2rem; }
  }
  .field {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    width: 100%;
    max-width: 100%;
    box-sizing: border-box;
  }
  .field.full-width {
    grid-column: span 2;
  }
  .field label {
    font-size: 0.9rem;
    font-weight: 700;
    color: var(--text-color);
    margin-left: 0.5rem;
  }
  .field input {
    background: var(--surface-variant);
    border: 2px solid transparent;
    padding: 1rem 1.5rem;
    border-radius: 1.25rem;
    color: var(--text-color);
    font-size: 1.05rem;
    transition: all 0.3s;
    width: 100%;
    max-width: 100%;
    box-sizing: border-box;
  }
  @media (max-width: 600px) {
    .field input {
      padding: 0.75rem 1rem;
    }
    .title-with-icon {
      gap: 1rem;
    }
  }

  .field input:focus {
    border-color: var(--primary-color);
    outline: none;
    background: rgba(var(--primary-rgb), 0.05);
    box-shadow: 0 0 0 5px rgba(var(--primary-rgb), 0.15);
  }
  .field input:disabled {
    opacity: 0.5;
    cursor: not-allowed;
  }
  .hint {
    font-size: 0.75rem;
    color: var(--text-muted);
    margin: 0 0 0 0.5rem;
  }

  .form-footer {
    display: flex;
    justify-content: flex-end;
    padding-top: 2.5rem;
    border-top: 1px solid var(--md-sys-color-outline-variant);
  }
  @media (max-width: 600px) {
    .form-footer {
      padding-top: 1.5rem;
    }
    .form-footer :global(button) {
      width: 100%;
    }
  }

  .magic-avatar-panel {
    margin-bottom: 3rem;
    padding: 2.5rem !important;
    border-radius: 2.5rem !important;
    display: grid;
    grid-template-columns: auto 1fr;
    gap: 3rem;
    border: 2px solid var(--primary-color) !important;
    background: rgba(var(--primary-rgb), 0.05) !important;
    animation: slideDown 0.4s cubic-bezier(0.4, 0, 0.2, 1);
    max-width: 900px;
    width: 100%;
    box-sizing: border-box;
  }
  @media (max-width: 1000px) {
    .magic-avatar-panel { 
      grid-template-columns: 1fr; 
      padding: 1.5rem !important; 
      gap: 2rem; 
      border-radius: 1.5rem !important; 
      text-align: center;
    }
  }
  @keyframes slideDown {
    from { opacity: 0; transform: translateY(-20px); }
    to { opacity: 1; transform: translateY(0); }
  }

  .dicebear-preview { display: flex; align-items: center; justify-content: center; }
  .preview-shield { position: relative; width: 140px; height: 140px; border-radius: 2rem; overflow: hidden; background: white; border: 4px solid var(--md-sys-color-outline-variant); }
  .preview-shield img { width: 100%; height: 100%; }
  .shuffle-btn { position: absolute; bottom: 0.5rem; right: 0.5rem; background: var(--primary-color); color: white; border: none; padding: 0.5rem; border-radius: 0.75rem; cursor: pointer; }

  .dicebear-settings { display: flex; flex-direction: column; gap: 1.5rem; }
  .styles-grid { display: flex; flex-wrap: wrap; gap: 0.5rem; }
  @media (max-width: 600px) {
    .styles-grid { justify-content: center; }
  }
  .style-chip { padding: 0.5rem 1rem; border-radius: 100px; border: 1px solid var(--md-sys-color-outline-variant); background: transparent; color: var(--text-muted); font-weight: 600; cursor: pointer; font-size: 0.85rem; text-transform: capitalize; }
  .style-chip.active { background: var(--primary-color); color: white; border-color: var(--primary-color); }

  .skeleton-form { height: 400px; background: var(--surface-variant); border-radius: 2rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0% { opacity: 0.3; } 50% { opacity: 0.6; } 100% { opacity: 0.3; } }
</style>
