import { browser } from '$app/environment';

export interface UserPreferences {
  repoTabs: {
    pr: boolean;
    issue: boolean;
    wiki: boolean;
    branch: boolean;
    commits: boolean;
    adr: boolean;
  };
  showSidebarOrgs: boolean;
}

const DEFAULT_PREFERENCES: UserPreferences = {
  repoTabs: {
    pr: true,
    issue: true,
    wiki: true,
    branch: true,
    commits: true,
    adr: true,
  },
  showSidebarOrgs: true
};

function createPreferences() {
  let initial = DEFAULT_PREFERENCES;
  if (browser) {
    const saved = localStorage.getItem('ggit_user_preferences');
    if (saved) {
      try {
        const parsed = JSON.parse(saved);
        initial = { 
          ...DEFAULT_PREFERENCES, 
          ...parsed,
          repoTabs: { ...DEFAULT_PREFERENCES.repoTabs, ...parsed.repoTabs }
        };
      } catch (e) {
        console.error('Failed to parse preferences', e);
      }
    }
  }

  let prefs = $state<UserPreferences>(initial);

  return {
    get current() { return prefs; },
    update(newPrefs: Partial<UserPreferences>) {
      Object.assign(prefs, newPrefs);
      if (browser) {
        localStorage.setItem('ggit_user_preferences', JSON.stringify(prefs));
      }
    },
    toggleRepoTab(tab: keyof UserPreferences['repoTabs']) {
      prefs.repoTabs[tab] = !prefs.repoTabs[tab];
      if (browser) {
        localStorage.setItem('ggit_user_preferences', JSON.stringify(prefs));
      }
    },
    save() {
      if (browser) {
        localStorage.setItem('ggit_user_preferences', JSON.stringify(prefs));
      }
    }
  };
}

export const userPrefs = createPreferences();
