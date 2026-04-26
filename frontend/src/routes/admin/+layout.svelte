<script lang="ts">
  import { page } from "$app/state";
  import { i18n } from "$lib/i18n.svelte";

  let { children } = $props();

  const adminSubItems = [
    { name: i18n.t("admin.users") || "Users", icon: "person", path: "/admin/users" },
    { name: i18n.t("admin.repositories") || "Repositories", icon: "book", path: "/admin/repos" },
    { name: i18n.t("admin.organizations") || "Organizations", icon: "corporate_fare", path: "/admin/orgs" },
    { name: i18n.t("admin.auths") || "Authentication", icon: "lock", path: "/admin/auth" },
    { name: i18n.t("admin.config") || "Configuration", icon: "settings", path: "/admin/config" },
    { name: i18n.t("admin.monitor") || "Monitor", icon: "monitoring", path: "/admin/monitor" },
    { name: i18n.t("admin.maintenance") || "Maintenance", icon: "build", path: "/admin/maintenance" },
    { name: i18n.t("admin.audit_log") || "Audit Log", icon: "history_edu", path: "/admin/audit_log" },
    { name: i18n.t("admin.notices") || "Notices", icon: "notifications", path: "/admin/notices" },
  ];
</script>

<div class="admin-layout-container">
  <nav class="admin-mobile-nav">
    {#each adminSubItems as item}
      <a href={item.path} class="mobile-nav-item" class:active={page.url.pathname === item.path}>
        <span class="material-symbols-rounded">{item.icon}</span>
        <span class="label">{item.name}</span>
      </a>
    {/each}
  </nav>

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
    gap: 0.5rem;
    overflow-x: auto;
    padding: 0.5rem 0 1.5rem;
    margin: 0 -1rem;
    padding-left: 1rem;
    scrollbar-width: none;
    -webkit-overflow-scrolling: touch;
  }
  .admin-mobile-nav::-webkit-scrollbar { display: none; }

  .mobile-nav-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    padding: 0.6rem 1rem;
    border-radius: 100px;
    background: rgba(255,255,255,0.05);
    color: var(--text-muted);
    text-decoration: none;
    font-size: 0.85rem;
    font-weight: 700;
    white-space: nowrap;
    border: 1px solid var(--md-sys-color-outline-variant);
    transition: all 0.2s;
  }

  .mobile-nav-item.active {
    background: var(--primary-color);
    color: white;
    border-color: var(--primary-color);
  }

  .mobile-nav-item span { font-size: 1.1rem; }

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
    }
    .admin-mobile-nav {
      display: flex;
    }
  }
</style>
