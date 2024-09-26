import type { IconName } from '$lib/types';
import { Box, Chip, Skeleton } from '@mui/material';
import { useCallback, useEffect, useState } from 'react';
import { createIconFromString, createIconOptions } from '../../lib/helpers/icons';
import { useLocalStorage } from '../../../../../lib/hooks/useLocalStorage';
import { useRouter } from 'next/router';
import useSafeRouterReplace from '../../lib/hooks/useSafeRouterReplace';
import debounce from '$lib/debounce';
import useReadLocalStorage from '$lib/hooks/useReadLocalStorage';

type Filters =
    | Partial<{
          [key in IconName]: {
              name: IconName;
              isActive: boolean;
          };
      }>
    | undefined;
type Props = {};
const Filter = ({}: Props) => {
    const { setQuery } = useSafeRouterReplace();
    const [, setFilters] = useLocalStorage<Filters>('filters', undefined);
    const filters = useReadLocalStorage<Filters>('filters') ?? undefined;
    const [toggleState, setToggleState] = useState<Filters>(filters);
    const chipOptions = createIconOptions().map((option) => ({
        ...option,
        isActive: toggleState?.[option.iconName]?.isActive ?? undefined,
    }));
    const handleClick = useCallback<(option: (typeof chipOptions)[number]) => void>((option) => {
        setToggleState((prev) => {
            return { ...prev, [option.iconName]: { name: option.iconName, isActive: !option.isActive } };
        });
    }, []);

    const debounceRouter = useCallback(
        debounce(() => {
            if (toggleState) {
                const test = Object.values(toggleState).map((ts) => ({
                    key: ts.name as string,
                    value: `${ts.isActive ?? 'false'}`,
                }));
                if (!test.every((t) => t.value === 'false')) {
                    setQuery(test);
                }
            }
        }, 500),
        [toggleState]
    );

    useEffect(() => {
        if (toggleState) {
            setFilters(toggleState);

            debounceRouter().catch((e) => console.warn(e));
        }
    }, [toggleState]);
    return (
        <Box
            sx={{
                display: 'grid',
                placeContent: 'center',
                marginBlock: '0.5rem',
                gridTemplateColumns: 'repeat(auto-fit, minmax(8.78rem, max-content))',
            }}
        >
            {chipOptions.map((option) =>
                option.isActive !== undefined ?
                    <Chip
                        key={option.label}
                        variant={option.isActive ? 'filled' : 'outlined'}
                        label={option.label}
                        color="primary"
                        icon={createIconFromString(option.iconName, option.isActive ? 'secondary' : 'primary')}
                        onClick={(e) => handleClick(option)}
                        sx={{ maxWidth: 'fit-content' }}
                    />
                :   <Skeleton
                        key={option.label}
                        variant="rounded"
                        width={145}
                        height={32}
                        sx={{ marginInline: '0.1rem' }}
                    />
            )}
        </Box>
    );
};

export default Filter;
