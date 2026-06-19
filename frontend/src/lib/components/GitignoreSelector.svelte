<script lang="ts">
  const OS_NAMES = ['macos', 'windows', 'linux'];
  const EDITOR_NAMES = ['visualstudiocode', 'jetbrains', 'vim', 'emacs', 'sublimetext', 'eclipse', 'netbeans', 'notepadpp', 'kate', 'cursor', 'zed', 'textmate', 'atom', 'sublime', 'code', 'clion', 'pycharm', 'webstorm', 'intellij'];

  let { gitignores = [], global_ignores = [], selectedGitignores = $bindable([]), gitignoreSearch = $bindable("") }: {
    gitignores: string[];
    global_ignores: string[];
    selectedGitignores: string[];
    gitignoreSearch: string;
  } = $props();

  let groupedGlobals = $derived.by(() => {
    const search = gitignoreSearch.toLowerCase();
    const filtered = global_ignores.filter(g => g.toLowerCase().includes(search));
    const security = filtered.filter(g => g.toLowerCase().includes('ggit'));
    const oss = filtered.filter(g => {
      const name = g.toLowerCase().replace('global/', '').replace('.gitignore', '');
      return OS_NAMES.some(os => name.includes(os)) && !g.toLowerCase().includes('ggit');
    });
    const editors = filtered.filter(g => {
      const name = g.toLowerCase().replace('global/', '').replace('.gitignore', '');
      return EDITOR_NAMES.some(ed => name.includes(ed)) && !g.toLowerCase().includes('ggit');
    });
    const others = filtered.filter(g => {
      const name = g.toLowerCase().replace('global/', '').replace('.gitignore', '');
      return !g.toLowerCase().includes('ggit') && !OS_NAMES.some(os => name.includes(os)) && !EDITOR_NAMES.some(ed => name.includes(ed));
    });
    return [
      { name: 'Seguridad y Estándar', icon: 'security', items: security, type: 'recommended', iconType: 'verified' },
      { name: 'Sistemas Operativos', icon: 'desktop_windows', items: oss, iconType: 'public' },
      { name: 'Editores e IDEs', icon: 'edit_square', items: editors, iconType: 'edit' },
      { name: 'Otros Globales', icon: 'extension', items: others, iconType: 'public' },
    ].filter(g => g.items.length > 0);
  });

  function toggleGitignore(g: string) {
    if (selectedGitignores.includes(g)) {
      selectedGitignores = selectedGitignores.filter(i => i !== g);
    } else {
      selectedGitignores = [...selectedGitignores, g];
    }
  }
</script>

<div class="init-item full">
  <div class="init-header-row">
    <div class="section-label">
      <span class="material-symbols-rounded">public</span>
      <strong>Templates Globales</strong>
    </div>
    <div class="search-box-mini">
      <span class="material-symbols-rounded">search</span>
      <input type="text" bind:value={gitignoreSearch} placeholder="Buscar templates..." />
    </div>
  </div>

  {#each groupedGlobals as group}
    <div class="gitignore-group">
      <div class="group-label">
        <span class="material-symbols-rounded">{group.icon}</span>
        {group.name}
      </div>
      <div class="gitignore-chips">
        {#each group.items as g}
          <button type="button" class="chip" title={g}
            class:active={selectedGitignores.includes(g)}
            class:recommended={group.type === 'recommended'}
            onclick={() => toggleGitignore(g)}>
            <span class="material-symbols-rounded">{group.iconType || 'public'}</span>
            {g.replace('global/', '').replace('.gitignore', '').replace(/_/g, ' ').trim()}
            {#if selectedGitignores.includes(g)}
              <span class="material-symbols-rounded">close</span>
            {/if}
          </button>
        {/each}
      </div>
    </div>
  {/each}

  <div class="gitignore-group" style="margin-top: 1.5rem;">
    <div class="group-label">
      <span class="material-symbols-rounded">terminal</span>
      <strong>Por Lenguaje</strong>
    </div>
    <div class="gitignore-chips">
      {#each gitignores.filter(g => g.toLowerCase().includes(gitignoreSearch.toLowerCase())) as g}
        <button type="button" class="chip" title={g}
          class:active={selectedGitignores.includes(g)}
          onclick={() => toggleGitignore(g)}>
          <span class="material-symbols-rounded">code</span>
          {g.replace('languages/', '').replace('.gitignore', '')}
          {#if selectedGitignores.includes(g)}
            <span class="material-symbols-rounded">close</span>
          {/if}
        </button>
      {/each}
    </div>
  </div>

  {#if selectedGitignores.length > 0}
    <div class="selection-summary">
      <span class="material-symbols-rounded">check</span>
      <strong>Mix seleccionado:</strong> {selectedGitignores.map(g => (g.split('/').pop() || '').replace('.gitignore', '')).join(' + ')}
    </div>
  {/if}
</div>

<style>
  .init-item.full { grid-column: 1 / -1; }
  .init-header-row { display: flex; justify-content: space-between; align-items: center; flex-wrap: wrap; gap: 0.5rem; margin-bottom: 1rem; }
  .section-label { display: flex; align-items: center; gap: 0.5rem; }
  .search-box-mini { display: flex; align-items: center; gap: 0.5rem; background: rgba(0,0,0,0.1); padding: 0.3rem 0.8rem; border-radius: 2rem; }
  .search-box-mini input { background: transparent; border: none; color: inherit; outline: none; width: 180px; font-family: inherit; }
  .gitignore-group { margin-bottom: 1rem; }
  .group-label { display: flex; align-items: center; gap: 0.5rem; font-size: 0.85rem; font-weight: 700; color: var(--text-muted); text-transform: uppercase; margin-bottom: 0.5rem; }
  .gitignore-chips { display: flex; flex-wrap: wrap; gap: 0.5rem; }
  .chip {
    display: flex; align-items: center; gap: 0.4rem;
    padding: 0.35rem 0.75rem; border-radius: 2rem;
    border: 1px solid var(--border-color); background: transparent;
    color: var(--text-muted); font-size: 0.8rem; cursor: pointer;
    transition: all 0.2s; font-family: inherit;
  }
  .chip:hover { border-color: var(--primary-color); color: var(--text-color); }
  .chip.active {
    background: var(--primary-color); color: var(--on-primary);
    border-color: var(--primary-color);
  }
  .chip.recommended {
    border-color: var(--primary-color);
    background: rgba(var(--primary-rgb),0.08);
  }
  .selection-summary {
    padding: 0.75rem; background: rgba(var(--primary-rgb),0.05);
    border-radius: 0.75rem; display: flex; align-items: center; gap: 0.5rem;
    font-size: 0.85rem;
  }
</style>
