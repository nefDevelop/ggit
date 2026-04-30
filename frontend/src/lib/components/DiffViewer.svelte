<script lang="ts">
  import {
    ChevronDown,
    ChevronRight,
    FileCode,
    Plus,
    Minus,
    Layers,
    FileText,
    Settings,
    MessageSquarePlus
  } from "lucide-svelte";
  import ReviewCommentComposer from "./ReviewCommentComposer.svelte";
  import MarkdownRenderer from "./ui/MarkdownRenderer.svelte";

  let { diff, allowComments = false, comments = [], onAddComment } = $props<{ 
    diff: any, 
    allowComments?: boolean, 
    comments?: any[],
    onAddComment?: (file: string, line: number, content: string) => void
  }>();

  let collapsedFiles = $state<Record<string, boolean>>({});
  let activeComposer = $state<{ file: string, line: number } | null>(null);

  function toggleFile(name: string) {
    collapsedFiles[name] = !collapsedFiles[name];
  }

  function getLineComments(fileName: string, lineNum: number) {
    return comments.filter((c: any) => c.tree_path === fileName && c.line === lineNum);
  }

  const getLineClass = (type: number) => {
    switch (type) {
      case 1:
        return "line-add";
      case 2:
        return "line-del";
      case 3:
        return "line-section";
      default:
        return "line-context";
    }
  };

  const getLineSign = (type: number) => {
    switch (type) {
      case 1:
        return "+";
      case 2:
        return "-";
      case 3:
        return "@@";
      default:
        return " ";
    }
  };
</script>

