import { Box, Chip } from '@mui/material';
import { useCallback, useEffect, useState } from 'react';
import { createIconFromString, createIconOptions } from '../../lib/helpers/icons';
import { useLocalStorage } from '$lib/hooks/useLocalStorage';
import type { Filters } from '../../lib/helpers/types/types';
import { DEFAULT_FILTERS } from '../../lib/hooks/useFilteredPoolEvents';

type Props = {};
const Filter = ({}: Props) => {
    const [filters, setFilters] = useLocalStorage<Filters>('filters', DEFAULT_FILTERS);
    const [toggleState, setToggleState] = useState<Filters>(filters);

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
                display: 'grid',
                placeContent: 'center',
                marginBlock: '0.5rem',
                gridTemplateColumns: 'repeat(auto-fit, minmax(8.78rem, max-content))',
            }}
        >
            {chipOptions.map((option) => (
                <Chip
                    key={option.label}
                    variant={option.isActive ? 'filled' : 'outlined'}
                    label={option.label}
                    color="primary"
                    icon={createIconFromString(option.iconName, option.isActive ? 'secondary' : 'primary')}
                    onClick={() => handleClick(option)}
                    sx={{ maxWidth: 'fit-content' }}
                />
            ))}
        </Box>
    );
};

export default Filter;
