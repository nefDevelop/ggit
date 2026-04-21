<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { authState } from '$lib/auth.svelte';
  import { i18n } from '$lib/i18n.svelte';
  import { page } from '$app/state';
  import ThemeSwitch from '$lib/components/ThemeSwitch.svelte';

  let showThemePanel = $state(false);

  const navItems = $derived([
    { name: i18n.t('dashboard'), icon: 'grid_view', path: '/dashboard' },
    { name: i18n.t('explore'), icon: 'search', path: '/explore/repos' },
    { name: i18n.t('organization'), icon: 'corporate_fare', path: '/explore/users' },
    ...(authState.user?.is_admin ? [
      { name: i18n.t('admin_panel'), icon: 'security', path: '/admin/users', isParent: true }
    ] : []),
    { name: i18n.t('settings'), icon: 'settings', path: '/user/settings/profile' },
  ]);

  const adminSubItems = $derived([
    { name: i18n.t('admin.users') || 'Users', icon: 'person', path: '/admin/users' },
    { name: i18n.t('admin.repositories') || 'Repositories', icon: 'book', path: '/admin/repos' },
    { name: i18n.t('admin.organizations') || 'Organizations', icon: 'corporate_fare', path: '/admin/orgs' },
    { name: i18n.t('admin.auths') || 'Authentication', icon: 'lock', path: '/admin/auth' },
    { name: i18n.t('admin.config') || 'Configuration', icon: 'settings', path: '/admin/config' },
    { name: i18n.t('admin.monitor') || 'Monitor', icon: 'monitoring', path: '/admin/monitor' },
    { name: i18n.t('admin.maintenance') || 'Maintenance', icon: 'build', path: '/admin/maintenance' },
    { name: i18n.t('admin.audit_log') || 'Audit Log', icon: 'history_edu', path: '/admin/audit_log' },
    { name: i18n.t('admin.notices') || 'Notices', icon: 'notifications', path: '/admin/notices' }
  ]);

  const isAdminPath = $derived(page.url.pathname.startsWith('/admin'));

  function toggleLanguage() {
    const next = i18n.locale === 'en-US' ? 'es-ES' : 'en-US';
    i18n.setLocale(next);
  }

  function getCsrfToken() {
    if (typeof document === 'undefined') return '';
    const match = document.cookie.match(new RegExp('(^| )XSRF-TOKEN=([^;]+)'));
    return match ? match[2] : '';
  }

  async function handleLogout() {
    await apiFetch('/api/v2/user/logout', { 
      method: 'POST'
    });
    window.location.href = '/user/login';
  }
</script>

