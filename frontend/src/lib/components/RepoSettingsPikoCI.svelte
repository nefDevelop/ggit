<script lang="ts">
  import type { RepoContext } from "$lib/types";
  import { apiFetch } from "$lib/api";
  import { onMount, getContext } from "svelte";
  import { toast } from "$lib/toast.svelte";
  import { dialog } from "$lib/dialog.svelte";
  import Card from "./ui/Card.svelte";
  import Button from "./ui/Button.svelte";
  import PageHeader from "./ui/PageHeader.svelte";

  import Rocket from "lucide-svelte/icons/rocket";
  import RefreshCw from "lucide-svelte/icons/refresh-cw";
  import Trash2 from "lucide-svelte/icons/trash-2";
  import CheckCircle2 from "lucide-svelte/icons/check-circle-2";
  import ShieldOff from "lucide-svelte/icons/shield-off";
  import ExternalLink from "lucide-svelte/icons/external-link";

  const repo = getContext<RepoContext>("repo");

  type State = {
    enabled: boolean;
    linked: boolean;
    pipeline?: string;
    public_url?: string;
  };

  let loading = $state(true);
  let busy = $state(false);
  let ci = $state<State>({ enabled: false, linked: false });

  const endpoint = $derived(`/api/v2/${repo.data.owner.username}/${repo.data.name}/settings/pikoci`);

  async function fetchState() {
    loading = true;
    try {
      const res = await apiFetch(endpoint);
      if (res.ok) {
        ci = res.data as State;
      }
    } catch (e) {
      console.error(e);
    } finally {
      loading = false;
    }
  }

  async function enable() {
    busy = true;
    try {
      const res = await apiFetch(endpoint, { method: "POST" });
      if (res.ok) {
        toast.success("CI activado y pipeline publicado");
        await fetchState();
      }
    } catch (e) {
      console.error(e);
    } finally {
      busy = false;
    }
  }

  async function disable() {
    if (!(await dialog.confirm("Se eliminará el pipeline en PikoCI y el vínculo con este repositorio. ¿Continuar?"))) {
      return;
    }

    busy = true;
    try {
      const res = await apiFetch(endpoint, { method: "DELETE" });
      if (res.ok) {
        toast.success("CI desactivado");
        await fetchState();
      }
    } catch (e) {
      console.error(e);
    } finally {
      busy = false;
    }
  }

  onMount(fetchState);
</script>

<div class="pikoci-settings">
  <PageHeader
    title="Integración continua"
    subtitle="PikoCI construye el repositorio y reporta el estado en los commits"
    icon={Rocket}
  >
    {#snippet actions()}
      <Button variant="tonal" icon={RefreshCw} onclick={fetchState} loading={loading} disabled={loading || busy}>
        Actualizar
      </Button>
    {/snippet}
  </PageHeader>

  <Card class="pikoci-card">
    {#if loading}
      <p class="muted">Cargando estado…</p>
    {:else}
      <div class="row">
        <div class="service">
          <div class="icon" class:on={ci.linked}>
            {#if ci.linked}
              <CheckCircle2 size={22} />
            {:else}
              <ShieldOff size={22} />
            {/if}
          </div>
          <div>
            <h3>{ci.linked ? "CI activo" : "CI desactivado"}</h3>
            <p class="muted">
              {#if ci.linked}
                Pipeline <code>{ci.pipeline}</code> publicado en PikoCI.
              {:else}
                El repositorio todavía no publica su pipeline.
              {/if}
            </p>
          </div>
        </div>

        <div class="actions">
          {#if !ci.enabled}
            <span class="muted">PikoCI está deshabilitado en este servidor.</span>
          {:else if ci.linked}
            <Button variant="outline" icon={RefreshCw} onclick={enable} loading={busy} disabled={busy}>
              Republicar
            </Button>
            <Button variant="danger" icon={Trash2} onclick={disable} loading={busy} disabled={busy}>
              Desactivar
            </Button>
          {:else}
            <Button variant="primary" icon={Rocket} onclick={enable} loading={busy} disabled={busy}>
              Activar CI
            </Button>
          {/if}
        </div>
      </div>

      {#if ci.public_url}
        <a class="link" href={ci.public_url} target="_blank" rel="noopener noreferrer">
          Abrir PikoCI
          <ExternalLink size={14} />
        </a>
      {/if}

      <p class="hint">
        El pipeline se define en <code>.ggit/ci.hcl</code> de la rama por defecto. Al cambiarlo en un push,
        ggit lo republica automáticamente.
      </p>
    {/if}
  </Card>
</div>

<style>
  .pikoci-settings {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }

  :global(.pikoci-card) {
    padding: 1.5rem !important;
    border-radius: 1rem !important;
    display: flex;
    flex-direction: column;
    gap: 1.25rem;
  }

  .row {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 1.5rem;
    flex-wrap: wrap;
  }

  .service {
    display: flex;
    gap: 1rem;
    align-items: flex-start;
  }

  .icon {
    display: flex;
    align-items: center;
    justify-content: center;
    width: 2.75rem;
    height: 2.75rem;
    border-radius: 0.75rem;
    background: var(--md-sys-color-surface-container-high);
    color: var(--md-sys-color-on-surface-variant);
  }

  .icon.on {
    background: color-mix(in srgb, var(--md-sys-color-primary), transparent 85%);
    color: var(--md-sys-color-primary);
  }

  h3 {
    margin: 0 0 0.25rem;
    font-size: 1.1rem;
    font-weight: 800;
  }

  .muted {
    margin: 0;
    color: var(--md-sys-color-on-surface-variant);
    font-size: 0.9rem;
    line-height: 1.5;
  }

  .actions {
    display: flex;
    gap: 0.75rem;
    align-items: center;
    flex-wrap: wrap;
  }

  .link {
    display: inline-flex;
    align-items: center;
    gap: 0.35rem;
    color: var(--md-sys-color-primary);
    font-weight: 700;
    text-decoration: none;
  }

  .link:hover {
    text-decoration: underline;
  }

  .hint {
    margin: 0;
    padding-top: 1rem;
    border-top: 1px solid var(--md-sys-color-outline-variant);
    color: var(--md-sys-color-on-surface-variant);
    font-size: 0.85rem;
    line-height: 1.5;
  }

  code {
    font-size: 0.85em;
  }
</style>
