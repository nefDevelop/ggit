<script lang="ts">
  import RepoFileBrowser from '$lib/components/RepoFileBrowser.svelte';
  import { page } from '$app/state';
</script>

<svelte:head>
  <title>{page.params.username}/{page.params.reponame} · Ggit</title>
</svelte:head>

<RepoFileBrowser path="" />
