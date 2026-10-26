<script lang="ts">
  import { onMount } from 'svelte';
  import { apiFetch } from "$lib/api";
  import { authState } from '$lib/auth.svelte';
  import { goto } from '$app/navigation';
  import { page } from '$app/state';
  import GitignoreSelector from '$lib/components/GitignoreSelector.svelte';

  let name = $state('');
  let description = $state('');
  let website = $state('');
  let projectFolder = $state('');
  let defaultBranch = $state('');
  let isPrivate = $state(false);
  let isUnlisted = $state(false);
  let autoInit = $state(true);
  let creationMode = $state<'new' | 'migrate'>(
    (page.url.searchParams.get('mode') as 'new' | 'migrate') || 'new'
  );

  let isArchived = $state(false);

  // Owners & Organizations
  let ownerID = $state(0);
  let organizations = $state<any[]>([]);

  // Migration fields
  let cloneAddr = $state('');
  let authUsername = $state('');
  let authPassword = $state('');

  let gitignores = $state<string[]>([]);
  let global_ignores = $state<string[]>([]);
  let licenses = $state<string[]>([]);
  let readmes = $state<string[]>([]);
  
  let selectedGitignores = $state<string[]>([]);
  let selectedLicense = $state('');
  
  let includeReadme = $state(true);
  let selectedReadme = $state('Default');
  
  let errorMsg = $state('');
  let loading = $state(false);
  let gitignoreSearch = $state('');

  async function loadTemplates() {
    try {
      const res = await apiFetch('/api/v2/repo/templates');
      if (res.ok) {
        const data = await res.json();
        gitignores = data.gitignores || [];
        global_ignores = data.global_ignores || [];
        licenses = data.licenses || [];
        readmes = data.readmes || [];
      }
    } catch (e) {
      console.error("Failed to load templates", e);
    }
  }

  async function loadOrganizations() {
    try {
      const res = await apiFetch('/api/v2/user/settings/orgs');
      if (res.ok) {
        organizations = await res.json();
      }
    } catch (e) {
      console.error("Failed to load orgs", e);
    }
  }

  onMount(() => {
    ownerID = authState.user?.id || 0;
    loadTemplates();
    loadOrganizations();
  });

  $effect(() => {
    // Auto-fill repo name from URL if empty and in migrate mode
    if (creationMode === 'migrate' && cloneAddr && !name) {
      const match = cloneAddr.match(/\/([^\/.]+)(\.git)?$/);
      if (match) {
        name = match[1];
      }
    }
  });


  async function handleSubmit() {
    if (!name) return;
    loading = true;
    errorMsg = '';
    
    try {
      let res;
      if (creationMode === 'new') {
        res = await apiFetch('/api/v2/user/repos', {
          method: 'POST',
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            uid: ownerID,
            name,
            description,
            website,
            project_folder: projectFolder,
            default_branch: defaultBranch,
            private: isPrivate,
            unlisted: isUnlisted,
            auto_init: autoInit,
            is_archived: isArchived,
            gitignores: selectedGitignores.join(','),
            license: selectedLicense,
            readme: includeReadme ? selectedReadme : ''
          })
        });
      } else {
        res = await apiFetch('/api/v2/user/repos/migrate', {
          method: 'POST',
          silent: true,
          headers: { 'Content-Type': 'application/json' },
          body: JSON.stringify({
            clone_addr: cloneAddr,
            auth_username: authUsername,
            auth_password: authPassword,
            uid: ownerID,
            repo_name: name,
            description,
            website,
            project_folder: projectFolder,
            default_branch: defaultBranch,
            private: isPrivate,
            unlisted: isUnlisted,
            is_archived: isArchived
          })
        });
      }
      
      const data = res.data || await res.json();
      if (res.ok) {
        goto(`/${authState.user.username}/${name}`);
      } else {
        if (res.status === 401) {
            errorMsg = 'El repositorio remoto requiere autenticación. Por favor, verifica el nombre de usuario y la contraseña.';
        } else {
            errorMsg = data.message || 'Failed to create/migrate repository';
        }
      }
    } catch (e) {
      errorMsg = 'An error occurred. Please try again.';
    } finally {
      loading = false;
    }
  }
