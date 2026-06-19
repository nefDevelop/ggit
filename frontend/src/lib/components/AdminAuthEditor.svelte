<script lang="ts">
  import { apiFetch } from "$lib/api";
  import Button from './ui/Button.svelte';
  import Input from './ui/Input.svelte';
  import Select from './ui/Select.svelte';
  import Checkbox from './ui/Checkbox.svelte';
  import Card from './ui/Card.svelte';
  
import { ArrowLeft, TriangleAlert } from "lucide-svelte";

  
  let { 
    source = null, 
    onSave = () => {}, 
    onCancel = () => {} 
  } = $props<{
    source?: any;
    onSave?: () => void;
    onCancel?: () => void;
  }>();

  let isNew = $derived(!source);
  let loading = $state(false);
  let error = $state('');

  // Form State
  let name = $state('');
  let type = $state(2); // Default LDAP
  let activated = $state(true);
  let isDefault = $state(false);

  // Config State (LDAP)
  let config = $state({
    Host: '',
    Port: 389,
    SecurityProtocol: 0,
    SkipVerify: false,
    BindDN: '',
    BindPassword: '',
    UserBase: '',
    UserDN: '',
    AttributeUsername: 'uid',
    AttributeName: 'givenName',
    AttributeSurname: 'sn',
    AttributeMail: 'mail',
    AttributesInBind: false,
    Filter: '(objectClass=inetOrgPerson)',
    AdminFilter: '',
    GroupEnabled: false,
    GroupDN: '',
    GroupFilter: '',
    GroupMemberUID: '',
    UserUID: ''
  });

  const typeOptions = [
    { value: 2, label: 'LDAP (via BindDN)' },
    { value: 5, label: 'LDAP (Simple Auth)' },
    { value: 6, label: 'OpenID Connect' }
  ];

  // Config State (OIDC) - We use the same 'config' object but with different fields
  let oidcConfig = $state({
    IssuerURL: '',
    ClientID: '',
    ClientSecret: '',
    Scopes: ['openid', 'profile', 'email'],
    RedirectURL: '',
    AttributeUsername: '',
    AttributeName: '',
    AttributeEmail: '',
    SkipVerify: false
  });

  $effect(() => {
    name = source?.name || '';
    type = source?.type || 2;
    activated = source?.is_actived ?? true;
    isDefault = source?.is_default ?? false;
    if (source?.RawConfig) {
        try { 
            const parsed = JSON.parse(source.RawConfig);
            if (type === 6) oidcConfig = { ...oidcConfig, ...parsed };
            else config = { ...config, ...parsed };
        } catch(e) {}
    }
  });

  const protocolOptions = [
    { value: 0, label: 'Unencrypted' },
    { value: 1, label: 'LDAPS' },
    { value: 2, label: 'StartTLS' }
  ];

  async function save() {
    if (!name) { error = 'El nombre es obligatorio'; return; }
    loading = true;
    error = '';
    
    try {
      const payload = {
        name,
        type: Number(type),
        activated,
        default: isDefault,
        config: type === 6 ? oidcConfig : config
      };

      const url = isNew ? '/api/v2/admin/auths' : `/api/v2/admin/auths/${source.id}`;
      const method = isNew ? 'POST' : 'PATCH';

      const res = await apiFetch(url, {
        method,
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(payload)
      });

      if (res.ok) {
        onSave();
      } else {
        const body = await res.json();
        error = body.message || 'Error al guardar la fuente de auth';
      }
    } catch (e) {
      error = 'Error de conexión';
    } finally {
      loading = false;
    }
  }
</script>

