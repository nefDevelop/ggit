import { render, screen, fireEvent } from '@testing-library/svelte';
import { describe, it, expect } from 'vitest';
import Input from '../Input.svelte';

describe('Input', () => {
  it('renders input with default type text', () => {
    render(Input, { props: { label: 'Name' } });
    const input = screen.getByLabelText('Name');
    expect(input).toBeInTheDocument();
    expect(input).toHaveAttribute('type', 'text');
  });

  it('renders with placeholder', () => {
    render(Input, { props: { placeholder: 'Enter value', label: 'Field' } });
    const input = screen.getByLabelText('Field');
    expect(input).toHaveAttribute('placeholder', 'Enter value');
  });

  it('renders with label', () => {
    render(Input, { props: { label: 'Username' } });
    expect(screen.getByText('Username')).toBeInTheDocument();
  });

  it('shows error text when error prop is set', () => {
    render(Input, { props: { label: 'Email', error: 'Invalid email' } });
    expect(screen.getByText('Invalid email')).toBeInTheDocument();
  });

  it('disables input when disabled prop is true', () => {
    render(Input, { props: { label: 'Name', disabled: true } });
    const input = screen.getByLabelText('Name');
    expect(input).toBeDisabled();
  });

  it('makes input required when required prop is true', () => {
    render(Input, { props: { label: 'Name', required: true } });
    const input = screen.getByLabelText('Name');
    expect(input).toBeRequired();
  });

  it('renders textarea when multiline is true', () => {
    render(Input, { props: { label: 'Bio', multiline: true } });
    const textarea = screen.getByLabelText('Bio');
    expect(textarea.tagName).toBe('TEXTAREA');
  });

  it('updates value on input', () => {
    render(Input, { props: { label: 'Name' } });
    const input = screen.getByLabelText('Name') as HTMLInputElement;
    fireEvent.input(input, { target: { value: 'new value' } });
    expect(input.value).toBe('new value');
  });

  it('has error class on container when error is set', () => {
    render(Input, { props: { label: 'Name', error: 'Error' } });
    const container = document.querySelector('.has-error');
    expect(container).toBeInTheDocument();
  });

  it('applies custom class', () => {
    render(Input, { props: { label: 'Name', class: 'custom-input' } });
    const field = document.querySelector('.field.custom-input');
    expect(field).toBeInTheDocument();
  });
});
