<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { goto } from '$app/navigation';
  
  let username = $state('');
  let fullName = $state('');
  let errorMsg = $state('');
  let loading = $state(false);

  async function handleSubmit() {
    if (!username) return;
    loading = true;
    errorMsg = '';
    
    try {
      const res = await apiFetch('/api/v2/user/settings/orgs', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          username,
          full_name: fullName
        })
      });
      
      const data = await res.json();
      if (res.ok) {
        goto(`/org/${username}`);
      } else {
        errorMsg = data.message || 'Failed to create organization';
      }
    } catch (e) {
      errorMsg = 'An error occurred. Please try again.';
    } finally {
      loading = false;
    }
  }
</script>

<div class="create-wrapper">
  <header class="hero-header">
    <div class="inner">
      <h1>Nueva Organización</h1>
      <p>Crea una organización para colaborar con tu equipo en múltiples repositorios.</p>
    </div>
  </header>

  <main class="create-content">
    <form class="create-form tonal island" onsubmit={(e) => { e.preventDefault(); handleSubmit(); }}>
      {#if errorMsg}
        <div class="error-toast">
          <span class="material-symbols-rounded">error</span>
          {errorMsg}
        </div>
      {/if}

      <div class="input-group">
        <label for="org-name">Nombre de la Organización</label>
        <input 
          type="text" 
          id="org-name" 
          bind:value={username} 
          placeholder="nombre-de-tu-org"
          required
          disabled={loading}
          pattern="[a-zA-Z0-9_\-\.]+"
        />
        <p class="help">{window.location.origin}/org/{username || 'nombre'}</p>
      </div>

      <div class="input-group">
        <label for="full-name">Nombre Completo <span class="optional">(opcional)</span></label>
        <input 
          type="text" 
          id="full-name" 
          bind:value={fullName} 
          placeholder="Nombre Público de la Organización"
          disabled={loading}
        />
      </div>

      <div class="form-actions">
        <button type="button" class="btn btn-ghost" onclick={() => history.back()}>Cancelar</button>
        <button type="submit" class="btn btn-primary" disabled={loading || !username}>
          {#if loading}
            <span class="material-symbols-rounded spin">sync</span> Creando...
          {:else}
            Crear Organización
          {/if}
        </button>
      </div>
    </form>
  </main>
</div>

<style>
  .create-wrapper { animation: fadeIn 0.4s ease-out; }
  
  .hero-header { background: var(--surface-variant); padding: 4rem 2rem; border-bottom: 1px solid var(--md-sys-color-outline-variant); margin-bottom: 2rem; }
  .hero-header .inner { max-width: 800px; margin: 0 auto; }
  .hero-header h1 { margin: 0; font-size: 2.5rem; font-weight: 900; }
  .hero-header p { margin: 0.5rem 0 0; color: var(--text-muted); font-size: 1.1rem; }

  .create-content { max-width: 600px; margin: 0 auto; padding: 2rem; }
  
  .create-form { padding: 2.5rem; display: flex; flex-direction: column; gap: 2rem; }
  
  .input-group { display: flex; flex-direction: column; gap: 0.75rem; }
  .input-group label { font-weight: 700; color: var(--text-color); }
  .input-group input { 
    background: var(--surface-variant); 
    border: 1px solid var(--md-sys-color-outline-variant); 
    padding: 1rem; 
    border-radius: 0.75rem; 
    color: var(--text-color);
    font-size: 1rem;
    font-weight: 600;
  }
  .input-group input:focus { border-color: var(--primary-color); outline: none; box-shadow: 0 0 0 3px rgba(var(--primary-rgb), 0.1); }
  .optional { font-weight: 400; opacity: 0.6; font-size: 0.85rem; }
  .help { font-size: 0.85rem; color: var(--text-muted); margin: 0; }

  .form-actions { display: flex; justify-content: flex-end; gap: 1rem; margin-top: 1rem; }
  
  .error-toast { background: rgba(239, 68, 68, 0.1); color: #ef4444; padding: 1rem; border-radius: 0.75rem; display: flex; align-items: center; gap: 0.75rem; font-weight: 600; }
  
  .spin { animation: spin 1s infinite linear; }
  @keyframes spin { from { transform: rotate(0deg); } to { transform: rotate(360deg); } }
  @keyframes fadeIn { from { opacity: 0; transform: translateY(10px); } to { opacity: 1; transform: translateY(0); } }
</style>
