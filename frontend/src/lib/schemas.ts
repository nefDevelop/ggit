import { z } from 'zod';

export const RepositorySchema = z.object({
  name: z.string()
    .min(1, 'El nombre es obligatorio')
    .max(100, 'El nombre es demasiado largo')
    .regex(/^[a-zA-Z0-9._-]+$/, 'El nombre contiene caracteres no válidos'),
  description: z.string().max(255, 'La descripción es demasiado larga').optional(),
  website: z.string().url('Debe ser una URL válida').or(z.literal('')).optional(),
  wiki_folder_name: z.string().min(1, 'El nombre de la carpeta wiki es obligatorio').default('docs'),
  private: z.boolean().default(false),
  unlisted: z.boolean().default(false),
  enable_issues: z.boolean().default(true),
  allow_public_issues: z.boolean().default(false),
  enable_external_tracker: z.boolean().default(false),
  external_tracker_url: z.string().url('URL inválida').or(z.literal('')).optional(),
  enable_pulls: z.boolean().default(true),
  pulls_ignore_whitespace: z.boolean().default(false),
  pulls_allow_rebase: z.boolean().default(false),
  pulls_allow_manual_merge: z.boolean().default(false),
  enable_wiki: z.boolean().default(true),
  allow_public_wiki: z.boolean().default(false),
  enable_external_wiki: z.boolean().default(false),
  external_wiki_url: z.string().url('URL inválida').or(z.literal('')).optional(),
  archived: z.boolean().default(false),
  default_branch: z.string().optional(),
});

export type RepositorySettings = z.infer<typeof RepositorySchema>;
