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
  import AdminSecurity from "./AdminSecurity.svelte";
  import PageHeader from "./ui/PageHeader.svelte";
  import ViewHeader from "./ui/ViewHeader.svelte";

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
    { id: "security", label: i18n.t('admin.security'), icon: "shield_check" },
    { id: "audit", label: i18n.t('admin.audit_log'), icon: "history_edu" },
    { id: "notices", label: i18n.t('admin.notices'), icon: "notifications" },
  ];
</script>

<div class="admin-container">
  <PageHeader 
    title={i18n.t('admin_panel') || "Panel de Administración"}
    subtitle="Administra los usuarios, organizaciones, repositorios y configuraciones globales"
    icon="security"
  />

  <ViewHeader 
    tabs={tabs} 
    bind:activeTab={activeTab} 
    showSearch={false}
  />

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
    {:else if activeTab === "security"}
      <AdminSecurity />
    {:else if activeTab === "audit"}
      <AdminAuditLog />
    {:else if activeTab === "notices"}
      <AdminNotices />
    {/if}
  </main>
</div>

<style>
  .admin-container {
    display: flex;
    flex-direction: column;
    gap: 0;
    min-height: 100vh;
  }

  .admin-content {
    display: flex;
    flex-direction: column;
    padding: 1.5rem 1.5rem 4rem;
    animation: contentFade 0.4s cubic-bezier(0.4, 0, 0.2, 1);
    min-width: 0;
  }

  @keyframes contentFade {
    from { opacity: 0; transform: translateY(10px); }
    to { opacity: 1; transform: translateY(0); }
  }
</style>
