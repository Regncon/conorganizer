'use client';
import type { EventDays } from '$app/(public)/page';
import { useSetCustomCssVariable } from '$lib/hooks/useSetCustomCssVariable';
import { type RemoveProperty, setCustomVariable } from '$lib/libClient';
import { Box, Typography, Link, type SxProps } from '@mui/material';
import { useEffect, useRef } from 'react';
import { unmountComponentAtNode } from 'react-dom';

type Props = {
    eventDays: EventDays;
};
const sxDayTypography: SxProps = {
    maxWidth: '5rem',
    minHeight: '4rem',
    display: 'grid',
    placeItems: 'center',
    padding: '0.5em',
};

const DaysHeader = ({ eventDays }: Props) => {
    const ref = useSetCustomCssVariable({ '--scroll-margin-top': 'height' });

    return (
        <>
            <Box
                component="header"
                sx={{
                    position: 'sticky',
                    top: '0',
                    backgroundColor: 'background.paper',
                    padding: '0.5rem',
                    zIndex: 1,
                }}
                ref={ref}
            >
                <Box sx={{ display: 'grid', placeContent: 'end', marginInline: '2rem', marginBlock: '0.5rem' }}>
                    FILTER
                </Box>
                <Box
                    sx={{
                        display: 'grid',
                        gridTemplateColumns: 'repeat(4, max-content)',
                        placeContent: 'end',
                        placeItems: 'center',
                    }}
                >
                    {Object.values(eventDays).map((day) => (
                        <Typography key={day} component={Link} href={`#${day}`} variant="h5" sx={sxDayTypography}>
                            {day}
                        </Typography>
                    ))}
                </Box>
            </Box>
        </>
    );
};

export default DaysHeader;
