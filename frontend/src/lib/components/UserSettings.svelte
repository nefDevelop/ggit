<script lang="ts">
  import { onMount } from 'svelte';
  import ThemeSelector from './ThemeSelector.svelte';
  import SettingsProfile from './SettingsProfile.svelte';
  import SettingsSecurity from './SettingsSecurity.svelte';
  import SettingsEmails from './SettingsEmails.svelte';
  import SettingsKeys from './SettingsKeys.svelte';
  import SettingsTokens from './SettingsTokens.svelte';
  import SettingsRepositories from './SettingsRepositories.svelte';
  import SettingsOrganizations from './SettingsOrganizations.svelte';

  let activeTab = $state('profile');

  const tabs = [
    { id: 'profile', label: 'Profile', icon: 'person' },
    { id: 'emails', label: 'Emails', icon: 'mail' },
    { id: 'security', label: 'Security', icon: 'security' },
    { id: 'keys', label: 'SSH Keys', icon: 'key' },
    { id: 'tokens', label: 'Access Tokens', icon: 'token' },
    { id: 'repos', label: 'Repositories', icon: 'book' },
    { id: 'orgs', label: 'Organizations', icon: 'corporate_fare' }
  ];
</script>

<div class="settings-container">
  <aside class="settings-sidebar island">
    <h3 class="sidebar-title">User Settings</h3>
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
    {#if activeTab === 'profile'}
      <SettingsProfile />
    {:else if activeTab === 'appearance'}
      <ThemeSelector />
    {:else if activeTab === 'emails'}
      <SettingsEmails />
    {:else if activeTab === 'security'}
      <SettingsSecurity />
    {:else if activeTab === 'keys'}
      <SettingsKeys />
    {:else if activeTab === 'tokens'}
      <SettingsTokens />
    {:else if activeTab === 'repos'}
      <SettingsRepositories />
    {:else if activeTab === 'orgs'}
      <SettingsOrganizations />
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

  .settings-content {
    display: flex;
    flex-direction: column;
    gap: 2rem;
  }
</style>
