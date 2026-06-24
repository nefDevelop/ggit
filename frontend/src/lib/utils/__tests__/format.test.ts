import { describe, it, expect } from 'vitest';

function formatSize(bytes: number) {
  if (bytes === 0) return '0 B';
  const k = 1024;
  const sizes = ['B', 'KB', 'MB', 'GB'];
  const i = Math.floor(Math.log(bytes) / Math.log(k));
  return parseFloat((bytes / Math.pow(k, i)).toFixed(1)) + ' ' + sizes[i];
}

function timeAgo(dateInput: any) {
  if (!dateInput) return '';
  const date = new Date(dateInput);
  const seconds = Math.floor((Date.now() - date.getTime()) / 1000);
  if (seconds < 60) return 'just now';
  const minutes = Math.floor(seconds / 60);
  if (minutes < 60) return `${minutes}m ago`;
  const hours = Math.floor(minutes / 60);
  if (hours < 24) return `${hours}h ago`;
  const days = Math.floor(hours / 24);
  if (days < 30) return `${days}d ago`;
  return date.toLocaleDateString();
}

function getHash(str: string) {
  if (!str) return 0;
  let hash = 0;
  for (let i = 0; i < str.length; i++) {
    hash = str.charCodeAt(i) + ((hash << 5) - hash);
  }
  return Math.abs(hash);
}

describe('formatSize', () => {
  it('returns "0 B" for 0 bytes', () => {
    expect(formatSize(0)).toBe('0 B');
  });

  it('formats bytes', () => {
    expect(formatSize(500)).toBe('500 B');
  });

  it('formats KB', () => {
    const result = formatSize(2048);
    expect(result).toMatch(/^2[\.\s]/);
    expect(result).toContain('KB');
  });

  it('formats MB', () => {
    const result = formatSize(1048576);
    expect(result).toMatch(/^1[\.\s]/);
    expect(result).toContain('MB');
  });

  it('formats GB', () => {
    const result = formatSize(1073741824);
    expect(result).toMatch(/^1[\.\s]/);
    expect(result).toContain('GB');
  });

  it('uses one decimal place', () => {
    expect(formatSize(1536)).toBe('1.5 KB');
  });
});

describe('timeAgo', () => {
  it('returns empty string for null input', () => {
    expect(timeAgo(null)).toBe('');
  });

  it('returns empty string for undefined input', () => {
    expect(timeAgo(undefined)).toBe('');
  });

  it('returns "just now" for recent timestamps', () => {
    expect(timeAgo(new Date().toISOString())).toBe('just now');
  });

  it('returns minutes ago', () => {
    const past = new Date(Date.now() - 5 * 60 * 1000);
    expect(timeAgo(past.toISOString())).toBe('5m ago');
  });

  it('returns hours ago', () => {
    const past = new Date(Date.now() - 3 * 60 * 60 * 1000);
    expect(timeAgo(past.toISOString())).toBe('3h ago');
  });

  it('returns days ago', () => {
    const past = new Date(Date.now() - 7 * 24 * 60 * 60 * 1000);
    expect(timeAgo(past.toISOString())).toBe('7d ago');
  });
});

describe('getHash', () => {
  it('returns 0 for empty string', () => {
    expect(getHash('')).toBe(0);
  });

  it('returns consistent hash for same input', () => {
    expect(getHash('hello')).toBe(getHash('hello'));
  });

  it('returns different hashes for different inputs', () => {
    expect(getHash('repo-one')).not.toBe(getHash('repo-two'));
  });

  it('returns positive number', () => {
    expect(getHash('test')).toBeGreaterThanOrEqual(0);
  });

  it('returns deterministic hash for numeric string', () => {
    expect(getHash('42')).toBe(getHash('42'));
  });

  it('handles long strings', () => {
    const long = 'a'.repeat(1000);
    expect(getHash(long)).toBeGreaterThan(0);
  });
});
