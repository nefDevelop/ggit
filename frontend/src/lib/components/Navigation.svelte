<script lang="ts">
  import { authState } from '$lib/auth.svelte';
  import { page } from '$app/state';

  const navItems = [
    { name: 'Home', icon: 'home', path: '/' },
    { name: 'Dashboard', icon: 'grid_view', path: '/dashboard' },
    { name: 'Explore', icon: 'search', path: '/explore/repos' },
    { name: 'Issues', icon: 'rule', path: '/issues' },
    { name: 'Settings', icon: 'settings', path: '/user/settings/profile' },
  ];

  function getCsrfToken() {
    if (typeof document === 'undefined') return '';
    const match = document.cookie.match(new RegExp('(^| )XSRF-TOKEN=([^;]+)'));
    return match ? match[2] : '';
  }

  async function handleLogout() {
    await fetch('/api/v2/user/logout', { 
      method: 'POST',
      headers: { 'X-Csrf-Token': getCsrfToken() }
    });
    window.location.href = '/user/login';
  }
</script>

<nav class="navigation-bar">
  <div class="brand-section">
    <div class="logo-circle">G</div>
    <span class="brand-name">Ggit</span>
  </div>

  <div class="nav-links">
    {#each navItems as item}
      <a 
        href={item.path} 
        class="nav-item" 
        class:active={page.url.pathname === item.path}
      >
        <span class="material-symbols-rounded">{item.icon}</span>
        <span class="item-label">{item.name}</span>
      </a>
    {/each}
  </div>

  {#if authState.isAuthenticated}
    <div class="user-section">
      <img src={authState.user?.avatar} alt="Avatar" class="mini-avatar" />
      <div class="user-info">
        <span class="username">{authState.user?.username}</span>
      </div>
      <button class="logout-btn" onclick={handleLogout} title="Logout">
        <span class="material-symbols-rounded">logout</span>
      </button>
    </div>
  {/if}
</nav>

<style>
  .navigation-bar {
    background-color: var(--md3-surface-container);
    display: flex;
    flex-direction: column;
    padding: 1.5rem 0.75rem;
    gap: 2rem;
    border: 1px solid var(--border-color);
    width: 280px;
    height: calc(100vh - 3rem);
    margin: 1.5rem 0 1.5rem 1.5rem;
    border-radius: var(--radius);
    box-sizing: border-box;
    z-index: 1000;
    box-shadow: var(--shadow-island);
  }

  .brand-section {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 0 1rem;
    margin-bottom: 1rem;
  }

  .logo-circle {
    width: 36px;
    height: 36px;
    background: linear-gradient(135deg, var(--primary-color), #7aa9f7);
    color: var(--on-primary);
    border-radius: 12px;
    display: flex;
    align-items: center;
    justify-content: center;
    font-weight: 800;
    box-shadow: 0 4px 10px rgba(0,0,0,0.2);
  }

  .brand-name {
    font-weight: 700;
    font-size: 1.4rem;
    letter-spacing: -0.03em;
  }

  .nav-links {
    display: flex;
    flex-direction: column;
    gap: 0.25rem;
    flex: 1;
  }

  .nav-item {
    display: flex;
    align-items: center;
    gap: 1.25rem;
    padding: 0.85rem 1rem;
    color: var(--text-muted);
    text-decoration: none;
    border-radius: 1rem;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    font-weight: 500;
  }

  .nav-item:hover {
    background-color: rgba(255, 255, 255, 0.06);
    color: var(--text-color);
    transform: translateX(4px);
  }

  .nav-item.active {
    background-color: var(--md3-tonal-primary);
    color: var(--primary-color);
    font-weight: 600;
  }

  .nav-item.active .material-symbols-rounded {
    font-variation-settings: 'FILL' 1;
  }

  .user-section {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.75rem;
    background: rgba(255, 255, 255, 0.04);
    border-radius: 1rem;
    border: 1px solid var(--border-color);
  }

  .mini-avatar {
    width: 36px;
    height: 36px;
    border-radius: 50%;
    border: 1px solid rgba(255,255,255,0.1);
  }

  .user-info {
    flex: 1;
    display: flex;
    flex-direction: column;
    overflow: hidden;
  }

  .username {
    font-size: 0.95rem;
    font-weight: 600;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }

  .logout-btn {
    background: transparent;
    border: none;
    color: var(--text-muted);
    cursor: pointer;
    padding: 0.5rem;
    border-radius: 0.5rem;
    display: flex;
    align-items: center;
    transition: all 0.2s;
  }

  .logout-btn:hover {
    color: #ef4444;
    background: rgba(239, 68, 68, 0.1);
  }

  /* Responsive Bottom Toolbar */
  @media (max-width: 768px) {
    .navigation-bar {
      width: 100vw;
      height: 80px;
      flex-direction: row;
      position: fixed;
      bottom: 0;
      left: 0;
      padding: 0 0.5rem;
      margin: 0;
      border-right: none;
      border-top: 1px solid var(--border-color);
      background: rgba(30, 30, 30, 0.95);
      backdrop-filter: blur(20px);
      justify-content: space-around;
      align-items: center;
      gap: 0;
      border-radius: 2rem 2rem 0 0;
      box-shadow: 0 -10px 30px rgba(0,0,0,0.3);
    }

    .brand-section, .user-section, .item-label {
      display: none;
    }

    .nav-links {
      flex-direction: row;
      justify-content: space-around;
      width: 100%;
      gap: 0;
    }

    .nav-item {
      flex-direction: column;
      padding: 0.75rem 0.5rem;
      gap: 0.25rem;
      flex: 1;
      border-radius: 0;
      transform: none !important;
    }

    .nav-item:hover {
      background: transparent;
    }

    .nav-item span {
      font-size: 28px;
    }

    .nav-item.active {
      background: transparent;
      color: var(--primary-color);
    }
    
    /* Active indicator bar for mobile */
    .nav-item.active::after {
      content: '';
      position: absolute;
      top: 0;
      width: 20px;
      height: 3px;
      background: var(--primary-color);
      border-radius: 0 0 4px 4px;
    }
  }
</style>
