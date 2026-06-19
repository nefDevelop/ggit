<script lang="ts">
  import type { RepoContext } from "$lib/types";
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

  

import Settings from "lucide-svelte/icons/settings";
import Users from "lucide-svelte/icons/users";
import GitBranch from "lucide-svelte/icons/git-branch";
import RefreshCw from "lucide-svelte/icons/refresh-cw";
import Webhook from "lucide-svelte/icons/webhook";
import Key from "lucide-svelte/icons/key";
import Terminal from "lucide-svelte/icons/terminal";
import Shield from "lucide-svelte/icons/shield";
import Sparkles from "lucide-svelte/icons/sparkles";



  const repo = getContext<RepoContext>('repo');
  console.log("DEBUG - Settings Page Context:", repo);
  console.log("DEBUG - Settings Page Repo Data:", repo?.data);
  
  let activeTab = $state('general');

  const tabs = [
    { id: 'general', label: 'General', icon: Settings },
    { id: 'collaboration', label: 'Collaboration', icon: Users },
    { id: 'branches', label: 'Branches', icon: GitBranch },
    { id: 'remotes', label: 'Remotes', icon: RefreshCw },
    { id: 'webhooks', label: 'Webhooks', icon: Webhook },
    { id: 'keys', label: 'Deploy Keys', icon: Key },
    { id: 'security', label: 'Seguridad', icon: Shield },
    { id: 'janitor', label: 'Janitor', icon: Sparkles }
  ];

  const visibleTabs = $derived(
    (authState.user?.is_admin || authState.user?.allow_git_hook)
      ? [...tabs, { id: 'git_hooks', label: 'Git Hooks', icon: Terminal }]
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
          {#if typeof tab.icon === 'string'}
            <span class="material-symbols-rounded">{tab.icon}</span>
          {:else if tab.icon}
            {@const Icon = tab.icon}
            <Icon size={18} />
          {/if}
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
      titleIcon={Settings}
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
    margin-top: 1rem;
  }

  .mobile-nav-wrapper {
    display: none;
  }

  @media (max-width: 1024px) {
    .settings-container {
      display: flex;
      flex-direction: column;
      gap: 1rem;
      padding-bottom: calc(160px + env(safe-area-inset-bottom));
    }
    .settings-sidebar {
      display: none;
    }
    .settings-content {
      width: 100%;
    }
    .mobile-nav-wrapper {
      display: block;
      position: fixed;
      bottom: calc(60px + env(safe-area-inset-bottom));
      left: 0;
      right: 0;
      width: 100vw;
      box-sizing: border-box;
      z-index: 999;
      background: var(--md-sys-color-surface-container);
      padding: 0.75rem 1rem;
      border-top: 1px solid var(--md-sys-color-outline-variant);
      box-shadow: 0 -8px 32px rgba(0, 0, 0, 0.1);
      margin-bottom: 0;
    }
  }
</style>
