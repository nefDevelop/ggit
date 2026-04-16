<script lang="ts">
  import { page } from '$app/state';
  import IssueDetail from '$lib/components/IssueDetail.svelte';
  
  let index = $derived(page.params.index);
</script>

<svelte:head>
  <title>#{index} · Issues · Ggit</title>
</svelte:head>

<IssueDetail {index} />
