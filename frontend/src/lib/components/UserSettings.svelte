<script lang="ts">
  import ThemeSelector from "./ThemeSelector.svelte";
  import SettingsProfile from "./SettingsProfile.svelte";
  import SettingsSecurity from "./SettingsSecurity.svelte";
  import SettingsEmails from "./SettingsEmails.svelte";
  import SettingsKeys from "./SettingsKeys.svelte";
  import SettingsTokens from "./SettingsTokens.svelte";
  import SettingsSessions from "./SettingsSessions.svelte";
  import SettingsBlockedUsers from "./SettingsBlockedUsers.svelte";
  import SettingsRepositories from "./SettingsRepositories.svelte";
  import SettingsOrganizations from "./SettingsOrganizations.svelte";
  import SettingsMyGgit from "./SettingsMyGgit.svelte";
  import { i18n } from "$lib/i18n.svelte";
  import ChipNavigation from "./ui/ChipNavigation.svelte";
  import { page } from "$app/state";
  import { goto } from "$app/navigation";

  let activeTab = $derived(page.url.searchParams.get("tab") || "profile");

  function setTab(id: string) {
    goto(`?tab=${id}`, { replaceState: true });
  }

  const tabs = $derived([
    { id: "profile", label: i18n.t("settings.profile"), icon: "person" },
    { id: "myggit", label: i18n.t("settings.my_ggit"), icon: "grid_view" },
    { id: "appearance", label: i18n.t("settings.appearance"), icon: "palette" },
    { id: "emails", label: i18n.t("settings.emails"), icon: "mail" },
    { id: "security", label: i18n.t("settings.security"), icon: "security" },
    { id: "sessions", label: "Sesiones", icon: "devices" },
    { id: "blocks", label: "Bloqueos", icon: "person_off" },
    { id: "keys", label: i18n.t("settings.ssh_keys"), icon: "key" },
    { id: "tokens", label: i18n.t("settings.access_tokens"), icon: "token" },
    { id: "repos", label: i18n.t("settings.repositories"), icon: "book" },
    { id: "orgs", label: i18n.t("settings.organizations"), icon: "corporate_fare" },
  ]);
</script>

<div class="settings-premium-container">
  <!-- Mobile Chip Rail -->
  <div class="mobile-nav-wrapper">
    <ChipNavigation 
      items={tabs} 
      activeId={activeTab} 
      title={i18n.t("settings.user_settings")} 
      titleIcon="manage_accounts"
      onSelect={(id) => setTab(id)} 
    />
  </div>

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
    {:else if activeTab === "sessions"}
      <SettingsSessions />
    {:else if activeTab === "blocks"}
      <SettingsBlockedUsers />
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
    display: flex;
    flex-direction: column;
    gap: 2rem;
    min-height: 80vh;
    max-width: 100%;
    box-sizing: border-box;
    padding: 1rem 0;
  }

  .mobile-nav-wrapper {
    display: none;
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
      padding: 0;
      padding-bottom: calc(160px + env(safe-area-inset-bottom));
    }
    .mobile-nav-wrapper {
      display: block;
      position: fixed;
      bottom: calc(70px + env(safe-area-inset-bottom));
      left: 0;
      right: 0;
      z-index: 999;
      background: var(--bg-color);
      padding: 0.5rem 0.5rem 0;
      border-top: 1px solid var(--md-sys-color-outline-variant);
      box-shadow: 0 -10px 30px rgba(0, 0, 0, 0.1);
    }
    .settings-content {
      gap: 2rem;
    }
  }

  @media (max-width: 600px) {
    .settings-premium-container {
      padding: 0 0.5rem;
      padding-bottom: calc(140px + env(safe-area-inset-bottom));
      gap: 1.5rem;
    }
    .settings-content {
      gap: 1.5rem;
    }
  }
</style>
