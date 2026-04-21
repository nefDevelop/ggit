<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { authState } from '$lib/auth.svelte';
  import { goto } from '$app/navigation';

  let name = $state('');
  let description = $state('');
  let isPrivate = $state(false);
  let autoInit = $state(true);

  let gitignores = $state<string[]>([]);
  let licenses = $state<string[]>([]);
  let readmes = $state<string[]>([]);
  
  let selectedGitignore = $state('');
  let selectedLicense = $state('');
  let selectedReadme = $state('Default');
  
  let errorMsg = $state('');
  let loading = $state(false);

  async function loadTemplates() {
    try {
      const res = await apiFetch('/api/v2/repo/templates');
      if (res.ok) {
        const data = await res.json();
        gitignores = data.gitignores || [];
        licenses = data.licenses || [];
        readmes = data.readmes || [];
      }
    } catch (e) {
      console.error("Failed to load templates", e);
    }
  }

  import { onMount } from 'svelte';
  onMount(loadTemplates);

  function getCsrfToken() {
    if (typeof document === 'undefined') return '';
    const match = document.cookie.match(new RegExp('(^| )XSRF-TOKEN=([^;]+)'));
    return match ? match[2] : '';
  }

  async function handleSubmit() {
    if (!name) return;
    loading = true;
    errorMsg = '';
    
    try {
      const res = await apiFetch('/api/v2/user/repos', {
        method: 'POST',
        headers: { 
          'Content-Type': 'application/json',
          'X-Csrf-Token': getCsrfToken()
        },
        body: JSON.stringify({
          name,
          description,
          private: isPrivate,
          auto_init: autoInit,
          gitignores: selectedGitignore,
          license: selectedLicense,
          readme: selectedReadme
        })
      });
      
      const data = await res.json();
      if (res.ok) {
        goto(`/${authState.user.username}/${name}`);
      } else {
        errorMsg = data.message || 'Failed to create repository';
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
        <h1>Crear repositorio</h1>
        <p>Un repositorio contiene todos los archivos de tu proyecto, incluyendo el historial de revisiones.</p>
      </div>
    </div>
  </header>

  <main class="create-grid">
    <!-- RIGHT COLUMN: FORM -->
    <form class="create-form-content" onsubmit={(e) => { e.preventDefault(); handleSubmit(); }}>
      {#if errorMsg}
        <div class="error-toast">
          <span class="material-symbols-rounded">error</span>
          {errorMsg}
        </div>
      {/if}

      <section class="form-section tonal">
        <div class="section-header">
          <span class="material-symbols-rounded">id_card</span>
          <h2>Información Básica</h2>
        </div>
        
        <div class="input-field">
          <label for="repo-name">Nombre del Repositorio</label>
          <div class="input-wrapper">
            <span class="prefix">{authState.user?.username} /</span>
            <input 
              type="text" 
              id="repo-name" 
              bind:value={name} 
              placeholder="mi-proyecto-increible"
            pattern="[a-zA-Z0-9_\-\.]+"
            title="Solo letras, números, guiones y puntos"
            disabled={loading}
            autofocus
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
      </section>

      <section class="form-section tonal">
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
      </section>

      <section class="form-section tonal">
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
            <strong>Generar README.md automático</strong>
            <span>Crea un archivo de presentación para poder clonar el repo de inmediato.</span>
          </div>
        </label>

        {#if autoInit}
          <div class="init-details">
            <div class="input-field">
              <label for="gitignore">.gitignore</label>
              <select id="gitignore" bind:value={selectedGitignore} disabled={loading}>
                <option value="">Ninguno</option>
                {#each gitignores as g}
                  <option value={g}>{g}</option>
                {/each}
              </select>
            </div>

            <div class="input-field">
              <label for="license">Licencia</label>
              <select id="license" bind:value={selectedLicense} disabled={loading}>
                <option value="">Ninguna</option>
                {#each licenses as l}
                  <option value={l}>{l}</option>
                {/each}
              </select>
            </div>

            <div class="input-field">
              <label for="readme">Plantilla README</label>
              <select id="readme" bind:value={selectedReadme} disabled={loading}>
                {#each readmes as r}
                  <option value={r}>{r}</option>
                {/each}
              </select>
            </div>
          </div>
        {/if}
      </section>

      <div class="form-actions">
        <button type="button" class="btn-text" onclick={() => goto('/dashboard')}>Cancelar</button>
        <button type="submit" class="btn-action-primary" disabled={loading || !name}>
          {#if loading}
            <span class="material-symbols-rounded spin">sync</span>
            Despegando...
          {:else}
            <span class="material-symbols-rounded">rocket_launch</span>
            Crear Repositorio
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
    -webkit-background-clip: text;
    -webkit-text-fill-color: transparent;
  }

  .hero-header p {
    font-size: 1.1rem;
    color: var(--text-muted);
    margin: 0.5rem 0 0 0;
  }

  /* GRID */
  .create-grid {
    max-width: 800px; /* Centramos el formulario para que se vea bien al quitar la columna */
    margin: 0 auto;
    padding-bottom: 5rem;
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

  .input-field label { font-weight: 700; font-size: 0.95rem; color: var(--text-secondary); }

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
    background: rgba(0,0,0,0.05);
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

  .init-details {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1.5rem;
    margin-top: 1.5rem;
    padding-top: 1.5rem;
    border-top: 1px solid var(--border-color);
  }

  .init-details .input-field:last-child {
    grid-column: span 2;
  }

  @media (max-width: 600px) {
    .init-details {
      grid-template-columns: 1fr;
    }
    .init-details .input-field:last-child {
      grid-column: span 1;
    }
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
    border-radius: 6px;
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

  .btn-hero {
    background: var(--primary-color);
    color: var(--on-primary);
    padding: 1rem 2.5rem;
    border-radius: 3rem;
    border: none;
    font-weight: 800;
    font-size: 1.1rem;
    display: flex;
    align-items: center;
    gap: 1rem;
    cursor: pointer;
    box-shadow: 0 10px 25px rgba(var(--primary-rgb), 0.2);
    transition: all 0.2s;
  }

  .btn-hero:hover:not(:disabled) {
    transform: translateY(-3px) scale(1.02);
    box-shadow: 0 15px 35px rgba(var(--primary-rgb), 0.3);
  }

  .btn-hero:disabled { opacity: 0.5; cursor: not-allowed; }

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
