<script lang="ts">
  import BookText from "lucide-svelte/icons/book-text";
  import Scale from "lucide-svelte/icons/scale";
  import Database from "lucide-svelte/icons/database";
  import Rocket from "lucide-svelte/icons/rocket";
  import ChevronRight from "lucide-svelte/icons/chevron-right";

  let { repoContext = null, repoData = null, username = "", reponame = "", langEntries = [], visibleLanguages = [], showAllLanguages = false, formatSize, onToggleLanguages }: {
    repoContext: any;
    repoData: any;
    // Optional because it has a default; SvelteKit types route params as
    // string | undefined, so a required string prop rejected every caller.
    username?: string;
    reponame?: string;
    langEntries: any[];
    visibleLanguages: any[];
    showAllLanguages: boolean;
    formatSize: (n: number) => string;
    onToggleLanguages: () => void;
  } = $props();
</script>

<aside class="side-column" class:hidden-sidebar={false}>
  <div class="side-section">
    <h4 class="side-title">About</h4>
    <p class="about-text">{repoContext?.data?.description || 'No description provided.'}</p>
    <div class="side-items-list">
      {#if repoContext?.full?.readme || repoData?.readme}
        <a href="#readme" class="side-item-row"><BookText size={16} /><span>Readme</span></a>
      {/if}
      {#if repoContext?.data?.license}
        <div class="side-item-row"><Scale size={16} /><span>{repoContext.data.license} License</span></div>
      {/if}
      {#if repoContext?.data?.size !== undefined}
        <div class="side-item-row"><Database size={16} /><span>{formatSize(repoContext.data.size)}</span></div>
      {/if}
    </div>
  </div>
  <div class="side-section">
    <h4 class="side-title">Releases</h4>
    {#if repoContext?.stats?.releases > 0}
      <a href={`/${username}/${reponame}/releases`} class="side-item-link">
        <div class="side-item-main"><Rocket size={18} class="icon-primary" /><span class="side-item-text"><b>{repoContext.stats.releases}</b> releases</span></div>
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
        <button class="lang-toggle" onclick={onToggleLanguages}>{showAllLanguages ? 'Show less' : `+ ${langEntries.length - 10} more`}</button>
      {/if}
    </div>
  {/if}
</aside>

<style>
  .side-column { display: flex; flex-direction: column; gap: 2.5rem; }
  .side-title { font-size: 0.75rem; font-weight: 800; text-transform: uppercase; letter-spacing: 0.1em; color: var(--text-muted); margin-bottom: 1rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); padding-bottom: 0.5rem; }
  .about-text { font-size: 0.95rem; line-height: 1.6; color: var(--text-color); opacity: 0.85; }
  .side-item-link { display: flex; align-items: center; justify-content: space-between; text-decoration: none; color: var(--text-color); padding: 0.75rem 1rem; background: rgba(var(--primary-rgb),0.05); border-radius: 12px; border: 1px solid var(--md-sys-color-outline-variant); transition: all 0.2s; }
  .side-item-link:hover { background: rgba(var(--primary-rgb),0.1); border-color: var(--primary-color); transform: translateY(-1px); }
  .side-item-main { display: flex; align-items: center; gap: 0.75rem; }
  :global(.icon-primary) { color: var(--primary-color); }
  .side-item-text { font-size: 0.9rem; }
  .side-item-row { display: flex; align-items: center; gap: 0.5rem; color: var(--text-muted); text-decoration: none; font-size: 0.9rem; }
  .side-items-list { display: flex; flex-direction: column; gap: 0.75rem; margin-top: 1rem; }
  .lang-bar { height: 12px; border-radius: 6px; overflow: hidden; display: flex; background: var(--sunken-2); margin-bottom: 1.25rem; box-shadow: inset 0 1px 3px rgba(0,0,0,0.3); }
  .lang-progress { height: 100%; transition: width 1s ease; }
  .lang-list { display: flex; flex-direction: column; gap: 0.6rem; }
  .lang-item { display: flex; align-items: center; gap: 0.75rem; font-size: 0.85rem; }
  .lang-item .dot { width: 10px; height: 10px; border-radius: 50%; }
  .lang-item .name { font-weight: 700; }
  .lang-item .percent { color: var(--text-muted); margin-left: auto; font-family: 'JetBrains Mono', monospace; font-size: 0.75rem; }
  .lang-toggle { background: transparent; border: none; color: var(--primary-color); font-size: 0.8rem; font-weight: 800; cursor: pointer; padding: 0.75rem 0; display: flex; align-items: center; gap: 0.5rem; }
  .lang-toggle:hover { text-decoration: underline; }
  .desktop-only { display: block; }
  @media (max-width: 900px) { .desktop-only { display: none; } }
</style>
