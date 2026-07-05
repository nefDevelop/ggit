<script lang="ts">
  import { apiFetch } from "$lib/api";
  import { onMount } from "svelte";
  import BranchSelector from "./BranchSelector.svelte";
  import Button from "./ui/Button.svelte";
  

import X from "lucide-svelte/icons/x";
import CloudUpload from "lucide-svelte/icons/cloud-upload";
import FileText from "lucide-svelte/icons/file-text";
import Trash2 from "lucide-svelte/icons/trash-2";



  let { 
    username, 
    reponame, 
    release = null, 
    onSaved = (data?: any) => {}, 
    onCancel = () => {} 
  } = $props<{
    username?: string;
    reponame?: string;
    release?: any;
    onSaved?: (data?: any) => void;
    onCancel?: () => void;
  }>();

  let title = $state("");
  let tagName = $state("");
  let target = $state("master");
  let note = $state("");
  let isDraft = $state(false);
  let isPrerelease = $state(false);
  let attachments = $state([] as any[]);

  $effect(() => {
    title = release?.title || "";
    tagName = release?.tag_name || "";
    target = release?.target || target;
    note = release?.note || "";
    isDraft = release?.is_draft || false;
    isPrerelease = release?.is_prerelease || false;
    attachments = (release?.attachments || []) as any[];
  });
  let branches = $state([] as any[]);
  let tags = $state([] as any[]);
  let uploading = $state(false);
  let saving = $state(false);

  async function fetchData() {
    try {
      const bRes = await apiFetch(`/api/v2/${username}/${reponame}/branches`);
      if (bRes.ok) {
        const data = await bRes.json();
        branches = [data.default, ...data.active, ...data.stale].map((b) => ({
          name: b.name,
        }));
        if (!release) target = data.default.name;
      }

      const tRes = await apiFetch(`/api/v2/${username}/${reponame}/tags`);
      if (tRes.ok) {
        const data = await tRes.json();
        tags = data.map((t: any) => ({ name: t.name }));
      }
    } catch (e) {
      console.error(e);
    }
  }

  async function handleFileUpload(e: Event) {
    const files = (e.target as HTMLInputElement).files;
    if (!files || files.length === 0) return;

    uploading = true;
    for (const file of Array.from(files)) {
      const formData = new FormData();
      formData.append("file", file);

      try {
        const res = await apiFetch(
          `/api/v2/${username}/${reponame}/releases/attachments`,
          {
            method: "POST",
            body: formData,
          },
        );

        if (res.ok) {
          const attachment = await res.json();
          attachments = [...attachments, attachment];
        }
      } catch (e) {
        console.error("Upload failed", e);
      }
    }
    uploading = false;
  }

  function removeAttachment(uuid: string) {
    attachments = attachments.filter((a) => a.uuid !== uuid);
  }

  async function save(publish = false) {
    if (saving || !tagName) return;
    saving = true;
    const method = release ? "PATCH" : "POST";
    const url = release
      ? `/api/v2/${username}/${reponame}/releases/${release.id}`
      : `/api/v2/${username}/${reponame}/releases`;

    const body = {
      tag_name: tagName,
      target: target,
      title: title,
      content: note,
      draft: !publish && isDraft,
      prerelease: isPrerelease,
      files: attachments.map((a) => a.uuid),
    };

    try {
      const res = await apiFetch(url, {
        method,
        headers: { "Content-Type": "application/json" },
        body: JSON.stringify(body),
      });

      if (res.ok) {
        const savedRelease = await res.json();
        onSaved(savedRelease);
      }
    } catch (e) {
      console.error(e);
    } finally {
      saving = false;
    }
  }

  onMount(fetchData);
</script>