<nav class="navigation-bar">
  <!-- USER AT TOP -->
  {#if authState.isAuthenticated}
    <div class="user-profile-summary">
      <div class="avatar-ring">
        <img src={authState.user?.avatar_url || authState.user?.avatar} alt="Avatar" class="sidebar-avatar" />
      </div>
      <div class="user-details">
        <span class="full-name">{authState.user?.full_name || authState.user?.username}</span>
        <span class="user-handle">@{authState.user?.username}</span>
      </div>
      <button class="minimal-logout" onclick={handleLogout} title="Logout">
        <span class="material-symbols-rounded">logout</span>
      </button>
    </div>
  {/if}

  <div class="nav-links">
    {#each navItems as item}
      <a 
        href={item.path} 
        class="nav-item" 
        class:active={page.url.pathname === item.path || (item.isParent && isAdminPath)}
      >
        <span class="material-symbols-rounded">{item.icon}</span>
        <span class="item-label">{item.name}</span>
      </a>

      {#if item.isParent && isAdminPath}
        <div class="sub-nav-links">
          {#each adminSubItems as sub}
            <a 
              href={sub.path} 
              class="sub-nav-item" 
              class:active={page.url.pathname === sub.path}
            >
              <span class="material-symbols-rounded">{sub.icon}</span>
              <span class="item-label">{sub.name}</span>
            </a>
          {/each}
        </div>
      {/if}
    {/each}
  </div>

  <!-- THEME SELECTOR -->
  <div class="theme-trigger-section">
    <button class="appearance-btn" onclick={() => showThemePanel = !showThemePanel}>
      <span class="material-symbols-rounded">palette</span>
      <span class="item-label">Apariencia</span>
    </button>
  </div>

  <!-- BRAND AT BOTTOM -->
  <div class="brand-section-bottom">
    <button class="lang-toggle-btn" onclick={toggleLanguage} title="Change Language">
        <span class="material-symbols-rounded">language</span>
        <span class="lang-text">{i18n.locale === 'en-US' ? 'EN' : 'ES'}</span>
    </button>
    <div class="brand-info">
        <span class="brand-name">Ggit</span>
        <span class="version">v2026.4.0</span>
    </div>
  </div>
</nav>

<ThemeSwitch bind:show={showThemePanel} />

<style>
  .navigation-bar {
    display: flex;
    flex-direction: column;
    padding: 2rem 1rem;
    width: 260px;
    height: 100vh;
    box-sizing: border-box;
    z-index: 1000;
    /* No borders, no islands, just flat integration with background */
  }

  /* USER PROFILE SUMMARY */
  .user-profile-summary {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 1rem;
    margin-bottom: 2rem;
    background: rgba(255, 255, 255, 0.03);
    border-radius: 1.5rem;
    border: 1px solid rgba(255, 255, 255, 0.05);
    transition: all 0.3s ease;
  }
  .user-profile-summary:hover {
    background: rgba(255, 255, 255, 0.06);
    transform: translateY(-2px);
  }

  .avatar-ring {
    display: flex;
    align-items: center;
    justify-content: center;
    position: relative;
    width: 48px;
    height: 48px;
    border-radius: 50%;
    background: linear-gradient(135deg, var(--primary-color), #8b5cf6);
    flex-shrink: 0;
  }

  .sidebar-avatar {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    object-fit: cover;
    border: 2px solid var(--bg-color);
  }

  .user-details {
    flex: 1;
    display: flex;
    flex-direction: column;
    min-width: 0;
  }

  .full-name {
    font-size: 0.9rem;
    font-weight: 800;
    color: var(--text-color);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .user-handle {
    font-size: 0.75rem;
    color: var(--text-muted);
    font-weight: 600;
  }

  .minimal-logout {
    background: transparent;
    border: none;
    color: var(--text-muted);
    cursor: pointer;
    padding: 0.5rem;
    display: flex;
    border-radius: 0.75rem;
    transition: all 0.2s;
  }
  .minimal-logout:hover {
    color: #ef4444;
    background: rgba(239, 68, 68, 0.1);
  }

  /* LINKS */
  .nav-links {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
    flex: 1;
  }

  .nav-item {
    display: flex;
    align-items: center;
    gap: 1.25rem;
    padding: 0.85rem 1.25rem;
    color: var(--text-muted);
    text-decoration: none;
    border-radius: 1rem;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    font-weight: 600;
    font-size: 0.95rem;
  }

  .nav-item:hover {
    color: var(--text-color);
    background: rgba(255, 255, 255, 0.03);
  }

  .nav-item.active {
    color: var(--primary-color);
    background: rgba(var(--primary-rgb), 0.08);
  }

  .nav-item.active .material-symbols-rounded {
    font-variation-settings: 'FILL' 1;
  }

  /* SUB NAV */
  .sub-nav-links {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    padding-left: 2rem;
    margin: 0.25rem 0;
    border-left: 1px solid var(--border-color);
    margin-left: 2rem;
  }

  .sub-nav-item {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 0.65rem 1rem;
    color: var(--text-muted);
    text-decoration: none;
    border-radius: 0.75rem;
    transition: all 0.2s;
    font-size: 0.85rem;
    font-weight: 600;
  }

  .sub-nav-item:hover {
    color: var(--text-color);
    background: rgba(255, 255, 255, 0.03);
  }

  .sub-nav-item.active {
    color: var(--primary-color);
    background: rgba(var(--primary-rgb), 0.08);
  }

  .sub-nav-item .material-symbols-rounded {
    font-size: 1.1rem;
  }

  /* BRAND SECTION BOTTOM */
  .brand-section-bottom {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 1rem;
    margin-top: auto;
    opacity: 0.6;
    transition: opacity 0.2s;
  }
  .brand-section-bottom:hover { opacity: 1; }

  .logo-circle-flat {
    width: 32px;
    height: 32px;
    background: var(--text-muted);
    color: var(--bg-color);
    border-radius: 8px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 900;
    font-size: 1rem;
  }

  .lang-toggle-btn {
    background: var(--surface-variant);
    border: 1px solid var(--border-color);
    color: var(--text-color);
    padding: 0.4rem 0.6rem;
    border-radius: 0.75rem;
    cursor: pointer;
    display: flex;
    align-items: center;
    gap: 0.4rem;
    transition: all 0.2s;
  }
  .lang-toggle-btn:hover { background: var(--surface-variant-hover); border-color: var(--primary-color); }
  .lang-toggle-btn .material-symbols-rounded { font-size: 1.1rem; }
  .lang-text { font-size: 0.75rem; font-weight: 800; }

  .brand-info { display: flex; flex-direction: column; }
  .brand-name { font-weight: 800; font-size: 1.1rem; letter-spacing: -0.02em; line-height: 1; }
  .version { font-size: 0.65rem; font-weight: 700; color: var(--text-muted); }

  /* Mobile stays similar but at bottom */
  @media (max-width: 1024px) {
    .navigation-bar {
        width: 100vw;
        height: 70px;
        flex-direction: row;
        position: fixed;
        bottom: 0;
        left: 0;
        z-index: 1000;
        padding: 0 1rem;
        background: var(--bg-color);
        border-top: 1px solid var(--border-color);
        align-items: center;
        padding-bottom: env(safe-area-inset-bottom);
        height: calc(70px + env(safe-area-inset-bottom));
        box-shadow: 0 -4px 20px rgba(0,0,0,0.1);
    }
    .user-profile-summary, .brand-section-bottom, .item-label, span:not(.material-symbols-rounded) { display: none; }
    .nav-links { flex-direction: row; width: 100%; justify-content: space-around; }
    .nav-item { padding: 0.5rem; }
    .theme-trigger-section { padding: 0; margin: 0; }
    .appearance-btn { width: auto; padding: 0.5rem; }
  }

  /* Theme Button Styles */
  .theme-trigger-section {
    padding: 0.5rem 0;
    margin-bottom: 0.5rem;
  }

  .appearance-btn {
    display: flex;
    align-items: center;
    gap: 1.25rem;
    width: 100%;
    padding: 0.85rem 1.25rem;
    color: var(--text-muted);
    background: transparent;
    border: none;
    cursor: pointer;
    border-radius: 1rem;
    transition: all 0.2s;
    font-weight: 600;
    font-size: 0.95rem;
  }

  .appearance-btn:hover {
    color: var(--text-color);
    background: rgba(255, 255, 255, 0.03);
  }

  .appearance-btn span.material-symbols-rounded {
    color: var(--primary-color);
  }
</style>
