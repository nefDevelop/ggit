<script lang="ts">
  import Palette from "lucide-svelte/icons/palette";

  let { showThemePanel = $bindable(false), pageLoadTime = 0, templateTime = 0 }: {
    showThemePanel: boolean;
    pageLoadTime: number;
    templateTime: number;
  } = $props();
</script>

<div class="brand-section-bottom">
  <div class="brand-info">
    <span class="brand-name">Ggit</span>
    <span class="version">v2026.4.0</span>
  </div>
  <button
    class="minimal-theme-btn"
    aria-expanded={showThemePanel} aria-haspopup="dialog"
    onclick={() => (showThemePanel = !showThemePanel)}
    title="Apariencia"
    aria-label="Abrir panel de apariencia"
  >
    <Palette size={20} />
  </button>
</div>

<div class="sidebar-perf-info">
  Page: {pageLoadTime}ms Template: {templateTime}ms
</div>

<style>
  .brand-section-bottom {
    display: flex; align-items: center; gap: 1rem;
    padding: 1rem; margin-top: auto; opacity: 0.6; transition: opacity 0.2s;
  }
  .brand-section-bottom:hover { opacity: 1; }
  .brand-info { display: flex; flex-direction: column; }
  .brand-name { font-weight: 800; font-size: 1.1rem; letter-spacing: -0.02em; line-height: 1; }
  .version { font-size: 0.65rem; font-weight: 700; color: var(--text-muted); }
  .minimal-theme-btn {
    background: transparent; border: none; color: var(--text-muted);
    cursor: pointer; padding: 0.5rem; display: flex;
    border-radius: 0.75rem; transition: all 0.2s; margin-left: auto;
  }
  .minimal-theme-btn:hover { color: var(--primary-color); background: rgba(var(--primary-rgb),0.1); }
  .sidebar-perf-info {
    position: absolute; bottom: 0.75rem; left: 2rem;
    font-size: 0.625rem; font-weight: 600; color: var(--text-muted);
    opacity: 0.4; font-family: monospace; white-space: nowrap;
    pointer-events: none;
  }
  @media (max-width: 1024px) {
    .sidebar-perf-info { display: none; }
    .brand-section-bottom { display: none; }
  }
</style>
