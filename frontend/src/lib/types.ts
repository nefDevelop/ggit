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
  // Not nullable on purpose. The [username]/[reponame] layout renders the page
  // through children() inside its {:else if repoData} branch, so a child can
  // never be mounted while this is null. Typing it as nullable only produced
  // "possibly null" noise at every read, and the obvious ways to quiet that,
  // a non-null assertion or any, would have moved the same trap to runtime.
  data: RepoData;
  full: any;
  entries: any[];
  stats: RepoStats | null;
  languages: any;
  loading: boolean;
}
