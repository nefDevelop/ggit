import { render, screen, fireEvent, waitFor } from '@testing-library/svelte';
import { describe, it, expect, vi, beforeEach } from 'vitest';
import RepoSettingsGeneralWrapper from './RepoSettingsGeneralWrapper.svelte';

const mockSettings = {
  id: 1,
  name: 'test-repo',
  description: 'Original description',
  website: '',
  private: false,
  unlisted: false,
  archived: false,
  enable_wiki: true,
  enable_issues: true,
  enable_pulls: true,
  enable_external_wiki: false,
  external_wiki_url: '',
  allow_public_wiki: false,
  enable_external_tracker: false,
  external_tracker_url: '',
  project_folder: '',
  wiki_folder_name: 'docs',
  default_branch: 'main',
  avatar_url: null,
  owner: { username: 'testuser' },
};

const repoContext = {
  data: {
    owner: { username: 'testuser' },
    name: 'test-repo',
    id: 1,
    private: false,
    description: 'A test repo',
  },
  full: null,
  entries: [],
  stats: { commits: 0, branches: 0, tags: 0, releases: 0 },
  languages: null,
  loading: false,
};

const mockApiFetch = vi.hoisted(() => vi.fn());
const mockToast = vi.hoisted(() => ({
  success: vi.fn(),
  error: vi.fn(),
  show: vi.fn(),
  info: vi.fn(),
  warning: vi.fn(),
}));
const mockDialog = vi.hoisted(() => ({
  confirm: vi.fn(async () => true),
  prompt: vi.fn(async () => null),
  alert: vi.fn(async () => undefined),
}));
const mockGoto = vi.hoisted(() => vi.fn());

vi.mock('$lib/api', () => ({ apiFetch: mockApiFetch }));
vi.mock('$lib/toast.svelte', () => ({ toast: mockToast }));
vi.mock('$lib/dialog.svelte', () => ({ dialog: mockDialog }));
vi.mock('$app/navigation', () => ({ goto: mockGoto }));

function setupDefaultMocks() {
  mockApiFetch.mockImplementation((_url: string, options?: any) => {
    if (!options || options.method === 'GET' || !options.method) {
      return Promise.resolve({
        ok: true,
        status: 200,
        json: async () => mockSettings,
        headers: new Headers({ 'content-type': 'application/json' }),
        text: async () => JSON.stringify(mockSettings),
      });
    }
    if (options.method === 'PATCH') {
      return Promise.resolve({
        ok: true,
        status: 200,
        json: async () => ({}),
        headers: new Headers({ 'content-type': 'application/json' }),
        text: async () => '{}',
      });
    }
    return Promise.resolve({
      ok: false,
      status: 404,
      json: async () => ({ message: 'Not found' }),
      headers: new Headers({ 'content-type': 'application/json' }),
      text: async () => 'Not found',
    });
  });
}

beforeEach(() => {
  vi.clearAllMocks();
  setupDefaultMocks();
  mockDialog.confirm.mockResolvedValue(true);
});

