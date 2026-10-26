<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from "svelte";
  import Card from "./ui/Card.svelte";
  import Button from "./ui/Button.svelte";
  import PageHeader from "./ui/PageHeader.svelte";

  import Rocket from "lucide-svelte/icons/rocket";
  import RefreshCw from "lucide-svelte/icons/refresh-cw";
  import CheckCircle2 from "lucide-svelte/icons/check-circle-2";
  import AlertTriangle from "lucide-svelte/icons/alert-triangle";
  import CloudOff from "lucide-svelte/icons/cloud-off";
  import ShieldAlert from "lucide-svelte/icons/shield-alert";
  import Info from "lucide-svelte/icons/info";
  import ExternalLink from "lucide-svelte/icons/external-link";

  type Health = {
    status: string;
    db: string;
    version: string;
    error?: string;
  };

  type PikoCIStatus = {
    enabled: boolean;
    mode: string;
    url: string;
    public_url?: string;
    reachable?: boolean;
    health?: Health;
    error?: string;
  };

  let loading = $state(true);
  let status = $state<PikoCIStatus | null>(null);
  let loadError = $state("");

  async function fetchStatus() {
    loading = true;
    loadError = "";
    try {
      const res = await apiFetch("/api/v2/admin/pikoci/status", { silent: true });
      if (res.ok) {
        status = res.data as PikoCIStatus;
      } else {
        loadError = `El servidor respondió ${res.status}`;
      }
    } catch (e) {
      console.error(e);
      loadError = "No se pudo consultar el estado de PikoCI";
    } finally {
      loading = false;
    }
  }

  onMount(fetchStatus);

  const summary = $derived.by(() => {
    if (!status) return { label: "Desconocido", tone: "muted" };
    if (!status.enabled) return { label: "Deshabilitado", tone: "muted" };
    if (status.reachable && status.health?.db === "connected") {
      return { label: "En marcha", tone: "ok" };
    }
    if (status.reachable) return { label: "Base de datos caída", tone: "warn" };
    return { label: "Sin conexión", tone: "error" };
  });
</script>

