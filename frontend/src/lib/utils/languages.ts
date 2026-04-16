export const languageColors: Record<string, string> = {
  "Go": "#00ADD8",
  "TypeScript": "#3178C6",
  "JavaScript": "#f1e05a",
  "Svelte": "#FF3E00",
  "Python": "#3572A5",
  "Rust": "#dea584",
  "C++": "#f34b7d",
  "C": "#555555",
  "HTML": "#e34c26",
  "CSS": "#563d7c",
  "Vue": "#41b883",
  "PHP": "#4F5D95",
  "Ruby": "#701516",
  "Java": "#b07219",
  "Shell": "#89e051",
  "Makefile": "#427819",
  "Dockerignore": "#384d54",
  "Dockerfile": "#384d54",
  "JSON": "#292929",
  "Markdown": "#083fa1",
  "Yaml": "#cb171e"
};

export function getLanguageColor(lang: string | undefined): string {
  if (!lang) return "var(--text-muted)";
  return languageColors[lang] || "#8b949e"; // Fallback color
}
