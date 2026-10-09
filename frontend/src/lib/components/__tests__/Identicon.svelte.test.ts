import { render } from '@testing-library/svelte';
import { describe, it, expect } from 'vitest';
import IdenticonWrapper from './IdenticonWrapper.svelte';

describe('Identicon', () => {
  it('renders an SVG element', () => {
    const { container } = render(IdenticonWrapper, { props: { seed: 'test' } });
    const svg = container.querySelector('svg.identicon');
    expect(svg).toBeInTheDocument();
  });

  it('renders with repo type', () => {
    const { container } = render(IdenticonWrapper, { props: { seed: 'my-repo', type: 'repo' } });
    const svg = container.querySelector('svg.repo-type');
    expect(svg).toBeInTheDocument();
  });

  it('renders with user type', () => {
    const { container } = render(IdenticonWrapper, { props: { seed: 'user', type: 'user' } });
    const svg = container.querySelector('svg.identicon');
    expect(svg).not.toHaveClass('repo-type');
  });

  it('produces deterministic output for same seed', () => {
    const { container: c1 } = render(IdenticonWrapper, { props: { seed: 'same-repo', type: 'repo' } });
    const { container: c2 } = render(IdenticonWrapper, { props: { seed: 'same-repo', type: 'repo' } });
    expect(c1.querySelector('svg')?.innerHTML).toBe(c2.querySelector('svg')?.innerHTML);
  });

  it('produces different output for different seeds', () => {
    const { container: c1 } = render(IdenticonWrapper, { props: { seed: 'repo-one', type: 'repo' } });
    const { container: c2 } = render(IdenticonWrapper, { props: { seed: 'repo-two', type: 'repo' } });
    expect(c1.querySelector('svg')?.innerHTML).not.toBe(c2.querySelector('svg')?.innerHTML);
  });

  it('renders repo type with rect elements (pixel art)', () => {
    const { container } = render(IdenticonWrapper, { props: { seed: 'pixel-repo', type: 'repo' } });
    const rects = container.querySelectorAll('rect');
    expect(rects.length).toBeGreaterThan(0);
  });

  it('renders user type with circles', () => {
    const { container } = render(IdenticonWrapper, { props: { seed: 'organic-user', type: 'user' } });
    const circles = container.querySelectorAll('circle');
    expect(circles.length).toBeGreaterThan(0);
  });

  it('handles empty seed gracefully', () => {
    const { container } = render(IdenticonWrapper, { props: { seed: '', type: 'repo' } });
    const svg = container.querySelector('svg');
    expect(svg).toBeInTheDocument();
  });

  it('respects custom size', () => {
    const { container } = render(IdenticonWrapper, { props: { seed: 'size-test', size: 64 } });
    const svg = container.querySelector('svg');
    expect(svg).toHaveAttribute('width', '64');
    expect(svg).toHaveAttribute('height', '64');
  });

  it('has a clipPath definition', () => {
    const { container } = render(IdenticonWrapper, { props: { seed: 'clip', type: 'repo' } });
    const clip = container.querySelector('clipPath');
    expect(clip).toBeInTheDocument();
  });

  it('falls back to palette when langColor is not a hex color', () => {
    // getLanguageColor devuelve "var(--...)" si el repo no tiene lenguaje.
    // Antes daba NaN: fondo transparente + bloques blancos invisibles en claro.
    const { container } = render(IdenticonWrapper, {
      props: { seed: 'infra', type: 'repo', langColor: 'var(--text-muted)' },
    });
    const svg = container.querySelector('svg') as SVGElement;
    expect(svg.style.background).not.toContain('NaN');
    expect(svg.style.background).toMatch(/^rgba\(\d+, \d+, \d+, 0\.15\)$/);
    const fills = [...container.querySelectorAll('g rect')].map((r) =>
      r.getAttribute('fill'),
    );
    expect(fills.length).toBeGreaterThan(0);
    for (const fill of fills) {
      expect(fill).toMatch(/^#[0-9a-fA-F]{6}$/);
    }
  });

  it('uses a valid hex langColor for background and blocks', () => {
    const { container } = render(IdenticonWrapper, {
      props: { seed: 'api-ejemplo', type: 'repo', langColor: '#00add8' },
    });
    const svg = container.querySelector('svg') as SVGElement;
    expect(svg.style.background).toBe('rgba(0, 173, 216, 0.7)');
  });
});
