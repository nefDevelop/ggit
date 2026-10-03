<script lang="ts">
  import { goto } from '$app/navigation';
  import { page } from '$app/state';
  import WikiView from '$lib/components/WikiView.svelte';

  // The index has no path of its own, so it forwards to the wiki root. Keeping the
  // redirect in a page component means /wiki works the same way the sidebar
  // links do, instead of needing a second copy of the view.
  const owner = $derived(page.params.username);
  const repo = $derived(page.params.reponame);

  $effect(() => {
    if (owner && repo) goto(`/${owner}/${repo}/wiki/Home`, { replaceState: true });
  });
</script>

<svelte:head>
  <title>Wiki · Ggit</title>
</svelte:head>

<div class="wiki-redirect">Abriendo la wiki…</div>

<style>
  .wiki-redirect {
    display: grid;
    place-items: center;
    min-height: 40vh;
    color: var(--md-sys-color-on-surface-variant);
  }
</style>
