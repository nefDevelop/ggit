<script lang="ts">
  
import CheckCircle2 from "lucide-svelte/icons/check-circle-2";
import XCircle from "lucide-svelte/icons/x-circle";
import AlertCircle from "lucide-svelte/icons/alert-circle";
import HelpCircle from "lucide-svelte/icons/help-circle";


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
  {@const Icon = info.icon}
  <div
    class="commit-status-icon {className}"
    style="--status-color: {info.color}; display: inline-flex;"
    title={info.label}
  >
    <Icon size={size} />
  </div>
{/if}

<style>
  .commit-status-icon {
    color: var(--status-color);
    flex-shrink: 0;
  }
</style>
