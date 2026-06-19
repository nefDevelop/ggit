import { render, screen, fireEvent } from '@testing-library/svelte';
import { describe, it, expect, vi } from 'vitest';
import ButtonWrapper from './ButtonWrapper.svelte';

describe('Button', () => {
  it('renders with default variant and size', () => {
    render(ButtonWrapper, { props: { childText: 'Click me' } });
    const btn = screen.getByRole('button', { name: /click me/i });
    expect(btn).toBeInTheDocument();
    expect(btn).toHaveClass('btn-primary', 'btn-md');
  });

  it('renders with custom variant and size', () => {
    render(ButtonWrapper, { props: { variant: 'danger', size: 'lg', childText: 'Delete' } });
    const btn = screen.getByRole('button', { name: /delete/i });
    expect(btn).toHaveClass('btn-danger', 'btn-lg');
  });

  it('renders as link when href is provided', () => {
    render(ButtonWrapper, { props: { href: '/test', childText: 'Link' } });
    const link = screen.getByRole('link', { name: /link/i });
    expect(link).toBeInTheDocument();
    expect(link).toHaveAttribute('href', '/test');
  });

  it('disables button when disabled prop is true', () => {
    render(ButtonWrapper, { props: { disabled: true, childText: 'Disabled' } });
    const btn = screen.getByRole('button', { name: /disabled/i });
    expect(btn).toBeDisabled();
  });

  it('shows spinner and disables button when loading', () => {
    render(ButtonWrapper, { props: { loading: true } });
    const btn = screen.getByRole('button');
    expect(btn).toBeDisabled();
    expect(btn.querySelector('.spinner')).toBeInTheDocument();
  });

  it('calls onclick when clicked', () => {
    const onclick = vi.fn();
    render(ButtonWrapper, { props: { onclick, childText: 'Click' } });
    fireEvent.click(screen.getByRole('button', { name: /click/i }));
    expect(onclick).toHaveBeenCalledOnce();
  });

  it('does not render children when loading', () => {
    render(ButtonWrapper, { props: { loading: true, childText: 'Hidden' } });
    expect(screen.getByRole('button')).not.toHaveTextContent(/hidden/i);
  });

  it('applies custom class', () => {
    render(ButtonWrapper, { props: { class: 'custom-class', childText: 'Styled' } });
    const btn = screen.getByRole('button', { name: /styled/i });
    expect(btn).toHaveClass('custom-class');
  });
});
