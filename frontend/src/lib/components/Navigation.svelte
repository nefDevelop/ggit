<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { authState } from "$lib/auth.svelte";
  import { i18n } from "$lib/i18n.svelte";
  import { page } from "$app/state";
  import ThemeSwitch from "$lib/components/ThemeSwitch.svelte";
  import Avatar from "$lib/components/Avatar.svelte";
  import { userPrefs } from "$lib/settings.svelte";

  let showThemePanel = $state(false);
  let userOrgs = $state<any[]>([]);

  async function fetchUserOrgs() {
    if (!authState.isAuthenticated) return;
    try {
      const res = await apiFetch("/api/v2/user/settings/orgs");
      if (res.ok) userOrgs = await res.json();
    } catch (e) {
      console.error("Failed to load sidebar orgs", e);
    }
  }

  $effect(() => {
    if (authState.isAuthenticated) fetchUserOrgs();
  });

  const navItems = $derived([
    { name: i18n.t("dashboard"), icon: "grid_view", path: "/dashboard" },
    { name: i18n.t("explore"), icon: "search", path: "/explore/repos", isParent: true, isExplore: true },
    ...(authState.user?.is_admin
      ? [
          {
            name: i18n.t("admin_panel"),
            icon: "security",
            path: "/admin/users",
            isParent: true,
          },
        ]
      : []),
    {
      name: i18n.t("settings_title"),
      icon: "manage_accounts",
      path: "/user/settings",
      isParent: true,
      isSettings: true,
    },
  ]);

  const adminSubItems = $derived([
    {
      name: i18n.t("admin.users") || "Users",
      icon: "person",
      path: "/admin/users",
    },
    {
      name: i18n.t("admin.repositories") || "Repositories",
      icon: "book",
      path: "/admin/repos",
    },
    {
      name: i18n.t("admin.organizations") || "Organizations",
      icon: "corporate_fare",
      path: "/admin/orgs",
    },
    {
      name: i18n.t("admin.auths") || "Authentication",
      icon: "lock",
      path: "/admin/auth",
    },
    {
      name: i18n.t("admin.config") || "Configuration",
      icon: "settings",
      path: "/admin/config",
    },
    {
      name: i18n.t("admin.monitor") || "Monitor",
      icon: "monitoring",
      path: "/admin/monitor",
    },
    {
      name: i18n.t("admin.maintenance") || "Maintenance",
      icon: "build",
      path: "/admin/maintenance",
    },
    {
      name: i18n.t("admin.audit_log") || "Audit Log",
      icon: "history_edu",
      path: "/admin/audit_log",
    },
    {
      name: i18n.t("admin.security") || "Seguridad",
      icon: "shield",
      path: "/admin/security",
    },
    {
      name: i18n.t("admin.janitor") || "Janitor",
      icon: "mop",
      path: "/admin/janitor",
    },
    {
      name: i18n.t("admin.notices") || "Notices",
      icon: "notifications",
      path: "/admin/notices",
    },
  ]);

  const settingsSubItems = $derived([
    { name: i18n.t("settings.profile") || "Profile", icon: "person", path: "/user/settings?tab=profile" },
    { name: i18n.t("settings.my_ggit") || "My Ggit", icon: "grid_view", path: "/user/settings?tab=myggit" },
    { name: i18n.t("settings.appearance") || "Appearance", icon: "palette", path: "/user/settings?tab=appearance" },
    { name: i18n.t("settings.emails") || "Emails", icon: "mail", path: "/user/settings?tab=emails" },
    { name: i18n.t("settings.security") || "Security", icon: "security", path: "/user/settings?tab=security" },
    { name: i18n.t("settings.ssh_keys") || "SSH Keys", icon: "key", path: "/user/settings?tab=keys" },
    { name: i18n.t("settings.access_tokens") || "Tokens", icon: "token", path: "/user/settings?tab=tokens" },
    { name: i18n.t("settings.repositories") || "Repositories", icon: "book", path: "/user/settings?tab=repos" },
    { name: i18n.t("settings.organizations") || "Organizations", icon: "corporate_fare", path: "/user/settings?tab=orgs" },
  ]);
  
  const exploreSubItems = $derived([
    { name: i18n.t("explore_repos"), icon: "book", path: "/explore/repos" },
    { name: i18n.t("explore_users"), icon: "person_search", path: "/explore/users" },
    { name: i18n.t("explore_organizations"), icon: "corporate_fare", path: "/explore/organizations" },
  ]);

  const isAdminPath = $derived(page.url.pathname.startsWith("/admin"));
  const isSettingsPath = $derived(page.url.pathname.startsWith("/user/settings"));
  const isExplorePath = $derived(page.url.pathname.startsWith("/explore"));

  function getCsrfToken() {
    if (typeof document === "undefined") return "";
    const match = document.cookie.match(new RegExp("(^| )XSRF-TOKEN=([^;]+)"));
    return match ? match[2] : "";
  }

  async function handleLogout() {
    await apiFetch("/api/v2/user/logout", {
      method: "POST",
    });
    window.location.href = "/user/login";
  }