<div class="release-editor shadow-xl">
  <header class="editor-header">
    <h3>{release ? "Edit Release" : "New Release"}</h3>
    <Button variant="ghost" icon={X} onclick={onCancel} aria-label="Close editor" />
  </header>

  <div class="editor-grid">
    <div class="main-form">
      <div class="field-row">
        <div class="field tag-field">
          <label for="tag-name">Tag Version</label>
          <input
            id="tag-name"
            type="text"
            list="tag-suggestions"
            bind:value={tagName}
            placeholder="e.g. v1.0.0"
            disabled={!!release}
          />
          <datalist id="tag-suggestions">
            {#each tags as tag}
              <option value={tag.name}></option>
            {/each}
          </datalist>
          <p class="help">
            Choose an existing tag, or create a new one on publish.
          </p>
        </div>

        <div class="field target-field">
          <BranchSelector
            label="Target Branch"
            {branches}
            bind:selected={target}
          />
        </div>
      </div>

      <div class="field">
        <label for="rel-title">Release Title</label>
        <input
          id="rel-title"
          type="text"
          bind:value={title}
          placeholder="Release title"
        />
      </div>

      <div class="field">
        <label for="rel-note">Release Notes</label>
        <textarea
          id="rel-note"
          bind:value={note}
          placeholder="Describe the changes in this version..."
        ></textarea>
      </div>

      <div class="attachments-zone">
        <span class="field-label">Binary Attachments</span>
        <div class="upload-area" class:uploading>
          <CloudUpload size={48} />
          <p>Drag and drop assets here, or <strong>browse</strong></p>
          <input id="release-attachments" type="file" multiple onchange={handleFileUpload} aria-label="Upload attachments" />
        </div>

        {#if attachments.length > 0}
          <div class="attachment-list">
            {#each attachments as a}
              <div class="attachment-item tonal">
                <FileText size={18} />
                <span class="name">{a.name}</span>
                <Button
                  variant="ghost"
                  size="sm"
                  icon={Trash2}
                  class="danger-btn"
                  onclick={() => removeAttachment(a.uuid)}
                  aria-label={`Remove ${a.name}`}
                />
              </div>
            {/each}
          </div>
        {/if}
      </div>
    </div>

    <aside class="editor-sidebar">
      <div class="options tonal">
        <h4>Options</h4>
        <label class="checkbox-row">
          <input type="checkbox" bind:checked={isPrerelease} />
          <div class="check-text">
            <strong>This is a pre-release</strong>
            <span
              >Inform users that this release is not ready for production.</span
            >
          </div>
        </label>

        {#if !release || (release && release.is_draft)}
          <label class="checkbox-row">
            <input type="checkbox" bind:checked={isDraft} />
            <div class="check-text">
              <strong>Save as draft</strong>
              <span>Only you and admins will see this until published.</span>
            </div>
          </label>
        {/if}
      </div>

      <div class="actions">
        <Button
          variant="primary"
          size="lg"
          onclick={() => save(!isDraft)}
          disabled={saving}
          loading={saving}
        >
          {release?.is_draft && !isDraft
            ? "Publish Release"
            : release
              ? "Update Release"
              : "Create Release"}
        </Button>
        <Button variant="tonal" onclick={onCancel}>Cancel</Button>
      </div>
    </aside>
  </div>
</div>

<style>
  .release-editor {
    width: 100%;
    max-width: 1100px;
    padding: 2.5rem;
    display: flex;
    flex-direction: column;
    gap: 2rem;
    background: var(--md-sys-color-surface);
    border-radius: 1.5rem;
    border: 1px solid var(--md-sys-color-outline-variant);
  }

  .editor-header {
    display: flex;
    justify-content: space-between;
    align-items: center;
  }
  .editor-header h3 {
    margin: 0;
    font-size: 1.5rem;
    font-weight: 800;
  }

  .editor-grid {
    display: grid;
    grid-template-columns: 1fr 350px;
    gap: 3rem;
    align-items: start;
  }

  .main-form {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }
  .field-row {
    display: grid;
    grid-template-columns: 1fr 1fr;
    gap: 1.5rem;
  }
  @media (max-width: 600px) {
    .field-row {
      grid-template-columns: 1fr;
    }
  }

  .field {
    display: flex;
    flex-direction: column;
    gap: 0.5rem;
  }
  .field label {
    font-size: 0.85rem;
    font-weight: 700;
    color: var(--text-muted);
    padding-left: 0.5rem;
  }
  .field input,
  .field textarea {
    background: var(--md-sys-color-surface-variant);
    border: 1px solid var(--md-sys-color-outline-variant);
    color: #fff;
    padding: 0.75rem 1rem;
    border-radius: 0.75rem;
    font-size: 1rem;
    outline: none;
    width: 100%;
    transition: all 0.2s;
  }
  .field input:focus,
  .field textarea:focus {
    border-color: var(--md-sys-color-primary);
    background: var(--md-sys-color-surface-container-high);
  }
  .tonal {
    background: var(--md-sys-color-surface-container);
    border-radius: 1.25rem;
    border: 1px solid var(--md-sys-color-outline-variant);
  }
  .field textarea {
    height: 350px;
    resize: vertical;
    font-family: inherit;
  }
  .help {
    font-size: 0.75rem;
    color: var(--text-muted);
    margin: 0;
    padding-left: 0.5rem;
  }

  .attachments-zone {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }
  .field-label {
    font-size: 0.85rem;
    font-weight: 700;
    color: var(--text-muted);
    padding-left: 0.5rem;
  }

  .upload-area {
    position: relative;
    border: 2px dashed var(--md-sys-color-outline-variant);
    border-radius: 1.5rem;
    padding: 2.5rem;
    text-align: center;
    display: flex;
    flex-direction: column;
    align-items: center;
    gap: 0.5rem;
    color: var(--text-muted);
    transition: all 0.2s;
    background: rgba(255, 255, 255, 0.01);
  }
  .upload-area:hover {
    border-color: var(--primary-color);
    background: rgba(var(--primary-rgb), 0.02);
    color: var(--text-color);
  }
  .upload-area input {
    position: absolute;
    inset: 0;
    opacity: 0;
    cursor: pointer;
  }
  .upload-area.uploading {
    opacity: 0.5;
    pointer-events: none;
  }

  .attachment-list {
    display: grid;
    grid-template-columns: repeat(auto-fill, minmax(200px, 1fr));
    gap: 1rem;
    margin-top: 1rem;
  }
  .attachment-item {
    padding: 0.75rem 1rem;
    display: flex;
    align-items: center;
    gap: 0.75rem;
    font-size: 0.9rem;
    position: relative;
  }
  .attachment-item .name {
    flex: 1;
    white-space: nowrap;
    overflow: hidden;
    text-overflow: ellipsis;
  }
  .editor-sidebar {
    display: flex;
    flex-direction: column;
    gap: 2rem;
    position: sticky;
    top: 1rem;
  }
  .options {
    padding: 1.5rem;
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
  }
  .options h4 {
    margin: 0;
    font-size: 1.1rem;
  }

  .checkbox-row {
    display: flex;
    gap: 1rem;
    cursor: pointer;
    transition: all 0.2s;
    user-select: none;
    padding: 0.5rem;
    border-radius: 0.75rem;
  }
  .checkbox-row:hover {
    background: rgba(255,255,255,0.05);
  }
  .checkbox-row:active {
    transform: scale(0.98);
  }
  .checkbox-row input {
    margin-top: 0.25rem;
  }
  .checkbox-row input {
    width: 1.25rem;
    height: 1.25rem;
    margin-top: 0.25rem;
    cursor: pointer;
  }
  .check-text {
    display: flex;
    flex-direction: column;
    gap: 0.2rem;
  }
  .check-text strong {
    font-size: 0.95rem;
  }
  .check-text span {
    font-size: 0.8rem;
    color: var(--text-muted);
  }

  .actions {
    display: flex;
    flex-direction: column;
    gap: 0.75rem;
  }

</style>
