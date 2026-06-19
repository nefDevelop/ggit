<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from "svelte";
  import { i18n } from '$lib/i18n.svelte';
  import Card from './ui/Card.svelte';
  import Button from './ui/Button.svelte';
  import AdminConfigEditor from './AdminConfigEditor.svelte';
  import PageHeader from './ui/PageHeader.svelte';
  
import { Settings, AppWindow, Server, Database, ShieldCheck, Book, Terminal, Mail, Edit2, ShieldAlert } from "lucide-svelte";


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
        if (config?.auth) {
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
    { key: "app", label: "Aplicación", icon: AppWindow },
    { key: "server", label: "Servidor", icon: Server },
    { key: "database", label: "Base de Datos", icon: Database },
    { key: "security", label: "Seguridad", icon: ShieldCheck },
    { key: "repository", label: "Repositorios", icon: Book },
    { key: "git", label: "Git System", icon: Terminal },
    { key: "email", label: "Servicio de Correo", icon: Mail },
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

<div class="admin-config-container">
    <PageHeader 
      title="Configuración del Sistema"
      subtitle="Estado detallado y ajustes del motor Ggit"
      icon={Settings}
    />

    <div class="admin-config-body">
      {#if loading}
        <div class="skeleton-config"></div>
      {:else}
        <div class="config-stack">
          {#each sections as section}
            {@const Icon = section.icon}
            <div class="config-modular-section">
              <div class="section-tag-wrapper">
                <div class="section-tag">
                  <Icon size={18} />
                  {section.label}
                </div>
                <Button variant="ghost" size="sm" icon={Edit2} onclick={() => editingSection = section.key}>
                  Configurar
                </Button>
              </div>
              <div class="config-grid">
                {#each Object.entries(config?.[section.key] || {}) as [label, value]}
                  <div class="config-tile">
                    <span class="tile-key">{label}</span>
                    <span class="tile-val">{typeof value === 'object' ? JSON.stringify(value) : (value || '—')}</span>
                  </div>
                {/each}
              </div>
            </div>
          {/each}

          <div class="special-action-card island">
            <div class="action-header">
               <ShieldAlert size={32} />
               <h4>Seguridad de Registro</h4>
            </div>
            <div class="action-body">
               <p>Configura una pregunta de seguridad para evitar registros automatizados de bots. El sistema cuenta además con un Honeypot integrado para bloquear el SPAM de manera silenciosa.</p>
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
                  <div class="form-actions">
                    <Button variant="primary" onclick={updateAuthSettings} disabled={savingSecurity}>
                      {savingSecurity ? 'Guardando...' : 'Guardar Ajustes de Seguridad'}
                    </Button>
                  </div>
               </div>
            </div>
          </div>
        </div>
      {/if}
    </div>
</div>

{#if editingSection}
  <AdminConfigEditor 
    section={editingSection} 
    data={config?.[editingSection]} 
    onSave={() => { editingSection = null; fetchConfig(); }} 
    onCancel={() => editingSection = null} 
  />
{/if}

<style>
  .admin-config-body {
    padding: 0 1.5rem 4rem;
    animation: fadeIn 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(4px); }
    to { opacity: 1; transform: translateY(0); }
  }
  .config-stack { display: flex; flex-direction: column; gap: 4rem; }

  .config-stack { display: flex; flex-direction: column; gap: 4rem; }

  .config-modular-section { display: flex; flex-direction: column; gap: 1.5rem; }
  
  .section-tag-wrapper { display: flex; justify-content: space-between; align-items: center; border-bottom: 2px solid var(--md-sys-color-outline-variant); padding-bottom: 0.75rem; gap: 1rem; }
  .section-tag { font-size: 0.8rem; font-weight: 950; text-transform: uppercase; color: var(--primary-color); letter-spacing: 0.15em; display: flex; align-items: center; gap: 0.75rem; border-bottom: none; padding-bottom: 0; width: fit-content; }

  .config-grid { display: grid; grid-template-columns: repeat(auto-fill, minmax(300px, 1fr)); gap: 1rem; }
  @media (max-width: 600px) {
    .config-grid { grid-template-columns: 1fr; }
    .config-tile { padding: 1rem; }
  }
  .config-tile { padding: 1.25rem 1.75rem; border-radius: 0.75rem; display: flex; flex-direction: column; gap: 0.25rem; border: 1px solid var(--md-sys-color-outline-variant); transition: all 0.2s; min-width: 0; }
  .config-tile:hover { transform: translateY(-3px); background: rgba(255,255,255,0.02); border-color: var(--primary-color); }
  .tile-key { font-size: 0.7rem; font-weight: 800; text-transform: uppercase; color: var(--text-muted); opacity: 0.7; }
  .tile-val { font-family: 'JetBrains Mono', monospace; font-size: 0.95rem; color: var(--text-color); overflow: hidden; text-overflow: ellipsis; white-space: nowrap; }

  .special-action-card { padding: 3rem; border-radius: 1.25rem; border: 1px solid var(--md-sys-color-outline-variant); background: rgba(var(--primary-rgb), 0.03); }
  @media (max-width: 600px) {
    .special-action-card { padding: 1.5rem; border-radius: 1.5rem; }
    .action-header h4 { font-size: 1.1rem; }
    .field-row { grid-template-columns: 1fr; }
    .form-actions :global(button) { width: 100%; }
  }
  .action-header { display: flex; align-items: center; gap: 1rem; margin-bottom: 1.5rem; color: var(--primary-color); }
  .action-header h4 { margin: 0; font-size: 1.5rem; font-weight: 900; }

  .action-form { display: flex; flex-direction: column; gap: 2.5rem; margin-top: 2rem; }
  .form-actions { display: flex; justify-content: flex-start; margin-top: 1.5rem; padding-top: 1.5rem; border-top: 1px solid var(--md-sys-color-outline-variant); }

  .field-row { display: grid; grid-template-columns: 1fr 1fr; gap: 1.5rem; }
  .field { display: flex; flex-direction: column; gap: 0.5rem; }
  .field label { font-size: 0.8rem; font-weight: 800; color: var(--text-color); }
  .field input { background: var(--surface-variant); border: 1px solid var(--md-sys-color-outline-variant); padding: 0.85rem 1.25rem; border-radius: 0.75rem; color: var(--text-color); }

  .premium-switch { display: flex; align-items: center; gap: 1rem; cursor: pointer; }
  .premium-switch input { display: none; }
  .switch-slider { width: 44px; height: 24px; background: var(--surface-variant); border-radius: 0.5rem; position: relative; border: 1px solid var(--md-sys-color-outline-variant); transition: 0.3s; }
  .switch-slider::before { content: ""; position: absolute; width: 16px; height: 16px; background: var(--text-muted); border-radius: 0.35rem; left: 3px; top: 3px; transition: 0.3s; }
  input:checked + .switch-slider { background: var(--primary-color); }
  input:checked + .switch-slider::before { transform: translateX(20px); background: white; }
  .switch-label { font-size: 0.95rem; font-weight: 700; color: var(--text-color); }

  .skeleton-config { height: 600px; background: var(--surface-variant); border-radius: 2rem; animation: pulse 1.5s infinite; }
</style>