</script>

<nav class="navigation-bar">
  <!-- USER AT TOP -->
  {#if authState.isAuthenticated}
    <div class="user-profile-summary">
      <div class="avatar-ring">
        <Avatar
          src={authState.user?.avatar_url}
          username={authState.user?.username}
          size={40}
        />
      </div>
      <div class="user-details">
        <span class="full-name"
          >{authState.user?.full_name || authState.user?.username}</span
        >
        <span class="user-handle">@{authState.user?.username}</span>
      </div>
      <button class="minimal-logout" onclick={handleLogout} title="Logout">
        <span class="material-symbols-rounded">logout</span>
      </button>
    </div>
  {/if}

  <div class="nav-links">
    {#each navItems as item}
      <a
        href={item.path}
        class="nav-item"
        class:active={page.url.pathname === item.path ||
          (item.isParent && ((isAdminPath && !item.isSettings && !item.isExplore) || (isSettingsPath && item.isSettings) || (isExplorePath && item.isExplore)))}
      >
        <span class="material-symbols-rounded">{item.icon}</span>
        <span class="item-label">{item.name}</span>
      </a>

      {#if item.isParent && ((isAdminPath && !item.isSettings && !item.isExplore) || (isSettingsPath && item.isSettings) || (isExplorePath && item.isExplore))}
        <div class="sub-nav-links">
          {#each (item.isSettings ? settingsSubItems : item.isExplore ? exploreSubItems : adminSubItems) as sub}
            <a
              href={sub.path}
              class="sub-nav-item"
              class:active={page.url.pathname === sub.path || (sub.path.includes('?') && page.url.search === sub.path.substring(sub.path.indexOf('?')))}
            >
              <span class="material-symbols-rounded">{sub.icon}</span>
              <span class="item-label">{sub.name}</span>
            </a>
          {/each}
        </div>
      {/if}
    {/each}

    <!-- ORGANIZATIONS SECTION -->
    {#if authState.isAuthenticated && userPrefs.current.showSidebarOrgs}
      <div class="nav-section-title">
        <span>{i18n.t("organizations") || "Organizaciones"}</span>
        <a href="/org/create" class="mini-add" title="Nueva Organización">
          <span class="material-symbols-rounded">add</span>
        </a>
      </div>
      <div class="nav-org-list">
        {#each userOrgs as org}
          <a
            href="/org/{org.username}"
            class="nav-item org-item"
            class:active={page.url.pathname === `/org/${org.username}`}
          >
            <Avatar src={org.avatar_url} username={org.username} size={24} />
            <span class="item-label">{org.full_name || org.username}</span>
          </a>
        {/each}
        {#if userOrgs.length === 0}
          <p class="nav-empty-hint">No eres parte de ninguna organización.</p>
        {/if}
      </div>
    {/if}
  </div>



  <!-- BRAND AT BOTTOM -->
  <div class="brand-section-bottom">
    <div class="brand-info">
      <span class="brand-name">Ggit</span>
      <span class="version">v2026.4.0</span>
    </div>
    <button
      class="minimal-theme-btn"
      onclick={() => (showThemePanel = !showThemePanel)}
      title="Apariencia"
    >
      <span class="material-symbols-rounded">palette</span>
    </button>
  </div>
</nav>

<ThemeSwitch bind:show={showThemePanel} />

<style>
  .navigation-bar {
    display: flex;
    flex-direction: column;
    padding: 2rem 1rem;
    width: 260px;
    height: 100vh;
    box-sizing: border-box;
    z-index: 1000;
    /* No borders, no s, just flat integration with background */
  }

  /* USER PROFILE SUMMARY */
  .user-profile-summary {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 1rem;
    margin-bottom: 2rem;
    background: var(--md-sys-color-secondary-container);
    color: var(--md-sys-color-on-secondary-container);
    border-radius: 1.5rem;
    border: 1px solid var(--md-sys-color-outline-variant);
    transition: all 0.3s ease;
  }
  .user-profile-summary:hover {
    background: rgba(255, 255, 255, 0.06);
    transform: translateY(-2px);
  }

  .avatar-ring {
    display: flex;
    align-items: center;
    justify-content: center;
    position: relative;
    width: 48px;
    height: 48px;
    border-radius: 50%;
    background: linear-gradient(135deg, var(--primary-color), var(--tertiary));
    flex-shrink: 0;
  }

  .user-details {
    flex: 1;
    display: flex;
    flex-direction: column;
    min-width: 0;
  }

  .full-name {
    font-size: 0.9rem;
    font-weight: 800;
    color: inherit;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .user-handle {
    font-size: 0.75rem;
    color: inherit;
    opacity: 0.8;
    font-weight: 600;
  }

  .minimal-logout {
    background: transparent;
    border: none;
    color: var(--text-muted);
    cursor: pointer;
    padding: 0.5rem;
    display: flex;
    border-radius: 0.75rem;
    transition: all 0.2s;
  }
  .minimal-logout:hover {
    color: #ef4444;
    background: rgba(239, 68, 68, 0.1);
  }

  /* LINKS */
  .nav-links {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    flex: 1;
  }

  .nav-item {
    display: flex;
    align-items: center;
    gap: 1.25rem;
    padding: 0.85rem 1.25rem;
    color: var(--text-muted);
    text-decoration: none;
    border-radius: 1rem;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    font-weight: 600;
    font-size: 0.95rem;
    min-width: 0;
  }

  .nav-item:hover {
    color: var(--text-color);
    background: rgba(255, 255, 255, 0.03);
  }

  .nav-item.active {
    color: var(--primary-color);
    background: rgba(var(--primary-rgb), 0.08);
  }

  .item-label {
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
    flex: 1;
  }

  /* ORGANIZATIONS */
  .nav-section-title {
    display: flex;
    align-items: center;
    justify-content: space-between;
    padding: 1.5rem 1.25rem 0.5rem;
    font-size: 0.75rem;
    font-weight: 800;
    text-transform: uppercase;
    color: var(--text-muted);
    letter-spacing: 0.05em;
  }

  .nav-org-list {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .org-item {
    padding: 0.65rem 1.25rem;
    font-size: 0.9rem;
  }

  .mini-add {
    color: var(--text-muted);
    text-decoration: none;
    transition: color 0.2s;
    display: flex;
  }
  .mini-add:hover {
    color: var(--primary-color);
  }
  .mini-add span {
    font-size: 1.25rem;
  }

  .nav-empty-hint {
    font-size: 0.75rem;
    color: var(--text-muted);
    padding: 0.5rem 1.25rem;
    opacity: 0.6;
    font-style: italic;
  }

  /* SUB NAV */
  .sub-nav-links {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    padding-left: 0.75rem;
    margin: 0.25rem 0;
    border-left: 1px solid var(--md-sys-color-outline-variant);
    margin-left: 1.25rem;
  }

  .sub-nav-item {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 0.65rem 1rem;
    color: var(--text-muted);
    text-decoration: none;
    border-radius: 0.75rem;
    transition: all 0.2s;
    font-size: 0.85rem;
    font-weight: 600;
  }

  .sub-nav-item:hover {
    color: var(--text-color);
    background: rgba(255, 255, 255, 0.03);
  }

  .sub-nav-item.active {
    color: var(--md-sys-color-on-secondary-container);
    background: var(--md-sys-color-secondary-container);
    box-shadow: var(--shadow-sm);
  }

  .sub-nav-item .material-symbols-rounded {
    font-size: 1.1rem;
  }

  /* BRAND SECTION BOTTOM */
  .brand-section-bottom {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 1rem;
    margin-top: auto;
    opacity: 0.6;
    transition: opacity 0.2s;
  }
  .brand-section-bottom:hover {
    opacity: 1;
  }

  .brand-info {
    display: flex;
    flex-direction: column;
  }
  .brand-name {
    font-weight: 800;
    font-size: 1.1rem;
    letter-spacing: -0.02em;
    line-height: 1;
  }
  .version {
    font-size: 0.65rem;
    font-weight: 700;
    color: var(--text-muted);
  }

  /* Mobile stays similar but at bottom */
  @media (max-width: 1024px) {
    .navigation-bar {
      width: 100vw;
      height: 70px;
      flex-direction: row;
      position: fixed;
      bottom: 0;
      left: 0;
      z-index: 1000;
      padding: 0 0.5rem;
      background: var(--bg-color);
      border-top: 1px solid var(--md-sys-color-outline-variant);
      align-items: center;
      padding-bottom: env(safe-area-inset-bottom);
      height: calc(70px + env(safe-area-inset-bottom));
      box-shadow: 0 -4px 20px rgba(0, 0, 0, 0.1);
      justify-content: space-around;
    }
    .user-profile-summary,
    .brand-section-bottom,
    .item-label,
    .nav-section-title,
    .sub-nav-links,
    .nav-empty-hint,
    span:not(.material-symbols-rounded) {
      display: none;
    }
    .nav-links {
      flex-direction: row;
      width: auto;
      flex: 1;
      justify-content: space-around;
      gap: 0;
    }
    .nav-item {
      padding: 0.75rem;
      gap: 0;
      justify-content: center;
      flex: 1;
    }
    .nav-org-list {
      display: none;
    }
    .appearance-btn {
      display: none;
    }
  }

  /* Theme Button Styles */
  .minimal-theme-btn {
    background: transparent;
    border: none;
    color: var(--text-muted);
    cursor: pointer;
    padding: 0.5rem;
    display: flex;
    border-radius: 0.75rem;
    transition: all 0.2s;
    margin-left: auto;
  }
  .minimal-theme-btn:hover {
    color: var(--primary-color);
    background: rgba(var(--primary-rgb), 0.1);
  }
  .minimal-theme-btn span {
    font-size: 1.25rem;
  }
</style>
