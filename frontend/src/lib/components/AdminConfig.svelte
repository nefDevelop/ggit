<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from "svelte";
  import { i18n } from '$lib/i18n.svelte';
  import Card from './ui/Card.svelte';
  import Button from './ui/Button.svelte';
  import AdminConfigEditor from './AdminConfigEditor.svelte';

  let config = $state<any>(null);
  let loading = $state(true);

  let regSecurity = $state({ enable: false, question: "", answer: "" });
  let savingSecurity = $state(false);
  let editingSection = $state<string | null>(null);

  async function fetchConfig() {
    loading = true;
    try {
      const res = await apiFetch("/api/v2/admin/config");
      if (res.ok) {
        config = await res.json();
        // Specific mapping for the shortcut card
        if (config.auth) {
          regSecurity.enable = config.auth.ENABLE_REGISTRATION_QUESTION === 'true';
          regSecurity.question = config.auth.REGISTRATION_QUESTION;
          regSecurity.answer = config.auth.REGISTRATION_ANSWER;
        }
      }
    } catch (e) { console.error(e); } finally { loading = false; }
  }

  async function updateAuthSettings() {
    savingSecurity = true;
    try {
      const res = await apiFetch("/api/v2/admin/config", {
        method: "PUT",
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify({
          section: "auth",
          data: {
            ENABLE_REGISTRATION_QUESTION: String(regSecurity.enable),
            REGISTRATION_QUESTION: regSecurity.question,
            REGISTRATION_ANSWER: regSecurity.answer,
          },
        }),
      });
      if (res.ok) fetchConfig();
    } catch (e) { console.error(e); } finally { savingSecurity = false; }
  }

  const sections = [
    { key: "app", label: "Aplicación", icon: "apps" },
    { key: "server", label: "Servidor", icon: "dns" },
    { key: "database", label: "Base de Datos", icon: "storage" },
    { key: "security", label: "Seguridad", icon: "security" },
    { key: "repository", label: "Repositorios", icon: "book" },
    { key: "git", label: "Git System", icon: "terminal" },
    { key: "email", label: "Servicio de Correo", icon: "mail" },
  ];

  let testEmail = $state("");
  let sendingTest = $state(false);

  async function sendTestEmail() {
    if (!testEmail) return;
    sendingTest = true;
    try {
      await apiFetch("/api/v2/admin/config/test_email", { method: "POST", headers: { "Content-Type": "application/json" }, body: JSON.stringify({ email: testEmail }) });
    } catch (e) { console.error(e); } finally { sendingTest = false; }
  }

  onMount(fetchConfig);
</script>