</script>

<div class="create-wrapper">
  <!-- HERO HEADER -->
  <header class="hero-header">
    <div class="hero-content">
      <div>
        <h1>{creationMode === 'new' ? 'Crear nuevo repositorio' : 'Importar repositorio'}</h1>
        <p>
          {creationMode === 'new' 
            ? 'Un repositorio contiene todos los archivos de tu proyecto, incluyendo el historial de revisiones.' 
            : 'Importa un repositorio existente desde otra URL para continuar tu trabajo aquí.'}
        </p>
      </div>
    </div>
  </header>

  <main class="create-grid">
    <div class="tabs-selector">
      <button 
        type="button" 
        class="tab-item" 
        class:active={creationMode === 'new'} 
        onclick={() => creationMode = 'new'}
      >
        <span class="material-symbols-rounded">add_circle</span>
        Nuevo Repositorio
      </button>
      <button 
        type="button" 
        class="tab-item" 
        class:active={creationMode === 'migrate'} 
        onclick={() => creationMode = 'migrate'}
      >
        <span class="material-symbols-rounded">cloud_download</span>
        Importar / Clonar
      </button>
    </div>

    <!-- RIGHT COLUMN: FORM -->
    <form class="create-form-content" onsubmit={(e) => { e.preventDefault(); handleSubmit(); }}>
      {#if errorMsg}
        <div class="error-toast">
          <span class="material-symbols-rounded">error</span>
          {errorMsg}
        </div>
      {/if}

      <section class="form-section island tonal">
        <div class="section-header">
          <span class="material-symbols-rounded">id_card</span>
          <h2>{creationMode === 'new' ? 'Información Básica' : 'Origen y Destino'}</h2>
        </div>
        
        {#if creationMode === 'migrate'}
          <div class="input-field">
            <label for="clone-addr">URL de Clonación (Git)</label>
            <div class="input-wrapper">
              <span class="prefix"><span class="material-symbols-rounded">link</span></span>
              <input 
                type="url" 
                id="clone-addr" 
                bind:value={cloneAddr} 
                placeholder="https://github.com/usuario/repo.git"
                required
              />
            </div>
          </div>

          <div class="auth-grid">
            <div class="input-field">
              <label for="auth-username">Usuario <span class="optional">(opcional)</span></label>
              <input type="text" id="auth-username" bind:value={authUsername} placeholder="Opcional" class="standard-input" />
            </div>
            <div class="input-field">
              <label for="auth-password">Contraseña / Token <span class="optional">(opcional)</span></label>
              <input type="password" id="auth-password" bind:value={authPassword} placeholder="Opcional" class="standard-input" />
            </div>
          </div>
          
          <div class="divider"></div>
        {/if}

        <div class="input-field">
          <label for="repo-owner">Propietario</label>
          <div class="input-wrapper select-only">
            <span class="prefix"><span class="material-symbols-rounded">person</span></span>
            <select id="repo-owner" bind:value={ownerID} disabled={loading}>
              <option value={authState.user?.id}>{authState.user?.username} (Mío)</option>
              {#each organizations as org}
                <option value={org.id}>{org.username} (Org)</option>
              {/each}
            </select>
          </div>
        </div>

        <div class="input-field">
          <label for="repo-name">Nombre del Repositorio</label>
          <div class="input-wrapper">
            <span class="prefix">{organizations.find(o => o.id === ownerID)?.username || authState.user?.username} /</span>
            <input 
              type="text" 
              id="repo-name" 
              bind:value={name} 
              placeholder="mi-proyecto-increible"
              pattern="[a-zA-Z0-9_\-\.]+"
              title="Solo letras, números, guiones y puntos"
              disabled={loading}
              required
            />
          </div>
        </div>

        <div class="input-field">
          <label for="repo-desc">Descripción <span class="optional">(opcional)</span></label>
          <textarea 
            id="repo-desc" 
            bind:value={description} 
            placeholder="¿De qué trata este repositorio?"
          disabled={loading}
          ></textarea>
        </div>

        <div class="input-grid">
          <div class="input-field">
            <label for="repo-site">Sitio Web <span class="optional">(opcional)</span></label>
            <div class="input-wrapper">
              <span class="prefix"><span class="material-symbols-rounded">language</span></span>
              <input type="url" id="repo-site" bind:value={website} placeholder="https://..." disabled={loading} />
            </div>
          </div>
          <div class="input-field">
            <label for="repo-branch">Rama por defecto <span class="optional">(opcional)</span></label>
            <div class="input-wrapper">
              <span class="prefix"><span class="material-symbols-rounded">fork_right</span></span>
              <input type="text" id="repo-branch" bind:value={defaultBranch} placeholder="main" disabled={loading} />
            </div>
          </div>
        </div>

        <div class="input-field">
          <label for="repo-folder">Carpeta / Grupo <span class="optional">(opcional)</span></label>
          <div class="input-wrapper">
            <span class="prefix"><span class="material-symbols-rounded">folder</span></span>
            <input type="text" id="repo-folder" bind:value={projectFolder} placeholder="ej. frontend, herramientas, docs" disabled={loading} />
          </div>
          <p class="helper-text">Agrupa este repositorio en una carpeta dentro del dashboard.</p>
        </div>
      </section>

      <section class="form-section island tonal">
        <div class="section-header">
          <span class="material-symbols-rounded">visibility</span>
          <h2>Visibilidad</h2>
        </div>
        
        <div class="visibility-selector">
          <label class="vm-option" class:selected={!isPrivate} class:disabled={loading}>
            <input type="radio" bind:group={isPrivate} value={false} disabled={loading} />
            <span class="material-symbols-rounded vm-icon">public</span>
            <div class="vm-info">
              <strong>Público</strong>
              <span>Visible para todos en el servidor.</span>
            </div>
            {#if !isPrivate}<span class="material-symbols-rounded dot">check_circle</span>{/if}
          </label>

          <label class="vm-option" class:selected={isPrivate} class:disabled={loading}>
            <input type="radio" bind:group={isPrivate} value={true} disabled={loading} />
            <span class="material-symbols-rounded vm-icon">lock</span>
            <div class="vm-info">
              <strong>Privado</strong>
              <span>Solo tú y tus colaboradores.</span>
            </div>
            {#if isPrivate}<span class="material-symbols-rounded dot">check_circle</span>{/if}
          </label>
        </div>

        <div class="unlisted-wrapper">
          <label class="fancy-checkbox" class:disabled={loading}>
            <input type="checkbox" bind:checked={isUnlisted} disabled={loading} />
            <div class="check-box">
              <span class="material-symbols-rounded check-icon">check</span>
            </div>
            <div class="check-text">
              <strong>Repositorio no listado (Unlisted)</strong>
              <span>Ocultar de la página de exploración pública.</span>
            </div>
          </label>
        </div>

        <div class="unlisted-wrapper">
          <label class="fancy-checkbox" class:disabled={loading}>
            <input type="checkbox" bind:checked={isArchived} disabled={loading} />
            <div class="check-box">
              <span class="material-symbols-rounded check-icon">check</span>
            </div>
            <div class="check-text">
              <strong>Archivar repositorio</strong>
              <span>Marcar como archivado desde el primer momento (solo lectura).</span>
            </div>
          </label>
        </div>
      </section>


      {#if creationMode === 'new'}
      <section class="form-section island tonal">
        <div class="section-header">
          <span class="material-symbols-rounded">settings</span>
          <h2>Inicialización</h2>
        </div>
        
        <label class="fancy-checkbox" class:disabled={loading}>
          <input type="checkbox" bind:checked={autoInit} disabled={loading} />
          <div class="check-box">
            <span class="material-symbols-rounded check-icon">check</span>
          </div>
          <div class="check-text">
            <strong>Inicializar repositorio</strong>
            <span>Crea el primer commit automáticamente con los archivos seleccionados.</span>
          </div>
        </label>

        {#if autoInit}
          <div class="init-details-v2">
            <!-- README OPTION -->
            <div class="init-item">
              <label class="item-toggle">
                <input type="checkbox" bind:checked={includeReadme} />
                <span class="material-symbols-rounded">description</span>
                <div>
                  <strong>Incluir README.md</strong>
                  {#if includeReadme}
                  <select bind:value={selectedReadme}>
                    {#each readmes as r}<option value={r}>{r}</option>{/each}
                  </select>
                  {/if}
                </div>
              </label>
            </div>

            <!-- LICENSE OPTION -->
            <div class="init-item">
              <label class="item-toggle">
                <input 
                  type="checkbox" 
                  checked={selectedLicense !== ''} 
                  onchange={(e) => {
                    const target = e.target as HTMLInputElement;
                    selectedLicense = target.checked ? (licenses[0] || '') : '';
                  }} 
                />
                <span class="material-symbols-rounded">balance</span>
                <div>
                  <strong>Incluir Licencia</strong>
                  {#if selectedLicense !== ''}
                  <select bind:value={selectedLicense}>
                    <option value="">Ninguna</option>
                    {#each licenses as l}<option value={l}>{l}</option>{/each}
                  </select>
                  {/if}
                </div>
              </label>
            </div>

            <!-- MULTI-GITIGNORE OPTION -->
            <GitignoreSelector {gitignores} {global_ignores} bind:selectedGitignores bind:gitignoreSearch />
          </div>
        {/if}
      </section>
      {/if}

      <div class="form-actions">
        <button type="button" class="btn-text" onclick={() => goto('/dashboard')}>Cancelar</button>
        <button type="submit" class="btn-action-primary" disabled={loading || !name}>
          {#if loading}
            <span class="material-symbols-rounded spin">sync</span>
            Despegando...
          {:else}
            <span class="material-symbols-rounded">{creationMode === 'new' ? 'rocket_launch' : 'cloud_sync'}</span>
            {creationMode === 'new' ? 'Crear Repositorio' : 'Importar Repositorio'}
          {/if}
        </button>
      </div>
    </form>
  </main>
</div>

<style>
  .create-wrapper {
    width: 100%;
    animation: fadeIn 0.4s ease-out;
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(10px); }
    to { opacity: 1; transform: translateY(0); }
  }

  .hero-header {
    background: var(--surface-variant);
    border-bottom: 1px solid var(--border-color);
    padding: 3.5rem 2rem;
    margin-bottom: 3rem;
  }

  .hero-content {
    display: flex;
    align-items: center;
    gap: 2rem;
    padding: 0 2rem;
  }

  .hero-header h1 {
    font-size: 2.5rem;
    font-weight: 900;
    margin: 0;
    background: linear-gradient(to right, var(--text-primary), var(--primary-color));
    background-clip: text;
    -webkit-background-clip: text;
    background-clip: text;
    -webkit-text-fill-color: transparent;
  }

  .hero-header p {
    font-size: 1.1rem;
    color: var(--text-muted);
    margin: 0.5rem 0 0 0;
  }

  /* GRID */
  .create-grid {
    max-width: 800px;
    margin: 0 auto;
    padding-bottom: 5rem;
    display: flex;
    flex-direction: column;
    gap: 2rem;
  }

  .tabs-selector {
    display: flex;
    gap: 0.5rem;
    padding: 0.25rem;
    background: var(--surface-variant);
    border-radius: 1.25rem;
    border: 1px solid var(--border-color);
    width: fit-content;
    margin: 0 auto;
  }
  
  .tab-item {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    padding: 0.75rem 1.5rem;
    border: none;
    background: transparent;
    color: var(--text-muted);
    font-weight: 700;
    cursor: pointer;
    border-radius: 1rem;
    transition: all 0.2s cubic-bezier(0.4, 0, 0.2, 1);
    font-size: 0.9rem;
  }
  
  .tab-item.active {
    background: var(--bg-color);
    color: var(--primary-color);
    box-shadow: var(--shadow-sm);
  }
  
  .tab-item span { font-size: 1.2rem; }

  .auth-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1.5rem;
    margin-top: 1.5rem;
  }

  .standard-input {
    background: var(--surface-variant);
    border: 1px solid var(--border-color);
    border-radius: 1rem;
    padding: 1rem;
    color: var(--text-color);
    font-weight: 600;
  }

  .divider {
    height: 1px;
    background: var(--border-color);
    margin: 2rem 0;
  }

  /* FORM SECTION */
  .form-section {
    margin-bottom: 2rem !important;
    padding: 2rem !important;
  }

  .section-header {
    display: flex;
    align-items: center;
    gap: 0.75rem;
    margin-bottom: 2rem;
    color: var(--primary-color);
  }

  .section-header h2 {
    font-size: 1.25rem;
    font-weight: 800;
    margin: 0;
    color: var(--text-primary);
  }

  .input-field {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
    margin-bottom: 1.5rem;
  }

  .input-wrapper.select-only {
    padding-right: 0.5rem;
  }
  
  .input-wrapper select {
    border: none;
    background-image: none;
    width: 100%;
    background: transparent;
    color: var(--text-color);
    font-weight: 600;
  }

  .unlisted-wrapper {
    margin-top: 1.5rem;
    padding-top: 1.5rem;
    border-top: 1px solid var(--border-color);
  }

  .input-field label { font-weight: 700; font-size: 0.95rem; color: var(--text-secondary); }

  .input-grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1.5rem;
    margin-top: 1.5rem;
  }

  .input-wrapper {
    display: flex;
    align-items: center;
    background: var(--surface-variant);
    border: 1px solid var(--border-color);
    border-radius: 1rem;
    overflow: hidden;
    transition: all 0.2s;
  }

  .input-wrapper:focus-within {
    border-color: var(--primary-color);
    box-shadow: 0 0 0 4px rgba(var(--primary-rgb), 0.1);
  }

  .prefix {
    background: var(--sunken-1);
    padding: 1rem;
    font-weight: 700;
    font-size: 0.9rem;
    color: var(--text-muted);
  }

  .input-wrapper input {
    flex: 1;
    background: transparent;
    border: none;
    padding: 1rem;
    color: var(--text-color);
    font-size: 1rem;
    font-weight: 600;
  }

  textarea {
    background: var(--surface-variant);
    border: 1px solid var(--border-color);
    border-radius: 1rem;
    padding: 1rem;
    min-height: 120px;
    color: var(--text-color);
    font-family: inherit;
    resize: vertical;
    transition: all 0.2s;
  }

  textarea:focus {
    border-color: var(--primary-color);
    box-shadow: 0 0 0 4px rgba(var(--primary-rgb), 0.1);
    outline: none;
  }
  
  textarea:disabled, input:disabled, select:disabled {
    opacity: 0.6;
    cursor: not-allowed;
  }

  select {
    background: var(--surface-variant);
    border: 1px solid var(--border-color);
    border-radius: 1rem;
    padding: 1rem;
    color: var(--text-color);
    font-family: inherit;
    font-size: 1rem;
    font-weight: 600;
    cursor: pointer;
    transition: all 0.2s;
    appearance: none;
    background-image: url("data:image/svg+xml,%3Csvg xmlns='http://www.w3.org/2000/svg' height='24' viewBox='0 -960 960 960' width='24'%3E%3Cpath d='M480-345 240-585l56-56 184 184 184-184 56 56-240 240Z'/%3E%3C/svg%3E");
    background-repeat: no-repeat;
    background-position: right 1rem center;
  }

  select:focus {
    border-color: var(--primary-color);
    box-shadow: 0 0 0 4px rgba(var(--primary-rgb), 0.1);
    outline: none;
  }

  .init-details-v2 {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    margin-top: 2rem;
    padding-top: 1.5rem;
    border-top: 1px solid var(--border-color);
  }

  .init-item {
    background: var(--surface-variant-hover);
    padding: 1.25rem;
    border-radius: 1.25rem;
    border: 1px solid var(--border-color);
  }

  .init-item.full { grid-column: span 2; }

  .item-toggle {
    display: flex;
    align-items: flex-start;
    gap: 1rem;
    cursor: pointer;
  }

  .item-toggle input { width: 20px; height: 20px; margin-top: 3px; }
  .item-toggle span { color: var(--primary-color); font-size: 1.5rem; }
  .item-toggle div { flex: 1; display: flex; flex-direction: column; gap: 0.5rem; }
  .item-toggle strong { font-size: 1rem; color: var(--text-primary); }

  .item-toggle select {
    padding: 0.5rem 1rem;
    font-size: 0.9rem;
    border-radius: 0.75rem;
    background: var(--surface-variant);
  }

  .init-header-row {
    display: flex;
    align-items: center;
    justify-content: space-between;
    margin-bottom: 1.5rem;
  }

  .init-header-row .section-label { margin-bottom: 0; }
  .init-header-row .search-box-mini { margin-left: 2rem; }

  .gitignore-group {
    margin-bottom: 1.5rem;
  }

  .group-label {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    font-size: 0.75rem;
    font-weight: 800;
    color: var(--text-muted);
    margin-bottom: 0.75rem;
    padding-left: 0.25rem;
    text-transform: uppercase;
    letter-spacing: 0.1em;
  }

  .group-label span {
    font-size: 1rem;
    color: var(--primary-color);
  }

  .gitignore-chips {
    display: flex;
    flex-wrap: wrap;
    gap: 0.5rem;
    max-height: 250px;
    overflow-y: auto;
    padding: 0.5rem;
    background: var(--sunken-1);
    border-radius: 1rem;
    min-height: 50px;
  }

  .search-box-mini {
    flex: 1;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    background: var(--surface-variant);
    border: 1px solid var(--border-color);
    padding: 0.25rem 0.75rem;
    border-radius: 2rem;
    margin-left: 1rem;
  }

  .search-box-mini span { font-size: 1.1rem; opacity: 0.5; }
  .search-box-mini input {
    background: transparent;
    border: none;
    font-size: 0.85rem;
    color: var(--text-color);
    width: 100%;
    outline: none;
  }

  .chip {
    padding: 0.5rem 1rem;
    border-radius: 2rem;
    background: var(--surface-variant);
    border: 1px solid var(--border-color);
    color: var(--text-muted);
    font-size: 0.85rem;
    font-weight: 600;
    cursor: pointer;
    display: flex;
    align-items: center;
    gap: 0.5rem;
    transition: all 0.2s;
  }

  .chip:hover { background: var(--border-color); color: var(--text-primary); }
  .chip.active {
    background: var(--primary-color);
    color: var(--on-primary);
    border-color: var(--primary-color);
  }
  .chip.recommended {
     background: rgba(var(--primary-rgb), 0.08);
     border-color: var(--primary-color);
     color: var(--primary-color);
     border-width: 2px;
  }
  .chip.recommended.active {
     background: var(--primary-color);
     color: var(--on-primary);
  }

  .selection-summary {
    margin-top: 1rem;
    font-size: 0.85rem;
    font-weight: 700;
    color: var(--primary-color);
    background: var(--primary-surface);
    padding: 0.5rem 1rem;
    border-radius: 1.5rem;
  }

  /* VISIBILITY SELECTOR */
  .visibility-selector {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1rem;
  }

  .vm-option {
    display: flex;
    align-items: center;
    gap: 1rem;
    padding: 1.5rem;
    background: var(--surface-variant);
    border: 1px solid var(--border-color);
    border-radius: 1.5rem;
    cursor: pointer;
    transition: all 0.2s;
    position: relative;
  }

  .vm-option:hover { background: var(--surface-variant-hover); transform: translateY(-2px); }
  .vm-option.selected { border-color: var(--primary-color); background: var(--primary-surface); }
  .vm-option input { display: none; }
  .vm-option.disabled { opacity: 0.6; cursor: not-allowed; pointer-events: none; }
  
  .vm-icon { font-size: 2rem; color: var(--primary-color); }
  .vm-info { display: flex; flex-direction: column; }
  .vm-info strong { font-size: 1.1rem; color: var(--text-primary); }
  .vm-info span { font-size: 0.85rem; color: var(--text-muted); }
  .dot { position: absolute; top: 1rem; right: 1rem; color: var(--primary-color); }

  @media (max-width: 600px) {
    .visibility-selector {
      grid-template-columns: 1fr;
    }
    .vm-option {
      padding: 1rem;
    }
    .hero-header {
      padding: 2rem 1rem;
    }
    .hero-header h1 {
      font-size: 1.8rem;
    }
    .form-section {
      padding: 1.25rem !important;
    }
  }

  /* FANCY CHECKBOX */
  .fancy-checkbox {
    display: flex;
    align-items: center;
    gap: 1.5rem;
    padding: 1rem;
    cursor: pointer;
  }

  .fancy-checkbox.disabled { opacity: 0.6; cursor: not-allowed; pointer-events: none; }
  .fancy-checkbox input { display: none; }
  
  .check-box {
    width: 24px;
    height: 24px;
    border: 2px solid var(--border-color);
    border-radius: 10px;
    transition: all 0.2s;
    display: flex;
    align-items: center;
    justify-content: center;
  }

  .check-icon {
    font-size: 16px;
    color: var(--on-primary);
    opacity: 0;
    transform: scale(0.5);
    transition: all 0.2s cubic-bezier(0.175, 0.885, 0.32, 1.275);
  }

  .fancy-checkbox input:checked + .check-box {
    background: var(--primary-color);
    border-color: var(--primary-color);
  }

  .fancy-checkbox input:checked + .check-box .check-icon {
    opacity: 1;
    transform: scale(1);
  }

  .check-text strong { display: block; font-size: 1rem; color: var(--text-primary); }
  .check-text span { font-size: 0.85rem; color: var(--text-muted); }

  /* ACTIONS */
  .form-actions {
    display: flex;
    justify-content: flex-end;
    align-items: center;
    gap: 2rem;
    margin-top: 2rem;
  }

  .btn-text {
    background: none;
    border: none;
    font-weight: 700;
    color: var(--text-muted);
    cursor: pointer;
  }

  .error-toast {
    background: var(--error-container);
    color: var(--on-error-container);
    padding: 1rem 1.5rem;
    border-radius: 1rem;
    margin-bottom: 2rem;
    display: flex;
    align-items: center;
    gap: 1rem;
    font-weight: 700;
    box-shadow: var(--shadow-sm);
    animation: shake 0.4s ease-out;
  }

  @keyframes shake {
    0%, 100% { transform: translateX(0); }
    25% { transform: translateX(-5px); }
    50% { transform: translateX(5px); }
    75% { transform: translateX(-5px); }
  }

  .spin { animation: spin 1s infinite linear; }
  @keyframes spin { from { rotate: 0deg; } to { rotate: 360deg; } }

  @media (max-width: 900px) {
    .hero-header { padding: 2rem 1rem; }
    .hero-header h1 { font-size: 1.75rem; }
  }
</style>
