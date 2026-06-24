<script lang="ts">
  import Avatar from '../Avatar.svelte';
  let { src = '', username = '?', seed = '', size = 40, type = 'user', class: className = '' }: Record<string, any> = $props();
</script>
<Avatar {src} {username} {seed} {size} {type} class={className} />
