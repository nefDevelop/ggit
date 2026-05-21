<script lang="ts">
  import { page } from "$app/state";
  import { i18n } from "$lib/i18n.svelte";
  import ChipNavigation from "$lib/components/ui/ChipNavigation.svelte";

  let { children } = $props();

  const adminSubItems = [
    { id: "/admin/users", name: i18n.t("admin.users") || "Users", icon: "person", path: "/admin/users" },
    { id: "/admin/repos", name: i18n.t("admin.repositories") || "Repositories", icon: "book", path: "/admin/repos" },
    { id: "/admin/orgs", name: i18n.t("admin.organizations") || "Organizations", icon: "corporate_fare", path: "/admin/orgs" },
    { id: "/admin/auth", name: i18n.t("admin.auths") || "Authentication", icon: "lock", path: "/admin/auth" },
    { id: "/admin/config", name: i18n.t("admin.config") || "Configuration", icon: "settings", path: "/admin/config" },
    { id: "/admin/monitor", name: i18n.t("admin.monitor") || "Monitor", icon: "monitoring", path: "/admin/monitor" },
    { id: "/admin/maintenance", name: i18n.t("admin.maintenance") || "Maintenance", icon: "build", path: "/admin/maintenance" },
    { id: "/admin/audit_log", name: i18n.t("admin.audit_log") || "Audit Log", icon: "history_edu", path: "/admin/audit_log" },
    { id: "/admin/security", name: "Seguridad", icon: "shield", path: "/admin/security" },
    { id: "/admin/janitor", name: "Janitor", icon: "mop", path: "/admin/janitor" },
    { id: "/admin/notices", name: i18n.t("admin.notices") || "Notices", icon: "notifications", path: "/admin/notices" },
  ].map(item => ({ ...item, label: item.name }));
</script>

<svelte:head>
  <title>{adminSubItems.find(i => i.id === page.url.pathname)?.name || 'Admin'} · Ggit</title>
</svelte:head>

<div class="admin-layout-container" class:edge-to-edge={page.url.pathname.startsWith('/admin')}>
  <main class="admin-content-full">
    {@render children()}
  </main>

  <div class="admin-mobile-nav">
    <ChipNavigation 
      items={adminSubItems} 
      activeId={page.url.pathname} 
      title={i18n.t("admin_panel") || "Panel de Admin"} 
      titleIcon="security"
    />
  </div>
</div>

<style>
  .admin-layout-container {
    padding: 0 1.5rem 4rem;
    width: 100%;
    box-sizing: border-box;
    margin: 0;
  }

  .admin-layout-container.edge-to-edge {
    padding: 0;
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
    .admin-layout-container, .admin-layout-container.edge-to-edge {
      padding: 1rem 0;
      display: block;
      min-height: auto;
    }
    .admin-content-full {
      width: 100%;
      padding: 0.5rem 1.5rem;
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
      border-top: 1px solid var(--md-sys-color-outline-variant);
      box-shadow: 0 -8px 32px rgba(0, 0, 0, 0.1);
      margin: 0;
      overflow-x: hidden;
    }
  }
</style>
