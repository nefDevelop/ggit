<script lang="ts">
  import ThemeSelector from "./ThemeSelector.svelte";
  import SettingsProfile from "./SettingsProfile.svelte";
  import SettingsSecurity from "./SettingsSecurity.svelte";
  import SettingsEmails from "./SettingsEmails.svelte";
  import SettingsKeys from "./SettingsKeys.svelte";
  import SettingsTokens from "./SettingsTokens.svelte";
  import SettingsRepositories from "./SettingsRepositories.svelte";
  import SettingsOrganizations from "./SettingsOrganizations.svelte";
  import SettingsMyGgit from "./SettingsMyGgit.svelte";
  import { i18n } from "$lib/i18n.svelte";

  let activeTab = $state("profile");

  const tabs = [
    { id: "profile", label: i18n.t("settings.profile"), icon: "person" },
    { id: "myggit", label: i18n.t("settings.my_ggit"), icon: "grid_view" },
    { id: "appearance", label: i18n.t("settings.appearance"), icon: "palette" },
    { id: "emails", label: i18n.t("settings.emails"), icon: "mail" },
    { id: "security", label: i18n.t("settings.security"), icon: "security" },
    { id: "keys", label: i18n.t("settings.ssh_keys"), icon: "key" },
    { id: "tokens", label: i18n.t("settings.access_tokens"), icon: "token" },
    { id: "repos", label: i18n.t("settings.repositories"), icon: "book" },
    { id: "orgs", label: i18n.t("settings.organizations"), icon: "corporate_fare" },
  ];
</script>

<div class="settings-premium-container">
  <aside class="settings-sidebar">
    <div class="sidebar-brand">
       <span class="material-symbols-rounded">manage_accounts</span>
       <h3 class="sidebar-title">{i18n.t("settings.user_settings")}</h3>
    </div>
    <nav class="settings-nav">
      {#each tabs as tab}
        <button
          class="nav-item"
          class:active={activeTab === tab.id}
          onclick={() => (activeTab = tab.id)}
        >
          <span class="material-symbols-rounded">{tab.icon}</span>
          <span class="label-text">{tab.label}</span>
          {#if activeTab === tab.id}
             <div class="active-indicator"></div>
          {/if}
        </button>
      {/each}
    </nav>
  </aside>

  <main class="settings-content">
    {#if activeTab === "profile"}
      <SettingsProfile />
    {:else if activeTab === "myggit"}
      <SettingsMyGgit />
    {:else if activeTab === "appearance"}
      <ThemeSelector />
    {:else if activeTab === "emails"}
      <SettingsEmails />
    {:else if activeTab === "security"}
      <SettingsSecurity />
    {:else if activeTab === "keys"}
      <SettingsKeys />
    {:else if activeTab === "tokens"}
      <SettingsTokens />
    {:else if activeTab === "repos"}
      <SettingsRepositories />
    {:else if activeTab === "orgs"}
      <SettingsOrganizations />
    {/if}
  </main>
</div>

<style>
  .settings-premium-container {
    display: grid;
    grid-template-columns: 320px minmax(0, 1fr);
    gap: 3rem;
    align-items: start;
    padding: 1rem 0;
    min-height: 80vh;
    max-width: 100%;
    box-sizing: border-box;
  }

  .settings-sidebar {
    position: sticky;
    top: 2.5rem;
    display: flex;
    flex-direction: column;
    gap: 3rem;
  }

  .sidebar-brand {
    display: flex;
    align-items: center;
    gap: 1.25rem;
    padding-left: 1rem;
  }
  .sidebar-brand span {
    font-size: 2.25rem;
    color: var(--primary-color);
  }

  .sidebar-title {
    font-size: 1.25rem;
    font-weight: 900;
    margin: 0;
    letter-spacing: -0.02em;
    color: var(--text-color);
  }

  .settings-nav {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }

  .nav-item {
    display: flex;
    align-items: center;
    gap: 1.25rem;
    padding: 1rem 1.5rem;
    background: transparent;
    border: none;
    color: var(--text-muted);
    font-size: 1.05rem;
    font-weight: 700;
    text-align: left;
    cursor: pointer;
    border-radius: 1.5rem;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
    width: 100%;
    position: relative;
    overflow: hidden;
    min-width: 0;
  }

  .nav-item:hover {
    background: rgba(255, 255, 255, 0.04);
    color: var(--text-color);
    transform: translateX(8px);
  }

  .nav-item.active {
    background: rgba(var(--primary-rgb), 0.1);
    color: var(--primary-color);
    box-shadow: 0 4px 15px rgba(var(--primary-rgb), 0.05);
  }

  .nav-item span.material-symbols-rounded {
    font-size: 1.5rem;
    flex-shrink: 0;
  }

  .label-text {
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    flex: 1;
  }

  .active-indicator {
    position: absolute;
    left: 0;
    top: 25%;
    bottom: 25%;
    width: 4px;
    background: var(--primary-color);
    border-radius: 0 4px 4px 0;
  }

  .settings-content {
    display: flex;
    flex-direction: column;
    gap: 3rem;
    animation: contentFade 0.5s ease-out;
    min-width: 0;
    max-width: 100%;
    width: 100%;
    box-sizing: border-box;
  }

  @keyframes contentFade {
    from { opacity: 0; transform: translateY(20px); }
    to { opacity: 1; transform: translateY(0); }
  }

  @media (max-width: 1100px) {
    .settings-premium-container {
      grid-template-columns: minmax(0, 1fr);
      padding: 0 1rem;
      gap: 2rem;
    }
    .settings-sidebar {
      position: relative;
      top: 0;
      gap: 1.5rem;
      z-index: 10;
    }
    .settings-nav {
      flex-direction: row;
      overflow-x: auto;
      padding-bottom: 0.5rem;
      margin: 0 -1rem;
      padding-left: 1rem;
      scrollbar-width: none;
      -webkit-overflow-scrolling: touch;
    }
    .settings-nav::-webkit-scrollbar { display: none; }
    .nav-item {
      width: auto;
      padding: 0.75rem 1.25rem;
      border-radius: 100px;
      flex-shrink: 0;
      gap: 0.5rem;
      font-size: 0.95rem;
    }
    .nav-item:hover { transform: none; }
    .active-indicator {
      display: none;
    }
    .settings-content {
      gap: 2rem;
    }
  }

  @media (max-width: 600px) {
    .sidebar-brand span { font-size: 1.75rem; }
    .sidebar-title { font-size: 1.1rem; }
    .settings-premium-container {
      padding: 0 0.5rem;
      gap: 1.5rem;
    }
    .settings-content {
      gap: 1.5rem;
    }
  }
</style>
