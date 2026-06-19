<script lang="ts">
  import { setContext } from 'svelte';
  import RepoSettingsGeneral from '../RepoSettingsGeneral.svelte';

  let { repoContext = { data: null, loading: false } }: Record<string, any> = $props();
  setContext('repo', repoContext);
</script>

<RepoSettingsGeneral />
