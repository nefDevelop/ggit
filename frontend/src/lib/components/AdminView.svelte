<script lang="ts">
  import { onMount } from "svelte";
  import { i18n } from '$lib/i18n.svelte';
  import AdminDashboard from "./AdminDashboard.svelte";
  import AdminUsers from "./AdminUsers.svelte";
  import AdminRepos from "./AdminRepos.svelte";
  import AdminAuths from "./AdminAuths.svelte";
  import AdminConfig from "./AdminConfig.svelte";
  import AdminNotices from "./AdminNotices.svelte";
  import AdminMonitor from "./AdminMonitor.svelte";
  import AdminOrgs from "./AdminOrgs.svelte";
  import AdminMaintenance from "./AdminMaintenance.svelte";
  import AdminAuditLog from "./AdminAuditLog.svelte";

  let activeTab = $state("dashboard");

  const tabs = [
    { id: "dashboard", label: i18n.t('dashboard'), icon: "dashboard" },
    { id: "users", label: i18n.t('admin.users'), icon: "person" },
    { id: "repos", label: i18n.t('admin.repositories'), icon: "book" },
    { id: "orgs", label: i18n.t('admin.organizations'), icon: "corporate_fare" },
    { id: "auths", label: i18n.t('admin.auths'), icon: "lock" },
    { id: "config", label: i18n.t('admin.config'), icon: "settings" },
    { id: "monitor", label: i18n.t('admin.monitor'), icon: "monitoring" },
    { id: "maintenance", label: i18n.t('admin.maintenance'), icon: "build" },
    { id: "audit", label: i18n.t('admin.audit_log'), icon: "history_edu" },
    { id: "notices", label: i18n.t('admin.notices'), icon: "notifications" },
  ];
</script>

<div class="admin-container">
  <aside class="admin-sidebar">
    <div class="sidebar-brand">
       <span class="material-symbols-rounded">security</span>
       <h3 class="sidebar-title">{i18n.t('admin_panel')}</h3>
    </div>
    <nav class="admin-nav">
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

  <main class="admin-content">
    {#if activeTab === "dashboard"}
      <AdminDashboard />
    {:else if activeTab === "users"}
      <AdminUsers />
    {:else if activeTab === "repos"}
      <AdminRepos />
    {:else if activeTab === "orgs"}
      <AdminOrgs />
    {:else if activeTab === "auths"}
      <AdminAuths />
    {:else if activeTab === "config"}
      <AdminConfig />
    {:else if activeTab === "monitor"}
      <AdminMonitor />
    {:else if activeTab === "maintenance"}
      <AdminMaintenance />
    {:else if activeTab === "audit"}
      <AdminAuditLog />
    {:else if activeTab === "notices"}
      <AdminNotices />
    {/if}
  </main>
</div>

<style>
  .admin-container {
    display: grid;
    grid-template-columns: 320px 1fr;
    gap: 3rem;
    align-items: start;
    padding: 2.5rem;
    min-height: 100vh;
  }

  .admin-sidebar {
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

  .admin-nav {
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

  .admin-content {
    display: flex;
    flex-direction: column;
    gap: 3rem;
    animation: contentFade 0.5s ease-out;
    min-width: 0;
  }

  @keyframes contentFade {
    from { opacity: 0; transform: translateY(20px); }
    to { opacity: 1; transform: translateY(0); }
  }

  @media (max-width: 1200px) {
    .admin-container {
      grid-template-columns: 1fr;
      padding: 1.5rem;
      gap: 2rem;
    }
    .admin-sidebar {
      position: relative;
      top: 0;
      gap: 1.5rem;
      z-index: 10;
    }
    .sidebar-brand { display: none; }
    .admin-nav {
      flex-direction: row;
      overflow-x: auto;
      padding-bottom: 0.5rem;
      margin: 0 -1.5rem;
      padding-left: 1.5rem;
      scrollbar-width: none;
      -webkit-overflow-scrolling: touch;
    }
    .admin-nav::-webkit-scrollbar { display: none; }
    .nav-item {
      width: auto;
      padding: 0.75rem 1.25rem;
      border-radius: 100px;
      flex-shrink: 0;
      gap: 0.5rem;
      font-size: 0.95rem;
    }
    .nav-item:hover { transform: none; }
    .nav-item span.material-symbols-rounded { font-size: 1.25rem; }
    .active-indicator {
      display: none;
    }
  }

  @media (max-width: 600px) {
    .sidebar-brand span { font-size: 1.25rem; }
    .sidebar-title { font-size: 1rem; }
  }
</style>
