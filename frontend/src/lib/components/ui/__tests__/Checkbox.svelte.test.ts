import { render, screen, fireEvent } from '@testing-library/svelte';
import { describe, it, expect, vi } from 'vitest';
import Checkbox from '../Checkbox.svelte';

describe('Checkbox', () => {
  it('renders with label', () => {
    render(Checkbox, { props: { checked: false, label: 'Accept terms' } });
    expect(screen.getByText('Accept terms')).toBeInTheDocument();
  });

  it('renders with description', () => {
    render(Checkbox, { props: { checked: false, label: 'Option', description: 'This is an option' } });
    expect(screen.getByText('This is an option')).toBeInTheDocument();
  });

  it('shows checkmark when checked is true', () => {
    render(Checkbox, { props: { checked: true, label: 'Enabled' } });
    const checkmark = document.querySelector('.checkmark.checked');
    expect(checkmark).toBeInTheDocument();
  });

  it('does not show checkmark when checked is false', () => {
    render(Checkbox, { props: { checked: false, label: 'Disabled' } });
    const checkmark = document.querySelector('.checkmark.checked');
    expect(checkmark).not.toBeInTheDocument();
  });

  it('toggles checked state on click', () => {
    render(Checkbox, { props: { checked: false, label: 'Toggle' } });
    const checkbox = screen.getByLabelText('Toggle') as HTMLInputElement;
    expect(checkbox.checked).toBe(false);
    fireEvent.click(checkbox);
    expect(checkbox.checked).toBe(true);
  });

  it('calls onchange with new value', () => {
    const onchange = vi.fn();
    render(Checkbox, { props: { checked: false, label: 'Test', onchange } });
    fireEvent.click(screen.getByLabelText('Test'));
    expect(onchange).toHaveBeenCalledWith(true);
  });

  it('does not render description when simple is true', () => {
    render(Checkbox, { props: { checked: false, label: 'Simple', description: 'Hidden', simple: true } });
    expect(screen.queryByText('Hidden')).not.toBeInTheDocument();
  });

  it('has simple class when simple prop is true', () => {
    render(Checkbox, { props: { checked: false, label: 'Simple', simple: true } });
    const container = document.querySelector('.checkbox-container.simple');
    expect(container).toBeInTheDocument();
  });
});
