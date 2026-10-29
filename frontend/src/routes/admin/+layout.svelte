<script lang="ts">
  import { page } from "$app/state";
  import { i18n } from "$lib/i18n.svelte";
  import ChipNavigation from "$lib/components/ui/ChipNavigation.svelte";
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
  import Rocket from "lucide-svelte/icons/rocket";

  let { children } = $props();

  const adminSubItems = [
    { id: "/admin/users", name: i18n.t("admin.users") || "Users", icon: User, path: "/admin/users" },
    { id: "/admin/repos", name: i18n.t("admin.repositories") || "Repositories", icon: BookText, path: "/admin/repos" },
    { id: "/admin/orgs", name: i18n.t("admin.organizations") || "Organizations", icon: Building2, path: "/admin/orgs" },
    { id: "/admin/auth", name: i18n.t("admin.auths") || "Authentication", icon: Lock, path: "/admin/auth" },
    { id: "/admin/config", name: i18n.t("admin.config") || "Configuration", icon: Settings, path: "/admin/config" },
    { id: "/admin/monitor", name: i18n.t("admin.monitor_label") || "Monitor", icon: Activity, path: "/admin/monitor" },
    { id: "/admin/maintenance", name: i18n.t("admin.maintenance") || "Maintenance", icon: Hammer, path: "/admin/maintenance" },
    { id: "/admin/audit_log", name: i18n.t("admin.audit_log") || "Audit Log", icon: ScrollText, path: "/admin/audit_log" },
    { id: "/admin/security", name: "Seguridad", icon: Shield, path: "/admin/security" },
    { id: "/admin/janitor", name: "Janitor", icon: Trash2, path: "/admin/janitor" },
    { id: "/admin/pikoci", name: "CI/CD", icon: Rocket, path: "/admin/pikoci" },
    { id: "/admin/notices", name: i18n.t("admin.notices") || "Notices", icon: Bell, path: "/admin/notices" },
  ].map(item => ({ ...item, label: item.name }));
</script>

<svelte:head>
  <title>{adminSubItems.find(i => i.id === page.url.pathname)?.name || 'Admin'} · Ggit</title>
</svelte:head>

<div class="admin-layout-container">
  <main class="admin-content-full">
    {@render children()}
  </main>

  <div class="admin-mobile-nav">
    <ChipNavigation 
      items={adminSubItems} 
      activeId={page.url.pathname} 
      title={i18n.t("admin_panel") || "Panel de Admin"} 
      titleIcon={Shield}
    />
  </div>
</div>

<style>
  .admin-layout-container {
    padding: 1.5rem 1.5rem 4rem;
    width: 100%;
    box-sizing: border-box;
    margin: 0;
  }

  .admin-mobile-nav {
    display: none;
    margin-bottom: 2rem;
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(10px); }
    to { opacity: 1; transform: translateY(0); }
  }

  .admin-content-full {
    display: flex;
    flex-direction: column;
    gap: 2rem;
    width: 100%;
    animation: fadeIn 0.3s ease-out;
  }

  /* Adjust for mobile if needed */
  @media (max-width: 1024px) {
    .admin-layout-container {
      padding: 1rem 1rem 5rem;
      display: block;
      min-height: auto;
    }
    .admin-content-full {
      width: 100%;
    }
    .admin-mobile-nav {
      display: flex;
      position: fixed;
      bottom: calc(60px + env(safe-area-inset-bottom));
      left: 0;
      right: 0;
      width: 100vw;
      box-sizing: border-box;
      z-index: 500;
      background: rgba(var(--surface-rgb), 0.7);
      backdrop-filter: blur(20px);
      -webkit-backdrop-filter: blur(20px);
      padding: 0.75rem 1rem;
      margin: 0;
      overflow-x: hidden;
    }
  }
</style>
