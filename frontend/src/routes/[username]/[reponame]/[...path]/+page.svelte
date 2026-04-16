<script lang="ts">
  import RepoFileBrowser from '$lib/components/RepoFileBrowser.svelte';
  import { page } from '$app/state';
  
  let path = $derived(page.params.path || '');
</script>

<RepoFileBrowser {path} />
