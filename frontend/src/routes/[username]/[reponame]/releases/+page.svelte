<script lang="ts">
  import RepoReleases from '$lib/components/RepoReleases.svelte';
  import { page } from '$app/state';

  let username = $derived(page.params.username);
  let reponame = $derived(page.params.reponame);
</script>

<svelte:head>
  <title>Releases · ggit</title>
</svelte:head>

<div class="page-container">
  <RepoReleases {username} {reponame} />
</div>

<style>
  .page-container {
    padding: 0;
  }
</style>
