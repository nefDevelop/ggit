<script lang="ts">
  import { page } from '$app/stores';
  import { goto } from '$app/navigation';
  import WebEditor from '$lib/components/WebEditor.svelte';

  const username = $page.params.username;
  const reponame = $page.params.reponame;
  const branch = $page.params.branch;
  const treePath = $page.params.path;

  function handleComplete() {
    goto(`/${username}/${reponame}/src/${branch}/${treePath}`);
  }
</script>

<svelte:head>
  <title>Editing {treePath} · {username}/{reponame}</title>
</svelte:head>

<div class="page-container">
  <WebEditor 
    {username} 
    {reponame} 
    {branch} 
    {treePath} 
    onComplete={handleComplete}
  />
</div>

<style>
  .page-container {
    max-width: 1400px;
    margin: 0 auto;
    padding: 2rem;
  }
</style>
