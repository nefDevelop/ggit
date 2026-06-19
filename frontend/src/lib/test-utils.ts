import { render as testingRender } from '@testing-library/svelte';
import type { RenderOptions, RenderResult } from '@testing-library/svelte';
import type { Component } from 'svelte';
import { setContext, getContext } from 'svelte';

type ComponentType = Component<any>;

interface RenderWithContextOptions extends RenderOptions {
  contexts?: [string, any][];
}

export function render<T extends Record<string, any>>(
  component: ComponentType,
  props?: T,
  options?: RenderWithContextOptions,
): RenderResult {
  const { contexts = [], ...renderOptions } = options || {};

  const Wrapper = createWrapper(component, props, contexts);

  return testingRender(Wrapper, renderOptions);
}

function createWrapper(
  Inner: ComponentType,
  props: any,
  contexts: [string, any][],
): ComponentType {
  return function Wrapper() {
    for (const [key, value] of contexts) {
      setContext(key, value);
    }

    return Inner as any;
  } as unknown as ComponentType;
}

export function mockApiFetch(overrides: Record<string, any> = {}) {
  const defaultMock = {
    ok: true,
    status: 200,
    json: vi.fn().mockResolvedValue({}),
    text: vi.fn().mockResolvedValue(''),
    headers: new Headers({ 'content-type': 'application/json' }),
  };

  vi.mock('$lib/api', () => ({
    apiFetch: vi.fn((url: string, _options?: any) => {
      const override = overrides[url];
      if (override) {
        return Promise.resolve({ ...defaultMock, ...override });
      }
      return Promise.resolve(defaultMock);
    }),
  }));
}

export const mockRepoContext = {
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
