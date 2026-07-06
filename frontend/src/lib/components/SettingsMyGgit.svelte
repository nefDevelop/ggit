<script lang="ts">
  import { userPrefs } from '$lib/settings.svelte';
  import { apiFetch } from '$lib/api';
  import { onMount } from 'svelte';
  import { authState } from '$lib/auth.svelte';
  import { i18n } from '$lib/i18n.svelte';
  import LanguageSelector from './LanguageSelector.svelte';
  import Card from './ui/Card.svelte';
  import Button from './ui/Button.svelte';

  type RepoTabKey = keyof (typeof userPrefs.current.repoTabs);

  const sections = $derived<{ id: RepoTabKey, label: string, icon: string, desc: string }[]>([
    { id: 'pr', label: i18n.t('pull_requests'), icon: 'merge', desc: i18n.t('settings.pr_desc') },
    { id: 'issue', label: i18n.t('issues'), icon: 'error', desc: i18n.t('settings.issue_desc') },
    { id: 'wiki', label: i18n.t('settings.wiki') || 'Wiki', icon: 'menu_book', desc: i18n.t('settings.wiki_desc') },
    { id: 'branch', label: i18n.t('settings.branches') || 'Branches', icon: 'fork_right', desc: i18n.t('settings.branch_desc') },
    { id: 'commits', label: i18n.t('settings.commits') || 'Commits', icon: 'history', desc: i18n.t('settings.commits_desc') },
    { id: 'tags', label: i18n.t('settings.tags') || 'Tags', icon: 'sell', desc: i18n.t('settings.tags_desc') || 'Gestiona versiones y etiquetas del código' },
    { id: 'milestones', label: i18n.t('settings.milestones') || 'Milestones', icon: 'event', desc: i18n.t('settings.milestones_desc') || 'Organiza el trabajo en hitos temporales' },
    { id: 'releases', label: i18n.t('settings.releases') || 'Releases', icon: 'rocket_launch', desc: i18n.t('settings.releases_desc') || 'Gestiona las versiones publicadas y binarios' },
    { id: 'labels', label: i18n.t('settings.labels') || 'Labels', icon: 'label', desc: i18n.t('settings.labels_desc') || 'Clasifica issues y PRs con etiquetas personalizadas' }
  ]);

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

<div class="settings-myggit-container">

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

        <div class="toggle-row">
          <div class="row-info">
            <div class="row-title">
              <span class="material-symbols-rounded">bar_chart</span>
              Barra de Lenguajes
            </div>
            <p class="row-desc">Muestra la barra de lenguajes en la lista de repositorios del dashboard.</p>
          </div>
          <label class="premium-switch">
            <input 
              type="checkbox" 
              checked={userPrefs.current.showLangBar} 
              onchange={() => {
                userPrefs.current.showLangBar = !userPrefs.current.showLangBar;
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
</div>

<style>
  


  .modular-stack { display: flex; flex-direction: column; gap: 3rem; }
  @media (max-width: 600px) {
    .modular-stack { gap: 2rem; }
  }
  .config-module { display: flex; flex-direction: column; gap: 1.5rem; }
  
  .module-tag { font-size: 0.8rem; font-weight: 950; text-transform: uppercase; color: var(--primary-color); letter-spacing: 0.15em; display: flex; align-items: center; gap: 0.75rem; border-bottom: 2px solid var(--md-sys-color-outline-variant); padding-bottom: 0.75rem; width: fit-content; opacity: 0.8; }
  .module-tag span { font-size: 1.25rem; }

  .tonal-list { display: flex; flex-direction: column; background: rgba(0,0,0,0.15); border-radius: 1.5rem; overflow: hidden; border: 1px solid var(--md-sys-color-outline-variant); }
  
  .toggle-row { display: flex; align-items: center; justify-content: space-between; padding: 2rem 2.5rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); gap: 3rem; transition: background 0.3s; }
  @media (max-width: 600px) { .toggle-row { flex-direction: column; align-items: stretch; padding: 1rem 1.25rem; gap: 1rem; } }
  .toggle-row:last-child { border-bottom: none; }
  .toggle-row:hover { background: rgba(255,255,255,0.02); }
  .toggle-row:active { background: rgba(255,255,255,0.05); }

  .row-info { flex: 1; display: flex; flex-direction: column; gap: 0.3rem; }
  .row-title { display: flex; align-items: center; gap: 1rem; font-size: 1.15rem; font-weight: 800; color: var(--text-color); }
  .row-title span { color: var(--primary-color); opacity: 0.7; font-size: 1.5rem; }
  .row-desc { font-size: 0.95rem; color: var(--text-muted); margin: 0; line-height: 1.5; }

  .premium-switch { position: relative; display: inline-block; width: 52px; height: 28px; flex-shrink: 0; cursor: pointer; transition: transform 0.2s; }
  .premium-switch:active { transform: scale(0.9); }
  .premium-switch input { opacity: 0; width: 0; height: 0; position: absolute; }
  .switch-slider { position: absolute; inset: 0; background: var(--surface-variant); border-radius: 100px; border: 1px solid var(--md-sys-color-outline-variant); transition: 0.4s cubic-bezier(0.4, 0, 0.2, 1); }
  .switch-slider:before { position: absolute; content: ""; height: 20px; width: 20px; left: 3px; bottom: 3px; background: var(--text-muted); border-radius: 50%; transition: 0.4s cubic-bezier(0.4, 0, 0.2, 1); }
  
  input:checked + .switch-slider { background: var(--primary-color); border-color: var(--primary-color); }
  input:checked + .switch-slider:before { transform: translateX(24px); background: white; }

  .language-card { padding: 1rem; }

  @media (max-width: 600px) {
    .premium-switch { align-self: flex-end; }
  }
</style>
