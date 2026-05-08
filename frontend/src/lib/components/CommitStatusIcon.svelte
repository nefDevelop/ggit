<script lang="ts">
  import { CheckCircle2, XCircle, AlertCircle, HelpCircle } from "lucide-svelte";

  let { status, size = 16, class: className = "" } = $props<{ 
    status: string | null | undefined, 
    size?: number,
    class?: string 
  }>();

  function getStatusInfo(s: string | null | undefined) {
    switch (s) {
      case 'success':
        return { icon: CheckCircle2, color: 'var(--success)', label: 'Success' };
      case 'failure':
      case 'error':
        return { icon: XCircle, color: 'var(--error)', label: 'Failure' };
      case 'pending':
        return { icon: AlertCircle, color: 'var(--warning)', label: 'Pending' };
      case 'warning':
        return { icon: AlertCircle, color: 'var(--warning)', label: 'Warning' };
      default:
        return null;
    }
  }

  let info = $derived(getStatusInfo(status));
</script>

{#if info}
  <div 
    class="commit-status-icon {className}" 
    style="--status-color: {info.color}; display: inline-flex;" 
    title={info.label}
  >
    <info.icon size={size} />
  </div>
{/if}

<style>
  .commit-status-icon {
    color: var(--status-color);
    flex-shrink: 0;
  }
</style>
