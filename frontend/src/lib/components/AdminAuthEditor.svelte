<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from 'svelte';
  
  let { authId = null, onSave, onCancel } = $props<{
    authId?: number | null;
    onSave: () => void;
    onCancel: () => void;
  }>();

  let name = $state('');
  let type = $state(2); // Default to LDAP
  let activated = $state(true);
  let isDefault = $state(false);
  
  // Specific configs
  let config = $state<any>({});
  
  let loading = $state(false);
  let submitting = $state(false);
  let error = $state('');

  const AUTH_TYPES = [
    { id: 2, name: 'LDAP (via BindDN)' },
    { id: 5, name: 'LDAP (simple auth)' },
    { id: 3, name: 'SMTP' },
    { id: 4, name: 'PAM' },
    { id: 6, name: 'GitHub OAuth' }
  ];

  function initConfig(t: number) {
    if (t === 2 || t === 5) {
      config = {
        security_protocol: 0,
        host: '',
        port: 389,
        bind_dn: '',
        bind_password: '',
        user_base: '',
        user_dn: '',
        filter: '',
        admin_filter: '',
        attribute_username: '',
        attribute_name: '',
        attribute_surname: '',
        attribute_mail: '',
        group_enabled: false,
        group_dn: '',
        group_filter: '',
        group_member_uid: '',
        user_uid: ''
      };
    } else if (t === 3) {
      config = {
        smtp_auth: 'PLAIN',
        smtp_host: '',
        smtp_port: 587,
        allowed_domains: '',
        tls: true,
        skip_verify: false
      };
    } else if (t === 4) {
      config = { pam_service_name: '' };
    } else if (t === 6) {
      config = { github_api_endpoint: 'https://api.github.com/' };
    }
  }

  $effect(() => {
    if (!authId) {
      initConfig(type);
    }
  });

  async function fetchAuth() {
    if (!authId) return;
    loading = true;
    try {
      const res = await apiFetch(`/api/v2/admin/auths/${authId}`);
      if (res.ok) {
        const data = await res.json();
        name = data.name;
        type = data.type;
        activated = data.is_activated;
        isDefault = data.is_default;
        config = data.config || {};
      }
    } catch (e) {
      error = 'Failed to load auth source.';
    } finally {
      loading = false;
    }
  }

  async function save() {
    submitting = true;
    error = '';
    const method = authId ? 'PATCH' : 'POST';
    const url = authId ? `/api/v2/admin/auths/${authId}` : '/api/v2/admin/auths';
    
    const payload = {
       name,
       type,
       activated,
       default: isDefault,
       config
    };

    try {
      const res = await apiFetch(url, {
        method,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      });
      if (res.ok) {
        onSave();
      } else {
        const data = await res.json();
        error = data.error || 'Failed to save authentication source.';
      }
    } catch (e) {
      error = 'A network error occurred.';
    } finally {
      submitting = false;
    }
  }

  onMount(fetchAuth);
</script>

