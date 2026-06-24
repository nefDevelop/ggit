import { render, screen, fireEvent } from '@testing-library/svelte';
import { describe, it, expect, vi } from 'vitest';
import ViewHeaderWrapper from './ViewHeaderWrapper.svelte';

const mockTabs = [
  { id: 'activity', label: 'Activity', icon: vi.fn() },
  { id: 'alphabetical', label: 'Alphabetical', icon: vi.fn() },
  { id: 'folders', label: 'Folders', icon: vi.fn() },
];

describe('ViewHeader', () => {
  it('renders search input', () => {
    const { container } = render(ViewHeaderWrapper, { props: { showSearch: true } });
    const input = container.querySelector('input');
    expect(input).toBeInTheDocument();
  });

  it('hides search when showSearch is false', () => {
    const { container } = render(ViewHeaderWrapper, { props: { showSearch: false } });
    const input = container.querySelector('input');
    expect(input).not.toBeInTheDocument();
  });

  it('renders tab pills when tabs are provided', () => {
    render(ViewHeaderWrapper, { props: { tabs: mockTabs } });
    expect(screen.getByText('Activity')).toBeInTheDocument();
    expect(screen.getByText('Alphabetical')).toBeInTheDocument();
    expect(screen.getByText('Folders')).toBeInTheDocument();
  });

  it('renders tab pills as buttons', () => {
    render(ViewHeaderWrapper, { props: { tabs: mockTabs } });
    const buttons = screen.getAllByRole('button');
    expect(buttons.length).toBe(3);
  });

  it('marks the active tab', () => {
    render(ViewHeaderWrapper, { props: { tabs: mockTabs, activeTab: 'alphabetical' } });
    const activePill = screen.getByText('Alphabetical')?.closest('button');
    const inactivePill = screen.getByText('Activity')?.closest('button');
    expect(activePill).toHaveClass('active');
    expect(inactivePill).not.toHaveClass('active');
  });

  it('updates active tab on click', async () => {
    let active = 'activity';
    const { component } = render(ViewHeaderWrapper, {
      props: { tabs: mockTabs, activeTab: active },
    });
    const alphabeticalBtn = screen.getByText('Alphabetical');
    await fireEvent.click(alphabeticalBtn);
  });

  it('does not render tabs section when tabs array is empty', () => {
    const { container } = render(ViewHeaderWrapper, { props: { tabs: [] } });
    expect(container.querySelector('.view-mode-selector')).not.toBeInTheDocument();
  });

  it('applies no-pills class when noPills is true', () => {
    render(ViewHeaderWrapper, { props: { tabs: mockTabs, noPills: true } });
    const pills = screen.getAllByRole('button');
    pills.forEach(pill => {
      expect(pill).toHaveClass('no-pills');
    });
  });

  it('renders with correct structure', () => {
    const { container } = render(ViewHeaderWrapper, { props: { tabs: mockTabs, showSearch: true } });
    expect(container.querySelector('.view-header-abstract')).toBeInTheDocument();
    expect(container.querySelector('.search-and-action-row')).toBeInTheDocument();
    expect(container.querySelector('.view-mode-selector-row')).toBeInTheDocument();
  });

  it('sets active tab class on the correct button', () => {
    render(ViewHeaderWrapper, { props: { tabs: mockTabs, activeTab: 'folders' } });
    const foldersBtn = screen.getByText('Folders')?.closest('button');
    expect(foldersBtn).toHaveClass('active');
    expect(screen.getByText('Activity')?.closest('button')).not.toHaveClass('active');
    expect(screen.getByText('Alphabetical')?.closest('button')).not.toHaveClass('active');
  });
});
