<script lang="ts">
  import BookText from "lucide-svelte/icons/book-text";
  import Scale from "lucide-svelte/icons/scale";
  import Database from "lucide-svelte/icons/database";
  import Rocket from "lucide-svelte/icons/rocket";
  import ChevronRight from "lucide-svelte/icons/chevron-right";

  let { repoContext = null, repoData = null, username = "", reponame = "", langEntries = [], visibleLanguages = [], showAllLanguages = false, formatSize, onToggleLanguages }: {
    repoContext: any;
    repoData: any;
    username: string;
    reponame: string;
    langEntries: any[];
    visibleLanguages: any[];
    showAllLanguages: boolean;
    formatSize: (n: number) => string;
    onToggleLanguages: () => void;
  } = $props();
</script>

<aside class="side-column">
  <div class="side-section">
    <h4 class="side-title">About</h4>
    <p class="about-text">{repoContext?.data?.description || 'No description provided.'}</p>
    <div class="side-items-list">
      {#if repoContext?.full?.readme || repoData?.readme}
        <a href="#readme" class="side-item-row">
          <BookText size={16} /><span>Readme</span>
        </a>
      {/if}
      {#if repoContext?.data?.license}
        <div class="side-item-row">
          <Scale size={16} /><span>{repoContext.data.license} License</span>
        </div>
      {/if}
      {#if repoContext?.data?.size !== undefined}
        <div class="side-item-row">
          <Database size={16} /><span>{formatSize(repoContext.data.size)}</span>
        </div>
      {/if}
    </div>
  </div>

  <div class="side-section">
    <h4 class="side-title">Releases</h4>
    {#if repoContext?.stats?.releases > 0}
      <a href={`/${username}/${reponame}/releases`} class="side-item-link">
        <div class="side-item-main">
          <Rocket size={18} class="icon-primary" />
          <span class="side-item-text"><b>{repoContext.stats.releases}</b> releases</span>
        </div>
        <ChevronRight size={18} />
      </a>
    {:else}
      <p class="about-text">No releases published.</p>
    {/if}
  </div>

  {#if langEntries.length > 0}
    <div class="side-section desktop-only">
      <h4 class="side-title">Languages</h4>
      <div class="lang-bar">
        {#each langEntries as lang}
          <div class="lang-progress" style="width: {lang.percent}%; background: {lang.color}" title="{lang.name}: {lang.percent}%"></div>
        {/each}
      </div>
      <div class="lang-list">
        {#each visibleLanguages as lang}
          <div class="lang-item">
            <span class="dot" style="background: {lang.color}"></span>
            <span class="name">{lang.name}</span>
            <span class="percent">{lang.percent}%</span>
          </div>
        {/each}
      </div>
      {#if langEntries.length > 10}
        <button class="lang-toggle" onclick={onToggleLanguages}>
          {showAllLanguages ? 'Show less' : `Show all (${langEntries.length})`}
        </button>
      {/if}
    </div>
  {/if}
</aside>

<style>
  .side-column { width: 280px; flex-shrink: 0; display: flex; flex-direction: column; gap: 1.5rem; }
  .side-section { padding: 0 1rem; }
  .side-title { font-size: 0.75rem; font-weight: 800; text-transform: uppercase; color: var(--text-muted); margin: 0 0 0.75rem; }
  .about-text { font-size: 0.9rem; color: var(--text-muted); margin: 0; line-height: 1.4; }
  .side-items-list { display: flex; flex-direction: column; gap: 0.75rem; margin-top: 1rem; }
  .side-item-row { display: flex; align-items: center; gap: 0.5rem; color: var(--text-muted); text-decoration: none; font-size: 0.9rem; }
  .side-item-link {
    display: flex; align-items: center; justify-content: space-between;
    text-decoration: none; color: inherit; padding: 0.5rem; border-radius: 0.5rem;
    transition: background 0.2s;
  }
  .side-item-link:hover { background: var(--md-sys-color-surface-container-high); }
  .side-item-main { display: flex; align-items: center; gap: 0.75rem; }
  .icon-primary { color: var(--primary-color); }
  .side-item-text { font-size: 0.9rem; }
  .lang-bar { display: flex; height: 8px; border-radius: 4px; overflow: hidden; background: var(--md-sys-color-surface-container); }
  .lang-progress { height: 100%; transition: width 0.3s; }
  .lang-list { display: flex; flex-direction: column; gap: 0.25rem; margin-top: 0.75rem; }
  .lang-item { display: flex; align-items: center; gap: 0.5rem; font-size: 0.85rem; }
  .dot { width: 10px; height: 10px; border-radius: 50%; flex-shrink: 0; }
  .name { flex: 1; font-weight: 600; }
  .percent { color: var(--text-muted); font-size: 0.8rem; }
  .lang-toggle { background: none; border: none; color: var(--primary-color); cursor: pointer; font-size: 0.8rem; font-weight: 700; padding: 0.5rem 0; }
  .desktop-only { display: block; }
  @media (max-width: 900px) { .side-column { display: none; } }
</style>
