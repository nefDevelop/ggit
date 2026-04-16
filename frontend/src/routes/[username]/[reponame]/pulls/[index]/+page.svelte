<script lang="ts">
  import { page } from '$app/state';
  import PullDetail from '$lib/components/PullDetail.svelte';
  
  let index = $derived(page.params.index);
</script>

<svelte:head>
  <title>#{index} · Pull Requests · Ggit</title>
</svelte:head>

<PullDetail {index} />
