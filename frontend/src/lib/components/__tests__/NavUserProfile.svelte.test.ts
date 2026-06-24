import { render, screen, fireEvent } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import NavUserProfileWrapper from './NavUserProfileWrapper.svelte';

const mockAuthState = vi.hoisted(() => ({
  isAuthenticated: true,
  user: {
    avatar_url: '',
    username: 'testuser',
    full_name: 'Test User',
    is_admin: false,
  },
}));

const mockToast = vi.hoisted(() => ({
  showNotices: false,
  hasUnread: false,
  clearUnread: vi.fn(),
}));

const mockApiFetch = vi.hoisted(() => vi.fn());

vi.mock('$lib/auth.svelte', () => ({ authState: mockAuthState }));
vi.mock('$lib/toast.svelte', () => ({ toast: mockToast }));
vi.mock('$lib/api', () => ({ apiFetch: mockApiFetch }));

describe('NavUserProfile', () => {
  beforeEach(() => {
    vi.clearAllMocks();
    mockAuthState.isAuthenticated = true;
    mockAuthState.user = { avatar_url: '', username: 'testuser', full_name: 'Test User', is_admin: false };
    mockToast.showNotices = false;
    mockToast.hasUnread = false;
  });

  it('renders when user is authenticated', () => {
    render(NavUserProfileWrapper);
    expect(screen.getByText('Test User')).toBeInTheDocument();
    expect(screen.getByText('@testuser')).toBeInTheDocument();
  });

  it('renders logout button', () => {
    render(NavUserProfileWrapper);
    const logoutBtn = screen.getByTitle('Logout');
    expect(logoutBtn).toBeInTheDocument();
  });

  it('uses full_name when available', () => {
    render(NavUserProfileWrapper);
    expect(screen.getByText('Test User')).toBeInTheDocument();
  });

  it('falls back to username when full_name is empty', () => {
    mockAuthState.user = { avatar_url: '', username: 'testuser', full_name: '', is_admin: false };
    render(NavUserProfileWrapper);
    expect(screen.getByText('testuser')).toBeInTheDocument();
  });

  it('shows unread dot when hasUnread is true', () => {
    mockToast.hasUnread = true;
    const { container } = render(NavUserProfileWrapper);
    expect(container.querySelector('.unread-dot')).toBeInTheDocument();
  });

  it('hides unread dot when hasUnread is false', () => {
    mockToast.hasUnread = false;
    const { container } = render(NavUserProfileWrapper);
    expect(container.querySelector('.unread-dot')).not.toBeInTheDocument();
  });

  it('toggles notices on click', async () => {
    mockToast.showNotices = false;
    render(NavUserProfileWrapper);
    const profile = document.querySelector('.user-profile-summary');
    expect(profile).toBeInTheDocument();
    await fireEvent.click(profile!);
    expect(mockToast.clearUnread).toHaveBeenCalled();
    expect(mockToast.showNotices).toBe(true);
  });

  it('does not render when user is not authenticated', () => {
    mockAuthState.isAuthenticated = false;
    const { container } = render(NavUserProfileWrapper);
    expect(container.querySelector('.user-profile-summary')).not.toBeInTheDocument();
  });

  it('calls logout API on logout button click', async () => {
    mockApiFetch.mockResolvedValueOnce({ ok: true });
    render(NavUserProfileWrapper);
    const logoutBtn = screen.getByTitle('Logout');
    await fireEvent.click(logoutBtn);
    expect(mockApiFetch).toHaveBeenCalledWith('/api/v2/user/logout', { method: 'POST' });
  });

  it('renders avatar ring', () => {
    const { container } = render(NavUserProfileWrapper);
    expect(container.querySelector('.avatar-ring')).toBeInTheDocument();
  });
});
