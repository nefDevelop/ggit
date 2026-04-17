<script lang="ts">
  import { onMount } from 'svelte';
  import AdminDashboard from './AdminDashboard.svelte';
  import AdminUsers from './AdminUsers.svelte';
  import AdminRepos from './AdminRepos.svelte';
  import AdminAuths from './AdminAuths.svelte';
  import AdminConfig from './AdminConfig.svelte';
  import AdminNotices from './AdminNotices.svelte';
  import AdminMonitor from './AdminMonitor.svelte';

  let activeTab = $state('dashboard');

  const tabs = [
    { id: 'dashboard', label: 'Dashboard', icon: 'dashboard' },
    { id: 'users', label: 'Users', icon: 'person' },
    { id: 'repos', label: 'Repositories', icon: 'book' },
    { id: 'auths', label: 'Authentication', icon: 'security' },
    { id: 'config', label: 'Configuration', icon: 'settings' },
    { id: 'monitor', label: 'Monitor', icon: 'monitoring' },
    { id: 'notices', label: 'System Notices', icon: 'notifications' }
  ];
</script>

<div class="admin-container">
  <aside class="admin-sidebar island">
    <h3 class="sidebar-title">Server Administration</h3>
    <nav class="admin-nav">
      {#each tabs as tab}
        <button 
          class="nav-item" 
          class:active={activeTab === tab.id}
          onclick={() => activeTab = tab.id}
        >
          <span class="material-symbols-rounded">{tab.icon}</span>
          {tab.label}
        </button>
      {/each}
    </nav>
  </aside>

  <main class="admin-content">
    {#if activeTab === 'dashboard'}
      <AdminDashboard />
    {:else if activeTab === 'users'}
      <AdminUsers />
    {:else if activeTab === 'repos'}
      <AdminRepos />
    {:else if activeTab === 'auths'}
      <AdminAuths />
    {:else if activeTab === 'config'}
      <AdminConfig />
    {:else if activeTab === 'monitor'}
      <AdminMonitor />
    {:else if activeTab === 'notices'}
      <AdminNotices />
    {/if}
  </main>
</div>

<style>
  .admin-container {
    display: grid;
    grid-template-columns: 280px 1fr;
    gap: 2rem;
    align-items: start;
    animation: fadeIn 0.3s ease-out;
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(10px); }
    to { opacity: 1; transform: translateY(0); }
  }

  .admin-sidebar {
    padding: 1.5rem;
    position: sticky;
    top: 2rem;
  }

  .sidebar-title {
    font-size: 0.75rem;
    font-weight: 800;
    text-transform: uppercase;
    letter-spacing: 0.05em;
    color: var(--text-muted);
    margin-bottom: 1.5rem;
    padding-left: 0.75rem;
  }

  .admin-nav {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
  }

  .nav-item {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.75rem 1rem;
    background: transparent;
    border: none;
    color: var(--text-muted);
    font-size: 0.95rem;
    font-weight: 600;
    text-align: left;
    cursor: pointer;
    border-radius: 0.75rem;
    transition: all 0.2s;
    width: 100%;
  }

  .nav-item:hover {
    background: rgba(255, 255, 255, 0.03);
    color: var(--text-color);
  }

  .nav-item.active {
    background: rgba(var(--primary-rgb), 0.1);
    color: var(--primary-color);
  }

  .nav-item span {
    font-size: 1.25rem;
  }

  .admin-content {
    display: flex;
    flex-direction: column;
    gap: 2rem;
  }
</style>
