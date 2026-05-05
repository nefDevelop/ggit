<script lang="ts">
  import type { Snippet } from 'svelte';
  
  interface Props {
    title: string;
    subtitle?: string;
    icon?: string;
    actions?: Snippet;
    avatar?: Snippet;
    compact?: boolean; // For smaller cards
  }
  
  let { title, subtitle, icon, actions, avatar, compact = false }: Props = $props();
</script>

<header class="page-header" class:compact>
  <div class="header-main">
    {#if avatar}
      <div class="avatar-slot">
        {@render avatar()}
      </div>
    {:else if icon}
      <div class="icon-slot">
        <span class="material-symbols-rounded">{icon}</span>
      </div>
    {/if}
    
    <div class="title-stack">
      <h2>{title}</h2>
      {#if subtitle}
        <p class="subtitle">{subtitle}</p>
      {/if}
    </div>
  </div>
  
  {#if actions}
    <div class="header-actions">
      {@render actions()}
    </div>
  {/if}
</header>

<style>
  .page-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
    gap: 2rem;
    padding: 2.5rem;
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
  }

  .page-header.compact {
    padding: 1.5rem 2rem;
  }

  .header-main {
    display: flex;
    align-items: center;
    gap: 2rem;
    min-width: 0;
  }

  .icon-slot {
    font-size: 3rem;
    color: var(--primary-color);
    background: var(--md-sys-color-primary-container);
    padding: 1.25rem;
    border-radius: 2rem;
    flex-shrink: 0;
    display: flex;
    align-items: center;
    justify-content: center;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .icon-slot span { font-size: inherit; }

  .avatar-slot {
    flex-shrink: 0;
    transition: all 0.3s cubic-bezier(0.4, 0, 0.2, 1);
  }

  .title-stack {
    display: flex;
    flex-direction: column;
    min-width: 0;
  }

  h2 {
    margin: 0;
    font-size: 2.25rem;
    font-weight: 950;
    letter-spacing: -0.04em;
    line-height: 1.1;
    overflow: hidden;
    text-overflow: ellipsis;
    white-space: nowrap;
    transition: all 0.3s ease;
  }

  @media (max-width: 1400px) {
    h2 { font-size: 1.85rem; }
    .icon-slot { font-size: 2.25rem; padding: 1rem; border-radius: 1.5rem; }
    .page-header { padding: 2rem; gap: 1.5rem; }
  }

  @media (max-width: 1100px) {
    h2 { font-size: 1.5rem; }
    .subtitle { font-size: 1rem; }
    .icon-slot { font-size: 1.75rem; padding: 0.8rem; border-radius: 1.25rem; }
  }

  .page-header.compact h2 {
    font-size: 1.75rem;
  }

  .subtitle {
    margin: 0.5rem 0 0 0;
    color: var(--text-muted);
    font-size: 1.1rem;
    line-height: 1.4;
  }

  .header-actions {
    display: flex;
    gap: 0.75rem;
    align-items: center;
    flex-shrink: 0;
  }

  @media (max-width: 800px) {
    .page-header {
      padding: 1.5rem;
      gap: 1rem;
    }

    .header-main {
      gap: 1rem;
    }

    h2 { font-size: 1.5rem; }
    .page-header.compact h2 { font-size: 1.25rem; }
    
    .subtitle { display: none; }

    .icon-slot {
      font-size: 1.5rem;
      padding: 0.6rem;
      border-radius: 1rem;
    }

    /* Scaling down avatars if they are wrapped in slots */
    .avatar-slot :global(.avatar) {
      width: 48px !important;
      height: 48px !important;
    }
    .avatar-slot :global(.avatar-shield-header) {
       width: 52px !important;
       height: 52px !important;
    }
  }

  @media (max-width: 500px) {
    .page-header {
      flex-direction: row; /* Still row to keep it compact */
      flex-wrap: wrap;
    }
    .header-actions {
      margin-left: auto;
    }
  }
</style>
