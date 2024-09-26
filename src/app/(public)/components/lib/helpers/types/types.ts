import type { IconName } from '$lib/types';

export type Filters =
    | Partial<{
          [key in IconName]: {
              name: IconName;
              isActive: boolean;
          };
      }>
    | undefined;
