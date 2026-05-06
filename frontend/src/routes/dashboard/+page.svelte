<script lang="ts">
  import DashboardView from '$lib/components/DashboardView.svelte';
  import { authState } from '$lib/auth.svelte';
  import { onMount } from 'svelte';
  import { goto } from '$app/navigation';

  $effect(() => {
    if (!authState.loading && !authState.isAuthenticated) {
      goto('/user/login');
    }
  });
</script>

<svelte:head>
  <title>Dashboard · Ggit</title>
</svelte:head>

{#if authState.isAuthenticated}
  <DashboardView />
{:else}
  <div class="loading-container">
    <div class="spinner"></div>
  </div>
{/if}

<style>
  .loading-container {
    display: flex;
    justify-content: center;
    align-items: center;
    height: 80vh;
  }
  .spinner {
    width: 40px;
    height: 40px;
    border: 3px solid rgba(255,255,255,0.1);
    border-radius: 50%;
    border-top-color: var(--primary-color);
    animation: spin 1s ease-in-out infinite;
  }
  @keyframes spin {
    to { transform: rotate(360deg); }
  }
</style>
