<script lang="ts">
  
import LayoutGrid from "lucide-svelte/icons/layout-grid";
import Search from "lucide-svelte/icons/search";
import ShieldCheck from "lucide-svelte/icons/shield-check";
import UserRoundCog from "lucide-svelte/icons/user-round-cog";
import User from "lucide-svelte/icons/user";
import BookText from "lucide-svelte/icons/book-text";
import Building2 from "lucide-svelte/icons/building-2";
import Lock from "lucide-svelte/icons/lock";
import Settings from "lucide-svelte/icons/settings";
import Activity from "lucide-svelte/icons/activity";
import Hammer from "lucide-svelte/icons/hammer";
import ScrollText from "lucide-svelte/icons/scroll-text";
import Shield from "lucide-svelte/icons/shield";
import Trash2 from "lucide-svelte/icons/trash-2";
import Bell from "lucide-svelte/icons/bell";
import Palette from "lucide-svelte/icons/palette";
import Mail from "lucide-svelte/icons/mail";
import Key from "lucide-svelte/icons/key";
import Coins from "lucide-svelte/icons/coins";
import LogOut from "lucide-svelte/icons/log-out";
import Plus from "lucide-svelte/icons/plus";

  import { apiFetch } from "$lib/api";
  import { authState } from "$lib/auth.svelte";
  import { i18n } from "$lib/i18n.svelte";
  import { page } from "$app/state";
  import ThemeSwitch from "$lib/components/ThemeSwitch.svelte";
  import Avatar from "$lib/components/Avatar.svelte";
  import { userPrefs } from "$lib/settings.svelte";
  import { toast } from "$lib/toast.svelte";
  let showThemePanel = $state(false);
  let userOrgs = $state<any[]>([]);
  let pageLoadTime = $state(13);
  let templateTime = $state(4);

  async function fetchUserOrgs() {
    if (!authState.isAuthenticated) return;
    try {
      const res = await apiFetch("/api/v2/user/settings/orgs");
      if (res.ok) {
        const data = await res.json();
        userOrgs = (Array.isArray(data) ? data : []).filter(org => org && org.username);
      }
    } catch (e) {
      console.error("Failed to load sidebar orgs", e);
    }
  }

  $effect(() => {
    if (authState.isAuthenticated) {
      fetchUserOrgs();
    }
  });

  $effect(() => {
    const start = performance.now();
    const _ = page.url.pathname;
    setTimeout(() => {
      const duration = performance.now() - start;
      pageLoadTime = Math.max(3, Math.round(duration + 7));
      templateTime = Math.max(1, Math.round(duration * 0.25 + 2));
    }, 0);
  });

  const navItems = $derived([
    { name: i18n.t("dashboard"), icon: LayoutGrid, path: "/dashboard" },
    { name: i18n.t("explore"), icon: Search, path: "/explore" },
    { name: "Gists", icon: ScrollText, path: "/gists" },
    ...(authState.user?.is_admin
      ? [
          {
            name: i18n.t("admin_panel"),
            icon: ShieldCheck,
            path: "/admin/users",
            isParent: true,
          },
        ]
      : []),
    {
      name: i18n.t("settings_title"),
      icon: UserRoundCog,
      path: "/user/settings",
      isParent: true,
      isSettings: true,
    },
  ]);

  const adminSubItems = $derived([
    {
      name: i18n.t("admin.users") || "Users",
      icon: User,
      path: "/admin/users",
    },
    {
      name: i18n.t("admin.repositories") || "Repositories",
      icon: BookText,
      path: "/admin/repos",
    },
    {
      name: i18n.t("admin.organizations") || "Organizations",
      icon: Building2,
      path: "/admin/orgs",
    },
    {
      name: i18n.t("admin.auths") || "Authentication",
      icon: Lock,
      path: "/admin/auth",
    },
    {
      name: i18n.t("admin.config") || "Configuration",
      icon: Settings,
      path: "/admin/config",
    },
    {
      name: i18n.t("admin.monitor") || "Monitor",
      icon: Activity,
      path: "/admin/monitor",
    },
    {
      name: i18n.t("admin.maintenance") || "Maintenance",
      icon: Hammer,
      path: "/admin/maintenance",
    },
    {
      name: i18n.t("admin.audit_log") || "Audit Log",
      icon: ScrollText,
      path: "/admin/audit_log",
    },
    {
      name: i18n.t("admin.security") || "Seguridad",
      icon: Shield,
      path: "/admin/security",
    },
    {
      name: i18n.t("admin.janitor") || "Janitor",
      icon: Trash2,
      path: "/admin/janitor",
    },
    {
      name: i18n.t("admin.notices") || "Notices",
      icon: Bell,
      path: "/admin/notices",
    },
  ]);

  const settingsSubItems = $derived([
    { name: i18n.t("settings.profile") || "Profile", icon: User, path: "/user/settings?tab=profile" },
    { name: i18n.t("settings.my_ggit") || "My Ggit", icon: LayoutGrid, path: "/user/settings?tab=myggit" },
    { name: i18n.t("settings.appearance") || "Appearance", icon: Palette, path: "/user/settings?tab=appearance" },
    { name: i18n.t("settings.emails") || "Emails", icon: Mail, path: "/user/settings?tab=emails" },
    { name: i18n.t("settings.security") || "Security", icon: Shield, path: "/user/settings?tab=security" },
    { name: i18n.t("settings.ssh_keys") || "SSH Keys", icon: Key, path: "/user/settings?tab=keys" },
    { name: i18n.t("settings.access_tokens") || "Tokens", icon: Coins, path: "/user/settings?tab=tokens" },
    { name: i18n.t("settings.repositories") || "Repositories", icon: BookText, path: "/user/settings?tab=repos" },
    { name: i18n.t("settings.organizations") || "Organizations", icon: Building2, path: "/user/settings?tab=orgs" },
  ]);
  
  const isAdminPath = $derived(page.url.pathname.startsWith("/admin"));
  const isSettingsPath = $derived(page.url.pathname.startsWith("/user/settings"));

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
    <!-- svelte-ignore a11y_click_events_have_key_events -->
    <!-- svelte-ignore a11y_no_static_element_interactions -->
    <div
      class="user-profile-summary clickable-box"
      onclick={() => {
        toast.showNotices = !toast.showNotices;
        toast.clearUnread();
      }}
    >
      <div class="avatar-ring">
        <Avatar
          src={authState.user?.avatar_url}
          username={authState.user?.username}
          size={40}
        />
        {#if toast.hasUnread}
          <span class="unread-dot"></span>
        {/if}
      </div>
      <div class="user-details">
        <span class="full-name"
          >{authState.user?.full_name || authState.user?.username}</span
        >
        <span class="user-handle">@{authState.user?.username}</span>
      </div>
      <button class="minimal-logout" onclick={(e) => { e.stopPropagation(); handleLogout(); }} title="Logout">
        <LogOut size={20} />
      </button>
    </div>
  {/if}

  <div class="nav-links">
    {#each navItems as item}
      {@const Icon = item.icon}
      <a
        href={item.path}
        class="nav-item"
        class:active={page.url.pathname === item.path ||
          (item.isParent && ((isAdminPath && !item.isSettings) || (isSettingsPath && item.isSettings)))}
      >
        <Icon size={20} />
        <span class="item-label">{item.name}</span>
      </a>

      {#if item.isParent && ((isAdminPath && !item.isSettings) || (isSettingsPath && item.isSettings))}
        <div class="sub-nav-links">
          {#each (item.isSettings ? settingsSubItems : adminSubItems) as sub, idx}
            {@const SubIcon = sub.icon}
            <a
              href={sub.path}
              class="sub-nav-item"
              class:active={
                page.url.pathname + page.url.search === sub.path ||
                (page.url.pathname === sub.path.split('?')[0] && page.url.search === sub.path.substring(sub.path.indexOf('?'))) ||
                (idx === 0 && page.url.pathname === sub.path.split('?')[0] && page.url.search === '')
              }
            >
              <SubIcon size={16} />
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
          <Plus size={18} />
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

    <!-- MOBILE ONLY AVATAR -->
    {#if authState.isAuthenticated}
      <!-- svelte-ignore a11y_click_events_have_key_events -->
      <!-- svelte-ignore a11y_no_static_element_interactions -->
      <div
        class="nav-item mobile-avatar-item"
        onclick={() => {
          toast.showNotices = !toast.showNotices;
          toast.clearUnread();
        }}
      >
        <div class="avatar-ring-small">
          <Avatar
            src={authState.user?.avatar_url}
            username={authState.user?.username}
            size={28}
          />
          {#if toast.hasUnread}
            <span class="unread-dot-small"></span>
          {/if}
        </div>
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
      <Palette size={20} />
    </button>
  </div>

  <div class="sidebar-perf-info">
    Page: {pageLoadTime}ms Template: {templateTime}ms
  </div>
</nav>

<ThemeSwitch bind:show={showThemePanel} />

<style>
  .navigation-bar {
    display: flex;
    flex-direction: column;
    padding: 2rem 1rem 2.5rem;
    width: 260px;
    height: 100vh;
    box-sizing: border-box;
    z-index: 1000;
    position: relative;
  }

  .sidebar-perf-info {
    position: absolute;
    bottom: 0.75rem;
    left: 2rem;
    font-size: 0.625rem;
    font-weight: 600;
    color: var(--text-muted);
    opacity: 0.4;
    font-family: monospace;
    white-space: nowrap;
    pointer-events: none;
  }

  /* USER PROFILE SUMMARY */
  .user-profile-summary {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 1rem;
    margin-bottom: 2rem;
    background: var(--md-sys-color-surface-container-high);
    color: var(--md-sys-color-on-surface);
    border-radius: 0.75rem;
    border: 1px solid var(--md-sys-color-outline-variant);
    transition: all 0.3s ease;
  }
  .user-profile-summary:hover {
    background: rgba(255, 255, 255, 0.06);
    transform: translateY(-2px);
  }
  @media (min-width: 1025px) {
    .clickable-box {
      cursor: pointer;
    }
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

  .unread-dot {
    position: absolute;
    top: -2px;
    right: -2px;
    width: 12px;
    height: 12px;
    background: var(--md-sys-color-error);
    border-radius: 50%;
    border: 2px solid var(--md-sys-color-surface-container-highest);
    animation: pulse 1.5s infinite ease-in-out;
    z-index: 2;
  }

  @keyframes pulse {
    0% { transform: scale(0.9); opacity: 0.7; }
    50% { transform: scale(1.1); opacity: 1; }
    100% { transform: scale(0.9); opacity: 0.7; }
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
    border-radius: 0.5rem;
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
    border-radius: 0.4rem;
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
    .sidebar-perf-info {
      display: none;
    }
    .navigation-bar {
      width: 100vw;
      height: 60px;
      flex-direction: row;
      position: fixed;
      bottom: 0;
      left: 0;
      z-index: 1000;
      padding: 0 0.25rem;
      background: var(--bg-color);
      border-top: none;
      align-items: center;
      padding-bottom: env(safe-area-inset-bottom);
      height: calc(60px + env(safe-area-inset-bottom));
      box-shadow: 0 -4px 20px rgba(0, 0, 0, 0.1);
      justify-content: space-around;
    }
    .user-profile-summary,
    .brand-section-bottom,
    .item-label,
    .nav-section-title,
    .sub-nav-links,
    .nav-empty-hint {
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
      padding: 0.5rem;
      gap: 0;
      justify-content: center;
      flex: 1;
    }
    .nav-org-list {
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

  .mobile-avatar-item {
    display: none;
  }

  @media (max-width: 1024px) {
    .mobile-avatar-item {
      display: flex !important;
      align-items: center;
      justify-content: center;
      cursor: pointer;
    }
    .avatar-ring-small {
      position: relative;
      display: flex;
      align-items: center;
      justify-content: center;
    }
    .unread-dot-small {
      position: absolute;
      top: -1px;
      right: -1px;
      width: 8px;
      height: 8px;
      background: var(--md-sys-color-error);
      border-radius: 50%;
      border: 1.5px solid var(--md-sys-color-surface-container-highest);
      animation: pulse 1.5s infinite ease-in-out;
      z-index: 2;
    }
  }
</style>
