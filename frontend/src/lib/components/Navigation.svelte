<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { authState } from '$lib/auth.svelte';
  import { page } from '$app/state';

  const navItems = $derived([
    { name: 'Dashboard', icon: 'grid_view', path: '/dashboard' },
    { name: 'Explore', icon: 'search', path: '/explore/repos' },
    { name: 'Organizations', icon: 'corporate_fare', path: '/explore/users' },
    ...(authState.user?.is_admin ? [
      { name: 'Admin Panel', icon: 'security', path: '/admin/users', isParent: true }
    ] : []),
    { name: 'Settings', icon: 'settings', path: '/user/settings/profile' },
  ]);

  const adminSubItems = [
    { name: 'Users', icon: 'person', path: '/admin/users' },
    { name: 'Repositories', icon: 'book', path: '/admin/repos' },
    { name: 'Organizations', icon: 'corporate_fare', path: '/admin/orgs' },
    { name: 'Authentication', icon: 'lock', path: '/admin/auth' },
    { name: 'Configuration', icon: 'settings', path: '/admin/config' },
    { name: 'Monitor', icon: 'monitoring', path: '/admin/monitor' },
    { name: 'Maintenance', icon: 'build', path: '/admin/maintenance' },
    { name: 'Notices', icon: 'notifications', path: '/admin/notices' }
  ];

  const isAdminPath = $derived(page.url.pathname.startsWith('/admin'));

  function getCsrfToken() {
    if (typeof document === 'undefined') return '';
    const match = document.cookie.match(new RegExp('(^| )XSRF-TOKEN=([^;]+)'));
    return match ? match[2] : '';
  }

  async function handleLogout() {
    await apiFetch('/api/v2/user/logout', { 
      method: 'POST',
      headers: { 'X-Csrf-Token': getCsrfToken() }
    });
    window.location.href = '/user/login';
  }
</script>

<nav class="navigation-bar">
  <!-- USER AT TOP -->
  {#if authState.isAuthenticated}
    <div class="user-section-top">
      <img src={authState.user?.avatar} alt="Avatar" class="mini-avatar" />
      <div class="user-info">
        <span class="username">{authState.user?.username}</span>
        <span class="user-role">Developer</span>
      </div>
      <button class="logout-btn-minimal" onclick={handleLogout} title="Logout">
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

  <!-- BRAND AT BOTTOM -->
  <div class="brand-section-bottom">
    <div class="logo-circle-flat">G</div>
    <div class="brand-info">
        <span class="brand-name">Ggit</span>
        <span class="version">v2026.4.0</span>
    </div>
  </div>
</nav>

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

  /* USER SECTION TOP */
  .user-section-top {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.75rem;
    margin-bottom: 2.5rem;
    border-radius: 1rem;
    transition: background 0.2s;
  }

  .mini-avatar {
    width: 40px;
    height: 40px;
    border-radius: 50%;
    border: 2px solid var(--border-color);
  }

  .user-info {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .username {
    font-size: 0.95rem;
    font-weight: 700;
    color: var(--text-color);
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .user-role {
    font-size: 0.7rem;
    color: var(--text-muted);
    font-weight: 500;
    text-transform: uppercase;
    letter-spacing: 0.05em;
  }

  .logout-btn-minimal {
    background: transparent;
    border: none;
    color: var(--text-muted);
    cursor: pointer;
    padding: 0.4rem;
    display: flex;
    border-radius: 50%;
    transition: all 0.2s;
  }
  .logout-btn-minimal:hover {
    color: #f87171;
    background: rgba(248, 113, 113, 0.1);
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

  .brand-info { display: flex; flex-direction: column; }
  .brand-name { font-weight: 800; font-size: 1.1rem; letter-spacing: -0.02em; line-height: 1; }
  .version { font-size: 0.65rem; font-weight: 700; color: var(--text-muted); }

  /* Mobile stays similar but at bottom */
  @media (max-width: 768px) {
    .navigation-bar {
        width: 100vw;
        height: 70px;
        flex-direction: row;
        position: fixed;
        bottom: 0;
        left: 0;
        padding: 0 1rem;
        background: var(--md3-surface-container);
        border-top: 1px solid var(--border-color);
        align-items: center;
    }
    .user-section-top, .brand-section-bottom, .item-label { display: none; }
    .nav-links { flex-direction: row; width: 100%; justify-content: space-around; }
    .nav-item { padding: 0.5rem; }
  }
</style>