<Card class="settings-premium-card" padding="lg">
  <header class="section-header">
    <div class="title-with-icon">
      <span class="material-symbols-rounded hero-icon">settings</span>
      <div>
        <h2>Configuración del Sistema</h2>
        <p class="subtitle">Estado detallado y ajustes del motor Ggit</p>
      </div>
    </div>
  </header>

  {#if loading}
     <div class="skeleton-config"></div>
  {:else}
    <div class="config-stack">
      {#each sections as section}
        <div class="config-modular-section">
          <div class="section-tag-wrapper">
            <div class="section-tag">
              <span class="material-symbols-rounded">{section.icon}</span>
              {section.label}
            </div>
            <Button variant="ghost" size="sm" icon="edit" onclick={() => editingSection = section.key}>
              Configurar
            </Button>
          </div>
          <div class="config-grid">
            {#each Object.entries(config[section.key] || {}) as [label, value]}
              <div class="config-tile tonal">
                <span class="tile-key">{label}</span>
                <span class="tile-val">{typeof value === 'object' ? JSON.stringify(value) : (value || '—')}</span>
              </div>
            {/each}
          </div>
        </div>
      {/each}

      <div class="special-action-card island">
        <div class="action-header">
           <span class="material-symbols-rounded">gpp_good</span>
           <h4>Seguridad de Registro</h4>
        </div>
        <div class="action-body">
           <p>Configura una pregunta de seguridad para evitar registros automatizados de bots.</p>
           <div class="action-form">
              <label class="premium-switch">
                 <input type="checkbox" bind:checked={regSecurity.enable} />
                 <span class="switch-slider"></span>
                 <span class="switch-label">Habilitar Pregunta de Seguridad</span>
              </label>
              
              {#if regSecurity.enable}
                <div class="field-row">
                  <div class="field">
                    <label for="reg-q">Pregunta</label>
                    <input id="reg-q" type="text" bind:value={regSecurity.question} placeholder="Ej: ¿Cuál es el nombre de este servicio?" />
                  </div>
                  <div class="field">
                    <label for="reg-a">Respuesta</label>
                    <input id="reg-a" type="text" bind:value={regSecurity.answer} placeholder="Ej: Ggit" />
                  </div>
                </div>
              {/if}
              <Button variant="primary" onclick={updateAuthSettings} disabled={savingSecurity}>
                {savingSecurity ? 'Guardando...' : 'Guardar Ajustes de Seguridad'}
              </Button>
           </div>
        </div>
      </div>
    </div>
  {/if}
</Card>

{#if editingSection}
  <AdminConfigEditor 
    section={editingSection} 
    data={config[editingSection]} 
    onSave={() => { editingSection = null; fetchConfig(); }} 
    onCancel={() => editingSection = null} 
  />
{/if}

<style>
  :global(.settings-premium-card) { border-radius: 3rem; background: rgba(var(--surface-rgb), 0.5); backdrop-filter: blur(24px); border: 1px solid var(--md-sys-color-outline-variant); }
  .section-header { margin-bottom: 4rem; }
  .header-flex { display: flex; justify-content: space-between; align-items: center; gap: 2rem; }
  @media (max-width: 600px) {
    .header-flex { flex-direction: column; align-items: stretch; text-align: center; }
    .title-with-icon { flex-direction: column; gap: 1rem; }
  }

  .title-with-icon { display: flex; align-items: center; gap: 2rem; }
  .hero-icon { font-size: 3rem; color: var(--primary-color); background: rgba(var(--primary-rgb), 0.1); padding: 1.25rem; border-radius: 2rem; flex-shrink: 0; }
  h2 { margin: 0; font-size: 2.25rem; font-weight: 950; letter-spacing: -0.04em; }
  .subtitle { margin: 0.5rem 0 0 0; color: var(--text-muted); font-size: 1.1rem; }

  .config-stack { display: flex; flex-direction: column; gap: 4rem; }

  .config-modular-section { display: flex; flex-direction: column; gap: 1.5rem; }
  
  .section-tag-wrapper { display: flex; justify-content: space-between; align-items: center; border-bottom: 2px solid var(--md-sys-color-outline-variant); padding-bottom: 0.75rem; gap: 1rem; }
  .section-tag { font-size: 0.8rem; font-weight: 950; text-transform: uppercase; color: var(--primary-color); letter-spacing: 0.15em; display: flex; align-items: center; gap: 0.75rem; border-bottom: none; padding-bottom: 0; width: fit-content; }
  .section-tag span { font-size: 1.25rem; }

  .config-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); gap: 1rem; }
  @media (max-width: 600px) {
    .config-grid { grid-template-columns: 1fr; }
    .config-tile { padding: 1rem; }
  }
  .config-tile { padding: 1.25rem 1.75rem; border-radius: 1.5rem; display: flex; flex-direction: column; gap: 0.25rem; border: 1px solid var(--md-sys-color-outline-variant); transition: all 0.2s; min-width: 0; }
  .config-tile:hover { transform: translateY(-3px); background: rgba(255,255,255,0.02); border-color: var(--primary-color); }
  .tile-key { font-size: 0.7rem; font-weight: 800; text-transform: uppercase; color: var(--text-muted); opacity: 0.7; }
  .tile-val { font-family: 'JetBrains Mono', monospace; font-size: 0.95rem; color: var(--text-color); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

  .special-action-card { padding: 3rem; border-radius: 2.5rem; border: 1px solid var(--md-sys-color-outline-variant); background: rgba(var(--primary-rgb), 0.03); }
  @media (max-width: 600px) {
    .special-action-card { padding: 1.5rem; border-radius: 1.5rem; }
    .action-header span { font-size: 2rem; }
    .action-header h4 { font-size: 1.25rem; }
    .field-row { grid-template-columns: 1fr; }
  }
  .field { display: flex; flex-direction: column; gap: 0.5rem; }
  .field label { font-size: 0.8rem; font-weight: 800; color: var(--text-color); }
  .field input { background: var(--surface-variant); border: 1px solid var(--md-sys-color-outline-variant); padding: 0.85rem 1.25rem; border-radius: 1rem; color: var(--text-color); }

  .premium-switch { display: flex; align-items: center; gap: 1rem; cursor: pointer; }
  .premium-switch input { display: none; }
  .switch-slider { width: 44px; height: 24px; background: var(--surface-variant); border-radius: 100px; position: relative; border: 1px solid var(--md-sys-color-outline-variant); transition: 0.3s; }
  .switch-slider::before { content: ""; position: absolute; width: 18px; height: 18px; background: var(--text-muted); border-radius: 50%; left: 2px; top: 2px; transition: 0.3s; }
  input:checked + .switch-slider { background: var(--primary-color); }
  input:checked + .switch-slider::before { transform: translateX(20px); background: white; }
  .switch-label { font-size: 0.95rem; font-weight: 700; color: var(--text-color); }

  .skeleton-config { height: 600px; background: var(--surface-variant); border-radius: 2rem; animation: pulse 1.5s infinite; }
</style>
