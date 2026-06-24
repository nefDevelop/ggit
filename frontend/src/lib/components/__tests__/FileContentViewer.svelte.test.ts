import { render, screen } from '@testing-library/svelte';
import { describe, it, expect, vi } from 'vitest';
import FileContentViewerWrapper from './FileContentViewerWrapper.svelte';

const mockFormatSize = (n: number) => n >= 1024 ? `${(n / 1024).toFixed(1)} KB` : `${n} B`;
const mockTimeAgo = (d: any) => d ? '5m ago' : '';

const mockCodeRows = [
  { num: 1, content: '<span class="hljs-keyword">import</span> { test }', blame: null },
  { num: 2, content: 'const x = 1;', blame: { author: 'Alice', time: 1700000000 } },
  { num: 3, content: 'console.log(x);', blame: null },
];

const mockRepoData = {
  name: 'test.js',
  size: 2048,
  is_text_file: true,
  content: 'import { test }\nconst x = 1;\nconsole.log(x);\n',
  branch: 'main',
};

const defaultProps = {
  repoData: mockRepoData,
  username: 'testuser',
  reponame: 'test-repo',
  realPath: 'src/test.js',
  showBlame: false,
  blameLoading: false,
  sidebarCollapsed: false,
  codeRows: mockCodeRows,
  formatSize: mockFormatSize,
  timeAgo: mockTimeAgo,
  onToggleBlame: vi.fn(),
  onToggleSidebar: vi.fn(),
};

describe('FileContentViewer', () => {
  it('renders file name', () => {
    render(FileContentViewerWrapper, { props: defaultProps });
    expect(screen.getByText('test.js')).toBeInTheDocument();
  });

  it('renders file size', () => {
    render(FileContentViewerWrapper, { props: defaultProps });
    expect(screen.getByText('2.0 KB')).toBeInTheDocument();
  });

  it('renders code rows', () => {
    const { container } = render(FileContentViewerWrapper, { props: defaultProps });
    const codeLines = container.querySelectorAll('.code-line-row');
    expect(codeLines.length).toBe(3);
  });

  it('renders line numbers', () => {
    const { container } = render(FileContentViewerWrapper, { props: defaultProps });
    const lineCells = container.querySelectorAll('.ln-cell');
    expect(lineCells[0]).toHaveTextContent('1');
    expect(lineCells[1]).toHaveTextContent('2');
    expect(lineCells[2]).toHaveTextContent('3');
  });

  it('renders code content with HTML highlighting', () => {
    const { container } = render(FileContentViewerWrapper, { props: defaultProps });
    const codeCells = container.querySelectorAll('.code-cell');
    expect(codeCells[0].querySelector('.hljs-keyword')).toBeInTheDocument();
  });

  it('renders empty state when codeRows is empty', () => {
    const { container } = render(FileContentViewerWrapper, {
      props: { ...defaultProps, codeRows: [] },
    });
    expect(container.querySelectorAll('.code-line-row').length).toBe(0);
  });

  it('does not render blame cells when showBlame is false', () => {
    const { container } = render(FileContentViewerWrapper, { props: defaultProps });
    expect(container.querySelectorAll('.blame-cell').length).toBe(0);
  });

  it('renders blame cells when showBlame is true', () => {
    const { container } = render(FileContentViewerWrapper, {
      props: { ...defaultProps, showBlame: true },
    });
    expect(container.querySelectorAll('.blame-cell').length).toBeGreaterThan(0);
  });

  it('shows blame author when blame data exists', () => {
    render(FileContentViewerWrapper, {
      props: { ...defaultProps, showBlame: true },
    });
    expect(screen.getByText('Alice')).toBeInTheDocument();
  });

  it('renders binary placeholder for non-text files', () => {
    const binaryData = { ...mockRepoData, is_text_file: false, content: '' };
    render(FileContentViewerWrapper, {
      props: { ...defaultProps, repoData: binaryData, codeRows: [] },
    });
    expect(screen.getByText(/Binary file/)).toBeInTheDocument();
  });

  it('renders Edit button for text files', () => {
    render(FileContentViewerWrapper, { props: defaultProps });
    const editLink = screen.getByTitle('Edit');
    expect(editLink).toBeInTheDocument();
    expect(editLink).toHaveAttribute('href', '/testuser/test-repo/_edit/main/src/test.js');
  });

  it('renders Raw link', () => {
    render(FileContentViewerWrapper, { props: defaultProps });
    const rawLink = screen.getByTitle('Raw');
    expect(rawLink).toBeInTheDocument();
    expect(rawLink).toHaveAttribute('href', '/api/v2/testuser/test-repo/raw/main/src/test.js');
  });

  it('renders Download link', () => {
    render(FileContentViewerWrapper, { props: defaultProps });
    const downloadLink = screen.getByTitle('Download');
    expect(downloadLink).toBeInTheDocument();
  });

  it('renders README section when readme_content exists', () => {
    const withReadme = { ...mockRepoData, readme_content: '# Readme' };
    render(FileContentViewerWrapper, {
      props: { ...defaultProps, repoData: withReadme },
    });
    expect(screen.getByText('README Preview')).toBeInTheDocument();
  });

  it('does not render README section when readme_content is absent', () => {
    render(FileContentViewerWrapper, { props: defaultProps });
    expect(screen.queryByText('README Preview')).not.toBeInTheDocument();
  });
});
