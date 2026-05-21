<script lang="ts">
  import { page } from '$app/stores';
  import RepoBranches from '$lib/components/RepoBranches.svelte';

  const username = $page.params.username;
  const reponame = $page.params.reponame;
</script>

<svelte:head>
  <title>Branches · {username}/{reponame}</title>
</svelte:head>

<RepoBranches {username} {reponame} />
