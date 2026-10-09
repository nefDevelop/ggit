<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import Input from '$lib/components/ui/Input.svelte';
  import Card from '$lib/components/ui/Card.svelte';
  import Select from '$lib/components/ui/Select.svelte';

  import ShieldCheck from "lucide-svelte/icons/shield-check";
  import Terminal from "lucide-svelte/icons/terminal";
  import AlertCircle from "lucide-svelte/icons/alert-circle";
  import Database from "lucide-svelte/icons/database";
  import Settings from "lucide-svelte/icons/settings";
  import WifiOff from "lucide-svelte/icons/wifi-off";

  let settings = $state<any>(null);
  let loading = $state(true);
  let alreadyInstalled = $state(false);
  let connectionError = $state(false);
  let form = $state<any>({});
  let submitting = $state(false);
  let error = $state('');

  async function fetchSettings() {
    loading = true;
    alreadyInstalled = false;
    connectionError = false;
    try {
      // silent para no disparar el toast de "sin permisos" cuando ya está instalado (403).
      const res = await apiFetch('/api/v2/install/settings', { silent: true });
      if (res.status === 403) {
        alreadyInstalled = true;
        return;
      }
      if (res.ok) {
        const data = await res.json();
        // Solo el cuerpo esperado abre el asistente; un 2xx con otra cosa
        // (fallback, proxy) se trata como ya instalado.
        if (data && data.default_settings) {
          settings = data;
          form = { ...data.default_settings };
        } else {
          alreadyInstalled = true;
        }
      }
    } catch (e) {
      connectionError = true;
    } finally {
      loading = false;
    }
  }

  async function install() {
    submitting = true;
    error = '';
    try {
      const res = await apiFetch('/api/v2/install', {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(form)
      });
      if (res.ok) {
        window.location.href = '/user/login';
      } else {
        // Si alguien instaló en paralelo, el backend responde 403.
        if (res.status === 403) {
          alreadyInstalled = true;
          settings = null;
          return;
        }
        const data = await res.json();
        error = data.error || data.message || 'La instalación ha fallado.';
      }
    } catch (e) {
      error = 'Error de red. Comprueba la conexión con el servidor.';
    } finally {
      submitting = false;
    }
  }

  onMount(fetchSettings);
</script>

