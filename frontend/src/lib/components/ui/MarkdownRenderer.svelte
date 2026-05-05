<script lang="ts">
  import { marked } from 'marked';
  import DOMPurify from 'dompurify';
  import hljs from 'highlight.js';

  let { content = '', repoContext = null } = $props();

  // Configure marked with Ggit-specific extensions
  const renderer = new marked.Renderer();
  
  // Custom code highlighting
  renderer.code = ({ text, lang }: { text: string; lang?: string }) => {
    const language = lang && hljs.getLanguage(lang) ? lang : 'plaintext';
    const highlighted = hljs.highlight(text, { language }).value;
    return `<pre><code class="hljs language-${language}">${highlighted}</code></pre>`;
  };

  // Extension for Mentions (@user), Issues (#123) and Commits (sha)
  const ggitExtension = {
    name: 'ggitSpecial',
    level: 'inline',
    start(src: string) { return src.match(/[@#\b[0-9a-f]{7,40}]/)?.index; },
    tokenizer(src: string) {
      // Wikilinks: [[Page|Label]]
      const wikiMatch = /^\[\[([^\]|]+)(?:\|([^\]]+))?\]\]/.exec(src);
      if (wikiMatch) {
        const dest = wikiMatch[1].trim().replace(/ /g, '-');
        const label = (wikiMatch[2] || wikiMatch[1]).trim();
        return {
          type: 'ggitSpecial',
          raw: wikiMatch[0],
          text: label,
          href: `./${dest}`
        };
      }
      // Mentions: @user
      const mentionMatch = /^@([a-zA-Z0-9-_\.]+)/.exec(src);
      if (mentionMatch) {
        return {
          type: 'ggitSpecial',
          raw: mentionMatch[0],
          text: mentionMatch[0],
          href: `/${mentionMatch[1]}`
        };
      }
      // Issues: #123
      const issueMatch = /^#([0-9]+)/.exec(src);
      if (issueMatch && repoContext) {
        return {
          type: 'ggitSpecial',
          raw: issueMatch[0],
          text: issueMatch[0],
          href: `/${repoContext.owner}/${repoContext.name}/issues/${issueMatch[1]}`
        };
      }
      // Commits: 7-40 hex chars
      const commitMatch = /^[0-9a-f]{7,40}/.exec(src);
      if (commitMatch && repoContext) {
        return {
          type: 'ggitSpecial',
          raw: commitMatch[0],
          text: commitMatch[0].slice(0, 7),
          href: `/${repoContext.owner}/${repoContext.name}/commit/${commitMatch[0]}`
        };
      }
      return;
    },
    renderer(token: any) {
      return `<a href="${token.href}" class="ggit-link">${token.text}</a>`;
    }
  };

  marked.use({ renderer, extensions: [ggitExtension] });

  // Derived content
  let contentHtml = $derived.by(() => {
    if (!content) return '';
    try {
      const rawHtml = marked.parse(content) as string;
      return DOMPurify.sanitize(rawHtml, {
        ADD_TAGS: ['Warning', 'Note', 'Tip', 'Important', 'Caution', 'info', 'notice'],
        ADD_ATTR: ['type']
      });
    } catch (e) {
      return '<p class="error">Error rendering markdown</p>';
    }
  });
</script>

<div class="markdown-body">
  {@html contentHtml}
</div>

<style>
  .markdown-body {
    color: var(--md-sys-color-on-surface);
    line-height: 1.6;
    font-size: 1rem;
  }

  :global(.markdown-body h1, .markdown-body h2) {
    border-bottom: 1px solid var(--md-sys-color-outline-variant);
    padding-bottom: 0.3em;
    margin-top: 1.5em;
    margin-bottom: 1em;
  }

  :global(.markdown-body p) {
    margin-bottom: 1rem;
  }

  :global(.markdown-body ul, .markdown-body ol) {
    padding-left: 2rem;
    margin-bottom: 1rem;
  }

  :global(.markdown-body li) {
    margin-bottom: 0.25rem;
  }

  :global(.markdown-body a) {
    color: var(--md-sys-color-primary);
    text-decoration: none;
  }

  :global(.markdown-body a:hover) {
    text-decoration: underline;
  }

  :global(.markdown-body .ggit-link) {
    color: var(--md-sys-color-primary);
    font-weight: 600;
    text-decoration: none;
    background: var(--md-sys-color-primary-container, rgba(0,0,0,0.05));
    padding: 0 0.4rem;
    border-radius: 4px;
    font-size: 0.9em;
    display: inline-block;
    transition: all 0.2s;
  }

  :global(.markdown-body .ggit-link:hover) {
    background: var(--md-sys-color-primary);
    color: var(--md-sys-color-on-primary);
    text-decoration: none;
    transform: translateY(-1px);
  }

  :global(.markdown-body img) {
    max-width: 100%;
    height: auto;
    border-radius: 8px;
    margin: 1rem 0;
    box-shadow: 0 4px 12px rgba(0,0,0,0.05);
  }

  /* Code Block Styling - Integrated with Ggit Theme */
  :global(.markdown-body pre) {
    background-color: var(--md-sys-color-surface-container-low);
    padding: 1.5rem;
    border-radius: 16px;
    overflow-x: auto;
    margin: 1.5rem 0;
    border: 1px solid var(--md-sys-color-outline-variant);
    font-family: 'JetBrains Mono', 'Fira Code', ui-monospace, monospace;
    font-size: 0.85rem;
    line-height: 1.6;
    box-shadow: inset 0 1px 4px rgba(0,0,0,0.02);
  }

  :global(.markdown-body code.hljs) {
    padding: 0;
    background: transparent;
  }

  /* MDX-style Callouts */
  :global(.markdown-body Warning, .markdown-body Note, .markdown-body Tip, .markdown-body Important, .markdown-body Caution) {
    display: block;
    padding: 1rem 1.5rem;
    margin: 1.5rem 0;
    border-radius: 12px;
    border-left: 5px solid;
    font-size: 0.95rem;
    line-height: 1.6;
    position: relative;
    box-shadow: 0 4px 12px rgba(0,0,0,0.03);
  }

  :global(.markdown-body Warning::before, .markdown-body Note::before, .markdown-body Tip::before, .markdown-body Important::before, .markdown-body Caution::before) {
    content: '';
    display: block;
    font-weight: 800;
    text-transform: uppercase;
    font-size: 0.75rem;
    letter-spacing: 0.05em;
    margin-bottom: 0.5rem;
    opacity: 0.8;
  }

  :global(.markdown-body Warning) { 
    background: rgba(255, 152, 0, 0.08); 
    border-color: #ff9800; 
    color: #e65100;
  }
  :global(.markdown-body Warning::before) { content: 'Warning'; }

  :global(.markdown-body Note) { 
    background: rgba(33, 150, 243, 0.08); 
    border-color: #2196f3; 
    color: #0d47a1;
  }
  :global(.markdown-body Note::before) { content: 'Note'; }

  :global(.markdown-body Tip) { 
    background: rgba(76, 175, 80, 0.08); 
    border-color: #4caf50; 
    color: #1b5e20;
  }
  :global(.markdown-body Tip::before) { content: 'Tip'; }

  :global(.markdown-body Important) { 
    background: rgba(156, 39, 176, 0.08); 
    border-color: #9c27b0; 
    color: #4a148c;
  }
  :global(.markdown-body Important::before) { content: 'Important'; }

  :global(.markdown-body Caution) { 
    background: rgba(244, 67, 54, 0.08); 
    border-color: #f44336; 
    color: #b71c1c;
  }
  :global(.markdown-body Caution::before) { content: 'Caution'; }
</style>
