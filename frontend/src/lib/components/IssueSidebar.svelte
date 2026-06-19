<script lang="ts">
  import Avatar from "$lib/components/Avatar.svelte";
  import Button from "$lib/components/ui/Button.svelte";
  import Select from "$lib/components/ui/Select.svelte";
  import Settings from "lucide-svelte/icons/settings";
  import X from "lucide-svelte/icons/x";

  let { issue = null, allLabels = [], allMilestones = [], allCollaborators = [], onUpdate }: {
    issue: any;
    allLabels: any[];
    allMilestones: any[];
    allCollaborators: any[];
    onUpdate?: (patch: any) => void;
  } = $props();

  let showLabelEditor = $state(false);
  let showMilestoneEditor = $state(false);
  let showAssigneeEditor = $state(false);

  const labelOptions = $derived(allLabels.map((l: any) => ({ value: l.id, label: l.name })));
  const milestoneOptions = $derived([
    { value: 0, label: 'No Milestone' },
    ...allMilestones.map(m => ({ value: m.id, label: m.title })),
  ]);
</script>

<aside class="sidebar">
  <div class="side-item">
    <div class="side-header">
      <h4 class="side-title">Assignee</h4>
      <Button variant="ghost" icon={showAssigneeEditor ? X : Settings} size="sm" onclick={() => showAssigneeEditor = !showAssigneeEditor} />
    </div>
    <div class="side-val">
      {#if showAssigneeEditor}
        <div class="editor-popover">
          <Select options={[{value: 0, label: 'Unassigned'}, ...allCollaborators.map(c => ({ value: c.id, label: c.username }))]} value={issue?.assignee?.id || 0}
            onchange={(id) => { onUpdate?.({ assignee_id: id === 0 ? -1 : id }); showAssigneeEditor = false; }} />
        </div>
      {/if}
      {#if issue?.assignee}
        <div class="user-pill">
          <Avatar src={issue.assignee.avatar_url} username={issue.assignee.username} size={24} />
          {issue.assignee.username}
        </div>
      {:else}
        <span class="muted">No one assigned</span>
      {/if}
    </div>
  </div>

  <div class="side-item">
    <div class="side-header">
      <h4 class="side-title">Labels</h4>
      <Button variant="ghost" icon={showLabelEditor ? X : Settings} size="sm" onclick={() => showLabelEditor = !showLabelEditor} />
    </div>
    <div class="side-val">
      {#if showLabelEditor}
        <div class="editor-popover">
          <Select options={[{value: 0, label: 'Add Label...'}, ...labelOptions.filter(opt => !issue?.labels?.some((l: any) => l.id === opt.value))]} value={0}
            onchange={(id) => { if (id === 0) return; const newLabels = [...(issue?.labels || []).map((l: any) => l.id), id]; onUpdate?.({ label_ids: newLabels }); }} />
        </div>
      {/if}
      <div class="labels-grid">
        {#each issue?.labels || [] as label}
          <span class="label-tag" style="background: {label.color}">
            {label.name}
            {#if showLabelEditor}
              <button class="remove-label" onclick={() => { const newLabels = issue.labels.filter((l: any) => l.id !== label.id).map((l: any) => l.id); onUpdate?.({ label_ids: newLabels }); }}>
                <X size={12} />
              </button>
            {/if}
          </span>
        {/each}
      </div>
    </div>
  </div>

  <div class="side-item">
    <div class="side-header">
      <h4 class="side-title">Milestone</h4>
      <Button variant="ghost" icon={showMilestoneEditor ? X : Settings} size="sm" onclick={() => showMilestoneEditor = !showMilestoneEditor} />
    </div>
    <div class="side-val">
      {#if showMilestoneEditor}
        <div class="editor-popover">
          <Select options={milestoneOptions} value={issue?.milestone?.id || 0}
            onchange={(id) => { onUpdate?.({ milestone_id: id }); showMilestoneEditor = false; }} />
        </div>
      {/if}
      {#if issue?.milestone}
        <span>{issue.milestone.title}</span>
      {:else}
        <span class="muted">No milestone</span>
      {/if}
    </div>
  </div>
</aside>

<style>
  .sidebar { width: 260px; flex-shrink: 0; display: flex; flex-direction: column; gap: 1.5rem; }
  .side-item { display: flex; flex-direction: column; gap: 0.5rem; }
  .side-header { display: flex; align-items: center; justify-content: space-between; }
  .side-title { margin: 0; font-size: 0.8rem; font-weight: 800; text-transform: uppercase; color: var(--text-muted); }
  .side-val { font-size: 0.9rem; }
  .muted { color: var(--text-muted); font-style: italic; }
  .user-pill { display: flex; align-items: center; gap: 0.5rem; font-weight: 600; }
  .editor-popover { margin: 0.5rem 0; }
  .labels-grid { display: flex; flex-wrap: wrap; gap: 0.25rem; margin-top: 0.5rem; }
  .label-tag {
    display: inline-flex; align-items: center; gap: 0.25rem;
    padding: 0.15rem 0.5rem; border-radius: 0.25rem;
    font-size: 0.75rem; font-weight: 700; color: #fff;
  }
  .remove-label {
    background: none; border: none; color: rgba(255,255,255,0.7);
    cursor: pointer; padding: 0; display: flex;
  }
  .remove-label:hover { color: #fff; }
  @media (max-width: 900px) { .sidebar { width: 100%; } }
</style>
