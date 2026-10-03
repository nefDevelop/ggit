import { vi } from 'vitest';
import '@testing-library/jest-dom/vitest';

const mockPage = {
  url: new URL('http://localhost:3005'),
  params: {},
  route: { id: null },
  status: 200,
  error: null,
  data: null,
};

vi.mock('$app/navigation', () => ({
  goto: vi.fn(),
}));

vi.mock('$app/state', () => ({
  page: mockPage,
}));

vi.mock('$app/stores', () => {
  const { writable } = require('svelte/store');
  return {
    page: writable(mockPage),
  };
});

vi.mock('$app/environment', () => ({
  browser: false,
  dev: true,
  building: false,
}));
