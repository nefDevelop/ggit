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
  import { page } from "$app/state";
  import { goto } from "$app/navigation";
  import PageHeader from "./ui/PageHeader.svelte";
  import ViewHeader from "./ui/ViewHeader.svelte";
  import ChipNavigation from "./ui/ChipNavigation.svelte";

  let activeTab = $state(page.url.searchParams.get("tab") || "profile");

  $effect(() => {
    const tabFromUrl = page.url.searchParams.get("tab") || "profile";
    if (tabFromUrl !== activeTab) {
      activeTab = tabFromUrl;
    }
  });

  $effect(() => {
    if (activeTab) {
      const currentTab = page.url.searchParams.get("tab") || "profile";
      if (activeTab !== currentTab) {
        goto(`?tab=${activeTab}`, { replaceState: true });
      }
    }
  });

  const tabs = $derived([
    { id: "profile", label: i18n.t("settings.profile") || "Perfil", icon: "person", description: "Configura tu información pública y avatar", path: "/user/settings?tab=profile" },
    { id: "myggit", label: i18n.t("settings.my_ggit") || "Mi Ggit", icon: "grid_view", description: "Personaliza tu página de inicio y preferencias del feed", path: "/user/settings?tab=myggit" },
    { id: "appearance", label: i18n.t("settings.appearance") || "Apariencia", icon: "palette", description: "Personaliza el tema, colores y tipografía de la interfaz", path: "/user/settings?tab=appearance" },
    { id: "emails", label: i18n.t("settings.emails") || "Correos", icon: "mail", description: "Gestiona las direcciones de correo electrónico vinculadas", path: "/user/settings?tab=emails" },
    { id: "security", label: i18n.t("settings.security") || "Seguridad", icon: "security", description: "Actualiza tu contraseña y configuración de doble factor", path: "/user/settings?tab=security" },
    { id: "sessions", label: "Sesiones", icon: "devices", description: "Controla los dispositivos y navegadores conectados", path: "/user/settings?tab=sessions" },
    { id: "blocks", label: "Bloqueos", icon: "person_off", description: "Gestiona los usuarios que tienes bloqueados", path: "/user/settings?tab=blocks" },
    { id: "keys", label: i18n.t("settings.ssh_keys") || "Claves SSH", icon: "key", description: "Gestiona las llaves públicas vinculadas a tu cuenta", path: "/user/settings?tab=keys" },
    { id: "tokens", label: i18n.t("settings.access_tokens") || "Tokens", icon: "token", description: "Claves de acceso seguro vía API o Git sin usar contraseña", path: "/user/settings?tab=tokens" },
    { id: "repos", label: i18n.t("settings.repositories") || "Repositorios", icon: "book", description: "Lista de todos tus proyectos personales y colaboraciones activas", path: "/user/settings?tab=repos" },
    { id: "orgs", label: i18n.t("settings.organizations") || "Organizaciones", icon: "corporate_fare", description: "Gestiona tu pertenencia a equipos y organizaciones", path: "/user/settings?tab=orgs" },
  ]);

  const activeTabObj = $derived(tabs.find(t => t.id === activeTab) || tabs[0]);
</script>

<div class="settings-premium-container">
  <PageHeader 
    title={activeTabObj.label}
    subtitle={activeTabObj.description}
    icon={activeTabObj.icon}
  />

  <div class="settings-premium-body">
    <div class="mobile-nav-wrapper">
      <ChipNavigation 
        items={tabs} 
        activeId={activeTab} 
        onSelect={(id) => activeTab = id}
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
</div>

<style>
  .settings-premium-container {
    display: flex;
    flex-direction: column;
    gap: 2rem;
    min-height: 80vh;
    max-width: 100%;
    box-sizing: border-box;
    padding: 0;
  }

  .settings-premium-body {
    padding: 0 1.5rem 4rem;
    display: flex;
    flex-direction: column;
    gap: 2rem;
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
      padding-bottom: calc(140px + env(safe-area-inset-bottom));
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
      border-top: none;
      box-shadow: 0 -10px 30px rgba(0, 0, 0, 0.1);
    }
    .settings-content {
      gap: 2rem;
    }
  }

  @media (max-width: 600px) {
    .settings-premium-container {
      padding: 0 0.5rem;
      padding-bottom: calc(130px + env(safe-area-inset-bottom));
      gap: 1.5rem;
    }
    .settings-content {
      gap: 1.5rem;
    }
  }
</style>