<div class="auth-editor">
  <div class="editor-header">
    <Button variant="ghost" icon={ArrowLeft} onclick={onCancel}>Volver</Button>
    <h2>{isNew ? 'Nueva Fuente de Autenticación' : `Editar: ${name}`}</h2>
  </div>

  {#if error}
    <div class="error-banner">
      <TriangleAlert size={20} />
      {error}
    </div>
  {/if}

  <div class="editor-grid">
    <section class="config-section">
      <h3>General</h3>
      <Card tonal padding="lg" class="form-card">
        <div class="form-row">
          <Input bind:value={name} label="Nombre del Origen" placeholder="Ej: LDAP Corporativo" />
          <Select bind:value={type} options={typeOptions} label="Tipo de Autenticación" />
        </div>
        <div class="form-row checks">
          <Checkbox bind:checked={activated} label="Activado" />
          <Checkbox bind:checked={isDefault} label="Origen por defecto" />
        </div>
      </Card>
    </section>

    {#if type === 6}
    <section class="config-section">
      <h3>Configuración OIDC</h3>
      <Card tonal padding="lg" class="form-card">
        <div class="form-row">
          <Input bind:value={oidcConfig.IssuerURL} label="URL del Emisor (Discovery)" placeholder="https://accounts.google.com" />
        </div>
        <div class="form-row">
          <Input bind:value={oidcConfig.ClientID} label="Client ID" />
          <Input bind:value={oidcConfig.ClientSecret} label="Client Secret" type="password" />
        </div>
        <div class="form-row">
          <Input bind:value={oidcConfig.RedirectURL} label="URL de Redirección (Callback)" placeholder="https://ggit.example.com/api/v2/user/login/oidc/CALLBACK" />
        </div>
        <div class="form-row">
          <div class="check-wrapper">
             <Checkbox bind:checked={oidcConfig.SkipVerify} label="Omitir verificación TLS" />
          </div>
        </div>
      </Card>
    </section>

    <section class="config-section">
      <h3>Mapeo de Atributos (Claims)</h3>
      <Card tonal padding="lg" class="form-card">
        <p class="section-hint">Claims para extraer información del usuario del ID Token</p>
        <div class="attribute-grid">
          <Input bind:value={oidcConfig.AttributeUsername} label="Claim Usuario" placeholder="preferred_username" />
          <Input bind:value={oidcConfig.AttributeName} label="Claim Nombre" placeholder="name" />
          <Input bind:value={oidcConfig.AttributeEmail} label="Claim Email" placeholder="email" />
        </div>
      </Card>
    </section>
    {:else}
    <section class="config-section">
      <h3>Conexión</h3>
      <Card tonal padding="lg" class="form-card">
        <div class="form-row">
          <Input bind:value={config.Host} label="Host" placeholder="ldap.ejemplo.com" />
          <Input type="number" bind:value={config.Port} label="Puerto" placeholder="389" />
        </div>
        <div class="form-row">
          <Select bind:value={config.SecurityProtocol} options={protocolOptions} label="Protocolo de Seguridad" />
          <div class="check-wrapper">
             <Checkbox bind:checked={config.SkipVerify} label="Omitir verificación TLS" />
          </div>
        </div>
      </Card>
    </section>

    <section class="config-section">
      <h3>Autenticación (Bind)</h3>
      <Card tonal padding="lg" class="form-card">
        <p class="section-hint">Deja vacío para búsqueda anónima (si el servidor lo permite)</p>
        <div class="form-row">
          <Input bind:value={config.BindDN} label="Bind DN" placeholder="cn=admin,dc=ejemplo,dc=com" />
          <Input type="password" bind:value={config.BindPassword} label="Bind Password" />
        </div>
        {#if type === 5}
          <Input bind:value={config.UserDN} label="User DN Template" placeholder="uid=%s,ou=Users,dc=ejemplo,dc=com" />
        {/if}
      </Card>
    </section>

    <section class="config-section">
      <h3>Búsqueda de Usuarios</h3>
      <Card tonal padding="lg" class="form-card">
        <div class="form-row">
          <Input bind:value={config.UserBase} label="User Search Base" placeholder="ou=Users,dc=ejemplo,dc=com" />
          <Input bind:value={config.Filter} label="User Filter" placeholder="(objectClass=inetOrgPerson)" />
        </div>
        <div class="form-row">
          <Input bind:value={config.AdminFilter} label="Admin Filter" placeholder="(memberOf=cn=admins,ou=Groups,dc=ejemplo,dc=com)" />
        </div>
        <div class="attribute-grid">
          <Input bind:value={config.AttributeUsername} label="Atributo Usuario" />
          <Input bind:value={config.AttributeName} label="Atributo Nombre" />
          <Input bind:value={config.AttributeSurname} label="Atributo Apellido" />
          <Input bind:value={config.AttributeMail} label="Atributo Email" />
        </div>
        <Checkbox bind:checked={config.AttributesInBind} label="Buscar atributos en contexto de Bind" />
      </Card>
    </section>

    <section class="config-section">
      <div class="section-title-row">
        <h3>Grupos</h3>
        <Checkbox bind:checked={config.GroupEnabled} label="Habilitar verificación de grupos" />
      </div>
      {#if config.GroupEnabled}
        <Card tonal padding="lg" class="form-card anim-in">
          <div class="form-row">
            <Input bind:value={config.GroupDN} label="Group Search Base" />
            <Input bind:value={config.GroupFilter} label="Group Filter" />
          </div>
          <div class="form-row">
            <Input bind:value={config.GroupMemberUID} label="Group Member Attribute" />
            <Input bind:value={config.UserUID} label="User UID Attribute" />
          </div>
        </Card>
      {/if}
    </section>
    {/if}
  </div>

  <div class="editor-footer">
    <Button variant="ghost" onclick={onCancel}>Cancelar</Button>
    <Button variant="primary" loading={loading} onclick={save}>
      {isNew ? 'Crear Origen' : 'Guardar Cambios'}
    </Button>
  </div>
</div>

<style>
  .auth-editor {
    animation: slideIn 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .editor-header {
    display: flex;
    align-items: center;
    gap: 2rem;
    margin-bottom: 3rem;
  }
  .editor-header h2 {
    margin: 0;
    font-size: 2rem;
    font-weight: 900;
  }

  .editor-grid {
    display: flex;
    flex-direction: column;
    gap: 3rem;
    max-width: 1000px;
  }

  .config-section h3 {
    margin: 0 0 1rem 0.5rem;
    font-size: 0.9rem;
    text-transform: uppercase;
    letter-spacing: 0.1em;
    color: var(--primary-color);
    font-weight: 900;
  }

  .section-title-row {
    display: flex;
    align-items: center;
    gap: 2rem;
    margin-bottom: 1rem;
  }
  .section-title-row h3 { margin-bottom: 0; }

  :global(.form-card) {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  .form-row {
    display: flex;
    gap: 1.5rem;
  }
  :global(.form-row > *) { flex: 1; }

  .checks { align-items: center; gap: 3rem; margin-top: 0.5rem; }
  .check-wrapper { display: flex; align-items: center; padding-top: 1rem; }

  .section-hint {
    margin: -0.5rem 0 0.5rem 0;
    font-size: 0.85rem;
    color: var(--text-muted);
    font-style: italic;
  }

  .attribute-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 1rem;
    background: rgba(0,0,0,0.2);
    padding: 1.5rem;
    border-radius: 0.75rem;
  }

  .error-banner {
    background: rgba(var(--error-rgb), 0.1);
    color: var(--error);
    padding: 1.25rem;
    border-radius: 0.75rem;
    display: flex;
    align-items: center;
    gap: 1rem;
    margin-bottom: 2rem;
    font-weight: 700;
    border: 1px solid rgba(var(--error-rgb), 0.2);
  }

  .editor-footer {
    margin-top: 4rem;
    padding-top: 2rem;
    border-top: 1px solid var(--md-sys-color-outline-variant);
    display: flex;
    justify-content: flex-end;
    gap: 1rem;
  }

  :global(.anim-in) {
    animation: expand 0.3s ease-out;
  }

  @keyframes slideIn {
    from { opacity: 0; transform: translateX(30px); }
    to { opacity: 1; transform: translateX(0); }
  }
  @keyframes expand {
    from { opacity: 0; transform: translateY(-10px); }
    to { opacity: 1; transform: translateY(0); }
  }

  @media (max-width: 768px) {
    .form-row { flex-direction: column; }
  }
</style>
