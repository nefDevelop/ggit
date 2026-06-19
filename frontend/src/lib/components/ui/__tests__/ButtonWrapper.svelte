<script lang="ts">
  import Button from '../Button.svelte';
  let { childText = '', ...rest }: Record<string, any> = $props();
</script>

<Button {...rest}>
  {childText}
</Button>
