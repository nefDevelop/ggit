import { render, screen } from '@testing-library/svelte';
import { describe, it, expect } from 'vitest';
import BadgeWrapper from './BadgeWrapper.svelte';

describe('Badge', () => {
  it('renders children content', () => {
    render(BadgeWrapper, { props: { childText: 'Status' } });
    expect(screen.getByText('Status')).toBeInTheDocument();
  });

  it('has badge class', () => {
    render(BadgeWrapper, { props: { childText: 'Badge' } });
    expect(document.querySelector('.badge')).toBeInTheDocument();
  });

  it('applies default variant (neutral) and size (md)', () => {
    render(BadgeWrapper, { props: { childText: 'Default' } });
    const badge = document.querySelector('.badge');
    expect(badge).toHaveClass('badge-neutral', 'badge-md');
  });

  it('applies custom variant', () => {
    render(BadgeWrapper, { props: { variant: 'success', childText: 'Success' } });
    expect(document.querySelector('.badge')).toHaveClass('badge-success');
  });

  it('applies custom size', () => {
    render(BadgeWrapper, { props: { size: 'sm', childText: 'Small' } });
    expect(document.querySelector('.badge')).toHaveClass('badge-sm');
  });

  it('applies custom class', () => {
    render(BadgeWrapper, { props: { class: 'custom-badge', childText: 'Styled' } });
    expect(document.querySelector('.badge')).toHaveClass('custom-badge');
  });
});
