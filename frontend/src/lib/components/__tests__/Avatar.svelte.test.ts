import { render } from '@testing-library/svelte';
import { describe, it, expect } from 'vitest';
import AvatarWrapper from './AvatarWrapper.svelte';

describe('Avatar', () => {
  it('renders an img element when src is provided', () => {
    const { container } = render(AvatarWrapper, { props: { src: 'https://example.com/avatar.png', username: 'user' } });
    const img = container.querySelector('img');
    expect(img).toBeInTheDocument();
    expect(img).toHaveAttribute('src', 'https://example.com/avatar.png');
  });

  it('falls back to Identicon when src is empty', () => {
    const { container } = render(AvatarWrapper, { props: { src: '', username: 'testuser', type: 'user' } });
    const svg = container.querySelector('svg');
    expect(svg).toBeInTheDocument();
    expect(container.querySelector('img')).not.toBeInTheDocument();
  });

  it('falls back to Identicon when src is null', () => {
    const { container } = render(AvatarWrapper, { props: { username: 'nobody', type: 'user' } });
    const svg = container.querySelector('svg');
    expect(svg).toBeInTheDocument();
  });

  it('renders with repo type', () => {
    const { container } = render(AvatarWrapper, { props: { username: 'my-repo', type: 'repo', seed: 'my-repo' } });
    const svg = container.querySelector('svg.repo-type');
    expect(svg).toBeInTheDocument();
  });

  it('uses seed over username for identicon fallback', () => {
    const { container: c1 } = render(AvatarWrapper, { props: { username: 'repo', seed: 'seed-A', type: 'repo' } });
    const { container: c2 } = render(AvatarWrapper, { props: { username: 'repo', seed: 'seed-B', type: 'repo' } });
    expect(c1.querySelector('svg')?.innerHTML).not.toBe(c2.querySelector('svg')?.innerHTML);
  });

  it('uses username as fallback seed when seed is empty', () => {
    const { container: c1 } = render(AvatarWrapper, { props: { username: 'unique-name', type: 'repo' } });
    const { container: c2 } = render(AvatarWrapper, { props: { username: 'unique-name', type: 'repo' } });
    expect(c1.querySelector('svg')?.innerHTML).toBe(c2.querySelector('svg')?.innerHTML);
  });

  it('applies repo-avatar class for repo type', () => {
    const { container } = render(AvatarWrapper, { props: { username: 'repo', type: 'repo' } });
    const avatarContainer = container.querySelector('.avatar-container');
    expect(avatarContainer).toHaveClass('repo-avatar');
  });

  it('does not apply repo-avatar class for user type', () => {
    const { container } = render(AvatarWrapper, { props: { username: 'user', type: 'user' } });
    const avatarContainer = container.querySelector('.avatar-container');
    expect(avatarContainer).not.toHaveClass('repo-avatar');
  });

  it('applies custom class name', () => {
    const { container } = render(AvatarWrapper, { props: { username: 'user', type: 'user', class: 'custom-class' } });
    const avatarContainer = container.querySelector('.avatar-container');
    expect(avatarContainer).toHaveClass('custom-class');
  });

  it('sets correct size style', () => {
    const { container } = render(AvatarWrapper, { props: { src: 'https://example.com/avatar.png', username: 'user', size: 48 } });
    const avatarContainer = container.querySelector('.avatar-container');
    expect(avatarContainer).toHaveStyle('max-width: 48px');
  });
});
