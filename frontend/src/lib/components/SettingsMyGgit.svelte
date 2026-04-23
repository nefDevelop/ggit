<script lang="ts">
  import { userPrefs } from '$lib/settings.svelte';
  import { apiFetch } from '$lib/api';
  import { onMount } from 'svelte';
  import { authState } from '$lib/auth.svelte';
  import { i18n } from '$lib/i18n.svelte';
  import LanguageSelector from './LanguageSelector.svelte';
  import Card from './ui/Card.svelte';
  import Button from './ui/Button.svelte';

  const sections = [
    { id: 'pr', label: i18n.t('pull_requests'), icon: 'merge', desc: i18n.t('settings.pr_desc') },
    { id: 'issue', label: i18n.t('issues'), icon: 'error', desc: i18n.t('settings.issue_desc') },
    { id: 'wiki', label: i18n.t('settings.wiki') || 'Wiki', icon: 'menu_book', desc: i18n.t('settings.wiki_desc') },
    { id: 'branch', label: i18n.t('settings.branches') || 'Branches', icon: 'fork_right', desc: i18n.t('settings.branch_desc') },
    { id: 'commits', label: i18n.t('settings.commits') || 'Commits', icon: 'history', desc: i18n.t('settings.commits_desc') }
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

<Card class="settings-premium-card" padding="lg">
  <header class="section-header">
    <div class="title-with-icon">
      <span class="material-symbols-rounded hero-icon">grid_view</span>
      <div>
        <h2>{i18n.t('settings.my_ggit')}</h2>
        <p class="subtitle">{i18n.t('settings.my_ggit_desc')}</p>
      </div>
    </div>
  </header>

  <div class="modular-stack">
    <!-- REPO MODULES -->
    <div class="config-module">
      <div class="module-tag">
        <span class="material-symbols-rounded">extension</span>
        {i18n.t('settings.repo_tab_modules')}
      </div>
      <div class="island tonal-list">
        {#each sections as section}
          <div class="toggle-row">
            <div class="row-info">
              <div class="row-title">
                <span class="material-symbols-rounded">{section.icon}</span>
                {section.label}
              </div>
              <p class="row-desc">{section.desc}</p>
            </div>
            <label class="premium-switch">
              <input 
                type="checkbox" 
                checked={userPrefs.current.repoTabs[section.id]} 
                onchange={() => userPrefs.toggleRepoTab(section.id)}
              />
              <span class="switch-slider"></span>
            </label>
          </div>
        {/each}
      </div>
    </div>

    <!-- DASHBOARD PREFS -->
    <div class="config-module">
      <div class="module-tag">
        <span class="material-symbols-rounded">dashboard_customize</span>
        {i18n.t('settings.dashboard_preferences')}
      </div>
      <div class="island tonal-list">
        <div class="toggle-row">
          <div class="row-info">
            <div class="row-title">
              <span class="material-symbols-rounded">timeline</span>
              {i18n.t('settings.activity_feed')}
            </div>
            <p class="row-desc">{i18n.t('settings.activity_feed_desc')}</p>
          </div>
          <label class="premium-switch">
            <input 
              type="checkbox" 
              checked={!profile.disable_news_feed} 
              onchange={toggleNewsFeed}
            />
            <span class="switch-slider"></span>
          </label>
        </div>

        <div class="toggle-row">
          <div class="row-info">
            <div class="row-title">
              <span class="material-symbols-rounded">corporate_fare</span>
              {i18n.t('settings.orgs_sidebar')}
            </div>
            <p class="row-desc">{i18n.t('settings.orgs_sidebar_desc')}</p>
          </div>
          <label class="premium-switch">
            <input 
              type="checkbox" 
              checked={userPrefs.current.showSidebarOrgs} 
              onchange={() => {
                userPrefs.current.showSidebarOrgs = !userPrefs.current.showSidebarOrgs;
                userPrefs.save();
              }}
            />
            <span class="switch-slider"></span>
          </label>
        </div>
      </div>
    </div>

    <!-- LANGUAGE -->
    <div class="config-module">
       <div class="module-tag">
         <span class="material-symbols-rounded">translate</span>
         Región e Idioma
       </div>
       <div class="island tonal-list language-card">
          <LanguageSelector />
       </div>
    </div>
  </div>
</Card>

<style>
  :global(.settings-premium-card) { border-radius: 3rem; background: rgba(var(--surface-rgb), 0.5); backdrop-filter: blur(24px); border: 1px solid var(--md-sys-color-outline-variant); }
  
  .section-header { margin-bottom: 4rem; }
  .title-with-icon { display: flex; align-items: center; gap: 2rem; }
  .hero-icon { font-size: 3rem; color: var(--primary-color); background: rgba(var(--primary-rgb), 0.1); padding: 1.25rem; border-radius: 2rem; }
  h2 { margin: 0; font-size: 2.25rem; font-weight: 950; letter-spacing: -0.04em; }
  .subtitle { margin: 0.5rem 0 0 0; color: var(--text-muted); font-size: 1.1rem; }

  .modular-stack { display: flex; flex-direction: column; gap: 4rem; }
  .config-module { display: flex; flex-direction: column; gap: 1.5rem; }
  
  .module-tag { font-size: 0.8rem; font-weight: 950; text-transform: uppercase; color: var(--primary-color); letter-spacing: 0.15em; display: flex; align-items: center; gap: 0.75rem; border-bottom: 2px solid var(--md-sys-color-outline-variant); padding-bottom: 0.75rem; width: fit-content; opacity: 0.8; }
  .module-tag span { font-size: 1.25rem; }

  .tonal-list { display: flex; flex-direction: column; background: rgba(0,0,0,0.15); border-radius: 2.25rem; overflow: hidden; border: 1px solid var(--md-sys-color-outline-variant); }
  
  .toggle-row { display: flex; align-items: center; justify-content: space-between; padding: 2rem 2.5rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); gap: 3rem; transition: background 0.3s; }
  .toggle-row:last-child { border-bottom: none; }
  .toggle-row:hover { background: rgba(255,255,255,0.02); }

  .row-info { flex: 1; display: flex; flex-direction: column; gap: 0.3rem; }
  .row-title { display: flex; align-items: center; gap: 1rem; font-size: 1.15rem; font-weight: 800; color: var(--text-color); }
  .row-title span { color: var(--primary-color); opacity: 0.7; font-size: 1.5rem; }
  .row-desc { font-size: 0.95rem; color: var(--text-muted); margin: 0; line-height: 1.5; }

  .premium-switch { position: relative; display: inline-block; width: 52px; height: 28px; flex-shrink: 0; cursor: pointer; }
  .premium-switch input { opacity: 0; width: 0; height: 0; }
  .switch-slider { position: absolute; inset: 0; background: var(--surface-variant); border-radius: 100px; border: 1px solid var(--md-sys-color-outline-variant); transition: 0.4s cubic-bezier(0.4, 0, 0.2, 1); }
  .switch-slider:before { position: absolute; content: ""; height: 20px; width: 20px; left: 3px; bottom: 3px; background: var(--text-muted); border-radius: 50%; transition: 0.4s cubic-bezier(0.4, 0, 0.2, 1); }
  
  input:checked + .switch-slider { background: var(--primary-color); border-color: var(--primary-color); }
  input:checked + .switch-slider:before { transform: translateX(24px); background: white; }

  .language-card { padding: 1rem; }

  @media (max-width: 600px) {
    .toggle-row { flex-direction: column; align-items: flex-start; gap: 1.5rem; padding: 1.5rem; }
    .premium-switch { align-self: flex-end; }
  }
</style>
