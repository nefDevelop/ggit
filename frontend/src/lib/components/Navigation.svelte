<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { authState } from "$lib/auth.svelte";
  import { i18n } from "$lib/i18n.svelte";
  import { page } from "$app/state";
  import { userPrefs } from "$lib/settings.svelte";
  import ThemeSwitch from "$lib/components/ThemeSwitch.svelte";
  import NavUserProfile from "$lib/components/NavUserProfile.svelte";
  import NavMainLinks from "$lib/components/NavMainLinks.svelte";
  import NavMobileAvatar from "$lib/components/NavMobileAvatar.svelte";
  import NavBrandFooter from "$lib/components/NavBrandFooter.svelte";
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
  import Rocket from "lucide-svelte/icons/rocket";
  import Bell from "lucide-svelte/icons/bell";
  import Palette from "lucide-svelte/icons/palette";
  import Mail from "lucide-svelte/icons/mail";
  import Key from "lucide-svelte/icons/key";
  import Coins from "lucide-svelte/icons/coins";
  import LogOut from "lucide-svelte/icons/log-out";
  import Plus from "lucide-svelte/icons/plus";

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
    if (authState.isAuthenticated) fetchUserOrgs();
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
    ...(authState.user?.is_admin ? [{ name: i18n.t("admin_panel"), icon: ShieldCheck, path: "/admin/users", isParent: true }] : []),
    { name: i18n.t("settings_title"), icon: UserRoundCog, path: "/user/settings", isParent: true, isSettings: true },
  ]);

  const adminSubItems = $derived([
    { name: i18n.t("admin.users") || "Users", icon: User, path: "/admin/users" },
    { name: i18n.t("admin.repositories") || "Repositories", icon: BookText, path: "/admin/repos" },
    { name: i18n.t("admin.organizations") || "Organizations", icon: Building2, path: "/admin/orgs" },
    { name: i18n.t("admin.auths") || "Authentication", icon: Lock, path: "/admin/auth" },
    { name: i18n.t("admin.config") || "Configuration", icon: Settings, path: "/admin/config" },
    { name: i18n.t("admin.monitor_label") || "Monitor", icon: Activity, path: "/admin/monitor" },
    { name: i18n.t("admin.maintenance") || "Maintenance", icon: Hammer, path: "/admin/maintenance" },
    { name: i18n.t("admin.audit_log") || "Audit Log", icon: ScrollText, path: "/admin/audit_log" },
    { name: i18n.t("admin.security") || "Seguridad", icon: Shield, path: "/admin/security" },
    { name: i18n.t("admin.janitor") || "Janitor", icon: Trash2, path: "/admin/janitor" },
    { name: "CI/CD", icon: Rocket, path: "/admin/pikoci" },
    { name: i18n.t("admin.notices") || "Notices", icon: Bell, path: "/admin/notices" },
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
</script>

<nav class="navigation-bar">
  <NavUserProfile />
  <NavMainLinks {navItems} {adminSubItems} {settingsSubItems} {userOrgs} showSidebarOrgs={userPrefs.current.showSidebarOrgs} />
  <NavMobileAvatar />
  <NavBrandFooter bind:showThemePanel {pageLoadTime} {templateTime} />
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
  @media (max-width: 1024px) {
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
      justify-content: space-around;
    }
  }
</style>
