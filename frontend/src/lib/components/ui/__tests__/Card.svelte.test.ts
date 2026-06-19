import { render, screen } from '@testing-library/svelte';
import { describe, it, expect } from 'vitest';
import CardWrapper from './CardWrapper.svelte';

describe('Card', () => {
  it('renders children content', () => {
    render(CardWrapper, { props: { childText: 'Card content' } });
    expect(screen.getByText('Card content')).toBeInTheDocument();
  });

  it('has card class', () => {
    render(CardWrapper, { props: { childText: 'Content' } });
    expect(document.querySelector('.card')).toBeInTheDocument();
  });

  it('applies default padding (md)', () => {
    render(CardWrapper, { props: { childText: 'Content' } });
    const card = document.querySelector('.card');
    expect(card).toHaveClass('p-md');
  });

  it('applies custom padding', () => {
    render(CardWrapper, { props: { padding: 'none', childText: 'Content' } });
    const card = document.querySelector('.card');
    expect(card).toHaveClass('p-0');
  });

  it('applies custom variant', () => {
    render(CardWrapper, { props: { variant: 'high', childText: 'Content' } });
    const card = document.querySelector('.card');
    expect(card).toHaveClass('bg-high');
  });

  it('has hoverable class when hover prop is true', () => {
    render(CardWrapper, { props: { hover: true, childText: 'Content' } });
    const card = document.querySelector('.card');
    expect(card).toHaveClass('hoverable');
  });

  it('has tonal class when tonal prop is true', () => {
    render(CardWrapper, { props: { tonal: true, childText: 'Content' } });
    const card = document.querySelector('.card');
    expect(card).toHaveClass('tonal');
  });

  it('applies custom class', () => {
    render(CardWrapper, { props: { class: 'extra-class', childText: 'Content' } });
    const card = document.querySelector('.card');
    expect(card).toHaveClass('extra-class');
  });
});
