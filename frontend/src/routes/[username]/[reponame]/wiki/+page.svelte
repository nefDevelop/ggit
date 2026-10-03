<script lang="ts">
  import { goto } from '$app/navigation';
  import { page } from '$app/state';
  import { apiFetch } from '$lib/api';
  import WikiView from '$lib/components/WikiView.svelte';

  const owner = $derived(page.params.username);
  const repo = $derived(page.params.reponame);
  let resolved = $state('');

  // The index forwards to the wiki's entry page, which is the first file in the
  // tree. It used to forward to "Home" unconditionally, so on a repository
  // without a wiki every visit to /wiki asked the API for a page that does not
  // exist and came back 404. With an empty tree there is nothing to forward to
  // and the view renders its own empty state.
  $effect(() => {
    const o = owner;
    const r = repo;
    if (!o || !r || resolved) return;

    let cancelled = false;
    (async () => {
      try {
        const res = await apiFetch(`/api/v2/${o}/${r}/wiki`);
        if (cancelled || !res.ok) return;
        const data = await res.json();
        const pages = Array.isArray(data) ? data : (data?.pages ?? []);
        const entry = pages.find((p: { name?: string; path?: string }) => p.name || p.path);
        const target = entry?.path || entry?.name;
        if (target) {
          resolved = target;
          await goto(`/${o}/${r}/wiki/${target}`, { replaceState: true });
        } else {
          // Nothing to open. Marked as resolved so the effect does not retry on
          // every re-render, and the view stays on its empty state.
          resolved = ' ';
        }
      } catch {
        if (!cancelled) resolved = ' ';
      }
    })();

    return () => { cancelled = true; };
  });
</script>

<svelte:head>
  <title>Wiki · Ggit</title>
</svelte:head>

{#if resolved}
  <WikiView path={resolved.trim()} />
{:else}
  <div class="wiki-redirect">Abriendo la wiki…</div>
{/if}

<style>
  .wiki-redirect {
    display: grid;
    place-items: center;
    min-height: 40vh;
    color: var(--md-sys-color-on-surface-variant);
  }
</style>