<div class="diff-viewer">
  <div class="diff-stats">
    <div class="stat-item">
      <Layers size={18} />
      <b>{diff?.Files?.length || 0}</b> archivos cambiados
    </div>
    <div class="stat-item success">
      <Plus size={16} />
      <b>{diff?.TotalAddition || 0}</b> adiciones
    </div>
    <div class="stat-item error">
      <Minus size={16} />
      <b>{diff?.TotalDeletion || 0}</b> eliminaciones
    </div>
  </div>

  <div class="diff-files">
    {#each diff?.Files || [] as file}
      <div class="diff-file" class:collapsed={collapsedFiles[file.Name]}>
        <header 
          class="file-header" 
          onclick={() => toggleFile(file.Name)} 
          onkeydown={(e) => e.key === 'Enter' && toggleFile(file.Name)}
          role="button"
          tabindex="0"
        >
          <div class="header-left">
            <button class="collapse-btn">
              {#if collapsedFiles[file.Name]}
                <ChevronRight size={18} />
              {:else}
                <ChevronDown size={18} />
              {/if}
            </button>
            <FileCode size={18} class="file-icon" />
            <span class="file-name">{file.Name}</span>
          </div>
          <div class="header-right">
            <span class="file-stats">
              <span class="add">+{file.Addition}</span>
              <span class="del">-{file.Deletion}</span>
            </span>
          </div>
        </header>

        {#if !collapsedFiles[file.Name]}
          <div class="file-content">
            {#each file.Sections as section}
              <div class="diff-section">
                {#each section.Lines as line}
                  <div class="diff-line-wrapper">
                    <div class="diff-line {getLineClass(line.Type)}">
                      <div class="line-num line-num-left">
                        {line.LeftLine || ""}
                      </div>
                      <div class="line-num line-num-right">
                        {#if allowComments && line.RightLine}
                          <button 
                            class="add-comment-btn" 
                            onclick={() => activeComposer = { file: file.Name, line: line.RightLine }}
                          >
                            <MessageSquarePlus size={14} />
                          </button>
                        {/if}
                        {line.RightLine || ""}
                      </div>
                      <div class="line-sign">{getLineSign(line.Type)}</div>
                      <div class="line-code">
                        <pre><code>{line.Content}</code></pre>
                      </div>
                    </div>

                    <!-- Existing Comments -->
                    {#if line.RightLine}
                      {#each getLineComments(file.Name, line.RightLine) as comment}
                        <div class="inline-comment-box tonal">
                          <header class="comment-header">
                             <strong>{comment.reviewer?.username || 'Usuario'}</strong>
                             <span class="time">{new Date(comment.created_at || Date.now()).toLocaleTimeString()}</span>
                          </header>
                          <div class="comment-body">
                             <MarkdownRenderer content={comment.content} />
                          </div>
                        </div>
                      {/each}
                    {/if}

                    <!-- New Comment Composer -->
                    {#if activeComposer?.file === file.Name && activeComposer?.line === line.RightLine}
                      <ReviewCommentComposer 
                        onSave={(content) => {
                          onAddComment?.(file.Name, line.RightLine, content);
                          activeComposer = null;
                        }}
                        onCancel={() => activeComposer = null}
                      />
                    {/if}
                  </div>
                {/each}
              </div>
            {/each}
          </div>
        {/if}
      </div>
    {/each}
  </div>
</div>

<style>
  .diff-viewer {
    display: flex;
    flex-direction: column;
    gap: 1.5rem;
    animation: fadeIn 0.3s ease-out;
  }

  .diff-stats {
    padding: 1rem 1.5rem;
    display: flex;
    gap: 2rem;
    align-items: center;
    font-size: 0.95rem;
  }

  .stat-item {
    display: flex;
    align-items: center;
    gap: 0.5rem;
    color: var(--md-sys-color-on-surface-variant);
  }

  .stat-item b {
    color: var(--md-sys-color-on-surface);
  }

  .stat-item.success b {
    color: var(--success);
  }
  .stat-item.error b {
    color: var(--error);
  }

  .diff-files {
    display: flex;
    flex-direction: column;
    gap: 1rem;
  }

  .diff-file {
    overflow: hidden;
    border: 1px solid var(--md-sys-color-outline-variant);
    transition: all 0.2s;
  }

  .file-header {
    background: var(--surface-variant);
    padding: 0.75rem 1.25rem;
    display: flex;
    justify-content: space-between;
    align-items: center;
    cursor: pointer;
    user-select: none;
  }

  .file-header:hover {
    background: var(--surface-variant-hover);
  }

  .header-left {
    display: flex;
    align-items: center;
    gap: 0.75rem;
  }

  .collapse-btn {
    background: transparent;
    border: none;
    color: var(--text-tertiary);
    display: flex;
    padding: 0;
    cursor: pointer;
  }

  .file-icon {
    color: var(--primary);
  }

  .file-name {
    font-weight: 600;
    font-size: 0.95rem;
    color: var(--md-sys-color-on-surface);
    font-family: "JetBrains Mono", monospace;
  }

  .file-stats {
    font-size: 0.85rem;
    font-weight: 700;
    display: flex;
    gap: 0.5rem;
  }

  .file-stats .add {
    color: var(--success);
  }
  .file-stats .del {
    color: var(--error);
  }

  .file-content {
    background: #0d0e10;
    overflow-x: auto;
  }

  .diff-line {
    display: flex;
    font-family: "JetBrains Mono", "Fira Code", monospace;
    font-size: 13px;
    line-height: 20px;
    white-space: pre;
  }

  .line-num {
    width: 45px;
    min-width: 45px;
    text-align: right;
    padding-right: 10px;
    color: var(--text-tertiary);
    background: rgba(255, 255, 255, 0.02);
    user-select: none;
    border-right: 1px solid var(--md-sys-color-outline-variant);
  }

  .line-sign {
    width: 25px;
    min-width: 25px;
    text-align: center;
    user-select: none;
    color: var(--text-tertiary);
  }

  .line-code {
    flex: 1;
    padding-left: 5px;
  }

  .line-code pre {
    margin: 0;
  }

  .diff-line-wrapper {
    display: flex;
    flex-direction: column;
  }

  .line-num-right {
    position: relative;
  }

  .add-comment-btn {
    position: absolute;
    right: 0;
    top: 50%;
    transform: translateY(-50%) translateX(50%);
    z-index: 10;
    background: var(--primary-color);
    color: white;
    border: none;
    border-radius: 4px;
    width: 24px;
    height: 24px;
    display: flex;
    align-items: center;
    justify-content: center;
    cursor: pointer;
    opacity: 0;
    transition: opacity 0.2s;
  }

  .diff-line:hover .add-comment-btn {
    opacity: 1;
  }

  .inline-comment-box {
    margin: 0.5rem 1rem 0.5rem 60px;
    padding: 1rem;
    border-radius: 0.75rem;
    border: 1px solid var(--md-sys-color-outline-variant);
    background: var(--md-sys-color-surface-container);
  }

  .comment-header {
    display: flex;
    justify-content: space-between;
    margin-bottom: 0.5rem;
    font-size: 0.85rem;
  }

  .comment-header strong {
    color: var(--secondary);
  }

  .comment-header .time {
    color: var(--text-muted);
    font-size: 0.75rem;
  }

  /* Colors */
  .line-add {
    background: rgba(103, 225, 163, 0.1);
  }
  .line-add .line-num {
    background: rgba(103, 225, 163, 0.15);
    color: #8ef9c1;
  }
  .line-add .line-sign {
    color: #8ef9c1;
  }

  .line-del {
    background: rgba(255, 180, 171, 0.1);
  }
  .line-del .line-num {
    background: rgba(255, 180, 171, 0.15);
    color: #ffcdc7;
  }
  .line-del .line-sign {
    color: #ffcdc7;
  }

  .line-section {
    background: rgba(168, 199, 250, 0.05);
    color: var(--primary);
    font-weight: 600;
  }
  .line-section .line-num {
    background: rgba(168, 199, 250, 0.1);
    color: var(--primary);
  }

  @keyframes fadeIn {
    from {
      opacity: 0;
      transform: translateY(5px);
    }
    to {
      opacity: 1;
      transform: translateY(0);
    }
  }
</style>