describe('RepoSettingsGeneral', () => {
  it('renders settings form after loading with repo data', async () => {
    render(RepoSettingsGeneralWrapper, { props: { repoContext } });

    expect(await screen.findByLabelText('Nombre del Repositorio')).toHaveValue('test-repo');
    expect(await screen.findByLabelText('Descripción')).toHaveValue('Original description');
  });

  it('renders all form fields', async () => {
    render(RepoSettingsGeneralWrapper, { props: { repoContext } });

    expect(await screen.findByLabelText('Nombre del Repositorio')).toBeInTheDocument();
    expect(await screen.findByLabelText('Descripción')).toBeInTheDocument();
    expect(await screen.findByLabelText('Página Web')).toBeInTheDocument();
    expect(await screen.findByText('Guardar Cambios Generales')).toBeInTheDocument();
  });

  it('renders danger zone', async () => {
    render(RepoSettingsGeneralWrapper, { props: { repoContext } });

    expect(await screen.findByText('Zona de Peligro')).toBeInTheDocument();
    expect(await screen.findByText('Eliminar este repositorio')).toBeInTheDocument();
    expect(await screen.findByText('Transferir Propiedad')).toBeInTheDocument();
  });

  it('fetches settings with correct URL', async () => {
    render(RepoSettingsGeneralWrapper, { props: { repoContext } });

    await waitFor(() => {
      const settingsCall = mockApiFetch.mock.calls.find((c: any[]) =>
        typeof c[0] === 'string' && c[0].includes('/settings'),
      );
      expect(settingsCall).toBeDefined();
      expect(settingsCall![0]).toContain('/testuser/test-repo/settings');
    });
  });

  it('saves description without changing name', async () => {
    render(RepoSettingsGeneralWrapper, { props: { repoContext } });

    await screen.findByLabelText('Nombre del Repositorio');
    const descInput = screen.getByLabelText('Descripción');
    await fireEvent.input(descInput, { target: { value: 'Updated description' } });

    fireEvent.click(screen.getByText('Guardar Cambios Generales'));

    await waitFor(() => {
      const patch = mockApiFetch.mock.calls.find(
        (c: any[]) => typeof c[0] === 'string' && c[0].includes('/settings') && c[1]?.method === 'PATCH',
      );
      expect(patch).toBeDefined();
      const body = JSON.parse(patch![1].body);
      expect(body.description).toBe('Updated description');
      expect(body.name).toBe('test-repo');
    });
  });

  it('shows success toast on save', async () => {
    render(RepoSettingsGeneralWrapper, { props: { repoContext } });

    await screen.findByText('Guardar Cambios Generales');
    fireEvent.click(screen.getByText('Guardar Cambios Generales'));

    await waitFor(() => {
      expect(mockToast.success).toHaveBeenCalled();
    });
  });

  it('does not show success toast when save fails', async () => {
    mockApiFetch.mockImplementation((url: string, options?: any) => {
      if (options?.method === 'PATCH') {
        return Promise.resolve({
          ok: false,
          status: 500,
          json: async () => ({ message: 'Server error' }),
          headers: new Headers({ 'content-type': 'application/json' }),
          text: async () => 'Server error',
        });
      }
      return Promise.resolve({
        ok: true,
        status: 200,
        json: async () => mockSettings,
        headers: new Headers({ 'content-type': 'application/json' }),
        text: async () => JSON.stringify(mockSettings),
      });
    });

    render(RepoSettingsGeneralWrapper, { props: { repoContext } });

    await screen.findByText('Guardar Cambios Generales');
    fireEvent.click(screen.getByText('Guardar Cambios Generales'));

    await waitFor(() => {
      expect(mockToast.success).not.toHaveBeenCalled();
    });
  });

  it('sends the confirmation name in the DELETE body', async () => {
    // El handler DeleteRepository exige {"name": "<repo>"} y responde 400 sin cuerpo.
    mockApiFetch.mockImplementation((_url: string, options?: any) => {
      if (options?.method === 'DELETE') {
        return Promise.resolve({
          ok: true,
          status: 200,
          json: async () => ({ status: 'ok' }),
          headers: new Headers({ 'content-type': 'application/json' }),
          text: async () => '{"status":"ok"}',
        });
      }
      return Promise.resolve({
        ok: true,
        status: 200,
        json: async () => mockSettings,
        headers: new Headers({ 'content-type': 'application/json' }),
        text: async () => JSON.stringify(mockSettings),
      });
    });

    render(RepoSettingsGeneralWrapper, { props: { repoContext } });

    const confirmInput = await screen.findByPlaceholderText(
      'Escribe el nombre del repositorio para confirmar',
    );
    await fireEvent.input(confirmInput, { target: { value: 'test-repo' } });

    const deleteBtn = await screen.findByText('Eliminar Repositorio');
    await waitFor(() => expect(deleteBtn).not.toBeDisabled());
    await fireEvent.click(deleteBtn);

    await waitFor(() => {
      const del = mockApiFetch.mock.calls.find(
        (c: any[]) => typeof c[0] === 'string' && c[1]?.method === 'DELETE' && c[0].includes('/settings'),
      );
      expect(del).toBeTruthy();
      expect(del![0]).toContain('/testuser/test-repo/settings');
      expect(JSON.parse(del![1].body)).toEqual({ name: 'test-repo' });
    });

    expect(mockGoto).toHaveBeenCalledWith('/dashboard');
  });

  it('does not call DELETE when the confirmation name does not match', async () => {
    render(RepoSettingsGeneralWrapper, { props: { repoContext } });

    const confirmInput = await screen.findByPlaceholderText(
      'Escribe el nombre del repositorio para confirmar',
    );
    await fireEvent.input(confirmInput, { target: { value: 'otro-repo' } });

    const deleteBtn = await screen.findByText('Eliminar Repositorio');
    expect(deleteBtn).toBeDisabled();
    await fireEvent.click(deleteBtn);

    expect(
      mockApiFetch.mock.calls.some((c: any[]) => c[1]?.method === 'DELETE'),
    ).toBe(false);
  });
});
