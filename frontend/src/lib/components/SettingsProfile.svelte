<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from "svelte";
  import { authState } from "$lib/auth.svelte";
  import Avatar from "./Avatar.svelte";
  import Checkbox from "./ui/Checkbox.svelte";
  import Input from "./ui/Input.svelte";
  import { i18n } from "$lib/i18n.svelte";
  import Card from "./ui/Card.svelte";
  import Button from "./ui/Button.svelte";
  import { getAvailableLanguages } from '$lib/locales';
  import PageHeader from "./ui/PageHeader.svelte";
  

import Upload from "lucide-svelte/icons/upload";
import Trash2 from "lucide-svelte/icons/trash-2";
import Wand2 from "lucide-svelte/icons/wand-2";
import ChevronDown from "lucide-svelte/icons/chevron-down";
import Dices from "lucide-svelte/icons/dices";



  const availableLanguages = getAvailableLanguages();

  function changeLanguage(event: Event) {
    const select = event.currentTarget as HTMLSelectElement;
    i18n.init(select.value);
  }

  let profile = $state<any>({
    username: "",
    full_name: "",
    email: "",
    website: "",
    location: "",
    avatar_email: "",
    disable_news_feed: false,
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
          avatar_email: profile.avatar_email,
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

<div class="settings-profile-container">
  <div class="profile-content-layout">
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
          <div class="avatar-form-row">
            <div class="avatar-shield-header">
              <Avatar
                src={profile.avatar_url}
                username={profile.username}
                size={80}
              />
              <div class="avatar-overlay">
                <label class="overlay-btn-mini" title="Subir foto">
                  <Upload size={18} />
                  <input
                    type="file"
                    accept="image/*"
                    onchange={handleAvatarUpload}
                    hidden
                  />
                </label>
                {#if profile.avatar_url && !profile.avatar_url.includes("magic")}
                    <button
                      class="overlay-btn-mini delete"
                      onclick={deleteAvatar}
                      title="Quitar foto"
                    >
                      <Trash2 size={18} />
                    </button>
                {/if}
              </div>
            </div>
            <div class="avatar-actions">
              <Button
                variant="ghost"
                icon={Wand2}
                onclick={() => (showDiceBear = !showDiceBear)}
              >
                <span class="avatar-magic-text">{showDiceBear ? "Cerrar selector" : "Avatar Mágico"}</span>
              </Button>
            </div>
          </div>
          <div class="form-grid">
            <Input label={i18n.t("username")} bind:value={profile.username} disabled hint={i18n.t("settings.username_helper")} />
            <Input label={i18n.t("settings.full_name")} bind:value={profile.full_name} placeholder="Tu nombre real" />
            <Input label={i18n.t("email")} type="email" bind:value={profile.email} disabled />
            <Input label={i18n.t("website")} type="url" bind:value={profile.website} placeholder="https://tuweb.com" />
            <Input label={i18n.t("settings.location")} bind:value={profile.location} placeholder="Ciudad, País" />
            <Input label="Correo de Avatar (Gravatar)" type="email" bind:value={profile.avatar_email} placeholder="email@example.com" />
            <div class="field">
              <label for="lang-select">{i18n.t("settings.language") || "Language"}</label>
              <div class="select-wrapper">
                <select id="lang-select" onchange={changeLanguage} value={i18n.locale}>
                    {#each availableLanguages as lang}
                        <option value={lang.code}>{lang.name}</option>
                    {/each}
                </select>
                <ChevronDown size={18} class="dropdown-icon" />
              </div>
              <p class="hint">{i18n.t("settings.language_helper") || "Choose your preferred language for the interface."}</p>
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
            <Dices size={18} />
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
</div>

<style>
  .avatar-form-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    gap: 1.5rem;
    margin-bottom: 2.5rem;
    padding-bottom: 2rem;
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
  }


  .avatar-shield-header {
    position: relative;
    width: 80px;
    height: 80px;
    border-radius: 1.5rem;
    overflow: hidden;
    border: 3px solid var(--md-sys-color-outline-variant);
    box-shadow: var(--shadow-lg);
    background: var(--md-sys-color-secondary-container);
    flex-shrink: 0;
  }
  .avatar-shield-header:hover .avatar-overlay {
    opacity: 1;
  }

  .avatar-overlay {
    position: absolute;
    inset: 0;
    background: rgba(0, 0, 0, 0.7);
    display: flex;
    align-items: center;
    justify-content: center;
    gap: 0.5rem;
    opacity: 0;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    backdrop-filter: blur(4px);
  }

  .overlay-btn-mini {
    width: 32px;
    height: 32px;
    border-radius: 50%;
    background: white;
    color: var(--primary-color);
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    border: none;
    transition: transform 0.2s;
  }
  .overlay-btn-mini:hover {
    transform: scale(1.1);
  }
  .overlay-btn-mini.delete {
    color: var(--error);
  }





  .profile-content-layout {
    width: 100%;
    box-sizing: border-box;
    padding: 2.5rem;
  }
  @media (max-width: 600px) {
    .profile-content-layout {
      padding: 1.25rem 1rem;
    }
    .avatar-magic-text {
      display: none;
    }
    .settings-form { gap: 1.5rem; }
    .form-grid { gap: 1rem; }
    .avatar-form-row { margin-bottom: 1.5rem; padding-bottom: 1.5rem; }
    .magic-avatar-panel { padding: 1rem !important; gap: 1rem; }
  }

  @media (max-width: 480px) {
    .profile-content-layout { padding: 0.75rem 0.5rem; }
    .settings-form { gap: 1rem; }
    .form-grid { gap: 0.75rem; }
    .field input { padding: 0.6rem 0.85rem; font-size: 1rem; border-radius: 1rem; }
    .select-wrapper select { padding: 0.6rem 0.85rem; font-size: 1rem; border-radius: 1rem; }
    .avatar-form-row { flex-direction: column; text-align: center; gap: 1rem; }
    .form-footer { padding-top: 1rem; }
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
  .field label {
    font-size: 0.9rem;
    font-weight: 700;
    color: var(--text-color);
    margin-left: 0.5rem;
  }
  .select-wrapper {
    position: relative;
    width: 100%;
  }
  .select-wrapper select {
    appearance: none;
    width: 100%;
    background: var(--md-sys-color-surface-container-high);
    border: 1px solid var(--md-sys-color-outline-variant);
    padding: 0.75rem 1rem;
    border-radius: 0.75rem;
    color: var(--text-color);
    font-size: 1rem;
    transition: all 0.2s;
    cursor: pointer;
    box-sizing: border-box;
  }
  .select-wrapper select:focus {
    outline: none;
    border-color: var(--md-sys-color-primary);
    background: var(--md-sys-color-surface-container-highest);
    box-shadow: 0 0 0 3px color-mix(in srgb, var(--md-sys-color-primary), transparent 85%);
  }
  :global(.dropdown-icon) {
    position: absolute;
    right: 1rem;
    top: 50%;
    transform: translateY(-50%);
    color: var(--text-muted);
    pointer-events: none;
    transition: transform 0.3s;
  }
  .select-wrapper select:focus + :global(.dropdown-icon) {
    transform: translateY(-50%) rotate(180deg);
    color: var(--primary-color);
  }
  .hint {
    font-size: 0.75rem;
    color: var(--text-muted);
    margin: 0 0 0 0.5rem;
  }

  .form-footer {
    display: flex;
    justify-content: flex-end;
    padding-top: 2rem;
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
    border-radius: 1.5rem !important;
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
  .preview-shield { position: relative; width: 140px; height: 140px; border-radius: 1.5rem; overflow: hidden; background: white; border: 4px solid var(--md-sys-color-outline-variant); }
  .preview-shield img { width: 100%; height: 100%; }
  .shuffle-btn { position: absolute; bottom: 0.5rem; right: 0.5rem; background: var(--primary-color); color: white; border: none; padding: 0.5rem; border-radius: 0.75rem; cursor: pointer; }

  .dicebear-settings { display: flex; flex-direction: column; gap: 1.5rem; }
  .styles-grid { display: flex; flex-wrap: wrap; gap: 0.5rem; }
  @media (max-width: 600px) {
    .styles-grid { justify-content: center; }
  }
  .style-chip { padding: 0.5rem 1rem; border-radius: 100px; border: 1px solid var(--md-sys-color-outline-variant); background: transparent; color: var(--text-muted); font-weight: 600; cursor: pointer; font-size: 0.85rem; text-transform: capitalize; }
  .style-chip.active { background: var(--primary-color); color: white; border-color: var(--primary-color); }

  .skeleton-form { height: 400px; background: var(--surface-variant); border-radius: 1.5rem; animation: pulse 1.5s infinite; }
  @keyframes pulse { 0% { opacity: 0.3; } 50% { opacity: 0.6; } 100% { opacity: 0.3; } }
</style>
