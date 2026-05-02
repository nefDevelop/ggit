<script lang="ts">
  import { getContext } from 'svelte';
  import RepoSettingsGeneral from '$lib/components/RepoSettingsGeneral.svelte';
  import RepoSettingsCollaboration from '$lib/components/RepoSettingsCollaboration.svelte';
  import RepoSettingsBranches from '$lib/components/RepoSettingsBranches.svelte';
  import RepoSettingsWebhooks from '$lib/components/RepoSettingsWebhooks.svelte';
  import RepoSettingsDeployKeys from '$lib/components/RepoSettingsDeployKeys.svelte';
  import RepoSettingsRemotes from '$lib/components/RepoSettingsRemotes.svelte';
  import RepoSettingsGitHooks from '$lib/components/RepoSettingsGitHooks.svelte';
  import RepoSettingsSecurity from '$lib/components/RepoSettingsSecurity.svelte';
  import RepoSettingsJanitor from '$lib/components/RepoSettingsJanitor.svelte';
  import ChipNavigation from '$lib/components/ui/ChipNavigation.svelte';
  import { authState } from '$lib/auth.svelte';

  const repo: any = getContext('repo');
  let activeTab = $state('general');

  const tabs = [
    { id: 'general', label: 'General', icon: 'settings' },
    { id: 'collaboration', label: 'Collaboration', icon: 'group' },
    { id: 'branches', label: 'Branches', icon: 'call_split' },
    { id: 'remotes', label: 'Remotes', icon: 'sync' },
    { id: 'webhooks', label: 'Webhooks', icon: 'webhook' },
    { id: 'keys', label: 'Deploy Keys', icon: 'vpn_key' },
    { id: 'security', label: 'Seguridad', icon: 'shield' },
    { id: 'janitor', label: 'Janitor', icon: 'mop' }
  ];

  const visibleTabs = $derived(
    (authState.user?.is_admin || authState.user?.allow_git_hook)
      ? [...tabs, { id: 'git_hooks', label: 'Git Hooks', icon: 'terminal' }]
      : tabs
  );
</script>

<svelte:head>
  <title>Settings · {repo?.data?.name} · Ggit</title>
</svelte:head>

<div class="settings-container">
  <aside class="settings-sidebar tonal">
    <h3 class="sidebar-title">Repository Settings</h3>
    <nav class="settings-nav">
      {#each visibleTabs as tab}
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
    {:else if activeTab === 'branches'}
        <RepoSettingsBranches />
    {:else if activeTab === 'remotes'}
        <RepoSettingsRemotes />
    {:else if activeTab === 'webhooks'}
        <RepoSettingsWebhooks />
    {:else if activeTab === 'keys'}
        <RepoSettingsDeployKeys />
    {:else if activeTab === 'git_hooks'}
        <RepoSettingsGitHooks />
    {:else if activeTab === 'security'}
        <RepoSettingsSecurity />
    {:else if activeTab === 'janitor'}
        <RepoSettingsJanitor />
    {/if}
  </main>

  <div class="mobile-nav-wrapper">
    <ChipNavigation 
      items={visibleTabs} 
      activeId={activeTab} 
      title="Repository Settings"
      titleIcon="settings"
      onSelect={(id) => activeTab = id} 
    />
  </div>
</div>

<style>
  .settings-container {
    display: grid;
    grid-template-columns: 280px 1fr;
    gap: 2rem;
    align-items: start;
    padding-bottom: 2rem;
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
    border-radius: 2rem;
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

  .nav-item span {
    font-size: 1.25rem;
  }

  .settings-content {
    display: flex;
    flex-direction: column;
    gap: 2rem;
  }

  .mobile-nav-wrapper {
    display: none;
  }

  @media (max-width: 1024px) {
    .settings-container {
      display: flex;
      flex-direction: column;
      gap: 1rem;
      padding-bottom: 2rem;
    }
    .settings-sidebar {
      display: none;
    }
    .settings-content {
      width: 100%;
    }
    .mobile-nav-wrapper {
      display: block;
      position: sticky;
      bottom: calc(70px + env(safe-area-inset-bottom));
      z-index: 999;
      background: var(--bg-color);
      padding: 0.5rem 0.5rem 0;
      border-top: none;
      box-shadow: 0 -10px 30px rgba(0, 0, 0, 0.1);
      margin-top: auto;
    }
  }
</style>
