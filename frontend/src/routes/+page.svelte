<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { authState } from '$lib/auth.svelte';
  import DashboardView from '$lib/components/DashboardView.svelte';

  // State to check backend while logged out
  let serverInfo = $state({ version: '...', status: 'Checking...' });

  async function checkBackend() {
    // Only check if not authenticated otherwise Dashboard will do its own fetching
    if (authState.isAuthenticated) return;
    
    try {
      const res = await apiFetch('/api/v2/explore/search');
      if (res.ok) {
        serverInfo.status = 'Online';
      } else {
        serverInfo.status = 'Authentication Required';
      }
    } catch (e) {
      serverInfo.status = 'Offline';
    }
  }

  $effect(() => {
    if (!authState.loading && !authState.isAuthenticated) {
      checkBackend();
    }
  });
</script>

{#if authState.loading}
  <div class="loading-screen">
    <div class="spinner"></div>
  </div>
{:else if authState.isAuthenticated}
  <DashboardView />
{:else}
  <div class="container landing-container">
    <header style="margin-bottom: 4rem; text-align: center;">
      <h1 class="hero-title text-gradient">Ggit</h1>
      <p style="color: var(--text-muted); font-size: 1.25rem;">
        The world's easiest self-hosted Git service, now faster and more beautiful.
      </p>
    </header>

    <div class="grid">
      <div class="glass">
        <h3>Welcome Back</h3>
        <p>Sign in to your account and manage your repositories.</p>
        <a href="/user/login" class="btn">Sign In</a>
      </div>

      <div class="glass">
        <h3>Explore</h3>
        <p>Discover interesting repositories and active developers.</p>
        <a href="/explore/repos" class="btn btn-tonal">Explore Repos</a>
      </div>

      <div class="glass">
        <h3>Server Status</h3>
        <div style="margin-top: 1rem; display: flex; align-items: center;">
          <span class="status-indicator {serverInfo.status === 'Online' ? 'online' : (serverInfo.status === 'Offline' ? 'offline' : 'auth')} animate"></span>
          <span>{serverInfo.status}</span>
        </div>
        <p style="font-size: 0.875rem; color: var(--text-muted); margin-top: 1rem;">
          Backend: <strong>Gin + GORM</strong>
        </p>
      </div>
    </div>

    <footer style="margin-top: 6rem; text-align: center; color: var(--text-muted);">
      <p>&copy; {new Date().getFullYear()} Ggit Modernization Project</p>
    </footer>
  </div>
{/if}

<style>
  .landing-container {
    padding-top: 4rem;
  }

  .hero-title {
    font-size: 4rem;
    margin-bottom: 1rem;
    background: linear-gradient(135deg, var(--on-primary), var(--primary-color));
    -webkit-background-clip: text;
    background-clip: text;
    -webkit-text-fill-color: transparent;
    font-weight: 800;
    letter-spacing: -0.03em;
  }

  .status-indicator {
    display: inline-block;
    width: 12px;
    height: 12px;
    border-radius: 50%;
    margin-right: 8px;
  }

  .status-indicator.online { background-color: var(--status-success); }
  .status-indicator.offline { background-color: var(--status-danger); }
  .status-indicator.auth { background-color: var(--color-warning); }

  @keyframes pulse {
    0% { opacity: 1; }
    50% { opacity: 0.4; }
    100% { opacity: 1; }
  }
  .animate {
    animation: pulse 2s infinite ease-in-out;
  }

  .loading-screen {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 100vh;
  }

  .spinner {
    width: 40px;
    height: 40px;
    border: 3px solid rgba(255,255,255,0.1);
    border-radius: 50%;
    border-top-color: var(--primary-color);
    animation: spin 1s ease-in-out infinite;
  }

  @keyframes spin {
    to { transform: rotate(360deg); }
  }
</style>

