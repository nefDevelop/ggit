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
  enable_issues: z.boolean().default(true),
  enable_pulls: z.boolean().default(true),
  enable_wiki: z.boolean().default(true),
  archived: z.boolean().default(false),
});

export type RepositorySettings = z.infer<typeof RepositorySchema>;
