import type { IconName } from '$lib/types';

export type Filters =
    | Partial<{
          [key in IconName]: {
              isActive: boolean;
          };
      }>
    | undefined;