<div class="install-page">
  <div class="background-decor">
    <div class="circle c1"></div>
    <div class="circle c2"></div>
    <div class="circle c3"></div>
  </div>

  <div class="install-container">
    <Card variant="highest" padding="none" class="install-card">
      <div class="install-layout">
        <div class="install-brand">
          <div class="brand-content">
            <Terminal size={64} class="logo-icon" />
            <h1>Ggit</h1>
            <p>Tu plataforma Git autoalojada, ligera y elegante, lista en minutos.</p>
          </div>
          <div class="brand-footer">
            <span>Instalación</span>
          </div>
        </div>

        <div class="install-form-area">
          {#if loading}
            <div class="spinner"></div>
          {:else if alreadyInstalled}
            <div class="state-block animate-in">
              <ShieldCheck size={48} class="state-icon ok" />
              <h2>Ya está instalado</h2>
              <p>Esta instancia de Ggit ya está configurada. Inicia sesión para continuar.</p>
              <Button variant="primary" href="/user/login" class="state-btn">
                Ir al inicio de sesión
              </Button>
            </div>
          {:else if connectionError}
            <div class="state-block animate-in">
              <WifiOff size={48} class="state-icon warn" />
              <h2>Sin conexión</h2>
              <p>No se puede contactar con el servidor. Verifica que esté encendido e inténtalo de nuevo.</p>
              <Button variant="outline" onclick={fetchSettings} class="state-btn">
                Reintentar
              </Button>
            </div>
          {:else if settings}
            <header class="form-header">
              <h2>Configuración inicial</h2>
              <p>Completa estos ajustes para poner en marcha tu instancia</p>
            </header>

            <form onsubmit={(e) => { e.preventDefault(); install(); }} class="install-form">
              <section class="form-section">
                <div class="section-title">
                  <Database size={20} />
                  <h3>Base de datos</h3>
                </div>
                <div class="grid">
                  <Select
                    label="Tipo"
                    bind:value={form.db_type}
                    options={(settings.db_types || []).map((t: string) => ({ value: t, label: t }))}
                  />
                  <Input id="db_host" label="Host" bind:value={form.db_host} placeholder="localhost" />
                  <Input id="db_user" label="Usuario" bind:value={form.db_user} placeholder="ggit" />
                  <Input id="db_name" label="Nombre" bind:value={form.db_name} placeholder="ggit" />
                  <Input id="db_path" label="Ruta (SQLite)" bind:value={form.db_path} placeholder="data/ggit.db" />
                </div>
              </section>

              <section class="form-section">
                <div class="section-title">
                  <Settings size={20} />
                  <h3>General</h3>
                </div>
                <div class="grid">
                  <Input id="app_name" label="Nombre de la aplicación" bind:value={form.app_name} />
                  <Input id="domain" label="Dominio" bind:value={form.domain} placeholder="localhost" />
                  <Input id="http_port" label="Puerto HTTP" bind:value={form.http_port} />
                  <Input id="app_url" label="URL pública" bind:value={form.app_url} placeholder="http://localhost:3000/" />
                  <Input id="run_user" label="Usuario de ejecución" bind:value={form.run_user} />
                  <Input id="default_branch" label="Rama por defecto" bind:value={form.default_branch} />
                </div>
              </section>

              <div class="info-banner">
                <ShieldCheck size={40} />
                <div class="info-text">
                  <strong>Administración automática</strong>
                  <p>La primera cuenta que se cree tras la instalación recibirá privilegios de administración automáticamente.</p>
                </div>
              </div>

              {#if error}
                <div class="error-banner animate-in">
                  <AlertCircle size={20} />
                  <p>{error}</p>
                </div>
              {/if}

              <Button type="submit" variant="primary" class="install-btn" loading={submitting}>
                Instalar Ggit
              </Button>
            </form>
          {/if}
        </div>
      </div>
    </Card>
  </div>
</div>

<style>
  .install-page {
    position: relative;
    min-height: 100%;
    width: 100%;
    display: flex;
    justify-content: center;
    align-items: center;
    overflow: hidden;
    background: transparent;
    padding: 2rem 0;
  }

  .background-decor {
    position: absolute;
    inset: 0;
    overflow: hidden;
    z-index: 0;
    filter: blur(100px);
    opacity: 0.35;
    pointer-events: none;
  }

  .circle {
    position: absolute;
    border-radius: 50%;
    animation: float 25s infinite alternate ease-in-out;
  }

  .c1 {
    width: 600px;
    height: 600px;
    background: radial-gradient(circle, var(--primary-color), var(--status-info));
    top: -150px;
    right: -150px;
  }

  .c2 {
    width: 500px;
    height: 500px;
    background: radial-gradient(circle, var(--status-accent), var(--status-accent));
    bottom: -150px;
    left: -100px;
    animation-delay: -5s;
    animation-duration: 30s;
  }

  .c3 {
    width: 400px;
    height: 400px;
    background: radial-gradient(circle, var(--status-info), var(--status-success));
    top: 50%;
    left: -200px;
    transform: translateY(-50%);
    animation-delay: -10s;
    animation-duration: 35s;
  }

  @keyframes float {
    0% { transform: translate(0, 0) rotate(0deg); }
    100% { transform: translate(60px, 40px) rotate(10deg); }
  }

  .install-container {
    position: relative;
    z-index: 10;
    width: 100%;
    max-width: 1000px;
    padding: 2rem;
    animation: slideUp 0.8s cubic-bezier(0.2, 0.8, 0.2, 1);
  }

  :global(.install-card) {
    position: relative;
    border-radius: 2.5rem !important;
    overflow: hidden !important;
    box-shadow:
      0 32px 80px rgba(0, 0, 0, 0.4),
      inset 0 1px 1px rgba(255, 255, 255, 0.1),
      inset 0 -1px 1px rgba(0, 0, 0, 0.2) !important;
    border: 1px solid rgba(0, 0, 0, 0.3) !important;
    background: var(--md-sys-color-surface-container-highest) !important;
  }

  .install-layout {
    display: grid;
    grid-template-columns: 1fr 1.6fr;
  }

  .install-brand {
    padding: 3.5rem;
    display: flex;
    flex-direction: column;
    justify-content: space-between;
    background: var(--md-sys-color-surface-container-low);
    border-right: 1px solid var(--md-sys-color-outline-variant);
  }

  .brand-content {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  :global(.logo-icon) {
    color: var(--primary-color);
  }

  .install-brand h1 {
    font-size: 3.5rem;
    font-weight: 950;
    margin: 0;
    letter-spacing: -0.05em;
    background: linear-gradient(to bottom, var(--text-color), var(--primary-color));
    -webkit-background-clip: text;
    background-clip: text;
    -webkit-text-fill-color: transparent;
  }

  .install-brand p {
    font-size: 1.1rem;
    line-height: 1.6;
    color: var(--text-muted);
    max-width: 300px;
    margin: 0;
  }

  .brand-footer {
    font-size: 0.8rem;
    font-weight: 800;
    color: var(--text-muted);
    opacity: 0.5;
    text-transform: uppercase;
    letter-spacing: 0.1em;
  }

  .install-form-area {
    padding: 3.5rem;
    display: flex;
    flex-direction: column;
    justify-content: center;
    background: var(--md-sys-color-surface-container-highest);
    min-height: 480px;
  }

  .form-header {
    margin-bottom: 2rem;
  }

  .form-header h2 {
    font-size: 2.25rem;
    font-weight: 900;
    margin: 0;
    letter-spacing: -0.02em;
  }

  .form-header p {
    color: var(--text-muted);
    margin-top: 0.5rem;
    font-size: 1.1rem;
  }

  .install-form {
    display: flex;
    flex-direction: column;
    gap: 2rem;
  }

  .form-section {
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
  }

  .section-title {
    display: flex;
    align-items: center;
    gap: 0.6rem;
    color: var(--primary-color);
  }

  .section-title h3 {
    margin: 0;
    font-size: 0.85rem;
    font-weight: 800;
    text-transform: uppercase;
    letter-spacing: 0.08em;
  }

  .grid {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1.25rem;
  }

  .info-banner {
    display: flex;
    gap: 1rem;
    align-items: center;
    padding: 1.5rem;
    border-radius: 1rem;
    border: 1px solid var(--primary-color);
    background: rgba(var(--primary-rgb), 0.05);
  }

  .info-banner :global(svg) {
    color: var(--primary-color);
    flex-shrink: 0;
  }

  .info-text strong {
    display: block;
    font-size: 1.1rem;
    margin-bottom: 0.25rem;
  }

  .info-text p {
    margin: 0;
    font-size: 0.9rem;
    line-height: 1.4;
    color: var(--text-muted);
  }

  .error-banner {
    background: rgba(var(--error-rgb), 0.1);
    color: var(--error);
    padding: 1rem 1.25rem;
    border-radius: 1rem;
    display: flex;
    align-items: center;
    gap: 0.75rem;
    font-weight: 600;
    font-size: 0.95rem;
  }

  .error-banner p { margin: 0; }

  :global(.install-btn) {
    height: 3.5rem !important;
    font-size: 1.1rem !important;
    font-weight: 800 !important;
    border-radius: 1.25rem !important;
  }

  .state-block {
    display: flex;
    flex-direction: column;
    align-items: center;
    text-align: center;
    gap: 1rem;
    padding: 2rem 0;
  }

  .state-block h2 {
    font-size: 1.75rem;
    font-weight: 900;
    margin: 0;
  }

  .state-block p {
    color: var(--text-muted);
    margin: 0 0 1rem 0;
    max-width: 380px;
    line-height: 1.6;
  }

  :global(.state-icon.ok) { color: var(--status-success); }
  :global(.state-icon.warn) { color: var(--status-accent); }

  :global(.state-btn) {
    min-width: 240px;
    height: 3.25rem !important;
    font-weight: 800 !important;
    border-radius: 1.25rem !important;
  }

  .spinner {
    margin: 3rem auto;
  }

  @keyframes slideUp {
    from { opacity: 0; transform: translateY(30px); }
    to { opacity: 1; transform: translateY(0); }
  }

  .animate-in {
    animation: fadeIn 0.3s ease-out;
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: scale(0.95); }
    to { opacity: 1; transform: scale(1); }
  }

  @media (max-width: 850px) {
    .install-layout {
      grid-template-columns: 1fr;
    }
    .install-brand {
      display: none;
    }
    .install-container {
      max-width: 560px;
    }
    .install-form-area {
      padding: 2.5rem;
    }
    .grid {
      grid-template-columns: 1fr;
    }
  }
</style>
