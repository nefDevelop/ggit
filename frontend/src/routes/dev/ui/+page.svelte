<script lang="ts">
  import { dev } from '$app/environment';
  import { onMount } from 'svelte';
  import Card from '$lib/components/ui/Card.svelte';
  import Button from '$lib/components/ui/Button.svelte';
  import Input from '$lib/components/ui/Input.svelte';
  import Select from '$lib/components/ui/Select.svelte';
  import Checkbox from '$lib/components/ui/Checkbox.svelte';
  import EmptyState from '$lib/components/ui/EmptyState.svelte';
  import Skeleton from '$lib/components/ui/Skeleton.svelte';
  import PageHeader from '$lib/components/ui/PageHeader.svelte';
  import PageActions from '$lib/components/ui/PageActions.svelte';
  import Palette from 'lucide-svelte/icons/palette';
  import Inbox from 'lucide-svelte/icons/inbox';

  const themes = [
    { value: 'default', label: 'Ggit Dynamic' },
    { value: 'ayu', label: 'Ayu Sunrise' },
    { value: 'monokai', label: 'Monokai Pro' },
    { value: 'nord', label: 'Nordic Frost' }
  ];

  const schemes = [
    { value: 'light', label: 'Claro' },
    { value: 'dark', label: 'Oscuro' }
  ];

  const paddings = ['none', 'sm', 'md', 'lg'] as const;

  let theme = $state('default');
  let scheme = $state('light');
  let text = $state('');
  let choice = $state('b');
  let checked = $state(true);

  const surfaces = [
    '--md-sys-color-background',
    '--md-sys-color-surface',
    '--md-sys-color-surface-container-lowest',
    '--md-sys-color-surface-container-low',
    '--md-sys-color-surface-container',
    '--md-sys-color-surface-container-high',
    '--md-sys-color-surface-container-highest',
    '--md-sys-color-surface-variant',
    '--sunken-1',
    '--sunken-2',
    '--sunken-3'
  ];

  onMount(() => {
    theme = localStorage.getItem('theme') || 'default';
    const saved = localStorage.getItem('color-scheme');
    scheme = saved === 'dark' ? 'dark' : 'light';
  });

  $effect(() => {
    if (!dev) return;
    document.documentElement.setAttribute('data-theme', theme);
    document.documentElement.setAttribute('data-color-scheme', scheme);
    localStorage.setItem('theme', theme);
    localStorage.setItem('color-scheme', scheme);
  });
</script>

<svelte:head>
  <title>UI Kit · Dev · Ggit</title>
</svelte:head>

