<script lang="ts">
  import { page } from '$app/state';
  import ReleaseEditor from '$lib/components/ReleaseEditor.svelte';

  let username = $derived(page.params.username);
  let reponame = $derived(page.params.reponame);
</script>

<svelte:head>
  <title>New Release · ggit</title>
</svelte:head>

<div class="page-container">
  <ReleaseEditor {username} {reponame} />
</div>

<style>
  .page-container {
    padding: 3rem 1rem;
    display: flex;
    justify-content: center;
  }
</style>
