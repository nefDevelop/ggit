<script lang="ts">
  // Refreshed preferences types
  import { userPrefs, type UserPreferences } from '$lib/settings.svelte';
  import { apiFetch } from '$lib/api';
  import { onMount } from 'svelte';
  import { authState } from '$lib/auth.svelte';
  import { i18n } from '$lib/i18n.svelte';
  import LanguageSelector from './LanguageSelector.svelte';

  const sections = [
    { id: 'pr', label: i18n.t('pull_requests'), icon: 'merge', desc: i18n.t('settings.pr_desc') || 'Display and manage code review requests.' },
    { id: 'issue', label: i18n.t('issues'), icon: 'error', desc: i18n.t('settings.issue_desc') || 'Track bugs, enhancements, and tasks.' },
    { id: 'wiki', label: i18n.t('settings.wiki') || 'Wiki', icon: 'menu_book', desc: i18n.t('settings.wiki_desc') || 'Host documentation for your project.' },
    { id: 'branch', label: i18n.t('settings.branches') || 'Branches', icon: 'fork_right', desc: i18n.t('settings.branch_desc') || 'View and manage development branches.' },
    { id: 'commits', label: i18n.t('settings.commits') || 'Commits', icon: 'history', desc: i18n.t('settings.commits_desc') || 'Browse the full history of changes.' }
  ] as const;

  let profile = $state<any>({ disable_news_feed: false });

  async function fetchProfile() {
    try {
      const res = await apiFetch('/api/v2/user/settings/profile');
      if (res.ok) profile = await res.json();
    } catch (e) {
      console.error(e);
    }
  }

  async function toggleNewsFeed() {
    profile.disable_news_feed = !profile.disable_news_feed;
    try {
      await apiFetch('/api/v2/user/settings/profile', {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
           name: authState.user?.username,
           disable_news_feed: profile.disable_news_feed
        })
      });
      if (authState.user) authState.user.disable_news_feed = profile.disable_news_feed;
    } catch (e) {
      console.error(e);
    }
  }

  onMount(fetchProfile);
</script>

<div class="settings-section">
  <header class="section-header-modern">
    <div class="header-info">
      <h3>{i18n.t('settings.my_ggit') || 'My Ggit'}</h3>
      <p>{i18n.t('settings.my_ggit_desc') || 'Customize your repository and dashboard experience.'}</p>
    </div>
  </header>

  <div class="switches-container">
    <div class="switches-group">
      <h4 class="group-title">{i18n.t('settings.repo_tab_modules') || 'Repository Tab Modules'}</h4>
      <div class="switches-grid">
        {#each sections as section}
          <div class="switch-row-item">
            <div class="item-info">
              <div class="title-with-icon">
                <span class="material-symbols-rounded mini">{section.icon}</span>
                <span class="label">{section.label}</span>
              </div>
              <p class="desc">{section.desc}</p>
            </div>
            
            <label class="precision-switch">
              <input 
                type="checkbox" 
                checked={userPrefs.current.repoTabs[section.id]} 
                onchange={() => userPrefs.toggleRepoTab(section.id)}
              />
              <span class="slider"></span>
            </label>
          </div>
        {/each}
      </div>
    </div>

    <div class="switches-group">
      <h4 class="group-title">{i18n.t('settings.dashboard_preferences') || 'Dashboard Preferences'}</h4>
      <div class="switches-grid">
        <div class="switch-row-item">
          <div class="item-info">
            <div class="title-with-icon">
              <span class="material-symbols-rounded mini">timeline</span>
              <span class="label">{i18n.t('settings.activity_feed') || 'Activity Feed'}</span>
            </div>
            <p class="desc">{i18n.t('settings.activity_feed_desc') || 'Show pushes, merges and relevant repository activity.'}</p>
          </div>
          
          <label class="precision-switch">
            <input 
              type="checkbox" 
              checked={!profile.disable_news_feed} 
              onchange={toggleNewsFeed}
            />
            <span class="slider"></span>
          </label>
        </div>

        <div class="switch-row-item">
          <div class="item-info">
            <div class="title-with-icon">
              <span class="material-symbols-rounded mini">corporate_fare</span>
              <span class="label">Mostrar Organizaciones en Sidebar</span>
            </div>
            <p class="desc">Muestra u oculta la sección de tus organizaciones en la navegación lateral.</p>
          </div>
          
          <label class="precision-switch">
            <input 
              type="checkbox" 
              checked={userPrefs.current.showSidebarOrgs} 
              onchange={() => {
                userPrefs.current.showSidebarOrgs = !userPrefs.current.showSidebarOrgs;
                userPrefs.save();
              }}
            />
            <span class="slider"></span>
          </label>
        </div>
      </div>
    </div>

    <LanguageSelector />
  </div>
</div>

<style>
  .settings-section { display: flex; flex-direction: column; gap: 2.5rem; }
  
  .section-header-modern { border-bottom: 1px solid var(--md-sys-color-outline-variant); padding-bottom: 1.5rem; }
  .section-header-modern h3 { margin: 0; font-size: 1.5rem; font-weight: 900; }
  .section-header-modern p { margin: 0.5rem 0 0 0; color: var(--text-muted); font-size: 0.95rem; }

  .switches-container { display: flex; flex-direction: column; gap: 3rem; }
  .switches-group { display: flex; flex-direction: column; gap: 1rem; }
  .group-title { font-size: 0.75rem; font-weight: 800; text-transform: uppercase; color: var(--text-muted); letter-spacing: 0.05em; padding-left: 0.5rem; }

  .switches-grid { display: flex; flex-direction: column; background: var(--surface-variant); border: 1px solid var(--md-sys-color-outline-variant); border-radius: 1.25rem; overflow: hidden; box-shadow: var(--shadow-sm); }
  
  .switch-row-item { display: flex; align-items: center; justify-content: space-between; padding: 1.5rem 2rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); gap: 2rem; }
  .switch-row-item:last-child { border-bottom: none; }
  .switch-row-item:hover { background: rgba(255,255,255,0.02); }

  .item-info { flex: 1; display: flex; flex-direction: column; gap: 0.25rem; }
  .title-with-icon { display: flex; align-items: center; gap: 1rem; font-weight: 700; color: var(--text-color); }
  .desc { font-size: 0.85rem; color: var(--text-muted); margin: 0; opacity: 0.8; }
  .mini { font-size: 1.25rem; color: var(--primary-color); }

  /* Precision Switch Styles */
  .precision-switch { position: relative; display: inline-block; width: 44px; height: 24px; flex-shrink: 0; }
  .precision-switch input { opacity: 0; width: 0; height: 0; }
  .slider { position: absolute; cursor: pointer; inset: 0; background-color: rgba(255,255,255,0.1); transition: .3s cubic-bezier(0.4, 0, 0.2, 1); border-radius: 20px; border: 1px solid var(--md-sys-color-outline-variant); }
  .slider:before { position: absolute; content: ""; height: 16px; width: 16px; left: 3px; bottom: 3px; background-color: var(--text-muted); transition: .3s cubic-bezier(0.4, 0, 0.2, 1); border-radius: 50%; }

  input:checked + .slider { background-color: var(--primary-color); border-color: var(--primary-color); }
  input:checked + .slider:before { transform: translateX(20px); background-color: var(--on-primary); }

  @media (max-width: 600px) {
    .switch-row-item { flex-direction: column; align-items: flex-start; gap: 1.5rem; }
    .precision-switch { align-self: flex-end; }
  }
</style>
