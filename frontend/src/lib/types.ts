export interface RepoOwner {
  username: string;
  [key: string]: any;
}

export interface RepoData {
  owner: RepoOwner;
  name: string;
  id: number;
  default_branch?: string;
  private?: boolean;
  description?: string;
  [key: string]: any;
}

export interface RepoStats {
  commits: number;
  branches: number;
  tags: number;
  releases: number;
}

export interface RepoContext {
  data: RepoData | null;
  full: any;
  entries: any[];
  stats: RepoStats | null;
  languages: any;
  loading: boolean;
}
