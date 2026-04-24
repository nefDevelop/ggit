<script lang="ts">
  import { page } from '$app/state';
  import WikiView from '$lib/components/WikiView.svelte';
  
  let path = $derived(page.params.path || 'Home');
</script>

<svelte:head>
  <title>{path} · Wiki · Ggit</title>
</svelte:head>

<WikiView {path} />
