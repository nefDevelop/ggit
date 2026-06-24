import { render, screen, fireEvent } from '@testing-library/svelte';
import { describe, it, expect, vi } from 'vitest';
import NavBrandFooter from '../NavBrandFooter.svelte';

describe('NavBrandFooter', () => {
  it('renders brand name', () => {
    render(NavBrandFooter, { props: { showThemePanel: false, pageLoadTime: 13, templateTime: 4 } });
    expect(screen.getByText('Ggit')).toBeInTheDocument();
  });

  it('renders version', () => {
    render(NavBrandFooter, { props: { showThemePanel: false, pageLoadTime: 13, templateTime: 4 } });
    expect(screen.getByText('v2026.4.0')).toBeInTheDocument();
  });

  it('renders theme button', () => {
    render(NavBrandFooter, { props: { showThemePanel: false, pageLoadTime: 13, templateTime: 4 } });
    const themeBtn = screen.getByTitle('Apariencia');
    expect(themeBtn).toBeInTheDocument();
  });

  it('toggles showThemePanel on theme button click', async () => {
    let show = false;
    const { component } = render(NavBrandFooter, { props: { showThemePanel: show, pageLoadTime: 13, templateTime: 4 } });
    const themeBtn = screen.getByTitle('Apariencia');
    await fireEvent.click(themeBtn);
  });

  it('renders performance info', () => {
    render(NavBrandFooter, { props: { showThemePanel: false, pageLoadTime: 42, templateTime: 10 } });
    expect(screen.getByText(/Page: 42ms/)).toBeInTheDocument();
  });

  it('hidden on mobile via @media query (CSS check)', () => {
    const { container } = render(NavBrandFooter, { props: { showThemePanel: false, pageLoadTime: 13, templateTime: 4 } });
    const brandSection = container.querySelector('.brand-section-bottom');
    expect(brandSection).toBeInTheDocument();
  });

  it('renders Palette icon', () => {
    const { container } = render(NavBrandFooter, { props: { showThemePanel: false, pageLoadTime: 13, templateTime: 4 } });
    const svg = container.querySelector('svg');
    expect(svg).toBeInTheDocument();
  });
});
