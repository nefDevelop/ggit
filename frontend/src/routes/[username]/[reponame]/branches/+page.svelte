<script lang="ts">
  import { page } from '$app/stores';
  import RepoBranches from '$lib/components/RepoBranches.svelte';

  const username = $page.params.username;
  const reponame = $page.params.reponame;
</script>

<svelte:head>
  <title>Branches · {username}/{reponame}</title>
</svelte:head>

<div class="page-container">
  <RepoBranches {username} {reponame} />
</div>

<style>
  .page-container {
    max-width: 1200px;
    margin: 0 auto;
    padding: 2rem;
  }
</style>
