import { Box, Chip, useMediaQuery } from '@mui/material';
import { useCallback, useEffect, useState } from 'react';
import { createIconFromString, createIconOptions } from '../../lib/helpers/icons';
import { useLocalStorage } from '$lib/hooks/useLocalStorage';
import type { Filters } from '../../lib/helpers/types/types';
import { DEFAULT_FILTERS } from '../../lib/hooks/useFilteredPoolEvents';

type Props = {};
const Filter = ({}: Props) => {
    const [filters, setFilters] = useLocalStorage<Filters>('filters', DEFAULT_FILTERS);
    const [toggleState, setToggleState] = useState<Filters>(filters);
    const isMobile = useMediaQuery('(max-width:633px)');

    const chipOptions = createIconOptions().map((option) => ({
        ...option,
        isActive: toggleState?.[option.iconName]?.isActive ?? false,
    }));
    const handleClick = useCallback<(option: (typeof chipOptions)[number]) => void>((option) => {
        setToggleState((prev) => {
            return { ...prev, [option.iconName]: { isActive: !option.isActive } };
        });
    }, []);

    useEffect(() => {
        if (toggleState) {
            setFilters(toggleState);
        }
    }, [toggleState]);
    return (
        <Box
            sx={{
                display: 'flex',
                flexWrap: 'wrap',
                marginBlock: '0.5rem',
                gap: isMobile ? '0.4rem' : '0.5rem',
            }}
        >
            {chipOptions.map((option) => (
                <Chip
                    key={option.label}
                    variant={option.isActive ? 'filled' : 'outlined'}
                    label={isMobile ? undefined : option.label}
                    color="primary"
                    icon={createIconFromString(
                        option.iconName,
                        option.isActive ? 'black' : 'primary',
                        undefined,
                        isMobile ? false : true
                    )}
                    onClick={() => handleClick(option)}
                    sx={{
                        maxWidth: 'fit-content',
                        '.MuiChip-label': {
                            display: isMobile ? 'none' : '',
                        },
                        paddingInline: isMobile ? '0.5rem' : '',
                    }}
                />
            ))}
        </Box>
    );
};

export default Filter;
