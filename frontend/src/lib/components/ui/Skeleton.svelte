<script lang="ts">
  interface Props {
    width?: string;
    height?: string;
    radius?: string;
    class?: string;
    circle?: boolean;
  }

  let { 
    width = '100%', 
    height = '1rem', 
    radius = '0.5rem', 
    class: className = '',
    circle = false
  }: Props = $props();

  const style = $derived(`
    width: ${width}; 
    height: ${height}; 
    border-radius: ${circle ? '50%' : radius};
  `);
</script>

<div class="skeleton {className}" {style}></div>

<style>
  .skeleton {
    background: var(--md-sys-color-surface-variant);
    animation: pulse 1.5s infinite ease-in-out;
    display: inline-block;
    vertical-align: middle;
  }

  @keyframes pulse {
    0% {
      opacity: 0.3;
    }
    50% {
      opacity: 0.7;
    }
    100% {
      opacity: 0.3;
    }
  }
</style>
