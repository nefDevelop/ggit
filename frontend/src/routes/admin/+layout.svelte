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
    { id: "/admin/notices", name: i18n.t("admin.notices") || "Notices", icon: "notifications", path: "/admin/notices" },
  ].map(item => ({ ...item, label: item.name }));
</script>

<div class="admin-layout-container">
  <div class="admin-mobile-nav">
    <ChipNavigation 
      items={adminSubItems} 
      activeId={page.url.pathname} 
      title={i18n.t("admin_panel") || "Panel de Admin"} 
      titleIcon="security"
    />
  </div>

  <main class="admin-content-full">
    {@render children()}
  </main>
</div>

<style>
  .admin-layout-container {
    animation: fadeIn 0.3s ease-out;
    padding: 2rem 4rem;
    max-width: 1600px;
    margin: 0 auto;
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
  }

  /* Adjust for mobile if needed */
  @media (max-width: 1024px) {
    .admin-layout-container {
      padding: 1rem;
      padding-bottom: calc(160px + env(safe-area-inset-bottom));
    }
    .admin-mobile-nav {
      display: flex;
      position: fixed;
      bottom: calc(70px + env(safe-area-inset-bottom));
      left: 0;
      right: 0;
      z-index: 999;
      background: var(--bg-color);
      padding: 0.5rem 0.5rem 0;
      border-top: 1px solid var(--md-sys-color-outline-variant);
      box-shadow: 0 -10px 30px rgba(0, 0, 0, 0.1);
      margin-bottom: 0;
    }
  }
</style>
