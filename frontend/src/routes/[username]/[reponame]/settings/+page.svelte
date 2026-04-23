<script lang="ts">
  import { getContext } from 'svelte';
  import RepoSettingsGeneral from '$lib/components/RepoSettingsGeneral.svelte';
  import RepoSettingsCollaboration from '$lib/components/RepoSettingsCollaboration.svelte';
  import RepoSettingsWebhooks from '$lib/components/RepoSettingsWebhooks.svelte';
  import RepoSettingsDeployKeys from '$lib/components/RepoSettingsDeployKeys.svelte';
  import RepoSettingsRemotes from '$lib/components/RepoSettingsRemotes.svelte';

  const repo: any = getContext('repo');
  let activeTab = $state('general');

  const tabs = [
    { id: 'general', label: 'General', icon: 'settings' },
    { id: 'collaboration', label: 'Collaboration', icon: 'group' },
    { id: 'remotes', label: 'Remotes', icon: 'sync' },
    { id: 'webhooks', label: 'Webhooks', icon: 'webhook' },
    { id: 'keys', label: 'Deploy Keys', icon: 'vpn_key' }
  ];
</script>

<svelte:head>
  <title>Settings · {repo?.data?.name} · Ggit</title>
</svelte:head>

<div class="settings-container">
  <aside class="settings-sidebar tonal">
    <h3 class="sidebar-title">Repository Settings</h3>
    <nav class="settings-nav">
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

  <main class="settings-content">
    {#if activeTab === 'general'}
      <RepoSettingsGeneral />
    {:else if activeTab === 'collaboration'}
        <RepoSettingsCollaboration />
    {:else if activeTab === 'remotes'}
        <RepoSettingsRemotes />
    {:else if activeTab === 'webhooks'}
        <RepoSettingsWebhooks />
    {:else if activeTab === 'keys'}
        <RepoSettingsDeployKeys />
    {/if}
  </main>
</div>

<style>
  .settings-container {
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

  .settings-sidebar {
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

  .settings-nav {
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

  .nav-item:hover:not(.disabled) {
    background: rgba(255, 255, 255, 0.03);
    color: var(--text-color);
  }

  .nav-item.active {
    background: rgba(var(--primary-rgb), 0.1);
    color: var(--primary-color);
  }

  .nav-item.disabled {
    opacity: 0.4;
    cursor: default;
    position: relative;
  }

  .nav-item .badge {
    margin-left: auto;
    background: rgba(255,255,255,0.1);
    font-size: 0.65rem;
    padding: 0.1rem 0.4rem;
    border-radius: 4px;
    font-weight: 800;
  }

  .nav-item span {
    font-size: 1.25rem;
  }

  .settings-content {
    display: flex;
    flex-direction: column;
    gap: 2rem;
  }
</style>
