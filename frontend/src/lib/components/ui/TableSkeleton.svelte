<script lang="ts">
  import Skeleton from "./Skeleton.svelte";

  interface Props {
    columns: number;
    rows?: number;
  }

  let { columns, rows = 5 }: Props = $props();
</script>

{#each Array(rows) as _}
  <tr class="skeleton-row">
    <td colspan={columns}>
      <div class="row-flex">
        <Skeleton height="3.5rem" radius="0.75rem" />
      </div>
    </td>
  </tr>
{/each}

<style>
  .skeleton-row td {
    padding: 1rem;
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
  }
  .row-flex {
    display: flex;
    align-items: center;
    gap: 1rem;
  }
</style>