{#if !dev}
  <div class="page">
    <Card>
      <p>Esta página solo está disponible en desarrollo.</p>
    </Card>
  </div>
{:else}
  <div class="page">
    <PageHeader
      title="UI Kit"
      subtitle="Superficies, componentes y temas en un solo sitio"
      icon={Palette}
    >
      {#snippet actions()}
        <Button variant="ghost" onclick={() => { theme = 'default'; scheme = 'light'; }}>
          Restaurar tema
        </Button>
      {/snippet}
    </PageHeader>

    <Card class="controls">
      <Select label="Tema" bind:value={theme} options={themes} class="control" />
      <Select label="Esquema" bind:value={scheme} options={schemes} class="control" />
      <p class="hint">
        Los cambios se aplican a toda la aplicación y se guardan en localStorage.
      </p>
    </Card>

    <section>
      <h2>Superficies</h2>
      <div class="swatches">
        {#each surfaces as token}
          <div class="swatch-cell">
            <div class="swatch" style={`background: var(${token})`}></div>
            <code>{token}</code>
          </div>
        {/each}
        <div class="swatch-cell">
          <div class="swatch scrim-sample"><span>scrim</span></div>
          <code>--scrim</code>
        </div>
      </div>
    </section>

    <section>
      <h2>Tonal</h2>
      <p class="hint">
        Canónico: <code>secondary-container</code>. Las tres variantes deben verse
        iguales; el tinte primary queda como referencia descartada.
      </p>
      <div class="tonal-grid">
        <div>
          <div class="tonal-sample tonal">.tonal global</div>
          <code>canónico</code>
        </div>
        <Card tonal>
          Card tonal
        </Card>
        <div>
          <div class="tonal-sample island tonal">.island.tonal</div>
          <code>island</code>
        </div>
        <div>
          <div class="tonal-sample primary-tint">primary 5%</div>
          <code>descartado</code>
        </div>
      </div>
    </section>

    <section>
      <h2>Botones</h2>
      <div class="row">
        <Button variant="primary">Primary</Button>
        <Button variant="secondary">Secondary</Button>
        <Button variant="tertiary">Tertiary</Button>
        <Button variant="tonal">Tonal</Button>
        <Button variant="outline">Outline</Button>
        <Button variant="ghost">Ghost</Button>
        <Button variant="danger">Danger</Button>
      </div>
      <div class="row">
        <Button size="sm">Small</Button>
        <Button size="md">Medium</Button>
        <Button size="lg">Large</Button>
        <Button disabled>Disabled</Button>
        <Button loading>Loading</Button>
      </div>
    </section>

    <section>
      <h2>Formularios</h2>
      <div class="form-grid">
        <Input label="Normal" placeholder="Escribe algo" bind:value={text} />
        <Input label="Con ayuda" hint="Texto de ayuda" placeholder="Opcional" />
        <Input label="Con error" error="Este campo es obligatorio" placeholder="Error" />
        <Select label="Select" bind:value={choice} options={[
          { value: 'a', label: 'Opción A' },
          { value: 'b', label: 'Opción B' }
        ]} hint="Mismo contrato que Input" />
      </div>
      <Checkbox bind:checked label="Checkbox con descripción" description="Una línea de contexto" />
    </section>

    <section>
      <h2>Tarjetas</h2>
      <div class="card-grid">
        {#each paddings as padding}
          <Card {padding}>
            <code>padding="{padding}"</code>
          </Card>
        {/each}
      </div>
    </section>

    <section>
      <h2>Cabeceras y acciones</h2>
      <Card>
        <PageActions>
          <Button variant="tonal" size="sm">Filtro abierto</Button>
          <Button variant="ghost" size="sm">Filtro cerrado</Button>
          <Button variant="primary" size="sm">Acción principal</Button>
        </PageActions>
      </Card>
    </section>

    <section>
      <h2>Estados</h2>
      <div class="state-grid">
        <EmptyState icon={Inbox} title="Nada por aquí" description="Este es el estado vacío compartido.">
          <Button variant="primary" size="sm">Acción</Button>
        </EmptyState>
        <Card>
          <div class="stack">
            <Skeleton variant="form" rows={2} />
            <Skeleton variant="list" rows={2} />
          </div>
        </Card>
      </div>
    </section>

    <section>
      <h2>Tipografía</h2>
      <Card class="type-card">
        <h1>Heading 1</h1>
        <h1 class="text-gradient">Gradient title</h1>
        <h2>Heading 2</h2>
        <h3>Heading 3</h3>
        <h4>Heading 4</h4>
        <p>Párrafo de cuerpo con <code>código inline</code> y un enlace.</p>
      </Card>
    </section>
  </div>
{/if}

<style>
  :global(.controls) {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
    gap: var(--space-5);
    align-items: end;
  }

  :global(.control) :global(.select-container) {
    width: 100%;
  }

  .hint {
    margin: 0;
    color: var(--md-sys-color-on-surface-variant);
    font-size: 0.85rem;
    line-height: 1.5;
  }

  section {
    display: flex;
    flex-direction: column;
    gap: var(--space-4);
  }

  section h2 {
    margin: 0;
  }

  .swatches {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
    gap: var(--space-4);
  }

  .swatch-cell {
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
  }

  .swatch {
    height: 64px;
    border-radius: var(--space-3);
    border: 1px solid var(--md-sys-color-outline-variant);
  }

  .scrim-sample {
    background: var(--scrim);
    display: flex;
    align-items: center;
    justify-content: center;
    color: #fff;
    font-size: 0.8rem;
  }

  .tonal-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(220px, 1fr));
    gap: var(--space-4);
  }

  .tonal-grid > div {
    display: flex;
    flex-direction: column;
    gap: var(--space-2);
  }

  .tonal-sample {
    display: flex;
    align-items: center;
    justify-content: center;
    min-height: 96px;
    padding: var(--space-5);
    border-radius: var(--radius);
    border: 1px solid var(--md-sys-color-outline-variant);
  }

  .primary-tint {
    background: color-mix(in srgb, var(--md-sys-color-primary), transparent 95%);
  }

  .row {
    display: flex;
    flex-wrap: wrap;
    gap: var(--space-3);
  }

  .form-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(240px, 1fr));
    gap: var(--space-5);
  }

  .card-grid {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(180px, 1fr));
    gap: var(--space-4);
  }

  .state-grid {
    display: grid;
    grid-template-columns: repeat(auto-fit, minmax(280px, 1fr));
    gap: var(--space-4);
    align-items: stretch;
  }

  .stack {
    display: flex;
    flex-direction: column;
    gap: var(--space-4);
  }

  :global(.type-card) h1,
  :global(.type-card) h2,
  :global(.type-card) h3,
  :global(.type-card) h4 {
    margin: 0 0 var(--space-2);
  }

  code {
    font-size: 0.8rem;
    color: var(--md-sys-color-on-surface-variant);
  }
</style>