<div class="auth-editor island island-interactive tonal">
  <header>
    <h3>{authId ? 'Edit Auth Source' : 'New Authentication Source'}</h3>
    <p>Configure how Ggit delegates user authentication to external systems.</p>
  </header>

  {#if loading}
    <div class="spinner-container"><div class="spinner"></div></div>
  {:else}
    <form onsubmit={(e) => { e.preventDefault(); save(); }}>
      <div class="form-grid">
        <div class="field">
          <label for="auth_name">Core Name</label>
          <input id="auth_name" type="text" bind:value={name} placeholder="e.g. Corporate LDAP" required />
        </div>
        <div class="field">
          <label for="auth_type">Source Type</label>
          <select id="auth_type" bind:value={type} disabled={!!authId}>
            {#each AUTH_TYPES as t}
              <option value={t.id}>{t.name}</option>
            {/each}
          </select>
        </div>
      </div>

      <div class="checkbox-group">
        <label class="checkbox">
          <input type="checkbox" bind:checked={activated} />
          <span>Activated</span>
        </label>
        <label class="checkbox">
          <input type="checkbox" bind:checked={isDefault} />
          <span>Default Auth Source</span>
        </label>
      </div>

      <div class="divider"></div>

      <!-- Provider Specific Config -->
      <section class="config-section">
        {#if type === 2 || type === 5}
          <div class="form-grid">
             <div class="field">
               <label>Protocol</label>
               <select bind:value={config.security_protocol}>
                 <option value={0}>Unencrypted</option>
                 <option value={1}>LDAPS (SSL)</option>
                 <option value={2}>STARTTLS</option>
               </select>
             </div>
             <div class="field">
               <label>Host</label>
               <input type="text" bind:value={config.host} placeholder="ldap.example.com" required />
             </div>
             <div class="field">
               <label>Port</label>
               <input type="number" bind:value={config.port} required />
             </div>
             {#if type === 2}
               <div class="field full">
                 <label>Bind DN</label>
                 <input type="text" bind:value={config.bind_dn} placeholder="cn=admin,dc=example,dc=com" />
               </div>
               <div class="field full">
                 <label>Bind Password</label>
                 <input type="password" bind:value={config.bind_password} />
               </div>
               <div class="field full">
                 <label>User Search Base</label>
                 <input type="text" bind:value={config.user_base} placeholder="ou=users,dc=example,dc=com" required />
               </div>
             {:else}
               <div class="field full">
                 <label>User DN Filter</label>
                 <input type="text" bind:value={config.user_dn} placeholder="uid=%s,ou=users,dc=example,dc=com" required />
               </div>
             {/if}
             <div class="field">
               <label>Search Filter</label>
               <input type="text" bind:value={config.filter} placeholder="(&(objectClass=posixAccount)(uid=%s))" required />
             </div>
             <div class="field">
               <label>Admin Filter</label>
               <input type="text" bind:value={config.admin_filter} />
             </div>
             <div class="field">
               <label>Mail Attribute</label>
               <input type="text" bind:value={config.attribute_mail} placeholder="mail" required />
             </div>
          </div>
        {:else if type === 3}
          <div class="form-grid">
            <div class="field">
              <label>SMTP Auth</label>
              <select bind:value={config.smtp_auth}>
                <option value="PLAIN">PLAIN</option>
                <option value="LOGIN">LOGIN</option>
                <option value="CRAM-MD5">CRAM-MD5</option>
              </select>
            </div>
            <div class="field">
              <label>SMTP Host</label>
              <input type="text" bind:value={config.smtp_host} required />
            </div>
            <div class="field">
              <label>SMTP Port</label>
              <input type="number" bind:value={config.smtp_port} required />
            </div>
            <div class="field full">
              <label>Allowed Domains</label>
              <input type="text" bind:value={config.allowed_domains} placeholder="example.com, company.local" />
            </div>
            <div class="checkbox-group">
               <label class="checkbox">
                 <input type="checkbox" bind:checked={config.tls} />
                 <span>Enable TLS</span>
               </label>
               <label class="checkbox">
                 <input type="checkbox" bind:checked={config.skip_verify} />
                 <span>Skip TLS Verify</span>
               </label>
            </div>
          </div>
        {:else if type === 4}
           <div class="field">
             <label>PAM Service Name</label>
             <input type="text" bind:value={config.pam_service_name} required />
           </div>
        {:else if type === 6}
           <div class="field">
             <label>GitHub API Endpoint</label>
             <input type="url" bind:value={config.github_api_endpoint} required />
             <small>Standard: https://api.github.com/</small>
           </div>
        {/if}
      </section>

      {#if error}
        <div class="error-msg">{error}</div>
      {/if}

      <div class="actions">
        <button type="button" class="btn btn-secondary" onclick={onCancel}>Cancel</button>
        <button type="submit" class="btn btn-primary" disabled={submitting}>
          {submitting ? 'Connecting...' : 'Save Source'}
        </button>
      </div>
    </form>
  {/if}
</div>

<style>
  .auth-editor { display: flex; flex-direction: column; gap: 1.5rem; max-width: 800px; padding: 2rem; }
  h3 { margin: 0; font-size: 1.5rem; font-weight: 800; }
  p { margin: 0; color: var(--text-muted); font-size: 0.95rem; }

  form { display: flex; flex-direction: column; gap: 1.5rem; }
  .form-grid { display: grid; grid-template-columns: 1fr 1fr; gap: 1.5rem; }
  .field.full { grid-column: span 2; }
  
  .field { display: flex; flex-direction: column; gap: 0.5rem; }
  .field label { font-size: 0.75rem; font-weight: 800; text-transform: uppercase; color: var(--text-muted); letter-spacing: 0.05em; }
  .field input, .field select { background: rgba(0,0,0,0.3); border: 1px solid var(--md-sys-color-outline-variant); color: #fff; padding: 0.75rem 1rem; border-radius: 0.5rem; font-size: 0.95rem; transition: all 0.2s; }
  .field input:focus { border-color: var(--primary-color); background: rgba(0,0,0,0.4); outline: none; }
  
  .checkbox-group { display: flex; gap: 2rem; margin: 0.5rem 0; }
  .checkbox { display: flex; align-items: center; gap: 0.75rem; cursor: pointer; font-size: 0.95rem; font-weight: 500; }
  .checkbox input { width: 1.25rem; height: 1.25rem; accent-color: var(--primary-color); }

  .divider { height: 1px; background: var(--md-sys-color-outline-variant); margin: 0.5rem 0; opacity: 0.5; }

  .config-section { background: rgba(255,255,255,0.02); padding: 1.5rem; border-radius: 1rem; border: 1px dashed var(--md-sys-color-outline-variant); }

  .error-msg { color: #ef4444; font-size: 0.9rem; font-weight: 700; background: rgba(239, 68, 68, 0.1); padding: 0.75rem 1rem; border-radius: 0.5rem; border-left: 4px solid #ef4444; }
  
  .actions { display: flex; justify-content: flex-end; gap: 1rem; margin-top: 1rem; }
  
  .spinner-container { display: flex; justify-content: center; padding: 3rem; }
  small { margin-top: 0.25rem; font-size: 0.8rem; color: var(--text-muted); }
</style>
