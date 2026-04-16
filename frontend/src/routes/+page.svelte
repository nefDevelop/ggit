<script>
  let serverInfo = $state({ version: '...', status: 'Checking...' });

  async function checkBackend() {
    try {
      const res = await fetch('/api/v1/repos/search');
      if (res.ok) {
        serverInfo.status = 'Online';
      } else {
        serverInfo.status = 'Authentication Required';
      }
    } catch (e) {
      serverInfo.status = 'Offline';
    }
  }

  checkBackend();
</script>

<div class="container">
  <header style="margin-bottom: 4rem; text-align: center;">
    <h1 style="font-size: 3.5rem; margin-bottom: 1rem; background: linear-gradient(to right, #6366f1, #3b82f6); -webkit-background-clip: text; -webkit-text-fill-color: transparent;">
      Gogs Modernized
    </h1>
    <p style="color: var(--text-muted); font-size: 1.25rem;">
      The world's easiest self-hosted Git service, now faster and more beautiful.
    </p>
  </header>

  <div class="grid">
    <div class="glass">
      <h3>Quick Start</h3>
      <p>Configure your server and start hosting your projects in minutes.</p>
      <a href="/install" class="btn">Go to Install</a>
    </div>

    <div class="glass">
      <h3>Explore</h3>
      <p>Discover interesting repositories and active developers.</p>
      <a href="/explore/repos" class="btn" style="background: transparent; border: 1px solid var(--primary-color);">Explore Repos</a>
    </div>

    <div class="glass">
      <h3>Server Status</h3>
      <div style="margin-top: 1rem;">
        <span style="display: inline-block; width: 12px; height: 12px; border-radius: 50%; margin-right: 8px; background-color: {serverInfo.status === 'Online' ? '#10b981' : (serverInfo.status === 'Offline' ? '#ef4444' : '#f59e0b')}; pulse animate"></span>
        <span>{serverInfo.status}</span>
      </div>
      <p style="font-size: 0.875rem; color: var(--text-muted); margin-top: 1rem;">
        Backend: <strong>Gin + GORM</strong>
      </p>
    </div>
  </div>

  <footer style="margin-top: 6rem; text-align: center; color: var(--text-muted);">
    <p>&copy; {new Date().getFullYear()} Gogs Modernization Project</p>
  </footer>
</div>

<style>
  @keyframes pulse {
    0% { opacity: 1; }
    50% { opacity: 0.5; }
    100% { opacity: 1; }
  }
  .animate {
    animation: pulse 2s infinite ease-in-out;
  }
</style>
