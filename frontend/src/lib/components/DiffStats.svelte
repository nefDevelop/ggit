<script lang="ts">
  import Layers from "lucide-svelte/icons/layers";
  import Plus from "lucide-svelte/icons/plus";
  import Minus from "lucide-svelte/icons/minus";

  let { files = 0, additions = 0, deletions = 0 }: {
    files: number;
    additions: number;
    deletions: number;
  } = $props();
</script>

<div class="diff-stats">
  <div class="stat-item">
    <Layers size={18} />
    <b>{files}</b> archivos cambiados
  </div>
  <div class="stat-item success">
    <Plus size={16} />
    <b>{additions}</b> adiciones
  </div>
  <div class="stat-item error">
    <Minus size={16} />
    <b>{deletions}</b> eliminaciones
  </div>
</div>

<style>
  .diff-stats {
    padding: 1rem 1.5rem;
    display: flex;
    gap: 2rem;
    align-items: center;
    font-size: 0.95rem;
  }
  .stat-item { display: flex; align-items: center; gap: 0.5rem; color: var(--md-sys-color-on-surface-variant); }
  .stat-item b { color: var(--md-sys-color-on-surface); }
  .stat-item.success b { color: var(--success); }
  .stat-item.error b { color: var(--error); }
  @media (max-width: 480px) { .diff-stats { padding: 0.75rem; flex-wrap: wrap; gap: 0.5rem; } }
</style>