<div class="admin-pikoci">
  <PageHeader
    title="CI/CD con PikoCI"
    subtitle="Estado del motor de integración continua y cómo llega ggit a él"
    icon={Rocket}
  >
    {#snippet actions()}
      <Button variant="tonal" icon={RefreshCw} onclick={fetchStatus} loading={loading} disabled={loading}>
        Actualizar
      </Button>
    {/snippet}
  </PageHeader>

  <div class="pikoci-body">
    {#if loading}
      <div class="skeleton-card"></div>
    {:else if loadError}
      <Card class="state-card">
        <AlertTriangle size={40} />
        <p>{loadError}</p>
      </Card>
    {:else if status}
      <Card class="status-card">
        <div class="status-head">
          <div class="service">
            <div class="service-icon" class:ok={summary.tone === "ok"} class:warn={summary.tone === "warn"}>
              <Rocket size={22} />
            </div>
            <div class="service-text">
              <h3>PikoCI</h3>
              <span class="mode">{status.mode}</span>
            </div>
          </div>

          <span
            class="badge"
            class:ok={summary.tone === "ok"}
            class:warn={summary.tone === "warn"}
            class:error={summary.tone === "error"}
          >
            {#if summary.tone === "ok"}
              <CheckCircle2 size={16} />
            {:else if summary.tone === "warn"}
              <AlertTriangle size={16} />
            {:else if summary.tone === "error"}
              <CloudOff size={16} />
            {:else}
              <Info size={16} />
            {/if}
            {summary.label}
          </span>
        </div>

        {#if status.mode === "subprocess"}
          <div class="notice warning">
            <ShieldAlert size={18} />
            <p>
              Modo subproceso: los pipelines se ejecutan con el usuario de ggit y pueden leer su
              configuración y su base de datos. En producción usa un contenedor o un usuario del
              sistema separado.
            </p>
          </div>
        {/if}

        {#if status.enabled && status.error}
          <div class="notice error">
            <CloudOff size={18} />
            <p>{status.error}</p>
          </div>
        {/if}

        {#if status.health?.error}
          <div class="notice error">
            <AlertTriangle size={18} />
            <p>{status.health.error}</p>
          </div>
        {/if}

        <div class="details">
          <div class="detail">
            <span class="label">API</span>
            <code>{status.url}</code>
          </div>
          <div class="detail">
            <span class="label">Versión</span>
            <span class="value">{status.health?.version || "—"}</span>
          </div>
          <div class="detail">
            <span class="label">Base de datos</span>
            <span class="value">{status.health?.db || "—"}</span>
          </div>
          <div class="detail">
            <span class="label">Interfaz</span>
            {#if status.public_url}
              <a class="link" href={status.public_url} target="_blank" rel="noopener noreferrer">
                Abrir PikoCI
                <ExternalLink size={14} />
              </a>
            {:else}
              <span class="value">—</span>
            {/if}
          </div>
        </div>
      </Card>

      {#if !status.enabled}
        <Card class="hint-card">
          <Info size={20} />
          <div>
            <strong>CI deshabilitado</strong>
            <p>
              Actívalo con <code>ENABLED = true</code> en la sección <code>[pikoci]</code> de
              <code>app.ini</code> y reinicia ggit.
            </p>
          </div>
        </Card>
      {/if}
    {/if}
  </div>
</div>

<style>
  .pikoci-body {
    padding: 1.5rem 1.5rem 4rem;
    animation: fadeIn 0.4s cubic-bezier(0.4, 0, 0.2, 1);
  }

  @keyframes fadeIn {
    from { opacity: 0; transform: translateY(4px); }
    to { opacity: 1; transform: translateY(0); }
  }

  :global(.status-card) {
    display: flex;
    flex-direction: column;
    gap: 1.5rem !important;
    padding: 2rem !important;
    border-radius: 1.5rem !important;
  }

  .status-head {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 1rem;
    flex-wrap: wrap;
  }

  .service {
    display: flex;
    align-items: center;
    gap: 1rem;
  }

  .service-icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 3rem;
    height: 3rem;
    border-radius: 0.85rem;
    background: var(--md-sys-color-surface-container-high);
    color: var(--md-sys-color-on-surface-variant);
  }

  .service-icon.ok {
    background: color-mix(in srgb, var(--md-sys-color-primary), transparent 85%);
    color: var(--md-sys-color-primary);
  }

  .service-icon.warn {
    background: color-mix(in srgb, var(--md-sys-color-error), transparent 88%);
    color: var(--md-sys-color-error);
  }

  .service-text h3 {
    margin: 0;
    font-size: 1.25rem;
    font-weight: 900;
    letter-spacing: -0.01em;
  }

  .mode {
    font-size: 0.8rem;
    text-transform: uppercase;
    letter-spacing: 0.08em;
    color: var(--md-sys-color-on-surface-variant);
    font-weight: 700;
  }

  .badge {
    display: inline-flex;
    align-items: center;
    gap: 0.45rem;
    padding: 0.5rem 0.9rem;
    border-radius: 999px;
    font-size: 0.85rem;
    font-weight: 800;
    background: var(--md-sys-color-surface-container-high);
    color: var(--md-sys-color-on-surface-variant);
  }

  .badge.ok {
    background: color-mix(in srgb, var(--md-sys-color-primary), transparent 85%);
    color: var(--md-sys-color-primary);
  }

  .badge.warn {
    background: color-mix(in srgb, var(--md-sys-color-error), transparent 88%);
    color: var(--md-sys-color-error);
  }

  .badge.error {
    background: color-mix(in srgb, var(--md-sys-color-error), transparent 88%);
    color: var(--md-sys-color-error);
  }

  .notice {
    display: flex;
    gap: 0.75rem;
    align-items: flex-start;
    padding: 1rem 1.25rem;
    border-radius: 0.75rem;
    border: 1px solid var(--md-sys-color-outline-variant);
    background: var(--md-sys-color-surface-container-low);
    font-size: 0.9rem;
    line-height: 1.5;
  }

  .notice p {
    margin: 0;
  }

  .notice.warning {
    border-color: color-mix(in srgb, var(--md-sys-color-error), transparent 70%);
    background: color-mix(in srgb, var(--md-sys-color-error), transparent 94%);
    color: var(--md-sys-color-error);
  }

  .notice.error {
    border-color: color-mix(in srgb, var(--md-sys-color-error), transparent 60%);
    background: color-mix(in srgb, var(--md-sys-color-error), transparent 92%);
    color: var(--md-sys-color-error);
  }

  .details {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(180px, 1fr));
    gap: 1.25rem;
    padding: 1.5rem;
    background: var(--md-sys-color-surface-container-low);
    border-radius: 0.75rem;
    border: 1px solid var(--md-sys-color-outline-variant);
  }

  .detail {
    display: flex;
    flex-direction: column;
    gap: 0.4rem;
    min-width: 0;
  }

  .label {
    font-size: 0.75rem;
    text-transform: uppercase;
    letter-spacing: 0.08em;
    color: var(--md-sys-color-on-surface-variant);
    font-weight: 800;
    opacity: 0.8;
  }

  .value,
  .link {
    font-size: 1rem;
    font-weight: 700;
    color: var(--md-sys-color-on-surface);
  }

  code {
    font-size: 0.9rem;
    font-weight: 600;
    color: var(--md-sys-color-on-surface);
    overflow-wrap: anywhere;
  }

  .link {
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    color: var(--md-sys-color-primary);
    text-decoration: none;
  }

  .link:hover {
    text-decoration: underline;
  }

  :global(.hint-card) {
    margin-top: 1.5rem;
    display: flex;
    gap: 1rem;
    align-items: flex-start;
    padding: 1.5rem !important;
    border-radius: 1rem !important;
  }

  :global(.hint-card strong) {
    display: block;
    margin-bottom: 0.25rem;
  }

  :global(.hint-card p) {
    margin: 0;
    color: var(--md-sys-color-on-surface-variant);
    font-size: 0.9rem;
    line-height: 1.5;
  }

  :global(.state-card) {
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 1rem;
    padding: 4rem !important;
    text-align: center;
    color: var(--md-sys-color-on-surface-variant);
  }

  .skeleton-card {
    height: 280px;
    border-radius: 1.5rem;
    background: var(--md-sys-color-surface-container-high);
    animation: pulse 1.5s infinite;
  }

  @keyframes pulse {
    0% { opacity: 1; }
    50% { opacity: 0.5; }
    100% { opacity: 1; }
  }
</style>
