<script lang="ts">
  import FileContentViewer from '../FileContentViewer.svelte';
  let { repoData = null, username = '', reponame = '', realPath = '', showBlame = false, blameLoading = false, sidebarCollapsed = false, codeRows = [], formatSize, timeAgo, onToggleBlame, onToggleSidebar }: Record<string, any> = $props();
</script>
<FileContentViewer {repoData} {username} {reponame} {realPath} {showBlame} {blameLoading} {sidebarCollapsed} {codeRows} {formatSize} {timeAgo} {onToggleBlame} {onToggleSidebar} />
